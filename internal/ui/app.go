package ui

import (
	"fmt"
	"strings"
	"time"

	"tuitodo/internal/domain"
	"tuitodo/internal/store"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const allCategoriesLabel = "Todas"

const (
	pendingSubtaskShortcuts   = " a adicionar  e editar  d apagar  espaço concluir"
	completedSubtaskShortcuts = " espaço reabrir"
)

var actionIDs = []string{"new-task", "new-category", "delete", "help", "quit"}

type App struct {
	application        *tview.Application
	pages              *tview.Pages
	tasksPane          *tview.Flex
	taskModeLine       *tview.TextView
	taskList           *tview.List
	categoryList       *tview.List
	subtasksPane       *tview.Flex
	parentTitle        *tview.TextView
	subtaskShortcuts   *tview.TextView
	detailList         *tview.List
	actionsBar         *tview.TextView
	taskRepository     store.TaskRepository
	categoryRepository store.CategoryRepository
	showingCompleted   bool
	selectedCategoryID *int64
	selectedParentID   int64
	selectedDay        time.Time
	completedDays      []domain.CompletedDay
	taskListEntries    []taskListEntry
	detailEntries      []detailEntry
	categories         []domain.Category
	reloading          bool
	modals             *ModalService
	actionsNavigating  bool
	selectedPane       paneID
	lastTopPane        paneID
	paneActive         bool
	selectingPane      bool
}

type taskListEntry struct {
	isDayHeader bool
	day         time.Time
	parentTask  domain.Task
}

type detailEntry struct {
	isParent bool
	task     domain.Task
}

func Run(taskRepository store.TaskRepository, categoryRepository store.CategoryRepository) error {
	app := &App{
		application:        tview.NewApplication(),
		taskRepository:     taskRepository,
		categoryRepository: categoryRepository,
	}
	app.build()
	if err := app.reload(); err != nil {
		return err
	}
	app.application.EnableMouse(true).SetRoot(app.pages, true)
	app.selectPane(paneTasks)
	return app.application.Run()
}

func (app *App) build() {
	app.taskModeLine = tview.NewTextView().SetDynamicColors(true).SetWrap(false)
	app.taskModeLine.SetText(taskModeLine(false))
	app.taskList = tview.NewList().ShowSecondaryText(true).SetHighlightFullLine(true).SetWrapAround(true)
	app.tasksPane = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(app.taskModeLine, 1, 0, false).
		AddItem(app.taskList, 0, 1, true)
	app.tasksPane.SetBorder(true).SetTitle(" Tarefas ")
	app.bindPaneFocus(paneTasks, app.tasksPane.Box)
	app.bindPaneFocus(paneTasks, app.taskList.Box)
	app.taskModeLine.SetScrollable(false)
	app.taskModeLine.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftDown {
			app.selectedPane = paneTasks
			app.lastTopPane = paneTasks
			app.enterPane()
			return tview.MouseConsumed, nil
		}
		return action, event
	})
	app.categoryList = tview.NewList().ShowSecondaryText(false).SetHighlightFullLine(true).SetWrapAround(true)
	app.categoryList.SetBorder(true).SetTitle(" Categorias ")
	app.bindPaneFocus(paneCategories, app.categoryList.Box)
	app.parentTitle = tview.NewTextView().SetText("Selecione uma tarefa").SetWrap(true)
	app.parentTitle.SetTextColor(tview.Styles.SecondaryTextColor)
	app.detailList = tview.NewList().ShowSecondaryText(false).SetHighlightFullLine(true).SetWrapAround(true)
	app.subtaskShortcuts = tview.NewTextView().SetText(pendingSubtaskShortcuts).SetWrap(false)
	app.subtaskShortcuts.SetTextColor(tview.Styles.SecondaryTextColor)
	app.subtasksPane = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(app.parentTitle, 1, 0, false).
		AddItem(app.detailList, 0, 1, true).
		AddItem(app.subtaskShortcuts, 1, 0, false)
	app.subtasksPane.SetBorder(true).SetTitle(" Subtarefas ")
	app.bindPaneFocus(paneDetail, app.subtasksPane.Box)
	app.bindPaneFocus(paneDetail, app.detailList.Box)
	app.actionsBar = tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetWrap(false)
	app.actionsBar.SetBorder(true).SetTitle(" Ações ")
	app.bindPaneFocus(paneActions, app.actionsBar.Box)
	app.actionsBar.SetText(` ["new-task"]a tarefa[""]   ["new-category"]c categoria[""]   ["delete"]d apagar[""]   ["help"]? atalhos[""]   ["quit"]q sair[""]`)
	app.actionsBar.SetScrollable(false)
	app.actionsBar.SetHighlightedFunc(func(added, removed, remaining []string) {
		if app.actionsNavigating || len(added) == 0 {
			return
		}
		app.runAction(added[0])
	})
	app.actionsBar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft, tcell.KeyUp:
			app.cycleAction(-1)
			return nil
		case tcell.KeyRight, tcell.KeyDown:
			app.cycleAction(1)
			return nil
		case tcell.KeyEnter:
			app.activateHighlightedAction()
			return nil
		}
		return event
	})

	bindListVimKeys(app.taskList)
	bindListVimKeys(app.categoryList)
	bindListVimKeys(app.detailList)

	app.taskList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		app.onTaskListChanged(index)
	})
	app.categoryList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		app.onCategoryChanged(index)
	})

	left := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(app.tasksPane, 0, 3, true).
		AddItem(app.categoryList, 0, 2, false)
	body := tview.NewFlex().
		AddItem(left, 0, 2, true).
		AddItem(app.subtasksPane, 0, 3, false)
	root := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(body, 0, 1, true).
		AddItem(app.actionsBar, 3, 0, false)

	app.pages = tview.NewPages().AddPage("main", root, true, true)
	app.modals = newModalService(app.application, app.pages, app.taskList)
	app.application.SetInputCapture(app.handleKeys)
}

