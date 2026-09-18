package ui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestModalServiceFilterKeyClosesOpenModal(t *testing.T) {
	pages := tview.NewPages().AddPage("main", tview.NewBox(), true, true)
	service := newModalService(tview.NewApplication(), pages, tview.NewBox())
	service.open = true
	pages.AddPage(modalPageName, tview.NewBox(), true, true)

	event := service.FilterKey(tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone))
	if event != nil {
		t.Fatal("expected Escape to be consumed")
	}
	if service.IsOpen() {
		t.Fatal("expected modal to close")
	}
}

func TestModalServiceFilterKeyIgnoresOtherKeys(t *testing.T) {
	pages := tview.NewPages().AddPage("main", tview.NewBox(), true, true)
	service := newModalService(tview.NewApplication(), pages, tview.NewBox())
	service.open = true

	event := tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone)
	got := service.FilterKey(event)
	if got != event {
		t.Fatal("expected non-Escape keys to pass through")
	}
	if !service.IsOpen() {
		t.Fatal("expected modal to stay open")
	}
}
