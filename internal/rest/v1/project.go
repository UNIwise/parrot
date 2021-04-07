package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/uniwise/parrot/pkg/poedit"
)

type Content struct {
	Extension, Type string
}

var (
	ContentMap = map[string]Content{
		"pot":             {Extension: "pot", Type: "text/plain; charset=utf-8"},
		"po":              {Extension: "po", Type: "text/plain; charset=utf-8"},
		"mo":              {Extension: "mo", Type: "application/octet-stream"},
		"xls":             {Extension: "xls", Type: "application/vnd.ms-excel"},
		"xlsx":            {Extension: "xlsx", Type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		"csv":             {Extension: "csv", Type: "text/csv; charset=utf-8"},
		"ini":             {Extension: "ini", Type: "text/plain; charset=utf-8"},
		"resw":            {Extension: "resw", Type: "application/xml"},
		"resx":            {Extension: "resx", Type: "application/xml"},
		"android_strings": {Extension: "xml", Type: "application/xml"},
		"apple_strings":   {Extension: "strings", Type: "text/plain; charset=utf-8"},
		"xliff":           {Extension: "xliff", Type: "application/xml"},
		"properties":      {Extension: "properties", Type: "text/plain; charset=utf-8"},
		"key_value_json":  {Extension: "json", Type: "application/json"},
		"json":            {Extension: "json", Type: "application/json"},
		"yml":             {Extension: "yml", Type: "text/plain; charset=utf-8"},
		"xlf":             {Extension: "xlf", Type: "application/xml"},
		"xmb":             {Extension: "xmb", Type: "application/xml"},
		"xtb":             {Extension: "xtb", Type: "application/xml"},
		"arb":             {Extension: "arb", Type: "application/json"},
		"rise_360_xliff":  {Extension: "xliff", Type: "application/xml"},
	}
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

	l := c.Log.WithFields(logrus.Fields{
		"project":  req.Project,
		"language": req.Language,
		"format":   req.Format,
	})
	if err := c.Validate(req); err != nil {
		l.WithError(err).Error("Error validating request")

		return echo.ErrBadRequest
	}

	format := "key_value_json"
	if req.Format != "" {
		format = req.Format
	}

	content, ok := ContentMap[format]
	if !ok {
		l.Error("No extension and content type found")

		return echo.ErrBadRequest
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
			l.WithError(err).Error("Error retrieving translation")

			return echo.ErrInternalServerError
		}
	}

	c.Response().Header().Add("Etag", trans.Checksum)
	c.Response().Header().Add("Cache-Control", fmt.Sprintf("max-age=%.0f", trans.TTL.Seconds()))
	c.Response().Header().Add("Content-Disposition", fmt.Sprintf("filename=%d-%s.%s", req.Project, req.Language, content.Extension))
	c.Response().Header().Add("Content-Transfer-Encoding", "8bit")

	return c.Stream(http.StatusOK, content.Type, bytes.NewReader(trans.Data))
}