func bindListVimKeys(list *tview.List) {
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		}
		return event
	})
}

func (app *App) handleKeys(event *tcell.EventKey) *tcell.EventKey {
	event = app.modals.FilterKey(event)
	if event == nil || app.modals.IsBlocking() {
		return event
	}

	event = app.handlePaneKeys(event)
	if event == nil {
		return nil
	}

	if event.Key() == tcell.KeyTab {
		app.cycleFocus(false)
		return nil
	}
	if event.Key() == tcell.KeyBacktab {
		app.cycleFocus(true)
		return nil
	}

	switch event.Rune() {
	case 'q':
		app.application.Stop()
		return nil
	case 'a':
		if app.showingCompleted {
			return nil
		}
		if app.paneActive && app.selectedPane == paneDetail {
			app.openNewSubtaskModal()
			return nil
		}
		app.openNewParentTaskModal()
		return nil
	case 'c':
		app.openNewCategoryModal()
		return nil
	case 'e':
		if app.paneActive {
			if app.showingCompleted && (app.selectedPane == paneDetail || app.selectedPane == paneTasks) {
				return nil
			}
			app.openEditModal()
		}
		return nil
	case 'd':
		if app.paneActive {
			if app.showingCompleted && (app.selectedPane == paneDetail || app.selectedPane == paneTasks) {
				return nil
			}
			app.openDeleteModal()
		}
		return nil
	case '1':
		app.showingCompleted = false
		_ = app.reload()
		return nil
	case '2':
		app.showingCompleted = true
		_ = app.reload()
		return nil
	case '?':
		app.openHelpModal()
		return nil
	case ' ':
		if app.paneActive {
			if app.showingCompleted && app.selectedPane == paneTasks {
				return nil
			}
			app.toggleFocusedCompletion()
		}
		return nil
	}
	return event
}

func (app *App) cycleFocus(backward bool) {
	order := []paneID{paneTasks, paneCategories, paneDetail, paneActions}
	index := 0
	for i, id := range order {
		if id == app.selectedPane {
			index = i
			break
		}
	}
	if backward {
		index = (index - 1 + len(order)) % len(order)
	} else {
		index = (index + 1) % len(order)
	}
	app.paneActive = false
	app.selectPane(order[index])
}

func (app *App) highlightAction(actionID string) {
	app.actionsNavigating = true
	if actionID == "" {
		app.actionsBar.Highlight()
	} else {
		app.actionsBar.Highlight(actionID)
	}
	app.actionsNavigating = false
}

func (app *App) cycleAction(offset int) {
	current := ""
	if highlights := app.actionsBar.GetHighlights(); len(highlights) > 0 {
		current = highlights[0]
	}
	index := 0
	found := false
	for i, actionID := range actionIDs {
		if actionID == current {
			index = i
			found = true
			break
		}
	}
	if found {
		index = (index + offset + len(actionIDs)) % len(actionIDs)
	} else if offset < 0 {
		index = len(actionIDs) - 1
	}
	app.highlightAction(actionIDs[index])
}

