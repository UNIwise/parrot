package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/uniwise/parrot/internal/project"
)

type getProjectItemResponse struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	NumberOfVersions int    `json:"numberOfVersions"`
	CreatedAt        string `json:"createdAt"`
}

type getAllProjectsResponse struct {
	Projects []getProjectItemResponse `json:"projects"`
}

func (h *Handlers) getAllProjects(ctx echo.Context, l *logrus.Entry) error {
	l.Debug("Retrieving projects")

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
	ID int `param:"id" validate:"required"`
}

type getProjectResponse struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	NumberOfVersions int    `json:"numberOfVersions"`
	CreatedAt        string `json:"createdAt"`
}

func (h *Handlers) getProject(ctx echo.Context, l *logrus.Entry) error {
	req := new(getProjectRequest)
	if err := ctx.Bind(req); err != nil {
		l.WithError(err).Error("Error binding request")

		return echo.ErrBadRequest
	}

	l = l.WithField("project", req.ID)

	l.Debug("Retrieving project")

	project, err := h.ProjectService.GetProjectByID(
		ctx.Request().Context(),
		req.ID,
	)
	if err != nil {
		if err.Error() == "failed to get project: not found" {
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
	ProjectID int `param:"id" validate:"required"`
}

type getProjectVersionsItemResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type getProjectVersionsResponse struct {
	Versions []getProjectVersionsItemResponse `json:"versions"`
}

func (h *Handlers) getProjectVersions(c echo.Context, l *logrus.Entry) error {
	req := new(getProjectVersionsRequest)
	if err := c.Bind(req); err != nil {
		l.WithError(err).Error("Error binding request")

		return echo.ErrBadRequest
	}

	l = l.WithField("project", req.ProjectID)

	l.Debug("Retrieving project versions")

	versions, err := h.ProjectService.GetProjectVersions(
		c.Request().Context(),
		req.ProjectID,
	)
	if err != nil {
		if err.Error() == "failed to get project versions: not found" {
			return echo.ErrNotFound
		}

		l.WithError(err).Error("Error retrieving project versions")

		return echo.ErrInternalServerError
	}

	response := h.newGetProjectVersionsResponse(versions)

	return c.JSON(http.StatusOK, response)
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

type deleteProjectVersionRequest struct {
	ProjectID uint   `param:"project_id" validate:"required,numeric"`
	VersionID string `param:"version_id" validate:"required,version,max=20"`
}

func (h *Handlers) deleteProjectVersion(c echo.Context, l *logrus.Entry) error {
	req := new(deleteProjectVersionRequest)
	if err := c.Bind(req); err != nil {
		l.WithError(err).Error("Error binding request")

		return echo.ErrBadRequest
	}

	l = l.WithFields(logrus.Fields{
		"project": req.ProjectID,
		"version": req.VersionID,
	})

	l.Debug("Deleting project version")

	err := h.ProjectService.DeleteProjectVersionByIDAndProjectID(
		c.Request().Context(),
		req.VersionID,
		req.ProjectID,
	)
	if err != nil {
		l.WithError(err).Error("Error deleting project version")

		return echo.ErrInternalServerError
	}

	l.Info("Project version deleted")

	return c.NoContent(http.StatusOK)
}

type postProjectVersionRequest struct {
	ID   int    `param:"id" validate:"required,numeric"`
	Name string `json:"name" validate:"required,version,max=20"`
}

func (h *Handlers) postProjectVersion(c echo.Context, l *logrus.Entry) error {
	req := new(postProjectVersionRequest)
	if err := c.Bind(req); err != nil {
		l.WithError(err).Error("Error binding request")

		return echo.ErrBadRequest
	}

	l = l.WithFields(logrus.Fields{
		"project": req.ID,
		"name":    req.Name,
	})

	l.Debug("Creating project version")

	if err := c.Validate(req); err != nil {
		l.WithError(err).Error("Error binding request")

		return echo.ErrBadRequest
	}

	l = l.WithFields(logrus.Fields{
		"project": req.ID,
		"name":    req.Name,
	})

	err := h.ProjectService.CreateLanguagesVersion(
		c.Request().Context(),
		req.ID,
		req.Name,
	)
	if err != nil {
		if errors.Is(err, project.ErrVersionAlreadyExist) {
			return echo.ErrBadRequest
		}
		
		l.WithError(err).Error("Error creating project version")

		return echo.ErrInternalServerError
	}

	l.Info("Creating project version")

	return c.NoContent(http.StatusCreated)
}
