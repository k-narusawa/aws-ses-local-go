package value

import "github.com/gofrs/uuid"

type UserID struct {
	ValueObject[string]
}

func NewUserID() UserID {
	userId, _ := uuid.NewV7()
	return UserID{ValueObject: ValueObject[string]{value: userId.String()}}
}

func (u UserID) String() string {
	return u.value
}

func OfUserID(value string) UserID {
	return UserID{ValueObject: ValueObject[string]{value: value}}
}