func (app *App) activateHighlightedAction() {
	highlights := app.actionsBar.GetHighlights()
	if len(highlights) == 0 {
		app.highlightAction(actionIDs[0])
		return
	}
	app.runAction(highlights[0])
}

func (app *App) runAction(actionID string) {
	switch actionID {
	case "new-task":
		if app.showingCompleted {
			return
		}
		app.openNewParentTaskModal()
	case "new-category":
		app.openNewCategoryModal()
	case "delete":
		app.openDeleteModal()
	case "help":
		app.openHelpModal()
	case "quit":
		app.application.Stop()
	}
}

func (app *App) reload() error {
	categories, err := app.categoryRepository.ListCategories()
	if err != nil {
		return fmt.Errorf("list categories: %w", err)
	}
	app.categories = categories

	var parentTasks []domain.Task
	if app.showingCompleted {
		parentTasks, err = app.taskRepository.ListCompletedParentTasks(app.selectedCategoryID)
	} else {
		parentTasks, err = app.taskRepository.ListPendingParentTasks(app.selectedCategoryID)
	}
	if err != nil {
		return fmt.Errorf("list parent tasks: %w", err)
	}

	app.reloading = true
	app.refreshCategoryList()
	app.refreshTaskList(parentTasks)
	app.reloading = false
	if app.taskList.GetItemCount() > 0 {
		app.onTaskListChanged(app.taskList.GetCurrentItem())
	} else {
		app.refreshDetail()
	}
	app.refreshTaskListTitle()
	return nil
}

func (app *App) refreshTaskListTitle() {
	app.taskModeLine.SetText(taskModeLine(app.showingCompleted))
}

func taskModeLine(showingCompleted bool) string {
	pending := "[::b]Pendentes 1[::-]"
	completed := "[::d]Concluídos 2[::-]"
	if showingCompleted {
		pending = "[::d]Pendentes 1[::-]"
		completed = "[::b]Concluídos 2[::-]"
	}
	return " " + pending + "  " + completed + " "
}

func (app *App) refreshCategoryList() {
	previousName := allCategoriesLabel
	if current := app.categoryList.GetCurrentItem(); current > 0 && current <= len(app.categories) {
		previousName = app.categories[current-1].Name
	} else if app.selectedCategoryID != nil {
		for _, category := range app.categories {
			if category.ID == *app.selectedCategoryID {
				previousName = category.Name
			}
		}
	}

	app.categoryList.Clear()
	app.categoryList.AddItem(allCategoriesLabel, "", 0, nil)
	selectedIndex := 0
	for i, category := range app.categories {
		app.categoryList.AddItem(category.Name, "", 0, nil)
		if category.Name == previousName {
			selectedIndex = i + 1
		}
	}
	if app.categoryList.GetItemCount() > 0 {
		app.categoryList.SetCurrentItem(selectedIndex)
	}
}

func (app *App) refreshTaskList(parentTasks []domain.Task) {
	app.taskList.Clear()
	app.taskListEntries = nil
	app.completedDays = nil

	app.taskList.ShowSecondaryText(!app.showingCompleted)
	if app.showingCompleted {
		app.completedDays = domain.GroupParentTasksByCompletedDay(parentTasks, time.Local)
		for _, day := range app.completedDays {
			app.taskListEntries = append(app.taskListEntries, taskListEntry{isDayHeader: true, day: day.Date})
			app.taskList.AddItem(day.Date.Format("02/01/2006"), "", 0, nil)
		}
	} else {
		for _, parentTask := range parentTasks {
			app.taskListEntries = append(app.taskListEntries, taskListEntry{parentTask: parentTask})
			app.taskList.AddItem(parentTask.Title, app.categoryName(parentTask.CategoryID), 0, nil)
		}
	}

	if len(app.taskListEntries) == 0 {
		app.selectedParentID = 0
		if app.showingCompleted {
			app.selectedDay = time.Time{}
		}
		return
	}

	index := -1
	if app.showingCompleted {
		index = app.indexOfSelectedDay()
		if index < 0 {
			index = 0
		}
	} else {
		index = app.indexOfParentTask(app.selectedParentID)
		if index < 0 {
			index = 0
			for i, entry := range app.taskListEntries {
				if !entry.isDayHeader {
					index = i
					break
				}
			}
		}
	}
	app.taskList.SetCurrentItem(index)
	app.onTaskListChanged(index)
}

func (app *App) indexOfSelectedDay() int {
	if app.selectedDay.IsZero() {
		return -1
	}
	for i, entry := range app.taskListEntries {
		if entry.isDayHeader && sameCalendarDay(entry.day, app.selectedDay) {
			return i
		}
	}
	return -1
}

