package rest

import (
	"go-app-template/domain"
	"go-app-template/domain/value"
	"go-app-template/usecase/todo"

	"go-app-template/internal/middleware"

	"github.com/labstack/echo/v4"
)

type ToDoService interface {
	Create(todo.CreateToDoInput) (*domain.ToDo, error)
	FindByUserId(value.UserID) ([]domain.ToDo, error)
	ChangeToDoDone(todo.ChangeStatusInput) (*domain.ToDo, error)
}

type ToDoHandler struct {
	ToDoService ToDoService
}

func NewToDoHandler(e *echo.Echo, todoService ToDoService) {
	handler := &ToDoHandler{
		ToDoService: todoService,
	}

	users := e.Group("/users/:user_id")

	users.POST("/todos", handler.Create)
	users.GET("/todos", handler.FindByUserId)
	users.PUT("/todos/:todo_id", handler.ChangeToDoDone)
}

func (h *ToDoHandler) Create(c echo.Context) error {
	in := todo.CreateToDoInput{}
	if err := c.Bind(&in); err != nil {
		return middleware.HandleError(c, err)
	}
	in.UserID = value.OfUserID(c.Param("user_id"))

	todo, err := h.ToDoService.Create(in)
	if err != nil {
		return middleware.HandleError(c, err)
	}

	return c.JSON(200, todo)
}

func (h *ToDoHandler) FindByUserId(c echo.Context) error {
	userId := value.OfUserID(c.Param("user_id"))

	todos, err := h.ToDoService.FindByUserId(userId)
	if err != nil {
		return middleware.HandleError(c, err)
	}

	return c.JSON(200, todos)
}

func (h *ToDoHandler) ChangeToDoDone(c echo.Context) error {
	userId := value.OfUserID(c.Param("user_id"))
	todoId := value.OfToDoID(c.Param("todo_id"))

	in := todo.ChangeStatusInput{}
	if err := c.Bind(&in); err != nil {
		return middleware.HandleError(c, err)
	}
	in.UserID = userId
	in.ToDoID = todoId

	todo, err := h.ToDoService.ChangeToDoDone(in)
	if err != nil {
		return middleware.HandleError(c, err)
	}

	return c.JSON(200, todo)
}
