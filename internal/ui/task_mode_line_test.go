package ui

import (
	"testing"

	"github.com/rivo/tview"
)

func TestTaskModeLine(t *testing.T) {
	todo := tview.Escape("A fazer [1]")
	completed := tview.Escape("Concluídos [2]")

	gotPending := taskModeLine(false)
	wantPending := " [::b]" + todo + "[::-]   [::d]" + completed + "[::-] "
	if gotPending != wantPending {
		t.Fatalf("taskModeLine(false) = %q, want %q", gotPending, wantPending)
	}

	gotCompleted := taskModeLine(true)
	wantCompleted := " [::d]" + todo + "[::-]   [::b]" + completed + "[::-] "
	if gotCompleted != wantCompleted {
		t.Fatalf("taskModeLine(true) = %q, want %q", gotCompleted, wantCompleted)
	}
}
