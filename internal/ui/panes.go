package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type paneID int

const (
	paneTasks paneID = iota
	paneCategories
	paneDetail
	paneActions
)

var (
	paneIdleColor     = tcell.ColorDarkGray
	paneSelectedColor = tcell.NewRGBColor(200, 200, 200)
	paneActiveColor   = tcell.NewRGBColor(80, 200, 120)
)

func nextPane(current, lastTop paneID, key tcell.Key) paneID {
	switch current {
	case paneTasks:
		switch key {
		case tcell.KeyDown:
			return paneCategories
		case tcell.KeyUp:
			return paneActions
		case tcell.KeyLeft, tcell.KeyRight:
			return paneDetail
		}
	case paneCategories:
		switch key {
		case tcell.KeyUp:
			return paneTasks
		case tcell.KeyDown:
			return paneActions
		case tcell.KeyLeft, tcell.KeyRight:
			return paneDetail
		}
	case paneDetail:
		switch key {
		case tcell.KeyLeft, tcell.KeyRight:
			if lastTop == paneCategories {
				return paneCategories
			}
			return paneTasks
		case tcell.KeyDown, tcell.KeyUp:
			return paneActions
		}
	case paneActions:
		switch key {
		case tcell.KeyUp:
			if lastTop == paneCategories || lastTop == paneDetail {
				return lastTop
			}
			return paneTasks
		case tcell.KeyDown:
			return paneTasks
		}
	}
	return current
}

func (app *App) panePrimitive(id paneID) tview.Primitive {
	switch id {
	case paneCategories:
		return app.categoryList
	case paneDetail:
		return app.subtasksPane
	case paneActions:
		return app.actionsBar
	default:
		return app.tasksPane
	}
}

func (app *App) paneBox(id paneID) *tview.Box {
	switch id {
	case paneCategories:
		return app.categoryList.Box
	case paneDetail:
		return app.subtasksPane.Box
	case paneActions:
		return app.actionsBar.Box
	default:
		return app.tasksPane.Box
	}
}

func (app *App) bindPaneFocus(id paneID, box *tview.Box) {
	box.SetFocusFunc(func() {
		app.selectedPane = id
		if id != paneActions {
			app.lastTopPane = id
		}
		app.refreshPaneBorders()
	})
	box.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftDown {
			app.selectedPane = id
			if id != paneActions {
				app.lastTopPane = id
			}
			app.enterPane()
		}
		return action, event
	})
}

func (app *App) selectPane(id paneID) {
	app.selectingPane = true
	app.selectedPane = id
	if id != paneActions {
		app.lastTopPane = id
		app.highlightAction("")
	} else {
		app.highlightAction(actionIDs[0])
	}
	primitive := app.panePrimitive(id)
	app.modals.Remember(primitive)
	app.application.SetFocus(primitive)
	app.selectingPane = false
	app.refreshPaneBorders()
}

func (app *App) enterPane() {
	app.paneActive = true
	switch app.selectedPane {
	case paneTasks:
		app.application.SetFocus(app.taskList)
	case paneDetail:
		app.application.SetFocus(app.detailList)
	case paneActions:
		app.highlightAction(actionIDs[0])
	}
	app.refreshPaneBorders()
}

func (app *App) leavePane() {
	app.paneActive = false
	if app.selectedPane == paneActions {
		app.highlightAction(actionIDs[0])
	}
	app.refreshPaneBorders()
}

func (app *App) movePane(key tcell.Key) {
	app.selectPane(nextPane(app.selectedPane, app.lastTopPane, key))
}

func (app *App) handlePaneKeys(event *tcell.EventKey) *tcell.EventKey {
	if app.paneActive {
		if isEscapeKey(event) {
			app.leavePane()
			return nil
		}
		return event
	}

	switch event.Key() {
	case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight:
		app.movePane(event.Key())
		return nil
	case tcell.KeyEnter:
		app.enterPane()
		return nil
	}

	switch event.Rune() {
	case 'k':
		app.movePane(tcell.KeyUp)
		return nil
	case 'j':
		app.movePane(tcell.KeyDown)
		return nil
	case 'h':
		app.movePane(tcell.KeyLeft)
		return nil
	case 'l':
		app.movePane(tcell.KeyRight)
		return nil
	}

	return event
}

func (app *App) refreshPaneBorders() {
	for _, id := range []paneID{paneTasks, paneCategories, paneDetail, paneActions} {
		box := app.paneBox(id)
		switch {
		case id == app.selectedPane && app.paneActive:
			box.SetBorderColor(paneActiveColor)
			box.SetBorderAttributes(tcell.AttrBold)
		case id == app.selectedPane:
			box.SetBorderColor(paneSelectedColor)
			box.SetBorderAttributes(tcell.AttrNone)
		default:
			box.SetBorderColor(paneIdleColor)
			box.SetBorderAttributes(tcell.AttrNone)
		}
	}
}
