package user_test

import (
	"go-app-template/domain"
	"go-app-template/domain/value"
	"go-app-template/usecase/mocks"
	"go-app-template/usecase/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Register(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	userService := user.NewService(mockUserRepo)

	t.Run("ユーザー登録", func(t *testing.T) {
		mockUserRepo.On("Store", mock.Anything).Return(nil).Once()

		user, err := userService.Register(user.RegisterUserInput{Name: "test"})
		assert.NoError(t, err)
		assert.NotNil(t, user)
	})
}

func TestService_Get(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	userService := user.NewService(mockUserRepo)

	mockUserIdStr := "01911bfa-4993-7b11-ae73-ffef34f92d62"

	mockUser := domain.User{
		UserID: value.OfUserID(mockUserIdStr),
		Name:   value.NewName("test"),
	}

	t.Run("ユーザー1件取得", func(t *testing.T) {
		mockUserRepo.On("FindByID", mock.Anything).Return(&mockUser, nil).Once()

		user, err := userService.Get(user.GetUserInput{UserID: value.OfUserID(mockUserIdStr)})
		assert.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("ユーザー1件取得: ユーザーが存在しない場合", func(t *testing.T) {
		mockUserRepo.On("FindByID", mock.Anything).Return(nil, domain.ErrUserNotFound).Once()

		user, err := userService.Get(user.GetUserInput{UserID: value.OfUserID(mockUserIdStr)})
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestService_GetAll(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	userService := user.NewService(mockUserRepo)

	mockUser := domain.User{
		UserID: value.OfUserID("01911bfa-4993-7b11-ae73-ffef34f92d62"),
		Name:   value.NewName("test"),
	}

	t.Run("ユーザー全件取得", func(t *testing.T) {
		mockUserRepo.On("FindAll").Return([]*domain.User{&mockUser}, nil).Once()

		users, err := userService.GetAll()
		assert.NoError(t, err)
		assert.NotNil(t, users)
	})
}

func TestService_Change(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	userService := user.NewService(mockUserRepo)

	mockUserIdStr := "01911bfa-4993-7b11-ae73-ffef34f92d62"

	mockUser := domain.User{
		UserID: value.OfUserID(mockUserIdStr),
		Name:   value.NewName("test"),
	}

	t.Run("ユーザー変更", func(t *testing.T) {
		mockUserRepo.On("FindByID", mock.Anything).Return(&mockUser, nil).Once()
		mockUserRepo.On("Update", mock.Anything).Return(nil).Once()

		user, err := userService.Change(user.ChangeUserInput{UserID: value.OfUserID(mockUserIdStr), Name: "test"})
		assert.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("ユーザー変更: ユーザーが存在しない場合", func(t *testing.T) {
		mockUserRepo.On("FindByID", mock.Anything).Return(nil, domain.ErrUserNotFound).Once()

		user, err := userService.Change(user.ChangeUserInput{UserID: value.OfUserID(mockUserIdStr), Name: "test"})
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
