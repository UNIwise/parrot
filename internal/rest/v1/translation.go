package v1

import (
	"bytes"
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/uniwise/parrot/internal/cache"
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

	b, err := c.Cacher.GetTranslation(ctx.Request().Context(), *req.Project, *req.Language, "key_value_json")
	if err == nil {
		return c.Stream(http.StatusOK, "application/json", bytes.NewReader(b))
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		c.Log.WithError(err).Error("Error fetching cached data")
		return echo.ErrInternalServerError
	}
	c.Log.Debug("Cache miss!")

	result, err := c.Client.FetchTerms(context.Background(), *req.Project, *req.Language, "key_value_json")
	if err != nil {
		c.Log.WithError(err).Error("Error fetching data from poeditor")
		return echo.ErrInternalServerError
	}
	if err := c.Cacher.SetTranslation(ctx.Request().Context(), *req.Project, *req.Language, "key_value_json", result); err != nil {
		c.Log.WithError(err).Error("Error caching translation data")
		return echo.ErrInternalServerError
	}

	return c.Stream(http.StatusOK, "application/json", bytes.NewReader(b))
}