func sameCalendarDay(left, right time.Time) bool {
	year, month, day := left.Date()
	otherYear, otherMonth, otherDay := right.Date()
	return year == otherYear && month == otherMonth && day == otherDay
}

func (app *App) indexOfParentTask(parentTaskID int64) int {
	if parentTaskID == 0 {
		return -1
	}
	for i, entry := range app.taskListEntries {
		if !entry.isDayHeader && entry.parentTask.ID == parentTaskID {
			return i
		}
	}
	return -1
}

func (app *App) categoryName(categoryID *int64) string {
	if categoryID == nil {
		return ""
	}
	for _, category := range app.categories {
		if category.ID == *categoryID {
			return category.Name
		}
	}
	return ""
}

func (app *App) onCategoryChanged(index int) {
	if app.reloading {
		return
	}
	if index <= 0 {
		app.selectedCategoryID = nil
	} else if index-1 < len(app.categories) {
		app.selectedCategoryID = &app.categories[index-1].ID
	}
	var parentTasks []domain.Task
	var err error
	if app.showingCompleted {
		parentTasks, err = app.taskRepository.ListCompletedParentTasks(app.selectedCategoryID)
	} else {
		parentTasks, err = app.taskRepository.ListPendingParentTasks(app.selectedCategoryID)
	}
	if err != nil {
		app.modals.ShowError(err.Error())
		return
	}
	app.refreshTaskList(parentTasks)
	app.refreshDetail()
}

func (app *App) onTaskListChanged(index int) {
	if app.reloading {
		return
	}
	if index < 0 || index >= len(app.taskListEntries) {
		app.selectedParentID = 0
		app.selectedDay = time.Time{}
		app.refreshDetail()
		return
	}
	entry := app.taskListEntries[index]
	if app.showingCompleted {
		app.selectedDay = entry.day
		app.refreshDetail()
		return
	}
	if entry.isDayHeader {
		return
	}
	app.selectedParentID = entry.parentTask.ID
	app.refreshDetail()
}

func (app *App) selectedParentTask() (domain.Task, bool) {
	for _, entry := range app.taskListEntries {
		if !entry.isDayHeader && entry.parentTask.ID == app.selectedParentID {
			return entry.parentTask, true
		}
	}
	return domain.Task{}, false
}

func (app *App) refreshDetail() {
	app.detailList.Clear()
	app.detailEntries = nil
	if app.showingCompleted {
		app.refreshCompletedDayDetail()
		return
	}
	app.subtaskShortcuts.SetText(pendingSubtaskShortcuts)
	parentTask, ok := app.selectedParentTask()
	if !ok {
		app.parentTitle.SetText("Selecione uma tarefa")
		return
	}
	app.parentTitle.SetText(parentTask.Title)

	subtasks, err := app.taskRepository.SubtasksByParentID(parentTask.ID)
	if err != nil {
		app.modals.ShowError(err.Error())
		return
	}
	for _, subtask := range subtasks {
		app.detailEntries = append(app.detailEntries, detailEntry{task: subtask})
		app.detailList.AddItem(checkboxLabel(subtask), "", 0, nil)
	}
}

func (app *App) refreshCompletedDayDetail() {
	app.subtaskShortcuts.SetText(completedSubtaskShortcuts)
	day, ok := app.completedDayByDate(app.selectedDay)
	if !ok {
		app.parentTitle.SetText("Selecione um dia")
		return
	}
	app.parentTitle.SetText(day.Date.Format("02/01/2006"))
	subtasksByParent := make(map[int64][]domain.Task, len(day.Tasks))
	for _, parentTask := range day.Tasks {
		subtasks, err := app.taskRepository.SubtasksByParentID(parentTask.ID)
		if err != nil {
			app.modals.ShowError(err.Error())
			return
		}
		subtasksByParent[parentTask.ID] = subtasks
	}
	for _, line := range buildCompletedDayLines(day, subtasksByParent) {
		app.detailEntries = append(app.detailEntries, detailEntry{isParent: line.isParent, task: line.task})
		app.detailList.AddItem(line.label, "", 0, nil)
	}
}

func (app *App) completedDayByDate(day time.Time) (domain.CompletedDay, bool) {
	if day.IsZero() {
		return domain.CompletedDay{}, false
	}
	for _, completedDay := range app.completedDays {
		if sameCalendarDay(completedDay.Date, day) {
			return completedDay, true
		}
	}
	return domain.CompletedDay{}, false
}

