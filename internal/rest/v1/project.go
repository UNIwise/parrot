package v1

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/uniwise/parrot/internal/poedit"
)

type getProjectLanguageRequest struct {
	Project  int    `param:"project" validate:"required"`
	Language string `param:"language" validate:"required,languageCode"`
}

func getProjectLanguage(ctx echo.Context) error {
	c := ctx.(*Context)

	req := new(getProjectLanguageRequest)
	if err := c.Bind(req); err != nil {
		c.Log.WithError(err).Error("Error binding request")

		return echo.ErrBadRequest
	}

	if err := c.Validate(req); err != nil {
		c.Log.Debug(req)
		c.Log.WithError(err).Error("Error validating request")

		return echo.ErrBadRequest
	}

	trans, err := c.ProjectService.GetTranslation(
		ctx.Request().Context(),
		req.Project,
		req.Language,
		"key_value_json",
	)
	if err != nil {
		switch err.(type) {
		case *poedit.ErrProjectPermissionDenied:
			return echo.ErrBadRequest
		case *poedit.ErrLanguageNotFound:
			return echo.ErrNotFound
		default:
			c.Log.WithError(err).Error("Error retrieving translation")

			return echo.ErrInternalServerError
		}
	}

	c.Response().Header().Add("Etag", trans.Checksum)
	c.Response().Header().Add("Cache-Control", fmt.Sprintf("max-age=%.0f", trans.TTL.Seconds()))

	return c.Stream(http.StatusOK, "application/json", bytes.NewReader(trans.Data))
}
