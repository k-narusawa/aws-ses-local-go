package mocks

import (
	"go-app-template/domain"
	"go-app-template/domain/value"
	"go-app-template/usecase/todo"

	"github.com/stretchr/testify/mock"
)

type ToDoService struct {
	mock.Mock
}

func (m *ToDoService) Create(input todo.CreateToDoInput) (*domain.ToDo, error) {
	ret := m.Called(input)

	var r0 *domain.ToDo
	if rf, ok := ret.Get(0).(func(todo.CreateToDoInput) *domain.ToDo); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ToDo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(todo.CreateToDoInput) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *ToDoService) FindByUserId(userId value.UserID) ([]domain.ToDo, error) {
	ret := m.Called(userId)

	var r0 []domain.ToDo
	if rf, ok := ret.Get(0).(func(value.UserID) []domain.ToDo); ok {
		r0 = rf(userId)
	} else {
		r0 = ret.Get(0).([]domain.ToDo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(value.UserID) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *ToDoService) ChangeToDoDone(input todo.ChangeStatusInput) (*domain.ToDo, error) {
	ret := m.Called(input)

	var r0 *domain.ToDo
	if rf, ok := ret.Get(0).(func(todo.ChangeStatusInput) *domain.ToDo); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ToDo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(todo.ChangeStatusInput) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
