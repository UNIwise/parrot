package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	gomock "go.uber.org/mock/gomock"
)

var (
	testCtx = context.Background()
	testKey = "test/key"
)

func TestService_PutObject(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	storage := NewMockStorage(ctrl)
	service := NewService(testCtx, storage)

	t.Run("PutObject should succeed", func(t *testing.T) {
		storage.EXPECT().PutObject(testCtx, testKey, nil, "").Times(1).Return(nil)

		err := service.PutObject(testCtx, testKey, nil, "")
		assert.NoError(t, err)
	})

	t.Run("PutObject should fail", func(t *testing.T) {
		storage.EXPECT().PutObject(testCtx, testKey, nil, "").Times(1).Return(errForTesting)

		err := service.PutObject(testCtx, testKey, nil, "")
		assert.ErrorIs(t, err, errForTesting)
	})
}

func TestService_DeleteObject(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	storage := NewMockStorage(ctrl)
	service := NewService(testCtx, storage)

	t.Run("DeleteObject should succeed", func(t *testing.T) {
		storage.EXPECT().DeleteObject(testCtx, testKey).Times(1).Return(nil)

		err := service.DeleteObject(testCtx, testKey)
		assert.NoError(t, err)
	})

	t.Run("DeleteObject should fail", func(t *testing.T) {
		storage.EXPECT().DeleteObject(testCtx, testKey).Times(1).Return(errForTesting)

		err := service.DeleteObject(testCtx, testKey)
		assert.ErrorIs(t, err, errForTesting)
	})
}

func TestService_GetObject(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	storage := NewMockStorage(ctrl)
	service := NewService(testCtx, storage)

	t.Run("GetObject should succeed", func(t *testing.T) {
		getObjectOutput := &s3.GetObjectOutput{}
		storage.EXPECT().GetObject(testCtx, testKey).Times(1).Return(getObjectOutput, nil)

		output, err := service.GetObject(testCtx, testKey)
		assert.NoError(t, err)
		assert.Equal(t, getObjectOutput, output)
	})

	t.Run("GetObject should fail", func(t *testing.T) {
		storage.EXPECT().GetObject(testCtx, testKey).Times(1).Return(nil, errForTesting)

		_, err := service.GetObject(testCtx, testKey)
		assert.ErrorIs(t, err, errForTesting)
	})
}
