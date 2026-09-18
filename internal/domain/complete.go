package domain

import "time"

func CompleteParentTask(parent Task, subtasks []Task, now time.Time) (Task, []Task, error) {
	if !parent.IsParent() {
		return Task{}, nil, ErrNotAParent
	}
	completedAt := now
	parent.CompletedAt = &completedAt
	updatedSubtasks := copyTasks(subtasks)
	for i := range updatedSubtasks {
		if updatedSubtasks[i].CompletedAt == nil {
			subtaskCompletedAt := now
			updatedSubtasks[i].CompletedAt = &subtaskCompletedAt
		}
	}
	return parent, updatedSubtasks, nil
}

func CompleteSubtask(parent Task, subtasks []Task, subtaskID int64, now time.Time) (Task, []Task, error) {
	if !parent.IsParent() {
		return Task{}, nil, ErrNotAParent
	}
	updatedSubtasks := copyTasks(subtasks)
	found := false
	allCompleted := true
	for i := range updatedSubtasks {
		if updatedSubtasks[i].ID == subtaskID {
			found = true
			if updatedSubtasks[i].CompletedAt == nil {
				subtaskCompletedAt := now
				updatedSubtasks[i].CompletedAt = &subtaskCompletedAt
			}
		}
		if updatedSubtasks[i].CompletedAt == nil {
			allCompleted = false
		}
	}
	if !found {
		return Task{}, nil, ErrSubtaskNotFound
	}
	if allCompleted {
		parentCompletedAt := now
		parent.CompletedAt = &parentCompletedAt
	}
	return parent, updatedSubtasks, nil
}

func ReopenSubtask(parent Task, subtasks []Task, subtaskID int64) (Task, []Task, error) {
	if !parent.IsParent() {
		return Task{}, nil, ErrNotAParent
	}
	updatedSubtasks := copyTasks(subtasks)
	found := false
	for i := range updatedSubtasks {
		if updatedSubtasks[i].ID == subtaskID {
			found = true
			updatedSubtasks[i].CompletedAt = nil
		}
	}
	if !found {
		return Task{}, nil, ErrSubtaskNotFound
	}
	parent.CompletedAt = nil
	return parent, updatedSubtasks, nil
}

func ReopenParentWithoutSubtasks(parent Task) (Task, error) {
	if !parent.IsParent() {
		return Task{}, ErrNotAParent
	}
	parent.CompletedAt = nil
	return parent, nil
}

func copyTasks(tasks []Task) []Task {
	copied := make([]Task, len(tasks))
	copy(copied, tasks)
	return copied
}
