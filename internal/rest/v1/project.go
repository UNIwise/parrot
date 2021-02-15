package v1

import (
	"bytes"
	"net/http"

	"github.com/labstack/echo/v4"
)

type getProjectLanguageRequest struct {
	Project  *int    `url:"project" validate:"required"`
	Language *string `url:"language" validate:"required,languageCode"`
}

func getProjectLanguage(ctx echo.Context) error {
	c := ctx.(*Context)

	req := new(getProjectLanguageRequest)
	if err := c.Bind(req); err != nil {
		c.Log.WithError(err).Error("Error binding request")

		return echo.ErrBadRequest
	}

	if err := c.Validate(req); err != nil {
		c.Log.WithError(err).Error("Error validating request")

		return echo.ErrBadRequest
	}

	b, err := c.ProjectService.GetTranslation(
		ctx.Request().Context(),
		*req.Project,
		*req.Language,
		"key_value_json",
	)
	if err != nil {
		c.Log.WithError(err).Error("Error retrieving translation")

		return echo.ErrInternalServerError
	}

	return c.Stream(http.StatusOK, "application/json", bytes.NewReader(b))
}
