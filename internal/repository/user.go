package repository

import (
	"go-app-template/domain"
	"go-app-template/domain/value"
)

var (
	users = map[value.UserID]*domain.User{}
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Store(user *domain.User) error {
	users[user.UserID] = user
	return nil
}

func (r *UserRepository) Update(user *domain.User) error {
	users[user.UserID] = user
	return nil
}

func (r *UserRepository) FindByID(userId value.UserID) (*domain.User, error) {
	user := users[userId]
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *UserRepository) FindAll() ([]*domain.User, error) {
	allUsers := make([]*domain.User, 0, len(users))
	for _, user := range users {
		allUsers = append(allUsers, user)
	}
	return allUsers, nil
}
