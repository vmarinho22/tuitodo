package store

import (
	"testing"
	"time"

	"tuitodo/internal/domain"
)

func openTestStore(t *testing.T) *SQLiteStore {
	t.Helper()
	sqliteStore, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() {
		_ = sqliteStore.Close()
	})
	return sqliteStore
}

func TestInsertAndListCategories(t *testing.T) {
	sqliteStore := openTestStore(t)
	category, err := sqliteStore.InsertCategory(domain.Category{Name: "trabalho", CreatedAt: time.Now()})
	if err != nil {
		t.Fatalf("InsertCategory: %v", err)
	}
	if category.ID == 0 {
		t.Fatal("expected category id")
	}
	categories, err := sqliteStore.ListCategories()
	if err != nil {
		t.Fatalf("ListCategories: %v", err)
	}
	if len(categories) != 1 || categories[0].Name != "trabalho" {
		t.Fatalf("categories = %+v", categories)
	}
}

func TestInsertCategoryRejectsDuplicateName(t *testing.T) {
	sqliteStore := openTestStore(t)
	if _, err := sqliteStore.InsertCategory(domain.Category{Name: "trabalho", CreatedAt: time.Now()}); err != nil {
		t.Fatalf("InsertCategory: %v", err)
	}
	if _, err := sqliteStore.InsertCategory(domain.Category{Name: "trabalho", CreatedAt: time.Now()}); err == nil {
		t.Fatal("expected duplicate name error")
	}
}

func TestListPendingParentTasksFiltersByCategory(t *testing.T) {
	sqliteStore := openTestStore(t)
	trabalho, _ := sqliteStore.InsertCategory(domain.Category{Name: "trabalho", CreatedAt: time.Now()})
	pessoal, _ := sqliteStore.InsertCategory(domain.Category{Name: "pessoal", CreatedAt: time.Now()})
	_, err := sqliteStore.InsertParentTask(domain.Task{CategoryID: &trabalho.ID, Title: "Relatório Q3", CreatedAt: time.Now()})
	if err != nil {
		t.Fatalf("InsertParentTask trabalho: %v", err)
	}
	_, err = sqliteStore.InsertParentTask(domain.Task{CategoryID: &pessoal.ID, Title: "Academia", CreatedAt: time.Now()})
	if err != nil {
		t.Fatalf("InsertParentTask pessoal: %v", err)
	}
	allPending, err := sqliteStore.ListPendingParentTasks(nil)
	if err != nil {
		t.Fatalf("ListPendingParentTasks todas: %v", err)
	}
	if len(allPending) != 2 {
		t.Fatalf("got %d pending, want 2", len(allPending))
	}
	trabalhoPending, err := sqliteStore.ListPendingParentTasks(&trabalho.ID)
	if err != nil {
		t.Fatalf("ListPendingParentTasks trabalho: %v", err)
	}
	if len(trabalhoPending) != 1 || trabalhoPending[0].Title != "Relatório Q3" {
		t.Fatalf("filtered = %+v", trabalhoPending)
	}
}

func TestInsertSubtaskRejectsNestedSubtask(t *testing.T) {
	sqliteStore := openTestStore(t)
	category, _ := sqliteStore.InsertCategory(domain.Category{Name: "trabalho", CreatedAt: time.Now()})
	parent, _ := sqliteStore.InsertParentTask(domain.Task{CategoryID: &category.ID, Title: "pai", CreatedAt: time.Now()})
	subtask, err := sqliteStore.InsertSubtask(domain.Task{ParentID: &parent.ID, Title: "sub", CreatedAt: time.Now()}, parent)
	if err != nil {
		t.Fatalf("InsertSubtask: %v", err)
	}
	_, err = sqliteStore.InsertSubtask(domain.Task{ParentID: &subtask.ID, Title: "neto", CreatedAt: time.Now()}, subtask)
	if err == nil {
		t.Fatal("expected nested subtask error")
	}
}

