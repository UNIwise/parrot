package v1

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/uniwise/parrot/internal/cache"
)

type getProjectLanguageRequest struct {
	Project  *int    `url:"project" validate:"required"`
	Language *string `url:"language" validate:"required"`
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

	b, err := c.Cacher.GetTranslation(*req.Project, *req.Language, "key_value_json")
	if err == nil {
		return c.Stream(http.StatusOK, "application/json", bytes.NewReader(b))
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		c.Log.WithError(err).Error("Error fetching cached data")
		return echo.ErrInternalServerError
	}
	c.Log.Debug("Cache miss!")

	return echo.NewHTTPError(http.StatusNotImplemented, "Not implemented")
}