type completedDayLine struct {
	isParent bool
	label    string
	task     domain.Task
}

func buildCompletedDayLines(day domain.CompletedDay, subtasksByParent map[int64][]domain.Task) []completedDayLine {
	lines := make([]completedDayLine, 0)
	for _, parentTask := range day.Tasks {
		lines = append(lines, completedDayLine{
			isParent: true,
			label:    completedParentLabel(parentTask),
			task:     parentTask,
		})
		for _, subtask := range subtasksByParent[parentTask.ID] {
			lines = append(lines, completedDayLine{
				label: completedSubtaskLabel(subtask),
				task:  subtask,
			})
		}
	}
	return lines
}

func completedParentLabel(task domain.Task) string {
	if task.CompletedAt == nil {
		return task.Title
	}
	return fmt.Sprintf("%s  %s", task.CompletedAt.Local().Format("15:04"), task.Title)
}

func completedSubtaskLabel(task domain.Task) string {
	return "  " + checkboxLabel(task)
}

func checkboxLabel(task domain.Task) string {
	mark := "[ ]"
	if task.IsCompleted() {
		mark = "[x]"
	}
	return mark + " " + task.Title
}

func (app *App) focusedDetailEntry() (detailEntry, bool) {
	index := app.detailList.GetCurrentItem()
	if index < 0 || index >= len(app.detailEntries) {
		return detailEntry{}, false
	}
	return app.detailEntries[index], true
}

func (app *App) toggleFocusedCompletion() {
	if app.showingCompleted {
		app.reopenFocusedCompletedEntry()
		return
	}

	parentTask, ok := app.selectedParentTask()
	if !ok {
		return
	}

	var entry detailEntry
	switch app.selectedPane {
	case paneTasks:
		entry = detailEntry{isParent: true, task: parentTask}
	case paneDetail:
		entry, ok = app.focusedDetailEntry()
		if !ok {
			return
		}
	default:
		return
	}

	subtasks, err := app.taskRepository.SubtasksByParentID(parentTask.ID)
	if err != nil {
		app.modals.ShowError(err.Error())
		return
	}
	now := time.Now()

	if entry.isParent {
		if len(subtasks) > 0 && parentTask.IsCompleted() {
			return
		}
		if parentTask.IsCompleted() {
			parentTask, err = domain.ReopenParentWithoutSubtasks(parentTask)
		} else {
			parentTask, subtasks, err = domain.CompleteParentTask(parentTask, subtasks, now)
		}
	} else if entry.task.IsCompleted() {
		parentTask, subtasks, err = domain.ReopenSubtask(parentTask, subtasks, entry.task.ID)
	} else {
		parentTask, subtasks, err = domain.CompleteSubtask(parentTask, subtasks, entry.task.ID, now)
	}
	if err != nil {
		app.modals.ShowError(err.Error())
		return
	}
	if err := app.taskRepository.SaveTaskCompletions(parentTask, subtasks); err != nil {
		app.modals.ShowError(err.Error())
		return
	}
	app.selectedParentID = parentTask.ID
	if err := app.reload(); err != nil {
		app.modals.ShowError(err.Error())
	}
}

func (app *App) reopenFocusedCompletedEntry() {
	if app.selectedPane != paneDetail {
		return
	}
	entry, ok := app.focusedDetailEntry()
	if !ok {
		return
	}

	parentTask := entry.task
	if !entry.isParent {
		if entry.task.ParentID == nil {
			return
		}
		var err error
		parentTask, err = app.taskRepository.ParentTaskByID(*entry.task.ParentID)
		if err != nil {
			app.modals.ShowError(err.Error())
			return
		}
	}

	subtasks, err := app.taskRepository.SubtasksByParentID(parentTask.ID)
	if err != nil {
		app.modals.ShowError(err.Error())
		return
	}

	if entry.isParent {
		if len(subtasks) > 0 {
			return
		}
		parentTask, err = domain.ReopenParentWithoutSubtasks(parentTask)
	} else {
		parentTask, subtasks, err = domain.ReopenSubtask(parentTask, subtasks, entry.task.ID)
	}
	if err != nil {
		app.modals.ShowError(err.Error())
		return
	}
	if err := app.taskRepository.SaveTaskCompletions(parentTask, subtasks); err != nil {
		app.modals.ShowError(err.Error())
		return
	}
	if err := app.reload(); err != nil {
		app.modals.ShowError(err.Error())
	}
}

func stringsTrim(value string) string {
	return strings.TrimSpace(value)
}
