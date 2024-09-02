package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/uniwise/parrot/internal/project"

	gomock "go.uber.org/mock/gomock"
)

var (
	errTest                     = errors.New("test error")
	errNotFoundTest             = errors.New("failed to get project: not found")
	testID               int64  = 1
	testVersionID        string = "test-id"
	testName             string = "testname"
	testNumberOfVersions int    = 3
	testCreatedAt               = time.Now()
	testCreatedAtString  string = testCreatedAt.String()
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
			CreatedAt:        testCreatedAtString,
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
					CreatedAt:        testCreatedAtString,
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
		CreatedAt:        testCreatedAtString,
	}}

	h := &Handlers{}

	response := h.newGetAllProjectsResponse(projects)

	assert.NotNil(t, response)
	assert.Len(t, response.Projects, 1)
	assert.Equal(t, response.Projects[0].ID, testID)
	assert.Equal(t, response.Projects[0].Name, testName)
	assert.Equal(t, response.Projects[0].NumberOfVersions, testNumberOfVersions)
	assert.Equal(t, response.Projects[0].CreatedAt, testCreatedAtString)
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
			CreatedAt:        testCreatedAtString,
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
			CreatedAt:        testCreatedAtString,
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
		CreatedAt:        testCreatedAtString,
	}

	h := &Handlers{}

	response := h.newGetProjectResponse(project)

	assert.NotNil(t, response)
	assert.Equal(t, response.ID, testID)
	assert.Equal(t, response.Name, testName)
	assert.Equal(t, response.NumberOfVersions, testNumberOfVersions)
	assert.Equal(t, response.CreatedAt, testCreatedAtString)
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
					ID:        testVersionID,
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
		CreatedAt: testCreatedAt.UTC(),
	}}

	h := &Handlers{}

	response := h.newGetProjectVersionsResponse(versions)

	assert.NotNil(t, response)
	assert.Len(t, response.Versions, 1)
	assert.Equal(t, response.Versions[0].ID, testVersionID)
	assert.Equal(t, response.Versions[0].Name, testName)
	assert.Equal(t, response.Versions[0].CreatedAt, testCreatedAt.UTC())
}

func TestDeleteProjectVersion(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	resp := httptest.NewRecorder()

	testCtx := e.NewContext(req, resp)

	testCtx.SetPath("/projects/:project_id/versions/:version_id")
	testCtx.SetParamNames("project_id", "version_id")
	testCtx.SetParamValues(fmt.Sprintf("%d", testID), fmt.Sprintf("%s", testVersionID))

	projectService := project.NewMockService(gomock.NewController(t))

	h := &Handlers{
		ProjectService: projectService,
	}
	t.Run("deleteProjectVersion, success", func(t *testing.T) {
		projectService.EXPECT().DeleteProjectVersionByIDAndProjectID(context.Background(), testVersionID, uint(testID)).Times(1).Return(nil)

		err := h.deleteProjectVersion(testCtx, logrus.NewEntry(logrus.New()))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("deleteProjectVersion, error", func(t *testing.T) {
		projectService.EXPECT().DeleteProjectVersionByIDAndProjectID(context.Background(), testVersionID, uint(testID)).Times(1).Return(errTest)

		err := h.deleteProjectVersion(testCtx, logrus.NewEntry(logrus.New()))
		assert.Error(t, err)
	})
}

type CustomValidator struct {
	validator *validator.Validate
}

// Validate implements the echo.Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestPostProjectVersion(t *testing.T) {
	t.Parallel()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	projectService := project.NewMockService(gomock.NewController(t))

	h := &Handlers{
		ProjectService: projectService,
	}

	t.Run("postProjectVersion, success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"id":   testID,
			"name": testName,
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		resp := httptest.NewRecorder()

		testCtx := e.NewContext(req, resp)

		testCtx.SetPath("/projects/:project_id/versions")
		testCtx.SetParamNames("project_id")
		testCtx.SetParamValues(fmt.Sprintf("%d", testID))

		projectService.EXPECT().CreateLanguagesVersion(context.Background(), int(testID), testName).Times(1).Return(nil)

		err := h.postProjectVersion(testCtx, logrus.NewEntry(logrus.New()))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("postProjectVersion, error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"id":   testID,
			"name": testName,
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		resp := httptest.NewRecorder()

		testCtx := e.NewContext(req, resp)

		testCtx.SetPath("/projects/:project_id/versions")
		testCtx.SetParamNames("project_id")
		testCtx.SetParamValues(fmt.Sprintf("%d", testID))

		projectService.EXPECT().CreateLanguagesVersion(context.Background(), int(testID), testName).Times(1).Return(errTest)

		err := h.postProjectVersion(testCtx, logrus.NewEntry(logrus.New()))
		assert.Error(t, err)
	})

	t.Run("postProjectVersion, validation error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"id":   testID,
			"name": "",
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		resp := httptest.NewRecorder()

		testCtx := e.NewContext(req, resp)

		testCtx.SetPath("/projects/:project_id/versions")
		testCtx.SetParamNames("project_id")
		testCtx.SetParamValues(fmt.Sprintf("%d", testID))

		err := h.postProjectVersion(testCtx, logrus.NewEntry(logrus.New()))
		assert.Error(t, err)
	})
}
