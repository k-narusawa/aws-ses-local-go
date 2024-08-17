package value

import "github.com/gofrs/uuid"

type ToDoID struct {
	ValueObject[string]
}

func NewToDoID() ToDoID {
	todoId, _ := uuid.NewV7()
	return ToDoID{ValueObject: ValueObject[string]{value: todoId.String()}}
}

func (t ToDoID) String() string {
	return t.value
}

func OfToDoID(value string) ToDoID {
	return ToDoID{ValueObject: ValueObject[string]{value: value}}
}
