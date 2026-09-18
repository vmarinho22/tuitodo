package ui

import (
	"testing"

	"tuitodo/internal/i18n"

	"github.com/rivo/tview"
)

func TestTaskModeLine(t *testing.T) {
	app := &App{catalog: i18n.NewCatalog(i18n.LocalePtBR)}
	todo := tview.Escape(app.t(i18n.KeyModeTodo))
	completed := tview.Escape(app.t(i18n.KeyModeDone))

	app.showingCompleted = false
	gotPending := app.taskModeLineText()
	wantPending := " [::b]" + todo + "[::-]   [::d]" + completed + "[::-] "
	if gotPending != wantPending {
		t.Fatalf("taskModeLineText() pending = %q, want %q", gotPending, wantPending)
	}

	app.showingCompleted = true
	gotCompleted := app.taskModeLineText()
	wantCompleted := " [::d]" + todo + "[::-]   [::b]" + completed + "[::-] "
	if gotCompleted != wantCompleted {
		t.Fatalf("taskModeLineText() completed = %q, want %q", gotCompleted, wantCompleted)
	}
}
