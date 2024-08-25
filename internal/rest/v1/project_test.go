package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/uniwise/parrot/internal/project"

	gomock "go.uber.org/mock/gomock"
)

var (
	errTest                     = errors.New("test error")
	errNotFoundTest             = errors.New("failed to get project: not found")
	testID               uint   = 1
	testVersionID        uint   = 1
	testName             string = "testname"
	testNumberOfVersions uint   = 3
	testCreatedAt               = time.Now()
)

func TestGetAllProjects(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	testCtx := e.NewContext(req, resp)

	testCtx.SetPath("/projects")

	projectService := project.NewMockService(gomock.NewController(t))

	h := &Handlers{
		ProjectService: projectService,
	}
	t.Run("GetAllProjects, success", func(t *testing.T) {
		projectService.EXPECT().GetAllProjects(context.Background()).Times(1).Return([]project.Project{{
			ID:               testID,
			Name:             testName,
			NumberOfVersions: testNumberOfVersions,
			CreatedAt:        testCreatedAt.UTC(),
		}}, nil)

		err := h.getAllProjects(testCtx, logrus.NewEntry(logrus.New()))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		var response getAllProjectsResponse
		err = json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, response, getAllProjectsResponse{
			Projects: []getProjectItemResponse{
				{
					ID:               testID,
					Name:             testName,
					NumberOfVersions: testNumberOfVersions,
					CreatedAt:        testCreatedAt.UTC(),
				},
			},
		})
	})

	t.Run("GetAllProjects, error", func(t *testing.T) {
		projectService.EXPECT().GetAllProjects(context.Background()).Times(1).Return(nil, errTest)

		err := h.getAllProjects(testCtx, logrus.NewEntry(logrus.New()))
		assert.Error(t, err)
	})
}

func TestHandlers_newGetAllProjectsResponse(t *testing.T) {
	t.Parallel()

	projects := []project.Project{{
		ID:               testID,
		Name:             testName,
		NumberOfVersions: testNumberOfVersions,
		CreatedAt:        testCreatedAt.UTC(),
	}}

	h := &Handlers{}

	response := h.newGetAllProjectsResponse(projects)

	assert.NotNil(t, response)
	assert.Len(t, response.Projects, 1)
	assert.Equal(t, response.Projects[0].ID, testID)
	assert.Equal(t, response.Projects[0].Name, testName)
	assert.Equal(t, response.Projects[0].NumberOfVersions, testNumberOfVersions)
	assert.Equal(t, response.Projects[0].CreatedAt, testCreatedAt.UTC())
}

func TestGetProject(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	testCtx := e.NewContext(req, resp)

	testCtx.SetPath("/projects/:id")
	testCtx.SetParamNames("id")
	testCtx.SetParamValues(fmt.Sprintf("%d", testID))

	projectService := project.NewMockService(gomock.NewController(t))

	h := &Handlers{
		ProjectService: projectService,
	}
	t.Run("GetProject, success", func(t *testing.T) {
		projectService.EXPECT().GetProjectByID(context.Background(), int(testID)).Times(1).Return(&project.Project{
			ID:               testID,
			Name:             testName,
			NumberOfVersions: testNumberOfVersions,
			CreatedAt:        testCreatedAt.UTC(),
		}, nil)

		err := h.getProject(testCtx, logrus.NewEntry(logrus.New()))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		var response getProjectResponse
		err = json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, response, getProjectResponse{
			ID:               testID,
			Name:             testName,
			NumberOfVersions: testNumberOfVersions,
			CreatedAt:        testCreatedAt.UTC(),
		})
	})

	t.Run("GetProject, error", func(t *testing.T) {
		projectService.EXPECT().GetProjectByID(context.Background(), int(testID)).Times(1).Return(nil, errTest)

		err := h.getProject(testCtx, logrus.NewEntry(logrus.New()))
		assert.Error(t, err)
	})
}

func TestHandlers_newGetProjectResponse(t *testing.T) {
	t.Parallel()

	project := project.Project{
		ID:               testID,
		Name:             testName,
		NumberOfVersions: testNumberOfVersions,
		CreatedAt:        testCreatedAt.UTC(),
	}

	h := &Handlers{}

	response := h.newGetProjectResponse(project)

	assert.NotNil(t, response)
	assert.Equal(t, response.ID, testID)
	assert.Equal(t, response.Name, testName)
	assert.Equal(t, response.NumberOfVersions, testNumberOfVersions)
	assert.Equal(t, response.CreatedAt, testCreatedAt.UTC())
}

func TestGetProjectVersions(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	testCtx := e.NewContext(req, resp)

	testCtx.SetPath("/projects/:id/versions")
	testCtx.SetParamNames("id")
	testCtx.SetParamValues(fmt.Sprintf("%d", testID))

	projectService := project.NewMockService(gomock.NewController(t))

	h := &Handlers{
		ProjectService: projectService,
	}
	t.Run("GetProjectVersions, success", func(t *testing.T) {
		projectService.EXPECT().GetProjectVersions(context.Background(), int(testID)).Times(1).Return([]project.Version{{
			ID:        testVersionID,
			Name:      testName,
			ProjectID: testID,
			CreatedAt: testCreatedAt.UTC(),
		}}, nil)

		err := h.getProjectVersions(testCtx, logrus.NewEntry(logrus.New()))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		var response getProjectVersionsResponse
		err = json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, response, getProjectVersionsResponse{
			Versions: []getProjectVersionsItemResponse{
				{
					ID:        testID,
					Name:      testName,
					CreatedAt: testCreatedAt.UTC(),
				},
			},
		})
	})

	t.Run("GetProjectVersions, error", func(t *testing.T) {
		projectService.EXPECT().GetProjectVersions(context.Background(), int(testID)).Times(1).Return(nil, errTest)

		err := h.getProjectVersions(testCtx, logrus.NewEntry(logrus.New()))
		assert.Error(t, err)
	})
}

func TestHandlers_newGetProjectVersionsResponse(t *testing.T) {
	t.Parallel()

	versions := []project.Version{{
		ID:        testVersionID,
		Name:      testName,
		ProjectID: testID,
		CreatedAt: testCreatedAt.UTC(),
	}}

	h := &Handlers{}

	response := h.newGetProjectVersionsResponse(versions)

	assert.NotNil(t, response)
	assert.Len(t, response.Versions, 1)
	assert.Equal(t, response.Versions[0].ID, testID)
	assert.Equal(t, response.Versions[0].Name, testName)
	assert.Equal(t, response.Versions[0].CreatedAt, testCreatedAt.UTC())
}
