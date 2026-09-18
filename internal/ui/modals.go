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
	app.openMessageModal("Atalhos", "Tab foca os painéis\n1 Pendentes  2 Concluídos\na nova tarefa  n nova sub  c categoria\ne editar  d apagar  espaço concluir/reabrir\nq sair  Esc fecha modal")
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
	focus := app.application.GetFocus()
	if focus == app.categoryList {
		app.openRenameCategoryModal()
		return
	}
	entry, ok := app.focusedDetailEntry()
	if !ok {
		parentTask, found := app.selectedParentTask()
		if !found {
			app.modals.ShowError("Nada selecionado para editar.")
			return
		}
		entry = detailEntry{isParent: true, task: parentTask}
	}
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Editar título ")
	styleForm(form)
	titleInput := addStyledInputField(form, "Título", entry.task.Title)
	form.AddButton("Salvar", func() {
		title := stringsTrim(titleInput.GetText())
		if err := app.taskRepository.UpdateTaskTitle(entry.task.ID, title); err != nil {
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
	focus := app.application.GetFocus()
	if focus == app.categoryList {
		app.openDeleteCategoryModal()
		return
	}
	parentTask, ok := app.selectedParentTask()
	if !ok {
		app.modals.ShowError("Selecione uma tarefa para apagar.")
		return
	}
	entry, hasDetail := app.focusedDetailEntry()
	target := parentTask
	label := parentTask.Title
	if hasDetail && !entry.isParent {
		target = entry.task
		label = entry.task.Title
	}
	app.openConfirmModal("Apagar", fmt.Sprintf("Apagar \"%s\"?", label), "Apagar", func() {
		if err := app.taskRepository.DeleteTaskByID(target.ID); err != nil {
			app.modals.Close()
			app.modals.ShowError(err.Error())
			return
		}
		if target.ID == app.selectedParentID {
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
