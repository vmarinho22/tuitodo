package domain

import "time"

type Category struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}

type Task struct {
	ID          int64
	ParentID    *int64
	CategoryID  *int64
	Title       string
	CompletedAt *time.Time
	CreatedAt   time.Time
}

func (task Task) IsParent() bool {
	return task.ParentID == nil
}

func (task Task) IsSubtask() bool {
	return task.ParentID != nil
}

func (task Task) IsCompleted() bool {
	return task.CompletedAt != nil
}
