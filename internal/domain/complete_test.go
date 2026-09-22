package domain

import (
	"testing"
	"time"
)

func TestCompleteParentTaskMarksParentAndPendingSubtasks(t *testing.T) {
	now := time.Date(2026, 9, 17, 14, 32, 0, 0, time.Local)
	alreadyDoneAt := now.Add(-time.Hour)
	categoryID := int64(1)
	parentID := int64(10)

	parent := Task{ID: parentID, CategoryID: &categoryID, Title: "Relatório Q3"}
	subtasks := []Task{
		{ID: 11, ParentID: &parentID, Title: "escrever intro"},
		{ID: 12, ParentID: &parentID, Title: "juntar dados", CompletedAt: &alreadyDoneAt},
	}

	completedParent, completedSubtasks, err := CompleteParentTask(parent, subtasks, now)
	if err != nil {
		t.Fatalf("CompleteParentTask: %v", err)
	}
	if completedParent.CompletedAt == nil || !completedParent.CompletedAt.Equal(now) {
		t.Fatalf("parent completed_at = %v, want %v", completedParent.CompletedAt, now)
	}
	if completedSubtasks[0].CompletedAt == nil || !completedSubtasks[0].CompletedAt.Equal(now) {
		t.Fatalf("pending subtask completed_at = %v, want %v", completedSubtasks[0].CompletedAt, now)
	}
	if completedSubtasks[1].CompletedAt == nil || !completedSubtasks[1].CompletedAt.Equal(alreadyDoneAt) {
		t.Fatalf("already completed subtask should keep original time, got %v", completedSubtasks[1].CompletedAt)
	}
}

func TestCompleteSubtaskCompletesParentWhenItIsTheLastPending(t *testing.T) {
	now := time.Date(2026, 9, 17, 15, 0, 0, 0, time.Local)
	alreadyDoneAt := now.Add(-time.Hour)
	categoryID := int64(1)
	parentID := int64(10)

	parent := Task{ID: parentID, CategoryID: &categoryID, Title: "Relatório Q3"}
	subtasks := []Task{
		{ID: 11, ParentID: &parentID, Title: "escrever intro"},
		{ID: 12, ParentID: &parentID, Title: "juntar dados", CompletedAt: &alreadyDoneAt},
	}

	completedParent, completedSubtasks, err := CompleteSubtask(parent, subtasks, 11, now)
	if err != nil {
		t.Fatalf("CompleteSubtask: %v", err)
	}
	if completedParent.CompletedAt == nil || !completedParent.CompletedAt.Equal(now) {
		t.Fatalf("parent completed_at = %v, want %v", completedParent.CompletedAt, now)
	}
	if completedSubtasks[0].CompletedAt == nil || !completedSubtasks[0].CompletedAt.Equal(now) {
		t.Fatalf("last subtask completed_at = %v, want %v", completedSubtasks[0].CompletedAt, now)
	}
	if !completedSubtasks[1].CompletedAt.Equal(alreadyDoneAt) {
		t.Fatalf("other subtask should keep original time, got %v", completedSubtasks[1].CompletedAt)
	}
}

func TestCompleteSubtaskDoesNotCompleteParentWhenOthersRemainPending(t *testing.T) {
	now := time.Date(2026, 9, 17, 15, 0, 0, 0, time.Local)
	categoryID := int64(1)
	parentID := int64(10)

	parent := Task{ID: parentID, CategoryID: &categoryID, Title: "Relatório Q3"}
	subtasks := []Task{
		{ID: 11, ParentID: &parentID, Title: "escrever intro"},
		{ID: 12, ParentID: &parentID, Title: "juntar dados"},
	}

	completedParent, completedSubtasks, err := CompleteSubtask(parent, subtasks, 11, now)
	if err != nil {
		t.Fatalf("CompleteSubtask: %v", err)
	}
	if completedParent.CompletedAt != nil {
		t.Fatalf("parent should stay pending, got %v", completedParent.CompletedAt)
	}
	if completedSubtasks[0].CompletedAt == nil {
		t.Fatal("focused subtask should be completed")
	}
	if completedSubtasks[1].CompletedAt != nil {
		t.Fatal("other subtask should stay pending")
	}
}

func TestReopenSubtaskClearsParentAndThatSubtaskOnly(t *testing.T) {
	now := time.Date(2026, 9, 17, 15, 0, 0, 0, time.Local)
	categoryID := int64(1)
	parentID := int64(10)

	parent := Task{ID: parentID, CategoryID: &categoryID, Title: "Relatório Q3", CompletedAt: &now}
	subtasks := []Task{
		{ID: 11, ParentID: &parentID, Title: "escrever intro", CompletedAt: &now},
		{ID: 12, ParentID: &parentID, Title: "juntar dados", CompletedAt: &now},
	}

	reopenedParent, reopenedSubtasks, err := ReopenSubtask(parent, subtasks, 11)
	if err != nil {
		t.Fatalf("ReopenSubtask: %v", err)
	}
	if reopenedParent.CompletedAt != nil {
		t.Fatal("parent should be pending")
	}
	if reopenedSubtasks[0].CompletedAt != nil {
		t.Fatal("reopened subtask should be pending")
	}
	if reopenedSubtasks[1].CompletedAt == nil {
		t.Fatal("other subtask should stay completed")
	}
}

