package project

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

var (
	errTest                        = errors.New("test error")
	testCtx                        = context.Background()
	testID               uint      = 1
	testName             string    = "testname"
	testNumberOfVersions uint      = 3
	testCreatedAt        time.Time = time.Now()
	testRenewalThreshold           = time.Hour
)

func TestServiceGetAllProjects(t *testing.T) {
	t.Parallel()

	repository := NewMockRepository(gomock.NewController(t))
	service := NewService(nil, nil, repository, nil, testRenewalThreshold, nil)

	t.Run("Success", func(t *testing.T) {
		repository.EXPECT().GetAllProjects(testCtx).Return(&[]Project{{
			ID:               testID,
			Name:             testName,
			NumberOfVersions: testNumberOfVersions,
			CreatedAt:        testCreatedAt,
		},
		}, nil)

		projects, err := service.GetAllProjects(testCtx)

		assert.NoError(t, err)
		assert.Equal(t, testID, (*projects)[0].ID)
		assert.Equal(t, testName, (*projects)[0].Name)
		assert.Equal(t, testNumberOfVersions, (*projects)[0].NumberOfVersions)
		assert.Equal(t, testCreatedAt, (*projects)[0].CreatedAt)
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
