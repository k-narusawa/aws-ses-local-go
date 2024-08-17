package web

import (
	"go-app-template/domain"
	"go-app-template/domain/value"
	"go-app-template/internal/middleware"
	"go-app-template/usecase/todo"
	"go-app-template/usecase/user"

	"github.com/labstack/echo/v4"
)

type UserService interface {
	Register(user.RegisterUserInput) (*domain.User, error)
	Get(user.GetUserInput) (*domain.User, error)
	GetAll() ([]*domain.User, error)
	Change(user.ChangeUserInput) (*domain.User, error)
}

type ToDoService interface {
	Create(todo.CreateToDoInput) (*domain.ToDo, error)
	FindByUserId(value.UserID) ([]domain.ToDo, error)
	ChangeToDoDone(todo.ChangeStatusInput) (*domain.ToDo, error)
}

type WebUserHandler struct {
	UserService UserService
	ToDoService ToDoService
}

func NewWebUserHandler(
	e *echo.Echo,
	userService UserService,
	todoService ToDoService,
) {
	handler := &WebUserHandler{
		UserService: userService,
		ToDoService: todoService,
	}

	e.GET("/", handler.AllUsers)
}

func (h *WebUserHandler) AllUsers(c echo.Context) error {
	users, err := h.UserService.GetAll()
	if err != nil {
		return middleware.HandleError(c, err)
	}

	return c.Render(200, "user.html", users)
}
