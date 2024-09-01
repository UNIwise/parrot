package project

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/jarcoal/httpmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/uniwise/parrot/internal/storage"
	"github.com/uniwise/parrot/pkg/poedit"
	gomock "go.uber.org/mock/gomock"
)

var (
	errTest                             = errors.New("test error")
	testCtx                             = context.Background()
	testID                    uint      = 1
	testName                  string    = "testname"
	testProjectID             uint      = 1
	testStorageKeyForDeletion string    = "1/test-uuid"
	testStorageKeyForListing  string    = "1/"
	testNumberOfVersions      int       = 3
	testCreatedAt             time.Time = time.Now()
	testRenewalThreshold                = time.Hour
	testUUID                            = "test-uuid"
	testGenerateUUID                    = func() (string, error) {
		return testUUID, nil
	}
	testGenerateTimestamp = func() int64 {
		return testCreatedAt.Unix()
	}
	testContentMetaMap = map[string]poedit.ContentMeta{
		"key_value_json":  {Extension: "json", Type: "application/json"},
	}
	testGetContentMetaMap = func() map[string]poedit.ContentMeta {
		return testContentMetaMap
	}
)

func TestServiceGetAllProjects(t *testing.T) {
	t.Parallel()

	storage := storage.NewMockStorage(gomock.NewController(t))
	client := poedit.NewMockClient(gomock.NewController(t))

	service := NewService(client, storage, nil, testRenewalThreshold, nil)

	listProjectsResponse := &poedit.ListProjectsResponse{
		Result: struct {
			Projects []struct {
				ID      int64  `json:"id"`
				Name    string `json:"name"`
				Public  int64  `json:"public"`
				Open    int64  `json:"open"`
				Created string `json:"created"`
			} `json:"projects"`
		}{
			Projects: []struct {
				ID      int64  `json:"id"`
				Name    string `json:"name"`
				Public  int64  `json:"public"`
				Open    int64  `json:"open"`
				Created string `json:"created"`
			}{
				{
					ID:      int64(testID),
					Name:    testName,
					Created: testCreatedAt.String(),
				},
			},
		},
	}

	prefix := fmt.Sprintf("%d/%s_%s_%d", testID, testUUID, testName, testGenerateTimestamp())

	listObjectsV2Output := &s3.ListObjectsV2Output{
		CommonPrefixes: []types.CommonPrefix{
			{
				Prefix: &prefix,
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		client.EXPECT().ListProjects(testCtx).Return(listProjectsResponse, nil)

		storage.EXPECT().ListObjects(testCtx, testStorageKeyForListing).Times(1).Return(listObjectsV2Output, nil)

		projects, err := service.GetAllProjects(testCtx)

		assert.NoError(t, err)
		assert.Equal(t, int64(testID), projects[0].ID)
		assert.Equal(t, testName, projects[0].Name)
		assert.Equal(t, 1, projects[0].NumberOfVersions)
		assert.Equal(t, testCreatedAt.String(), projects[0].CreatedAt)
	})

	t.Run("Failure, client fail", func(t *testing.T) {
		client.EXPECT().ListProjects(testCtx).Return(nil, errTest)

		projects, err := service.GetAllProjects(testCtx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errTest)
		assert.Nil(t, projects)
	})

	t.Run("Failure, storage fail", func(t *testing.T) {
		client.EXPECT().ListProjects(testCtx).Return(listProjectsResponse, nil)

		storage.EXPECT().ListObjects(testCtx, testStorageKeyForListing).Times(1).Return(nil, errTest)

		projects, err := service.GetAllProjects(testCtx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errTest)
		assert.Nil(t, projects)
	})
}

func TestServiceGetProjectById(t *testing.T) {
	t.Parallel()

	storage := storage.NewMockStorage(gomock.NewController(t))
	client := poedit.NewMockClient(gomock.NewController(t))
	service := NewService(client, storage, nil, testRenewalThreshold, nil)

	projectResponse := &poedit.ViewProjectResponse{
		Result: struct {
			Project struct {
				ID                int64  `json:"id"`
				Name              string `json:"name"`
				Description       string `json:"description"`
				Public            int64  `json:"public"`
				Open              int64  `json:"open"`
				ReferenceLanguage string `json:"reference_language"` // nolint:tagliatelle
				Terms             int64  `json:"terms"`
				Created           string `json:"created"`
			} `json:"project"`
		}{
			Project: struct {
				ID                int64  `json:"id"`
				Name              string `json:"name"`
				Description       string `json:"description"`
				Public            int64  `json:"public"`
				Open              int64  `json:"open"`
				ReferenceLanguage string `json:"reference_language"` // nolint:tagliatelle
				Terms             int64  `json:"terms"`
				Created           string `json:"created"`
			}{
				ID:      int64(testID),
				Name:    testName,
				Created: testCreatedAt.String(),
			},
		},
	}

	prefix := fmt.Sprintf("%d/%s_%s_%d", testID, testUUID, testName, testGenerateTimestamp())

	listObjectsV2Output := &s3.ListObjectsV2Output{
		CommonPrefixes: []types.CommonPrefix{
			{
				Prefix: &prefix,
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		client.EXPECT().ViewProject(testCtx, poedit.ViewProjectRequest{
			ID: int(testID),
		}).Return(projectResponse, nil)

		storage.EXPECT().ListObjects(testCtx, testStorageKeyForListing).Times(1).Return(listObjectsV2Output, nil)

		project, err := service.GetProjectByID(testCtx, int(testID))

		assert.NoError(t, err)
		assert.Equal(t, int64(testID), (*project).ID)
		assert.Equal(t, testName, (*project).Name)
		assert.Equal(t, 1, (*project).NumberOfVersions)
		assert.Equal(t, testCreatedAt.String(), (*project).CreatedAt)
	})

	t.Run("Failure, client fail", func(t *testing.T) {
		client.EXPECT().ViewProject(testCtx, poedit.ViewProjectRequest{
			ID: int(testID),
		}).Return(nil, errTest)

		project, err := service.GetProjectByID(testCtx, int(testID))

		assert.Error(t, err)
		assert.Nil(t, project)
	})

	t.Run("Failure, storage fail", func(t *testing.T) {
		client.EXPECT().ViewProject(testCtx, poedit.ViewProjectRequest{
			ID: int(testID),
		}).Return(projectResponse, nil)

		storage.EXPECT().ListObjects(testCtx, testStorageKeyForListing).Times(1).Return(nil, errTest)

		project, err := service.GetProjectByID(testCtx, int(testID))

		assert.Error(t, err)
		assert.ErrorIs(t, err, errTest)
		assert.Nil(t, project)
	})
}

func TestServiceGetProjectVersions(t *testing.T) {
	t.Parallel()

	storage := storage.NewMockStorage(gomock.NewController(t))
	service := NewService(nil, storage, nil, testRenewalThreshold, nil)

	prefix := fmt.Sprintf("%d/%s_%s_%d", testID, testUUID, testName, testGenerateTimestamp())

	listObjectsV2Output := &s3.ListObjectsV2Output{
		CommonPrefixes: []types.CommonPrefix{
			{
				Prefix: &prefix,
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		storage.EXPECT().ListObjects(testCtx, testStorageKeyForListing).Times(1).Return(listObjectsV2Output, nil)

		versions, err := service.GetProjectVersions(testCtx, int(testProjectID))

		assert.NoError(t, err)
		assert.Equal(t, testUUID, versions[0].ID)
		assert.Equal(t, testName, versions[0].Name)
		assert.Equal(t, testCreatedAt.Unix(), versions[0].CreatedAt.Unix())
	})

	t.Run("Failure", func(t *testing.T) {
		storage.EXPECT().ListObjects(testCtx, testStorageKeyForListing).Times(1).Return(nil, errTest)

		versions, err := service.GetProjectVersions(testCtx, int(testProjectID))

		assert.Error(t, err)
		assert.ErrorIs(t, err, errTest)
		assert.Nil(t, versions)
	})
}

func TestServiceDeleteProjectVersionByIDAndVersionID(t *testing.T) {
	t.Parallel()

	storage := storage.NewMockStorage(gomock.NewController(t))
	service := NewService(nil, storage, nil, testRenewalThreshold, nil)

	t.Run("fail, fails to delete objects in S3", func(t *testing.T) {
		storage.EXPECT().DeleteObjects(testCtx, testStorageKeyForDeletion).Times(1).Return(errTest)

		err := service.DeleteProjectVersionByIDAndProjectID(testCtx, testUUID, testProjectID)

		assert.ErrorIs(t, err, errTest)
		assert.ErrorContains(t, err, "Failed to delete project version in S3")
	})

	t.Run("success", func(t *testing.T) {
		storage.EXPECT().DeleteObjects(testCtx, testStorageKeyForDeletion).Times(1).Return(nil)

		err := service.DeleteProjectVersionByIDAndProjectID(testCtx, testUUID, testProjectID)

		assert.NoError(t, err)
	})
}

func TestServiceCreateLanguagesVersion(t *testing.T) {
	t.Parallel()

	storage := storage.NewMockStorage(gomock.NewController(t))
	poeditClient := poedit.NewMockClient(gomock.NewController(t))

	service := NewService(poeditClient, storage, nil, testRenewalThreshold, nil)
	service.generateUUID = testGenerateUUID
	service.generateTimestamp = testGenerateTimestamp
	service.getContentMetaMap = testGetContentMetaMap

	listProjectLanguagesRequest := poedit.ListProjectLanguagesRequest{
		ID: int(testProjectID),
	}

	listAvailableLanguagesResponse := &poedit.ListProjectLanguagesResponse{
		Result: struct {
			Languages []struct {
				Name         string  `json:"name"`
				Code         string  `json:"code"`
				Translations int64   `json:"translations"`
				Percentage   float64 `json:"percentage"`
				Updated      string  `json:"updated"`
			} `json:"languages"`
		}{
			Languages: []struct {
				Name         string  `json:"name"`
				Code         string  `json:"code"`
				Translations int64   `json:"translations"`
				Percentage   float64 `json:"percentage"`
				Updated      string  `json:"updated"`
			}{
				{
					Name: "English",
					Code: "en",
				},
			},
		},
	}

	exportProjectRequest := poedit.ExportProjectRequest{
		ID:       int(testProjectID),
		Language: "en",
		Type:     "key_value_json",
		Filters:  []string{"translated"},
	}

	exportProjectResponse := &poedit.ExportProjectResponse{
		Result: struct {
			URL string `json:"url"`
		}{
			URL: "http://example.com/file.json",
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fileName := fmt.Sprintf("%d/%s_%s_%d/%s.json", testID, testUUID, testName, testGenerateTimestamp(), "en")

	t.Run("Success", func(t *testing.T) {
		poeditClient.EXPECT().ListProjectLanguages(testCtx, listProjectLanguagesRequest).Return(listAvailableLanguagesResponse, nil)
		poeditClient.EXPECT().ExportProject(testCtx, exportProjectRequest).Return(exportProjectResponse, nil)

		httpmock.RegisterResponder("GET", "http://example.com/file.json",
			httpmock.NewStringResponder(200, `{"key":"value"}`))

		storage.EXPECT().PutObject(testCtx, fileName, gomock.Any(), gomock.Any(), "application/json").Return(nil)

		err := service.CreateLanguagesVersion(testCtx, int(testProjectID), testName)

		assert.NoError(t, err)
	})

	t.Run("Fail, fails to list project languages", func(t *testing.T) {
		poeditClient.EXPECT().ListProjectLanguages(testCtx, listProjectLanguagesRequest).Return(nil, errTest)

		err := service.CreateLanguagesVersion(testCtx, int(testProjectID), testName)

		assert.ErrorIs(t, err, errTest)
	})

	t.Run("Fail, fails to export project", func(t *testing.T) {
		poeditClient.EXPECT().ListProjectLanguages(testCtx, listProjectLanguagesRequest).Return(listAvailableLanguagesResponse, nil)
		poeditClient.EXPECT().ExportProject(testCtx, exportProjectRequest).Return(nil, errTest)

		err := service.CreateLanguagesVersion(testCtx, int(testProjectID), testName)

		assert.ErrorIs(t, err, errTest)
	})

	t.Run("Fail, fails to download file", func(t *testing.T) {
		poeditClient.EXPECT().ListProjectLanguages(testCtx, listProjectLanguagesRequest).Return(listAvailableLanguagesResponse, nil)
		poeditClient.EXPECT().ExportProject(testCtx, exportProjectRequest).Return(exportProjectResponse, nil)

		httpmock.RegisterResponder("GET", "http://example.com/file.json",
			httpmock.NewErrorResponder(errTest))

		err := service.CreateLanguagesVersion(testCtx, int(testProjectID), testName)

		assert.ErrorIs(t, err, errTest)
	})

	t.Run("Fail, fails to upload file to S3", func(t *testing.T) {
		poeditClient.EXPECT().ListProjectLanguages(testCtx, listProjectLanguagesRequest).Return(listAvailableLanguagesResponse, nil)
		poeditClient.EXPECT().ExportProject(testCtx, exportProjectRequest).Return(exportProjectResponse, nil)

		httpmock.RegisterResponder("GET", "http://example.com/file.json",
			httpmock.NewStringResponder(200, `{"key":"value"}`))

		storage.EXPECT().PutObject(testCtx, fileName, gomock.Any(), gomock.Any(), "application/json").Return(errTest)

		err := service.CreateLanguagesVersion(testCtx, int(testProjectID), testName)

		assert.ErrorIs(t, err, errTest)
	})
}
