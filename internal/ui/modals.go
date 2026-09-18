package ui

import (
	"fmt"
	"time"

	"tuitodo/internal/domain"
	"tuitodo/internal/store"

	"github.com/rivo/tview"
)

func (app *App) openMessageModal(title, message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.closeModal()
		})
	modal.SetTitle(" " + title + " ")
	app.rememberFocus()
	app.modalOpen = true
	app.pages.AddPage("modal", modal, true, true)
	app.application.SetFocus(modal)
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
			app.showError("O nome da categoria não pode ser vazio.")
			return
		}
		_, err := app.categoryRepository.InsertCategory(domain.Category{
			Name:      name,
			CreatedAt: time.Now(),
		})
		if err != nil {
			app.showError("Não foi possível criar a categoria. O nome já existe?")
			return
		}
		app.closeModal()
		if err := app.reload(); err != nil {
			app.showError(err.Error())
		}
	})
	form.AddButton("Cancelar", func() { app.closeModal() })
	app.openPage(form)
}

func (app *App) openNewParentTaskModal() {
	if len(app.categories) == 0 {
		app.showError("Crie uma categoria primeiro (c).")
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
	categoryDropDown := addStyledDropDown(form, "Categoria", names, selected)
	form.AddButton("Salvar", func() {
		title := stringsTrim(titleInput.GetText())
		if err := domain.ValidateTitle(title); err != nil {
			app.showError("O título não pode ser vazio.")
			return
		}
		categoryIndex, _ := categoryDropDown.GetCurrentOption()
		if categoryIndex < 0 || categoryIndex >= len(app.categories) {
			app.showError("Escolha uma categoria.")
			return
		}
		categoryID := app.categories[categoryIndex].ID
		parentTask, err := app.taskRepository.InsertParentTask(domain.Task{
			CategoryID: &categoryID,
			Title:      title,
			CreatedAt:  time.Now(),
		})
		if err != nil {
			app.showError(err.Error())
			return
		}
		app.selectedParentID = parentTask.ID
		app.showingCompleted = false
		app.closeModal()
		if err := app.reload(); err != nil {
			app.showError(err.Error())
		}
	})
	form.AddButton("Cancelar", func() { app.closeModal() })
	app.openPage(form)
}

func (app *App) openNewSubtaskModal() {
	parentTask, ok := app.selectedParentTask()
	if !ok {
		app.showError("Selecione uma tarefa pai.")
		return
	}
	if parentTask.IsCompleted() {
		app.showError("Reabra a tarefa antes de adicionar subtarefas.")
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
				app.showError("O título não pode ser vazio.")
				return
			}
			app.showError(err.Error())
			return
		}
		app.closeModal()
		if err := app.reload(); err != nil {
			app.showError(err.Error())
		}
	})
	form.AddButton("Cancelar", func() { app.closeModal() })
	app.openPage(form)
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
			app.showError("Nada selecionado para editar.")
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
				app.showError("O título não pode ser vazio.")
				return
			}
			app.showError(err.Error())
			return
		}
		app.closeModal()
		if err := app.reload(); err != nil {
			app.showError(err.Error())
		}
	})
	form.AddButton("Cancelar", func() { app.closeModal() })
	app.openPage(form)
}

func (app *App) openRenameCategoryModal() {
	index := app.categoryList.GetCurrentItem()
	if index <= 0 || index-1 >= len(app.categories) {
		app.showError("Selecione uma categoria para renomear.")
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
			app.showError("Não foi possível renomear a categoria.")
			return
		}
		app.closeModal()
		if err := app.reload(); err != nil {
			app.showError(err.Error())
		}
	})
	form.AddButton("Cancelar", func() { app.closeModal() })
	app.openPage(form)
}

func (app *App) openDeleteModal() {
	focus := app.application.GetFocus()
	if focus == app.categoryList {
		app.openDeleteCategoryModal()
		return
	}
	parentTask, ok := app.selectedParentTask()
	if !ok {
		app.showError("Selecione uma tarefa para apagar.")
		return
	}
	entry, hasDetail := app.focusedDetailEntry()
	target := parentTask
	label := parentTask.Title
	if hasDetail && !entry.isParent {
		target = entry.task
		label = entry.task.Title
	}
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Apagar \"%s\"?", label)).
		AddButtons([]string{"Apagar", "Cancelar"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Apagar" {
				if err := app.taskRepository.DeleteTaskByID(target.ID); err != nil {
					app.closeModal()
					app.showError(err.Error())
					return
				}
				if target.ID == app.selectedParentID {
					app.selectedParentID = 0
				}
			}
			app.closeModal()
			if err := app.reload(); err != nil {
				app.showError(err.Error())
			}
		})
	app.rememberFocus()
	app.modalOpen = true
	app.pages.AddPage("modal", modal, true, true)
	app.application.SetFocus(modal)
}

func (app *App) openDeleteCategoryModal() {
	index := app.categoryList.GetCurrentItem()
	if index <= 0 || index-1 >= len(app.categories) {
		app.showError("Selecione uma categoria para apagar.")
		return
	}
	category := app.categories[index-1]
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Apagar a categoria \"%s\"?", category.Name)).
		AddButtons([]string{"Apagar", "Cancelar"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Apagar" {
				if err := app.categoryRepository.DeleteCategory(category.ID); err != nil {
					app.closeModal()
					if err == store.ErrCategoryInUse {
						app.showError("Essa categoria ainda tem tarefas.")
						return
					}
					app.showError(err.Error())
					return
				}
				app.selectedCategoryID = nil
			}
			app.closeModal()
			if err := app.reload(); err != nil {
				app.showError(err.Error())
			}
		})
	app.rememberFocus()
	app.modalOpen = true
	app.pages.AddPage("modal", modal, true, true)
	app.application.SetFocus(modal)
}
