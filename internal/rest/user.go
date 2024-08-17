package rest

import (
	"go-app-template/domain"
	"go-app-template/domain/value"
	"go-app-template/internal/middleware"
	"go-app-template/usecase/user"

	"github.com/labstack/echo/v4"
)

type UserService interface {
	Register(user.RegisterUserInput) (*domain.User, error)
	Get(user.GetUserInput) (*domain.User, error)
	GetAll() ([]*domain.User, error)
	Change(user.ChangeUserInput) (*domain.User, error)
}

type UserHandler struct {
	UserService UserService
}

func NewUserHandler(e *echo.Echo, userService UserService) {
	handler := &UserHandler{
		UserService: userService,
	}

	e.POST("/users", handler.Register)
	e.GET("/users/:user_id", handler.Get)
	e.GET("/users", handler.GetAll)
	e.PUT("/users/:user_id", handler.Change)
}

func (h *UserHandler) Register(c echo.Context) error {
	in := new(user.RegisterUserInput)
	if err := c.Bind(&in); err != nil {
		return middleware.HandleError(c, err)
	}

	user, err := h.UserService.Register(*in)
	if err != nil {
		return middleware.HandleError(c, err)
	}

	return c.JSON(201, user)
}

func (h *UserHandler) Get(c echo.Context) error {
	userID := c.Param("user_id")
	in := user.GetUserInput{
		UserID: value.OfUserID(userID),
	}
	user, err := h.UserService.Get(in)
	if err != nil {
		return middleware.HandleError(c, err)
	}

	return c.JSON(200, user)
}

func (h *UserHandler) GetAll(c echo.Context) error {
	users, err := h.UserService.GetAll()
	if err != nil {
		return middleware.HandleError(c, err)
	}

	return c.JSON(200, users)
}

func (h *UserHandler) Change(c echo.Context) error {
	in := user.ChangeUserInput{}
	userId := c.Param("user_id")
	if err := c.Bind(&in); err != nil {
		return err
	}
	in.UserID = value.OfUserID(userId)

	user, err := h.UserService.Change(in)
	if err != nil {
		return middleware.HandleError(c, err)
	}

	return c.JSON(200, user)
}
