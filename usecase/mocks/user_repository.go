package mocks

import (
	"go-app-template/domain"
	"go-app-template/domain/value"

	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Store(user *domain.User) error {
	ret := m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m *UserRepository) Update(user *domain.User) error {
	ret := m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m *UserRepository) FindByID(userID value.UserID) (*domain.User, error) {
	ret := m.Called(userID)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(value.UserID) *domain.User); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(value.UserID) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *UserRepository) FindAll() ([]*domain.User, error) {
	ret := m.Called()

	var r0 []*domain.User
	if rf, ok := ret.Get(0).(func() []*domain.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
