package domain

import (
	"go-app-template/domain/value"
)

type ToDo struct {
	ID     value.ToDoID `json:"id"`
	Title  value.Title  `json:"title"`
	Done   value.Done   `json:"done"`
	UserID value.UserID `json:"user_id"`
}

func NewToDo(title string, userId value.UserID) *ToDo {
	return &ToDo{
		ID:     value.NewToDoID(),
		Title:  value.NewTitle(title),
		Done:   value.OfDone(false),
		UserID: userId,
	}
}

func (t *ToDo) ChangeDone(done value.Done) {
	t.Done = done
}
