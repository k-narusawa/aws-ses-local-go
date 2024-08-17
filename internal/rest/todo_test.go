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

func TestToDoHandler_Create(t *testing.T) {
	mockToDoService := new(mocks.ToDoService)
	mockToDoService.On("Create", mock.Anything).Return(&domain.ToDo{}, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users/01911bfa-4993-7b11-ae73-ffef34f92d62/todos", strings.NewReader(`{"title":"test"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/todos")
	c.SetParamNames("user_id")
	c.SetParamValues("01911bfa-4993-7b11-ae73-ffef34f92d62")

	var handler = rest.ToDoHandler{
		ToDoService: mockToDoService,
	}

	err := handler.Create(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestToDoHandler_FindByUserId(t *testing.T) {
	mockToDoService := new(mocks.ToDoService)
	mockUserIdStr := "01911bfa-4993-7b11-ae73-ffef34f92d62"

	mockToDo := domain.ToDo{
		UserID: value.OfUserID(mockUserIdStr),
		Title:  value.NewTitle("test"),
	}
	mockToDoService.On("FindByUserId", mock.Anything).Return([]domain.ToDo{mockToDo}, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/"+mockUserIdStr+"/todos", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/todos")
	c.SetParamNames("user_id")
	c.SetParamValues(mockUserIdStr)

	var handler = rest.ToDoHandler{
		ToDoService: mockToDoService,
	}

	err := handler.FindByUserId(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestToDoHandler_ChangeToDoDone(t *testing.T) {
	mockToDoService := new(mocks.ToDoService)
	mockToDoService.On("ChangeToDoDone", mock.Anything).Return(&domain.ToDo{}, nil).Once()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users/01911bfa-4993-7b11-ae73-ffef34f92d62/todos/01911bfa-4993-7b11-ae73-ffef34f92d62", strings.NewReader(`{"done":true}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:user_id/todos/:todo_id")
	c.SetParamNames("user_id", "todo_id")
	c.SetParamValues("01911bfa-4993-7b11-ae73-ffef34f92d62", "01911bfa-4993-7b11-ae73-ffef34f92d62")

	var handler = rest.ToDoHandler{
		ToDoService: mockToDoService,
	}

	err := handler.ChangeToDoDone(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
