package rest_test

import (
	"go-app-template/domain"
	"go-app-template/domain/value"
	"go-app-template/internal/rest"
	"go-app-template/internal/rest/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_Register(t *testing.T) {
	mockUserService := new(mocks.UserService)
	mockUserService.On("Register", mock.Anything).Return(&domain.User{}, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"test"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var handler = rest.UserHandler{
		UserService: mockUserService,
	}

	err := handler.Register(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestUserHandler_Get(t *testing.T) {
	mockUserService := new(mocks.UserService)
	mockUserIdStr := "01911bfa-4993-7b11-ae73-ffef34f92d62"

	mockUser := domain.User{
		UserID: value.OfUserID(mockUserIdStr),
		Name:   value.NewName("test"),
	}
	mockUserService.On("Get", mock.Anything).Return(&mockUser, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/"+mockUserIdStr, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id")
	c.SetParamNames("user_id")
	c.SetParamValues(mockUserIdStr)

	var handler = rest.UserHandler{
		UserService: mockUserService,
	}

	err := handler.Get(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUserHandler_Get_NotFound(t *testing.T) {
	mockUserService := new(mocks.UserService)
	mockUserIdStr := "01911bfa-4993-7b11-ae73-ffef34f92d62"

	mockUserService.On("Get", mock.Anything).Return(nil, domain.ErrUserNotFound).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/"+mockUserIdStr, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id")
	c.SetParamNames("user_id")
	c.SetParamValues(mockUserIdStr)

	var handler = rest.UserHandler{
		UserService: mockUserService,
	}

	err := handler.Get(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUserHandler_GetAll(t *testing.T) {
	mockUserService := new(mocks.UserService)
	mockUserService.On("GetAll").Return([]*domain.User{}, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var handler = rest.UserHandler{
		UserService: mockUserService,
	}

	err := handler.GetAll(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUserHandler_Change(t *testing.T) {
	mockUserService := new(mocks.UserService)
	mockUserIdStr := "01911bfa-4993-7b11-ae73-ffef34f92d62"

	mockUser := domain.User{
		UserID: value.OfUserID(mockUserIdStr),
		Name:   value.NewName("test"),
	}
	mockUserService.On("Change", mock.Anything).Return(&mockUser, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users/"+mockUserIdStr, strings.NewReader(`{"name":"test"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id")
	c.SetParamNames("user_id")
	c.SetParamValues(mockUserIdStr)

	var handler = rest.UserHandler{
		UserService: mockUserService,
	}

	err := handler.Change(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
