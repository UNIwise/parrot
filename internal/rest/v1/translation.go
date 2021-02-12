package v1

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type getProjectLanguageRequest struct {
	Project  *int    `url:"project" validate:"required"`
	Language *string `url:"language" validate:"required"`
}

func getProjectLanguage(c echo.Context) error {
	return errors.New("Not implemented")
}
