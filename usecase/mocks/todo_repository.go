package mocks

import (
	"go-app-template/domain"
	"go-app-template/domain/value"

	"github.com/stretchr/testify/mock"
)

type ToDoRepository struct {
	mock.Mock
}

func (m *ToDoRepository) FindAll() ([]domain.ToDo, error) {
	ret := m.Called()

	var r0 []domain.ToDo
	if rf, ok := ret.Get(0).(func() []domain.ToDo); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).([]domain.ToDo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *ToDoRepository) FindByID(todoID value.ToDoID) (*domain.ToDo, error) {
	ret := m.Called(todoID)

	var r0 *domain.ToDo
	if rf, ok := ret.Get(0).(func(value.ToDoID) *domain.ToDo); ok {
		r0 = rf(todoID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ToDo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(value.ToDoID) error); ok {
		r1 = rf(todoID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *ToDoRepository) FindByUserID(userID value.UserID) ([]domain.ToDo, error) {
	ret := m.Called(userID)

	var r0 []domain.ToDo
	if rf, ok := ret.Get(0).(func(value.UserID) []domain.ToDo); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).([]domain.ToDo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(value.UserID) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *ToDoRepository) Store(todo domain.ToDo) error {
	ret := m.Called(todo)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.ToDo) error); ok {
		r0 = rf(&todo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m *ToDoRepository) Update(todo domain.ToDo) error {
	ret := m.Called(todo)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.ToDo) error); ok {
		r0 = rf(&todo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m *ToDoRepository) Delete(todoID value.ToDoID) error {
	ret := m.Called(todoID)

	var r0 error
	if rf, ok := ret.Get(0).(func(value.ToDoID) error); ok {
		r0 = rf(todoID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
