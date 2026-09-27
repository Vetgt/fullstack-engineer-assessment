package task

import "errors"

var (
	ErrNotFound       = errors.New("task not found")
	ErrDuplicateTitle = errors.New("task title already exists")
)
