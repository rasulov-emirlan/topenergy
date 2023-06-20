package tasks

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
