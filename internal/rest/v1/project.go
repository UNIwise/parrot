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

type getProjectRequest struct {
	Project int `param:"project" validate:"required"`
}

type getProjectResponse struct {
	Id        int                          `json:"id"`
	Languages []getProjectResponseLanguage `json:"languages"`
}

type getProjectResponseLanguage struct {
	Updated    int64   `json:"updated"`
	Code       string  `json:"code"`
	Percentage float64 `json:"percentage"`
}

func getProject(ctx echo.Context) error {
	c := ctx.(*Context)

	req := new(getProjectRequest)
	if err := c.Bind(req); err != nil {
		c.Log.WithError(err).Error("Error binding request")

		return echo.ErrBadRequest
	}

	l := c.Log.WithFields(logrus.Fields{
		"project": req.Project,
	})
	if err := c.Validate(req); err != nil {
		l.WithError(err).Error("Error validating request")

		return echo.ErrBadRequest
	}

	project, err := c.ProjectService.GetProjectMeta(
		ctx.Request().Context(),
		req.Project,
	)

	if errors.Is(err, context.Canceled) {
		return echo.NewHTTPError(499, "client closed request")
	}

	if err != nil {
		switch err.(type) {
		case *poedit.ErrProjectPermissionDenied:
			return echo.ErrBadRequest
		default:
			l.WithError(err).Error("Error retrieving project information")

			return echo.ErrInternalServerError
		}
	}

	languages := make([]getProjectResponseLanguage, len(project.Meta.Languages))
	for i, l := range project.Meta.Languages {
		languages[i] = getProjectResponseLanguage{
			Updated:    l.Updated.Unix(),
			Code:       l.Code,
			Percentage: l.Percentage,
		}
	}

	c.Response().Header().Add("Etag", project.Checksum)
	c.Response().Header().Add("Cache-Control", fmt.Sprintf("max-age=%.0f", project.TTL.Seconds()))

	return ctx.JSON(http.StatusOK, &getProjectResponse{
		Id:        req.Project,
		Languages: languages,
	})
}

type getProjectLanguageRequest struct {
	Project  int    `param:"project" validate:"required"`
	Language string `param:"language" validate:"required,languageCode"`
	Format   string `query:"format" validate:"default=key_value_json,oneof=po pot mo xls xlsx csv ini resw resx android_strings apple_strings xliff properties key_value_json json yml xlf xmb xtb arb rise_360_xliff"`
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

	contentMeta, err := poedit.GetContentMeta(format)
	if err != nil {
		l.Error("No extension and content type found")

		return echo.ErrBadRequest
	}

	lang, err := c.ProjectService.GetLanguage(
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

	c.Response().Header().Add("Etag", lang.Checksum)
	c.Response().Header().Add("Cache-Control", fmt.Sprintf("max-age=%.0f", lang.TTL.Seconds()))
	c.Response().Header().Add("Content-Disposition", fmt.Sprintf("filename=%d-%s.%s", req.Project, req.Language, contentMeta.Extension))
	c.Response().Header().Add("Content-Transfer-Encoding", "8bit")

	return c.Stream(http.StatusOK, contentMeta.Type, bytes.NewReader(lang.Data))
}
