package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/uniwise/parrot/internal/project"
	"github.com/uniwise/parrot/pkg/poedit"
)

type getProjectLanguageRequest struct {
	Project  int    `param:"project" validate:"required"`
	Language string `param:"language" validate:"required,languageCode"`
	Format   string `query:"format" validate:"omitempty,oneof=po pot mo xls xlsx csv ini resw resx android_strings apple_strings xliff properties key_value_json json yml xlf xmb xtb arb rise_360_xliff"`
}

func (h *Handlers) getProjectLanguage(ctx echo.Context, l *logrus.Entry) error {
	req := new(getProjectLanguageRequest)
	if err := ctx.Bind(req); err != nil {
		l.WithError(err).Error("Error binding request")

		return echo.ErrBadRequest
	}

	l = l.WithFields(logrus.Fields{
		"project":  req.Project,
		"language": req.Language,
		"format":   req.Format,
	})

	if err := ctx.Validate(req); err != nil {
		l.WithError(err).Error("Error validating request")

		return echo.ErrBadRequest
	}

	format := "key_value_json"
	if req.Format != "" {
		format = req.Format
	}

	contentMeta, err := poedit.GetContentMeta(format)
	if err != nil {
		l.WithError(err).Errorf("No content meta found for format %s", format)

		return echo.ErrBadRequest
	}

	trans, err := h.ProjectService.GetTranslation(
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

	if ctx.Request() == nil {
		l.Error("Request is nil")

		return errors.New("request is nil")
	}

	if ctx.Request().Header.Get("If-None-Match") == trans.Checksum {
		return ctx.NoContent(http.StatusNotModified)
	}

	ctx.Response().Header().Add("Etag", trans.Checksum)
	ctx.Response().Header().Add("Cache-Control", fmt.Sprintf("max-age=%.0f", trans.TTL.Seconds()))
	ctx.Response().Header().Add("Content-Disposition", fmt.Sprintf("filename=%d-%s.%s", req.Project, req.Language, contentMeta.Extension))
	ctx.Response().Header().Add("Content-Transfer-Encoding", "8bit")

	return ctx.Stream(http.StatusOK, contentMeta.Type, bytes.NewReader(trans.Data))
}

type getProjectItemResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	NumberOfVersions uint      `json:"numberOfVersions"`
	CreatedAt        time.Time `json:"createdAt"`
}

type getAllProjectsResponse struct {
	Projects []getProjectItemResponse `json:"projects"`
}

func (h *Handlers) getAllProjects(ctx echo.Context, l *logrus.Entry) error {
	projects, err := h.ProjectService.GetAllProjects(ctx.Request().Context())
	if err != nil {
		l.WithError(err).Error("Error retrieving projects")

		return echo.ErrInternalServerError
	}

	response := h.newGetAllProjectsResponse(projects)

	return ctx.JSON(http.StatusOK, response)
}

func (h *Handlers) newGetAllProjectsResponse(
	projects []project.Project,
) *getAllProjectsResponse {
	response := &getAllProjectsResponse{
		Projects: make([]getProjectItemResponse, len(projects)),
	}

	for i, project := range projects {

		response.Projects[i] = getProjectItemResponse{
			ID:               project.ID,
			Name:             project.Name,
			NumberOfVersions: project.NumberOfVersions,
			CreatedAt:        project.CreatedAt,
		}
	}

	return response
}

type getProjectRequest struct {
	ID  int    `param:"id" validate:"required"`
}

type getProjectResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	NumberOfVersions uint      `json:"numberOfVersions"`
	CreatedAt        time.Time `json:"createdAt"`
}

func (h *Handlers) getProject(ctx echo.Context, l *logrus.Entry) error {
	req := new(getProjectRequest)
	if err := ctx.Bind(req); err != nil {
		l.WithError(err).Error("Error binding request")

		return echo.ErrBadRequest
	}

	l = l.WithField("project", req.ID)

	project, err := h.ProjectService.GetProjectByID(
		ctx.Request().Context(),
		req.ID,
	)
	if err != nil {
		if (err.Error() == "failed to get project: not found") {
			return echo.ErrNotFound
		}

		l.WithError(err).Error("Error retrieving project")

		return echo.ErrInternalServerError
	}

	response := h.newGetProjectResponse(*project)

	return ctx.JSON(http.StatusOK, response)
}

func (h *Handlers) newGetProjectResponse(
	project project.Project,
) *getProjectResponse {
	return &getProjectResponse{
		ID:               project.ID,
		Name:             project.Name,
		NumberOfVersions: project.NumberOfVersions,
		CreatedAt:        project.CreatedAt,
	}
}

type getProjectVersionsRequest struct {
	ProjectID  int    `param:"id" validate:"required"`
}

type getProjectVersionsItemResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	CreatedAt        time.Time `json:"createdAt"`
}

type getProjectVersionsResponse struct {
	Versions []getProjectVersionsItemResponse `json:"versions"`
}

func (h *Handlers) getProjectVersions(ctx echo.Context, l *logrus.Entry) error {
	req := new(getProjectVersionsRequest)
	if err := ctx.Bind(req); err != nil {
		l.WithError(err).Error("Error binding request")

		return echo.ErrBadRequest
	}

	l = l.WithField("project", req.ProjectID)

	versions, err := h.ProjectService.GetProjectVersions(
		ctx.Request().Context(),
		req.ProjectID,
	)
	if err != nil {
		if (err.Error() == "failed to get project versions: not found") {
			return echo.ErrNotFound
		}

		l.WithError(err).Error("Error retrieving project versions")

		return echo.ErrInternalServerError
	}

	response := h.newGetProjectVersionsResponse(versions)

	return ctx.JSON(http.StatusOK, response)
}

func (h *Handlers) newGetProjectVersionsResponse(
	versions []project.Version,
) *getProjectVersionsResponse {
	response := &getProjectVersionsResponse{
		Versions: make([]getProjectVersionsItemResponse, len(versions)),
	}

	for i, version := range versions {
		response.Versions[i] = getProjectVersionsItemResponse{
			ID:        version.ID,
			Name:      version.Name,
			CreatedAt: version.CreatedAt,
		}
	}

	return response
}

