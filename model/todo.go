package model

import (
	uuid "github.com/nu7hatch/gouuid"
)

// Todo type
type Todo struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	IsComplete bool   `json:"isComplete"`
}

// TodoBuffer type
type TodoBuffer struct {
	Title      string `json:"title"`
	IsComplete bool   `json:"isComplete"`
}

// NewTodo - instantiates Todo from buffer
func (b *TodoBuffer) NewTodo() *Todo {
	u, _ := uuid.NewV4()
	return &Todo{
		ID:         u.String(),
		Title:      b.Title,
		IsComplete: b.IsComplete,
	}
}
