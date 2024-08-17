package repository

import (
	"aws-ses-local-go/domain"
	"aws-ses-local-go/domain/value"
)

type UserRepository interface {
	Store(user *domain.User) error
	Update(user *domain.User) error
	FindByID(userID value.UserID) (*domain.User, error)
	FindAll() ([]*domain.User, error)
}
