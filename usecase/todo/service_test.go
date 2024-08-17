package todo_test

import (
	"go-app-template/domain"
	"go-app-template/domain/value"
	"go-app-template/usecase/mocks"
	"go-app-template/usecase/todo"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Create(t *testing.T) {
	mockToDoRepo := new(mocks.ToDoRepository)
	mockUserRepo := new(mocks.UserRepository)

	todoService := todo.NewService(mockToDoRepo, mockUserRepo)

	mockUserIdStr := "01911bfa-4993-7b11-ae73-ffef34f92d62"

	mockUser := domain.User{
		UserID: value.OfUserID(mockUserIdStr),
		Name:   value.NewName("test"),
	}

	t.Run("ToDo作成", func(t *testing.T) {
		mockToDoRepo.On("Store", mock.Anything).Return(nil).Once()
		mockUserRepo.On("FindByID", mock.Anything).Return(&mockUser, nil).Once()

		todo, err := todoService.Create(todo.CreateToDoInput{Title: "test"})
		assert.NoError(t, err)
		assert.NotNil(t, todo)
	})

	t.Run("ToDo作成: ユーザーが存在しない場合", func(t *testing.T) {
		mockToDoRepo.On("Store", mock.Anything).Return(nil).Once()
		mockUserRepo.On("FindByID", mock.Anything).Return(nil, domain.ErrUserNotFound).Once()

		todo, err := todoService.Create(todo.CreateToDoInput{Title: "test"})
		assert.Error(t, err)
		assert.Nil(t, todo)
	})
}

func TestService_FindByUserId(t *testing.T) {
	mockToDoRepo := new(mocks.ToDoRepository)
	mockUserRepo := new(mocks.UserRepository)

	todoService := todo.NewService(mockToDoRepo, mockUserRepo)

	mockUserIdStr := "01911bfa-4993-7b11-ae73-ffef34f92d62"
	mockToDoIdStr := "01911bfa-4993-7b11-ae73-ffef34f92d62"

	mockToDo := domain.ToDo{
		ID:     value.OfToDoID(mockToDoIdStr),
		Title:  value.NewTitle("test"),
		UserID: value.OfUserID(mockUserIdStr),
	}

	t.Run("ToDo取得", func(t *testing.T) {
		mockToDoRepo.On("FindByUserID", mock.Anything).Return([]domain.ToDo{mockToDo}, nil).Once()

		todos, err := todoService.FindByUserId(value.OfUserID(mockUserIdStr))
		assert.NoError(t, err)
		assert.NotNil(t, todos)
	})
}

func TestService_ChangeToDoDone(t *testing.T) {
	mockToDoRepo := new(mocks.ToDoRepository)
	mockUserRepo := new(mocks.UserRepository)

	todoService := todo.NewService(mockToDoRepo, mockUserRepo)

	mockUserIdStr := "01911bfa-4993-7b11-ae73-ffef34f92d62"
	mockToDoIdStr := "01911bfa-4993-7b11-ae73-ffef34f92d62"

	mockToDo := domain.ToDo{
		ID:     value.OfToDoID(mockToDoIdStr),
		Title:  value.NewTitle("test"),
		UserID: value.OfUserID(mockUserIdStr),
	}

	t.Run("ToDoのステータス変更", func(t *testing.T) {
		mockToDoRepo.On("FindByID", mock.Anything).Return(&mockToDo, nil).Once()
		mockToDoRepo.On("Update", mock.Anything).Return(nil).Once()

		todo, err := todoService.ChangeToDoDone(todo.ChangeStatusInput{ToDoID: value.OfToDoID(mockToDoIdStr), Done: true})
		assert.NoError(t, err)
		assert.NotNil(t, todo)
	})

	t.Run("ToDoのステータス変更: ToDoが存在しない場合", func(t *testing.T) {
		mockToDoRepo.On("FindByID", mock.Anything).Return(nil, domain.ErrToDoNotFound).Once()

		todo, err := todoService.ChangeToDoDone(todo.ChangeStatusInput{ToDoID: value.OfToDoID(mockToDoIdStr), Done: true})
		assert.Error(t, err)
		assert.Nil(t, todo)
	})
}
