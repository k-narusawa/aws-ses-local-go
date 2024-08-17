package middleware

import (
	"errors"
	"go-app-template/domain"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func HandleError(c echo.Context, err error) error {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}

	switch err {
	case domain.ErrUserNotFound:
		log.Warn(err)
		domainError := domain.ToDomainError(err)
		return c.JSON(domain.StatusMap[domainError.Code], domainError)
	case domain.ErrToDoNotFound:
		log.Warn(err)
		domainError := domain.ToDomainError(err)
		return c.JSON(domain.StatusMap[domainError.Code], domainError)
	default:
		log.Error(err)
		return c.JSON(500, domain.ToDomainError(err))
	}
}
