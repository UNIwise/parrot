package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/uniwise/parrot/internal/poedit"
)

type getProjectLanguageRequest struct {
	Project  int    `param:"project" validate:"required"`
	Language string `param:"language" validate:"required,languageCode"`
	Format   string `query:"format" validate:"omitempty,oneof=po pot mo xls xlsx csv ini resw resx android_strings apple_strings xliff properties key_value_json json yml xlf xmb xtb arb rise_360_xliff"`
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

	format := "key_value_json"
	if req.Format != "" {
		format = req.Format
	}

	trans, err := c.ProjectService.GetTranslation(
		ctx.Request().Context(),
		req.Project,
		req.Language,
		format,
	)
	if errors.Is(err, context.Canceled) {
		return echo.NewHTTPError(499, "client closed request")
	}

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
