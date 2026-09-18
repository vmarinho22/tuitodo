package ui

import (
	"time"

	"tuitodo/internal/domain"
	"tuitodo/internal/i18n"
	"tuitodo/internal/store"

	"github.com/rivo/tview"
)

func (app *App) openMessageModal(title, message string) {
	form := newStyledForm(title)
	addStyledMessage(form, message)
	form.AddButton(app.t(i18n.KeyButtonOK), app.modals.Close)
	focusFormButtons(form)
	app.modals.OpenPage(form)
}

func (app *App) openHelpModal() {
	app.openMessageModal(app.t(i18n.KeyHelpTitle), app.t(i18n.KeyHelpBody))
}

func (app *App) openConfigModal() {
	locales := i18n.Locales()
	names := make([]string, len(locales))
	selected := 0
	previous := i18n.LocaleEnUS
	if app.catalog != nil {
		previous = app.catalog.Locale()
	}
	for i, locale := range locales {
		names[i] = locale.NativeLabel()
		if locale == previous {
			selected = i
		}
	}
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" " + app.t(i18n.KeyConfigTitle) + " ")
	styleForm(form)
	picker := addStyledCategoryPicker(form, app.t(i18n.KeyConfigLanguage), names, selected)
	saved := false
	picker.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		if index < 0 || index >= len(locales) {
			return
		}
		app.previewConfigLocale(locales[index], form, picker)
	})
	form.AddButton(app.t(i18n.KeyButtonSave), func() {
		index := picker.GetCurrentItem()
		if index < 0 || index >= len(locales) {
			return
		}
		locale := locales[index]
		if err := app.settingsRepository.SetSetting(store.SettingLocale, string(locale)); err != nil {
			app.showError(err)
			return
		}
		saved = true
		app.catalog.SetLocale(locale)
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.showError(err)
		}
	})
	form.AddButton(app.t(i18n.KeyButtonCancel), app.modals.Close)
	app.modals.onClose = func() {
		if saved || app.catalog == nil || app.catalog.Locale() == previous {
			return
		}
		app.catalog.SetLocale(previous)
		_ = app.reload()
	}
	app.modals.OpenPage(form)
}

func (app *App) previewConfigLocale(locale i18n.Locale, form *tview.Form, picker *categoryPicker) {
	if app.catalog == nil || app.catalog.Locale() == locale {
		return
	}
	app.catalog.SetLocale(locale)
	if err := app.reload(); err != nil {
		app.showError(err)
		return
	}
	form.SetTitle(" " + app.t(i18n.KeyConfigTitle) + " ")
	picker.SetCaption(app.t(i18n.KeyConfigLanguage))
	if form.GetButtonCount() >= 2 {
		form.GetButton(0).SetLabel(app.t(i18n.KeyButtonSave))
		form.GetButton(1).SetLabel(app.t(i18n.KeyButtonCancel))
	}
}

func (app *App) openNewCategoryModal() {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" " + app.t(i18n.KeyFormNewCategory) + " ")
	styleForm(form)
	nameInput := addStyledInputField(form, app.t(i18n.KeyFormName), "")
	form.AddButton(app.t(i18n.KeyButtonSave), func() {
		name := stringsTrim(nameInput.GetText())
		if err := domain.ValidateTitle(name); err != nil {
			app.modals.ShowError(app.t(i18n.KeyEmptyName))
			return
		}
		_, err := app.categoryRepository.InsertCategory(domain.Category{
			Name:      name,
			CreatedAt: time.Now(),
		})
		if err != nil {
			app.modals.ShowError(app.t(i18n.KeyCategoryExists))
			return
		}
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.showError(err)
		}
	})
	form.AddButton(app.t(i18n.KeyButtonCancel), func() { app.modals.Close() })
	app.modals.OpenPage(form)
}

func (app *App) openNewParentTaskModal() {
	if app.showingCompleted {
		return
	}
	if len(app.categories) == 0 {
		app.modals.ShowError(app.t(i18n.KeyCreateCategoryFirst))
		return
	}
	names := make([]string, 0, len(app.categories))
	selected := 0
	for i, category := range app.categories {
		names = append(names, category.Name)
		if app.selectedCategoryID != nil && category.ID == *app.selectedCategoryID {
			selected = i
		}
	}
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" " + app.t(i18n.KeyFormNewTask) + " ")
	styleForm(form)
	titleInput := addStyledInputField(form, app.t(i18n.KeyFormTitle), "")
	categoryPicker := addStyledCategoryPicker(form, app.t(i18n.KeyFormCategory), names, selected)
	form.AddButton(app.t(i18n.KeyButtonSave), func() {
		title := stringsTrim(titleInput.GetText())
		if err := domain.ValidateTitle(title); err != nil {
			app.modals.ShowError(app.t(i18n.KeyEmptyTitle))
			return
		}
		categoryIndex := categoryPicker.GetCurrentItem()
		if categoryIndex < 0 || categoryIndex >= len(app.categories) {
			app.modals.ShowError(app.t(i18n.KeyChooseCategory))
			return
		}
		categoryID := app.categories[categoryIndex].ID
		parentTask, err := app.taskRepository.InsertParentTask(domain.Task{
			CategoryID: &categoryID,
			Title:      title,
			CreatedAt:  time.Now(),
		})
		if err != nil {
			app.showError(err)
			return
		}
		app.selectedParentID = parentTask.ID
		app.showingCompleted = false
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.showError(err)
		}
	})
	form.AddButton(app.t(i18n.KeyButtonCancel), func() { app.modals.Close() })
	app.modals.OpenPage(form)
}

