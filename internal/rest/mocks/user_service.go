package mocks

import (
	"go-app-template/domain"
	"go-app-template/usecase/user"

	"github.com/stretchr/testify/mock"
)

type UserService struct {
	mock.Mock
}

func (m *UserService) Register(input user.RegisterUserInput) (*domain.User, error) {
	ret := m.Called(input)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(user.RegisterUserInput) *domain.User); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(user.RegisterUserInput) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *UserService) Get(input user.GetUserInput) (*domain.User, error) {
	ret := m.Called(input)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(user.GetUserInput) *domain.User); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(user.GetUserInput) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *UserService) GetAll() ([]*domain.User, error) {
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

func (m *UserService) Change(input user.ChangeUserInput) (*domain.User, error) {
	ret := m.Called(input)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(user.ChangeUserInput) *domain.User); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(user.ChangeUserInput) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