func TestSaveTaskCompletionsAndListCompleted(t *testing.T) {
	sqliteStore := openTestStore(t)
	now := time.Date(2026, 9, 17, 14, 32, 0, 0, time.Local)
	category, _ := sqliteStore.InsertCategory(domain.Category{Name: "trabalho", CreatedAt: now})
	parent, _ := sqliteStore.InsertParentTask(domain.Task{CategoryID: &category.ID, Title: "Relatório Q3", CreatedAt: now})
	subtask, _ := sqliteStore.InsertSubtask(domain.Task{ParentID: &parent.ID, Title: "escrever intro", CreatedAt: now}, parent)
	completedParent, completedSubtasks, err := domain.CompleteParentTask(parent, []domain.Task{subtask}, now)
	if err != nil {
		t.Fatalf("CompleteParentTask: %v", err)
	}
	if err := sqliteStore.SaveTaskCompletions(completedParent, completedSubtasks); err != nil {
		t.Fatalf("SaveTaskCompletions: %v", err)
	}
	pending, err := sqliteStore.ListPendingParentTasks(nil)
	if err != nil {
		t.Fatalf("ListPendingParentTasks: %v", err)
	}
	if len(pending) != 0 {
		t.Fatalf("pending after complete = %+v", pending)
	}
	completed, err := sqliteStore.ListCompletedParentTasks(nil)
	if err != nil {
		t.Fatalf("ListCompletedParentTasks: %v", err)
	}
	if len(completed) != 1 || completed[0].CompletedAt == nil {
		t.Fatalf("completed = %+v", completed)
	}
}

func TestDeleteParentTaskDeletesSubtasks(t *testing.T) {
	sqliteStore := openTestStore(t)
	category, _ := sqliteStore.InsertCategory(domain.Category{Name: "trabalho", CreatedAt: time.Now()})
	parent, _ := sqliteStore.InsertParentTask(domain.Task{CategoryID: &category.ID, Title: "pai", CreatedAt: time.Now()})
	_, _ = sqliteStore.InsertSubtask(domain.Task{ParentID: &parent.ID, Title: "sub", CreatedAt: time.Now()}, parent)
	if err := sqliteStore.DeleteTaskByID(parent.ID); err != nil {
		t.Fatalf("DeleteTaskByID: %v", err)
	}
	subtasks, err := sqliteStore.SubtasksByParentID(parent.ID)
	if err != nil {
		t.Fatalf("SubtasksByParentID: %v", err)
	}
	if len(subtasks) != 0 {
		t.Fatalf("subtasks after delete = %+v", subtasks)
	}
}

func TestSettingMissingReturnsEmpty(t *testing.T) {
	sqliteStore := openTestStore(t)
	value, err := sqliteStore.Setting(SettingLocale)
	if err != nil {
		t.Fatalf("Setting: %v", err)
	}
	if value != "" {
		t.Fatalf("Setting = %q, want empty", value)
	}
}

func TestSetAndGetSetting(t *testing.T) {
	sqliteStore := openTestStore(t)
	if err := sqliteStore.SetSetting(SettingLocale, "pt-BR"); err != nil {
		t.Fatalf("SetSetting: %v", err)
	}
	value, err := sqliteStore.Setting(SettingLocale)
	if err != nil {
		t.Fatalf("Setting: %v", err)
	}
	if value != "pt-BR" {
		t.Fatalf("Setting = %q, want pt-BR", value)
	}
	if err := sqliteStore.SetSetting(SettingLocale, "en-US"); err != nil {
		t.Fatalf("SetSetting overwrite: %v", err)
	}
	value, err = sqliteStore.Setting(SettingLocale)
	if err != nil {
		t.Fatalf("Setting after overwrite: %v", err)
	}
	if value != "en-US" {
		t.Fatalf("Setting = %q, want en-US", value)
	}
}

func TestDeleteCategoryFailsWhenParentTasksExist(t *testing.T) {
	sqliteStore := openTestStore(t)
	category, _ := sqliteStore.InsertCategory(domain.Category{Name: "trabalho", CreatedAt: time.Now()})
	_, _ = sqliteStore.InsertParentTask(domain.Task{CategoryID: &category.ID, Title: "pai", CreatedAt: time.Now()})
	if err := sqliteStore.DeleteCategory(category.ID); err == nil {
		t.Fatal("expected delete category error")
	}
}
