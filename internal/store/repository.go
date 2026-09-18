package store

import "tuitodo/internal/domain"

type CategoryRepository interface {
	InsertCategory(category domain.Category) (domain.Category, error)
	RenameCategory(categoryID int64, name string) error
	DeleteCategory(categoryID int64) error
	CategoryByID(categoryID int64) (domain.Category, error)
	ListCategories() ([]domain.Category, error)
}

type TaskRepository interface {
	InsertParentTask(task domain.Task) (domain.Task, error)
	InsertSubtask(subtask domain.Task, parent domain.Task) (domain.Task, error)
	UpdateTaskTitle(taskID int64, title string) error
	SaveTaskCompletions(parent domain.Task, subtasks []domain.Task) error
	DeleteTaskByID(taskID int64) error
	ParentTaskByID(taskID int64) (domain.Task, error)
	SubtasksByParentID(parentID int64) ([]domain.Task, error)
	ListPendingParentTasks(categoryID *int64) ([]domain.Task, error)
	ListCompletedParentTasks(categoryID *int64) ([]domain.Task, error)
}
