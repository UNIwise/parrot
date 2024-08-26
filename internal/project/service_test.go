package project

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/uniwise/parrot/internal/storage"
	"github.com/uniwise/parrot/pkg/connectors/database"
	gomock "go.uber.org/mock/gomock"
)

var (
	errTest                        = errors.New("test error")
	testCtx                        = context.Background()
	testID               uint      = 1
	testName             string    = "testname"
	testProjectID        uint      = 1
	testStorageKey       string    = "testkey"
	testNumberOfVersions uint      = 3
	testCreatedAt        time.Time = time.Now()
	testRenewalThreshold           = time.Hour
)

func TestServiceGetAllProjects(t *testing.T) {
	t.Parallel()

	repository := NewMockRepository(gomock.NewController(t))
	service := NewService(nil, nil, repository, nil, testRenewalThreshold, nil)

	t.Run("Success", func(t *testing.T) {
		repository.EXPECT().GetAllProjects(testCtx).Return([]Project{{
			ID:               testID,
			Name:             testName,
			NumberOfVersions: testNumberOfVersions,
			CreatedAt:        testCreatedAt,
		},
		}, nil)

		projects, err := service.GetAllProjects(testCtx)

		assert.NoError(t, err)
		assert.Equal(t, testID, projects[0].ID)
		assert.Equal(t, testName, projects[0].Name)
		assert.Equal(t, testNumberOfVersions, projects[0].NumberOfVersions)
		assert.Equal(t, testCreatedAt, projects[0].CreatedAt)
	})

	t.Run("Not found", func(t *testing.T) {
		repository.EXPECT().GetAllProjects(testCtx).Return(nil, ErrNotFound)

		projects, err := service.GetAllProjects(testCtx)

		assert.Error(t, err)
		assert.Nil(t, projects)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("Failure", func(t *testing.T) {
		repository.EXPECT().GetAllProjects(testCtx).Return(nil, errTest)

		projects, err := service.GetAllProjects(testCtx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errTest)
		assert.Nil(t, projects)
	})
}

func TestServiceGetProjectById(t *testing.T) {
	t.Parallel()

	repository := NewMockRepository(gomock.NewController(t))
	service := NewService(nil, nil, repository, nil, testRenewalThreshold, nil)

	t.Run("Success", func(t *testing.T) {
		repository.EXPECT().GetProjectByID(testCtx, int(testID)).Return(&Project{
			ID:               testID,
			Name:             testName,
			NumberOfVersions: testNumberOfVersions,
			CreatedAt:        testCreatedAt,
		}, nil)

		project, err := service.GetProjectByID(testCtx, int(testID))

		assert.NoError(t, err)
		assert.Equal(t, testID, (*project).ID)
		assert.Equal(t, testName, (*project).Name)
		assert.Equal(t, testNumberOfVersions, (*project).NumberOfVersions)
		assert.Equal(t, testCreatedAt, (*project).CreatedAt)
	})

	t.Run("Not found", func(t *testing.T) {
		repository.EXPECT().GetProjectByID(testCtx, int(testID)).Return(nil, ErrNotFound)

		project, err := service.GetProjectByID(testCtx, int(testID))

		assert.Error(t, err)
		assert.Nil(t, project)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("Failure", func(t *testing.T) {
		repository.EXPECT().GetProjectByID(testCtx, int(testID)).Return(nil, errTest)

		project, err := service.GetProjectByID(testCtx, int(testID))

		assert.Error(t, err)
		assert.ErrorIs(t, err, errTest)
		assert.Nil(t, project)
	})
}

func TestServiceGetProjectVersions(t *testing.T) {
	t.Parallel()

	repository := NewMockRepository(gomock.NewController(t))
	service := NewService(nil, nil, repository, nil, testRenewalThreshold, nil)

	t.Run("Success", func(t *testing.T) {
		repository.EXPECT().GetProjectVersions(testCtx, int(testProjectID)).Return([]Version{{
			ID:        testID,
			Name:      testName,
			ProjectID: testProjectID,
			CreatedAt: testCreatedAt,
		},
		}, nil)

		versions, err := service.GetProjectVersions(testCtx, int(testProjectID))

		assert.NoError(t, err)
		assert.Equal(t, testID, versions[0].ID)
		assert.Equal(t, testName, versions[0].Name)
		assert.Equal(t, testProjectID, versions[0].ProjectID)
		assert.Equal(t, testCreatedAt, versions[0].CreatedAt)
	})

	t.Run("Not found", func(t *testing.T) {
		repository.EXPECT().GetProjectVersions(testCtx, int(testProjectID)).Return(nil, ErrNotFound)

		versions, err := service.GetProjectVersions(testCtx, int(testProjectID))

		assert.Error(t, err)
		assert.Nil(t, versions)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("Failure", func(t *testing.T) {
		repository.EXPECT().GetProjectVersions(testCtx, int(testProjectID)).Return(nil, errTest)

		versions, err := service.GetProjectVersions(testCtx, int(testProjectID))

		assert.Error(t, err)
		assert.ErrorIs(t, err, errTest)
		assert.Nil(t, versions)
	})
}

func TestServiceDeleteProjectVersionByIDAndVersionID(t *testing.T) {
	t.Parallel()

	repository := NewMockRepository(gomock.NewController(t))
	storage := storage.NewMockStorage(gomock.NewController(t))
	service := NewService(nil, storage, repository, nil, testRenewalThreshold, nil)

	version := &Version{
		ID:         testID,
		StorageKey: testStorageKey,
		ProjectID:  testProjectID,
		Name:       testName,
		CreatedAt:  testCreatedAt,
	}

	t.Run("fail, fails to get version", func(t *testing.T) {
		repository.EXPECT().GetVersionByIDAndProjectID(testCtx, testID, testProjectID).Times(1).Return(nil, errTest)

		err := service.DeleteProjectVersionByIDAndProjectID(testCtx, testID, testProjectID)

		assert.ErrorIs(t, err, errTest)
		assert.ErrorContains(t, err, fmt.Sprintf("Failed to retrieve project version with ID %d and project ID %d", testID, testProjectID))
	})

	t.Run("fail, version not found", func(t *testing.T) {
		repository.EXPECT().GetVersionByIDAndProjectID(testCtx, testID, testProjectID).Times(1).Return(nil, ErrNotFound)

		err := service.DeleteProjectVersionByIDAndProjectID(testCtx, testID, testProjectID)

		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("fail, fails to begin transaction", func(t *testing.T) {
		repository.EXPECT().GetVersionByIDAndProjectID(testCtx, testID, testProjectID).Times(1).Return(version, nil)
		repository.EXPECT().DeleteVersionByIDTransaction(testCtx, testID).Times(1).Return(nil, errTest)

		err := service.DeleteProjectVersionByIDAndProjectID(testCtx, testID, testProjectID)

		assert.ErrorIs(t, err, errTest)
		assert.ErrorContains(t, err, "Failed to begin delete project version transaction")
	})

	t.Run("fail, none deleted", func(t *testing.T) {
		repository.EXPECT().GetVersionByIDAndProjectID(testCtx, testID, testProjectID).Times(1).Return(version, nil)
		repository.EXPECT().DeleteVersionByIDTransaction(testCtx, testID).Times(1).Return(nil, ErrNotDeleted)

		err := service.DeleteProjectVersionByIDAndProjectID(testCtx, testID, testProjectID)

		assert.ErrorIs(t, err, ErrNotDeleted)
	})

	t.Run("fail, fails to delete objects in S3", func(t *testing.T) {
		db, sql := database.NewMockClient(t)
		sql.ExpectBegin()
		sql.ExpectRollback()
		testTx := db.WithContext(testCtx).Begin()

		repository.EXPECT().GetVersionByIDAndProjectID(testCtx, testID, testProjectID).Times(1).Return(version, nil)
		repository.EXPECT().DeleteVersionByIDTransaction(testCtx, testID).Times(1).Return(testTx, nil)

		storage.EXPECT().DeleteObject(testCtx, testStorageKey).Times(1).Return(errTest)

		err := service.DeleteProjectVersionByIDAndProjectID(testCtx, testID, testProjectID)

		assert.ErrorIs(t, err, errTest)
		assert.ErrorContains(t, err, "Failed to delete project version in S3")
	})

	t.Run("fail, fails to commit transaction", func(t *testing.T) {
		db, sql := database.NewMockClient(t)
		sql.ExpectBegin()
		sql.ExpectRollback()
		testTx := db.WithContext(testCtx).Begin()

		repository.EXPECT().GetVersionByIDAndProjectID(testCtx, testID, testProjectID).Times(1).Return(version, nil)
		repository.EXPECT().DeleteVersionByIDTransaction(testCtx, testID).Times(1).Return(testTx, nil)

		storage.EXPECT().DeleteObject(testCtx, testStorageKey).Times(1).Return(nil)

		err := service.DeleteProjectVersionByIDAndProjectID(testCtx, testID, testProjectID)

		assert.ErrorContains(t, err, "Failed to commit delete project version transaction")
	})

	t.Run("success", func(t *testing.T) {
		db, sql := database.NewMockClient(t)
		sql.ExpectBegin()
		sql.ExpectCommit()
		testTx := db.WithContext(testCtx).Begin()

		repository.EXPECT().GetVersionByIDAndProjectID(testCtx, testID, testProjectID).Times(1).Return(version, nil)
		repository.EXPECT().DeleteVersionByIDTransaction(testCtx, testID).Times(1).Return(testTx, nil)

		storage.EXPECT().DeleteObject(testCtx, testStorageKey).Times(1).Return(nil)

		err := service.DeleteProjectVersionByIDAndProjectID(testCtx, testID, testProjectID)

		assert.NoError(t, err)
	})
}
