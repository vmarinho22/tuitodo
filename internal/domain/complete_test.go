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

func TestGroupParentTasksByCompletedDay(t *testing.T) {
	loc := time.Local
	morning := time.Date(2026, 9, 17, 9, 10, 0, 0, loc)
	afternoon := time.Date(2026, 9, 17, 14, 32, 0, 0, loc)
	yesterday := time.Date(2026, 9, 16, 18, 0, 0, 0, loc)

	days := GroupParentTasksByCompletedDay([]Task{
		{ID: 1, Title: "tarde", CompletedAt: &afternoon},
		{ID: 2, Title: "manhã", CompletedAt: &morning},
		{ID: 3, Title: "ontem", CompletedAt: &yesterday},
	}, loc)

	if len(days) != 2 {
		t.Fatalf("got %d days, want 2", len(days))
	}
	if days[0].Tasks[0].Title != "tarde" || days[0].Tasks[1].Title != "manhã" {
		t.Fatalf("17/09 order = %v, %v", days[0].Tasks[0].Title, days[0].Tasks[1].Title)
	}
	if days[1].Tasks[0].Title != "ontem" {
		t.Fatalf("16/09 first = %s", days[1].Tasks[0].Title)
	}
}