func TestReopenParentWithoutSubtasks(t *testing.T) {
	now := time.Date(2026, 9, 17, 15, 0, 0, 0, time.Local)
	categoryID := int64(1)
	parent := Task{ID: 10, CategoryID: &categoryID, Title: "Academia", CompletedAt: &now}

	reopenedParent, err := ReopenParentWithoutSubtasks(parent)
	if err != nil {
		t.Fatalf("ReopenParentWithoutSubtasks: %v", err)
	}
	if reopenedParent.CompletedAt != nil {
		t.Fatal("parent should be pending")
	}
}

func TestValidateSubtaskRejectsNestedSubtask(t *testing.T) {
	categoryID := int64(1)
	parentID := int64(10)
	subtaskID := int64(11)
	parent := Task{ID: parentID, CategoryID: &categoryID, Title: "pai"}
	nestedParent := Task{ID: subtaskID, ParentID: &parentID, Title: "sub"}

	err := ValidateSubtask(Task{ParentID: &subtaskID, Title: "neto"}, nestedParent)
	if err == nil {
		t.Fatal("expected error for nested subtask")
	}

	err = ValidateSubtask(Task{ParentID: &parentID, Title: "sub"}, parent)
	if err != nil {
		t.Fatalf("valid subtask: %v", err)
	}
}

func TestValidateTitleRejectsEmpty(t *testing.T) {
	if err := ValidateTitle("   "); err == nil {
		t.Fatal("expected error for empty title")
	}
	if err := ValidateTitle("ok"); err != nil {
		t.Fatalf("valid title: %v", err)
	}
}

func TestGroupCompletedActivityByDayPendingParentWithSub(t *testing.T) {
	loc := time.Local
	dayA := time.Date(2026, 9, 18, 10, 0, 0, 0, loc)
	parent := Task{ID: 1, Title: "pai"}
	sub := Task{ID: 11, ParentID: ptrInt64(1), Title: "sub A", CompletedAt: &dayA}

	days := GroupCompletedActivityByDay([]Task{parent}, map[int64][]Task{1: {sub}}, loc)
	if len(days) != 1 {
		t.Fatalf("days = %d, want 1", len(days))
	}
	if len(days[0].Entries) != 1 || days[0].Entries[0].Parent.ID != 1 {
		t.Fatalf("entries = %+v", days[0].Entries)
	}
	if len(days[0].Entries[0].Subtasks) != 1 || days[0].Entries[0].Subtasks[0].ID != 11 {
		t.Fatalf("subs = %+v", days[0].Entries[0].Subtasks)
	}
}

func TestGroupCompletedActivityByDaySplitsSubsAcrossDays(t *testing.T) {
	loc := time.Local
	dayA := time.Date(2026, 9, 18, 10, 0, 0, 0, loc)
	dayB := time.Date(2026, 9, 19, 11, 0, 0, 0, loc)
	parent := Task{ID: 1, Title: "pai"}
	subA := Task{ID: 11, ParentID: ptrInt64(1), Title: "a", CompletedAt: &dayA}
	subB := Task{ID: 12, ParentID: ptrInt64(1), Title: "b", CompletedAt: &dayB}

	days := GroupCompletedActivityByDay([]Task{parent}, map[int64][]Task{1: {subA, subB}}, loc)
	if len(days) != 2 {
		t.Fatalf("days = %d, want 2", len(days))
	}
	if !sameDay(days[0].Date, dayB) || days[0].Entries[0].Subtasks[0].ID != 12 {
		t.Fatalf("newest day = %+v", days[0])
	}
	if !sameDay(days[1].Date, dayA) || days[1].Entries[0].Subtasks[0].ID != 11 {
		t.Fatalf("older day = %+v", days[1])
	}
}

func TestGroupCompletedActivityByDayParentCompletionDayKeepsOnlyThatDaySubs(t *testing.T) {
	loc := time.Local
	dayA := time.Date(2026, 9, 18, 10, 0, 0, 0, loc)
	dayC := time.Date(2026, 9, 20, 15, 0, 0, 0, loc)
	parent := Task{ID: 1, Title: "pai", CompletedAt: &dayC}
	subA := Task{ID: 11, ParentID: ptrInt64(1), Title: "a", CompletedAt: &dayA}
	subC := Task{ID: 12, ParentID: ptrInt64(1), Title: "c", CompletedAt: &dayC}

	days := GroupCompletedActivityByDay([]Task{parent}, map[int64][]Task{1: {subA, subC}}, loc)
	if len(days) != 2 {
		t.Fatalf("days = %d, want 2", len(days))
	}
	var dayCEntry CompletedDayEntry
	for _, day := range days {
		if sameDay(day.Date, dayC) {
			dayCEntry = day.Entries[0]
		}
	}
	if dayCEntry.Parent.CompletedAt == nil || len(dayCEntry.Subtasks) != 1 || dayCEntry.Subtasks[0].ID != 12 {
		t.Fatalf("day C entry = %+v", dayCEntry)
	}
}

func sameDay(left, right time.Time) bool {
	y1, m1, d1 := left.Date()
	y2, m2, d2 := right.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func ptrInt64(value int64) *int64 {
	return &value
}
