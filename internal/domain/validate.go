package domain

import (
	"errors"
	"strings"
)

var (
	ErrEmptyTitle      = errors.New("title is empty")
	ErrNestedSubtask   = errors.New("subtask cannot have a subtask as parent")
	ErrParentHasNoID   = errors.New("parent task has no id")
	ErrSubtaskNotFound = errors.New("subtask not found")
	ErrNotAParent      = errors.New("task is not a parent")
)

func ValidateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmptyTitle
	}
	return nil
}

func ValidateSubtask(subtask Task, parent Task) error {
	if err := ValidateTitle(subtask.Title); err != nil {
		return err
	}
	if !parent.IsParent() {
		return ErrNestedSubtask
	}
	if parent.ID == 0 {
		return ErrParentHasNoID
	}
	if subtask.ParentID == nil || *subtask.ParentID != parent.ID {
		return ErrNestedSubtask
	}
	return nil
}
