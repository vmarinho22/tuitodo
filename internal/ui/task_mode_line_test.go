package ui

import "testing"

func TestTaskModeLine(t *testing.T) {
	pending := taskModeLine(false)
	wantPending := " [::b]Pendentes 1[::-]  [::d]Concluídos 2[::-] "
	if pending != wantPending {
		t.Fatalf("taskModeLine(false) = %q, want %q", pending, wantPending)
	}

	completed := taskModeLine(true)
	wantCompleted := " [::d]Pendentes 1[::-]  [::b]Concluídos 2[::-] "
	if completed != wantCompleted {
		t.Fatalf("taskModeLine(true) = %q, want %q", completed, wantCompleted)
	}
}
