package ui

import (
	"fmt"
	"time"

	"tuitodo/internal/domain"
	"tuitodo/internal/store"

	"github.com/rivo/tview"
)

func (app *App) openMessageModal(title, message string) {
	form := newStyledForm(title)
	addStyledMessage(form, message)
	form.AddButton("OK", app.modals.Close)
	focusFormButtons(form)
	app.modals.OpenPage(form)
}

func (app *App) openHelpModal() {
	app.openMessageModal("Atalhos", "Setas escolhem o painel  Enter entra  Esc volta\n1 A fazer  2 Concluídos\na nova tarefa  c categoria\nEm Tarefas: e editar  d apagar  espaço concluir\nEm Subtarefas: a adicionar  e editar  d apagar  espaço concluir\nEm Concluídos: escolha o dia; espaço reabre\nq sair")
}

func (app *App) openNewCategoryModal() {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Nova categoria ")
	styleForm(form)
	nameInput := addStyledInputField(form, "Nome", "")
	form.AddButton("Salvar", func() {
		name := stringsTrim(nameInput.GetText())
		if err := domain.ValidateTitle(name); err != nil {
			app.modals.ShowError("O nome da categoria não pode ser vazio.")
			return
		}
		_, err := app.categoryRepository.InsertCategory(domain.Category{
			Name:      name,
			CreatedAt: time.Now(),
		})
		if err != nil {
			app.modals.ShowError("Não foi possível criar a categoria. O nome já existe?")
			return
		}
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.modals.ShowError(err.Error())
		}
	})
	form.AddButton("Cancelar", func() { app.modals.Close() })
	app.modals.OpenPage(form)
}

func (app *App) openNewParentTaskModal() {
	if app.showingCompleted {
		return
	}
	if len(app.categories) == 0 {
		app.modals.ShowError("Crie uma categoria primeiro (c).")
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
	form.SetBorder(true).SetTitle(" Nova tarefa ")
	styleForm(form)
	titleInput := addStyledInputField(form, "Título", "")
	categoryPicker := addStyledCategoryPicker(form, "Categoria", names, selected)
	form.AddButton("Salvar", func() {
		title := stringsTrim(titleInput.GetText())
		if err := domain.ValidateTitle(title); err != nil {
			app.modals.ShowError("O título não pode ser vazio.")
			return
		}
		categoryIndex := categoryPicker.GetCurrentItem()
		if categoryIndex < 0 || categoryIndex >= len(app.categories) {
			app.modals.ShowError("Escolha uma categoria.")
			return
		}
		categoryID := app.categories[categoryIndex].ID
		parentTask, err := app.taskRepository.InsertParentTask(domain.Task{
			CategoryID: &categoryID,
			Title:      title,
			CreatedAt:  time.Now(),
		})
		if err != nil {
			app.modals.ShowError(err.Error())
			return
		}
		app.selectedParentID = parentTask.ID
		app.showingCompleted = false
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.modals.ShowError(err.Error())
		}
	})
	form.AddButton("Cancelar", func() { app.modals.Close() })
	app.modals.OpenPage(form)
}

func (app *App) openNewSubtaskModal() {
	if app.showingCompleted {
		return
	}
	parentTask, ok := app.selectedParentTask()
	if !ok {
		app.modals.ShowError("Selecione uma tarefa pai.")
		return
	}
	if parentTask.IsCompleted() {
		app.modals.ShowError("Reabra a tarefa antes de adicionar subtarefas.")
		return
	}
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Nova subtarefa ")
	styleForm(form)
	titleInput := addStyledInputField(form, "Título", "")
	form.AddButton("Salvar", func() {
		title := stringsTrim(titleInput.GetText())
		subtask := domain.Task{ParentID: &parentTask.ID, Title: title, CreatedAt: time.Now()}
		if _, err := app.taskRepository.InsertSubtask(subtask, parentTask); err != nil {
			if err == domain.ErrEmptyTitle {
				app.modals.ShowError("O título não pode ser vazio.")
				return
			}
			app.modals.ShowError(err.Error())
			return
		}
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.modals.ShowError(err.Error())
		}
	})
	form.AddButton("Cancelar", func() { app.modals.Close() })
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
	form.SetBorder(true).SetTitle(" Editar título ")
	styleForm(form)
	titleInput := addStyledInputField(form, "Título", task.Title)
	form.AddButton("Salvar", func() {
		title := stringsTrim(titleInput.GetText())
		if err := app.taskRepository.UpdateTaskTitle(task.ID, title); err != nil {
			if err == domain.ErrEmptyTitle {
				app.modals.ShowError("O título não pode ser vazio.")
				return
			}
			app.modals.ShowError(err.Error())
			return
		}
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.modals.ShowError(err.Error())
		}
	})
	form.AddButton("Cancelar", func() { app.modals.Close() })
	app.modals.OpenPage(form)
}

func (app *App) openRenameCategoryModal() {
	index := app.categoryList.GetCurrentItem()
	if index <= 0 || index-1 >= len(app.categories) {
		app.modals.ShowError("Selecione uma categoria para renomear.")
		return
	}
	category := app.categories[index-1]
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Renomear categoria ")
	styleForm(form)
	nameInput := addStyledInputField(form, "Nome", category.Name)
	form.AddButton("Salvar", func() {
		name := stringsTrim(nameInput.GetText())
		if err := app.categoryRepository.RenameCategory(category.ID, name); err != nil {
			app.modals.ShowError("Não foi possível renomear a categoria.")
			return
		}
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.modals.ShowError(err.Error())
		}
	})
	form.AddButton("Cancelar", func() { app.modals.Close() })
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
		app.modals.ShowError("Selecione uma tarefa para apagar.")
		return
	}
	app.confirmDeleteTask(parentTask, true)
}

func (app *App) confirmDeleteTask(target domain.Task, isParent bool) {
	app.openConfirmModal("Apagar", fmt.Sprintf("Apagar \"%s\"?", target.Title), "Apagar", func() {
		if err := app.taskRepository.DeleteTaskByID(target.ID); err != nil {
			app.modals.Close()
			app.modals.ShowError(err.Error())
			return
		}
		if isParent && target.ID == app.selectedParentID {
			app.selectedParentID = 0
		}
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.modals.ShowError(err.Error())
		}
	})
}

func (app *App) openDeleteCategoryModal() {
	index := app.categoryList.GetCurrentItem()
	if index <= 0 || index-1 >= len(app.categories) {
		app.modals.ShowError("Selecione uma categoria para apagar.")
		return
	}
	category := app.categories[index-1]
	app.openConfirmModal("Apagar", fmt.Sprintf("Apagar a categoria \"%s\"?", category.Name), "Apagar", func() {
		if err := app.categoryRepository.DeleteCategory(category.ID); err != nil {
			app.modals.Close()
			if err == store.ErrCategoryInUse {
				app.modals.ShowError("Essa categoria ainda tem tarefas.")
				return
			}
			app.modals.ShowError(err.Error())
			return
		}
		app.selectedCategoryID = nil
		app.modals.Close()
		if err := app.reload(); err != nil {
			app.modals.ShowError(err.Error())
		}
	})
}

func (app *App) openConfirmModal(title, message, confirmLabel string, onConfirm func()) {
	form := newStyledForm(title)
	addStyledMessage(form, message)
	form.AddButton(confirmLabel, onConfirm)
	form.AddButton("Cancelar", app.modals.Close)
	focusFormButtons(form)
	app.modals.OpenCompact(form)
}