func (app *App) openNewSubtaskModal() {
	if app.showingCompleted {
		return
	}
	parentTask, ok := app.selectedParentTask()
	if !ok {
		app.modals.ShowError(app.t(i18n.KeySelectParent))
		return
	}
	if parentTask.IsCompleted() {
		app.modals.ShowError(app.t(i18n.KeyReopenBeforeSub))
		return
	}
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" " + app.t(i18n.KeyFormNewSubtask) + " ")
	styleForm(form)
	titleInput := addStyledInputField(form, app.t(i18n.KeyFormTitle), "")
	form.AddButton(app.t(i18n.KeyButtonSave), func() {
		title := stringsTrim(titleInput.GetText())
		subtask := domain.Task{ParentID: &parentTask.ID, Title: title, CreatedAt: time.Now()}
		if _, err := app.taskRepository.InsertSubtask(subtask, parentTask); err != nil {
			if err == domain.ErrEmptyTitle {
				app.modals.ShowError(app.t(i18n.KeyEmptyTitle))
				return
			}
			app.showError(err)
			return
		}
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.showError(err)
		}
	})
	form.AddButton(app.t(i18n.KeyButtonCancel), func() { app.modals.Close() })
	app.modals.OpenPage(form)
}

func (app *App) openEditModal() {
	if app.showingCompleted {
		if app.selectedPane == paneCategories {
			app.openRenameCategoryModal()
		}
		return
	}
	switch app.selectedPane {
	case paneCategories:
		app.openRenameCategoryModal()
	case paneTasks:
		parentTask, ok := app.selectedParentTask()
		if !ok {
			return
		}
		app.openEditTitleModal(parentTask)
	case paneDetail:
		entry, ok := app.focusedDetailEntry()
		if !ok {
			return
		}
		app.openEditTitleModal(entry.task)
	}
}

func (app *App) openEditTitleModal(task domain.Task) {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" " + app.t(i18n.KeyFormEditTitle) + " ")
	styleForm(form)
	titleInput := addStyledInputField(form, app.t(i18n.KeyFormTitle), task.Title)
	form.AddButton(app.t(i18n.KeyButtonSave), func() {
		title := stringsTrim(titleInput.GetText())
		if err := app.taskRepository.UpdateTaskTitle(task.ID, title); err != nil {
			if err == domain.ErrEmptyTitle {
				app.modals.ShowError(app.t(i18n.KeyEmptyTitle))
				return
			}
			app.showError(err)
			return
		}
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.showError(err)
		}
	})
	form.AddButton(app.t(i18n.KeyButtonCancel), func() { app.modals.Close() })
	app.modals.OpenPage(form)
}

func (app *App) openRenameCategoryModal() {
	index := app.categoryList.GetCurrentItem()
	if index <= 0 || index-1 >= len(app.categories) {
		app.modals.ShowError(app.t(i18n.KeySelectCategoryRename))
		return
	}
	category := app.categories[index-1]
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" " + app.t(i18n.KeyFormRenameCategory) + " ")
	styleForm(form)
	nameInput := addStyledInputField(form, app.t(i18n.KeyFormName), category.Name)
	form.AddButton(app.t(i18n.KeyButtonSave), func() {
		name := stringsTrim(nameInput.GetText())
		if err := app.categoryRepository.RenameCategory(category.ID, name); err != nil {
			app.modals.ShowError(app.t(i18n.KeyRenameFailed))
			return
		}
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.showError(err)
		}
	})
	form.AddButton(app.t(i18n.KeyButtonCancel), func() { app.modals.Close() })
	app.modals.OpenPage(form)
}

func (app *App) openDeleteModal() {
	if app.showingCompleted {
		if app.paneActive && app.selectedPane == paneCategories {
			app.openDeleteCategoryModal()
		}
		return
	}
	if app.paneActive && app.selectedPane == paneCategories {
		app.openDeleteCategoryModal()
		return
	}
	if app.paneActive && app.selectedPane == paneDetail {
		entry, ok := app.focusedDetailEntry()
		if !ok {
			return
		}
		app.confirmDeleteTask(entry.task, false)
		return
	}
	parentTask, ok := app.selectedParentTask()
	if !ok {
		app.modals.ShowError(app.t(i18n.KeySelectTaskDelete))
		return
	}
	app.confirmDeleteTask(parentTask, true)
}

func (app *App) confirmDeleteTask(target domain.Task, isParent bool) {
	app.openConfirmModal(app.t(i18n.KeyButtonDelete), app.t(i18n.KeyDeleteTaskConfirm, target.Title), app.t(i18n.KeyButtonDelete), func() {
		if err := app.taskRepository.DeleteTaskByID(target.ID); err != nil {
			app.modals.Close()
			app.showError(err)
			return
		}
		if isParent && target.ID == app.selectedParentID {
			app.selectedParentID = 0
		}
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.showError(err)
		}
	})
}

func (app *App) openDeleteCategoryModal() {
	index := app.categoryList.GetCurrentItem()
	if index <= 0 || index-1 >= len(app.categories) {
		app.modals.ShowError(app.t(i18n.KeySelectCategoryDelete))
		return
	}
	category := app.categories[index-1]
	app.openConfirmModal(app.t(i18n.KeyButtonDelete), app.t(i18n.KeyDeleteCategoryConfirm, category.Name), app.t(i18n.KeyButtonDelete), func() {
		if err := app.categoryRepository.DeleteCategory(category.ID); err != nil {
			app.modals.Close()
			app.showError(err)
			return
		}
		app.selectedCategoryID = nil
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.showError(err)
		}
	})
}

func (app *App) openConfirmModal(title, message, confirmLabel string, onConfirm func()) {
	form := newStyledForm(title)
	addStyledMessage(form, message)
	form.AddButton(confirmLabel, onConfirm)
	form.AddButton(app.t(i18n.KeyButtonCancel), app.modals.Close)
	focusFormButtons(form)
	app.modals.OpenCompact(form)
}
