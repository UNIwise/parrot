package storage

import (
	"context"
	"errors"
	"testing"
	"strings"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

var errForTesting = errors.New("this is an error for testing")

type MockS3APIMethods struct {
	DeleteObjectFunction func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	PutObjectFunction    func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObjectFunction    func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func (m *MockS3APIMethods) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return m.DeleteObjectFunction(ctx, params, optFns...)
}

func (m *MockS3APIMethods) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return m.PutObjectFunction(ctx, params, optFns...)
}

func (m *MockS3APIMethods) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m.GetObjectFunction(ctx, params, optFns...)
}

func TestNewS3Client(t *testing.T) {
	t.Parallel()

	client, err := NewS3Client(context.Background(), S3StorageConfig{})

	assert.NoError(t, err)
	assert.NotNil(t, client.client)
	assert.NotNil(t, client.config)
}

func TestDeleteObject(t *testing.T) {
	t.Parallel()

	var (
		ctx    context.Context
		bucket string = "test-bucket"
		region string = "test-region"
		key    string = "test/keyOne"
	)

	t.Run("Delete object, fail", func(t *testing.T) {
		mockClient := &MockS3APIMethods{

			DeleteObjectFunction: func(
				_ctx context.Context,
				_params *s3.DeleteObjectInput,
				_optFns ...func(*s3.Options),
			) (*s3.DeleteObjectOutput, error) {
				assert.Equal(t, ctx, _ctx)
				assert.Equal(t, bucket, *_params.Bucket)
				return nil, errForTesting
			},
		}

		storage := &S3ClientImpl{
			config: S3StorageConfig{Bucket: bucket, Region: region},
			client: mockClient,
		}

		actualError := storage.DeleteObject(ctx, key)

		assert.Error(t, actualError)
	})
	t.Run("Delete object, success", func(t *testing.T) {
		mockClient := &MockS3APIMethods{

			DeleteObjectFunction: func(
				_ctx context.Context,
				_params *s3.DeleteObjectInput,
				_optFns ...func(*s3.Options),
			) (*s3.DeleteObjectOutput, error) {
				assert.Equal(t, ctx, _ctx)
				assert.Equal(t, bucket, *_params.Bucket)
				return &s3.DeleteObjectOutput{}, nil
			},
		}

		storage := &S3ClientImpl{
			config: S3StorageConfig{Bucket: bucket, Region: region},
			client: mockClient,
		}

		actualError := storage.DeleteObject(ctx, key)

		assert.NoError(t, actualError)
	})
}

func TestPutObject(t *testing.T) {
	t.Parallel()

	var (
		ctx    context.Context
		bucket string = "test-bucket"
		region string = "test-region"
		key    string = "test/keyOne"
		payloadReader io.Reader = strings.NewReader("Hello, Reader!")
	)

	t.Run("Put object, fail", func(t *testing.T) {
		mockClient := &MockS3APIMethods{

			PutObjectFunction: func(
				_ctx context.Context,
				_params *s3.PutObjectInput,
				_optFns ...func(*s3.Options),
			) (*s3.PutObjectOutput, error) {
				assert.Equal(t, ctx, _ctx)
				assert.Equal(t, bucket, *_params.Bucket)
				assert.Equal(t, payloadReader, _params.Body)
				return nil, errForTesting
			},
		}

		storage := &S3ClientImpl{
			config: S3StorageConfig{Bucket: bucket, Region: region},
			client: mockClient,
		}

		actualError := storage.PutObject(ctx, key, payloadReader, "txt")

		assert.Error(t, actualError)
	})
	t.Run("Put object, success", func(t *testing.T) {
		mockClient := &MockS3APIMethods{

			PutObjectFunction: func(
				_ctx context.Context,
				_params *s3.PutObjectInput,
				_optFns ...func(*s3.Options),
			) (*s3.PutObjectOutput, error) {
				assert.Equal(t, ctx, _ctx)
				assert.Equal(t, bucket, *_params.Bucket)
				assert.Equal(t, payloadReader, _params.Body)
				return &s3.PutObjectOutput{}, nil
			},
		}

		storage := &S3ClientImpl{
			config: S3StorageConfig{Bucket: bucket, Region: region},
			client: mockClient,
		}

		actualError := storage.PutObject(ctx, key, payloadReader, "txt")

		assert.NoError(t, actualError)
	})
}

func TestGetObject(t *testing.T) {
	t.Parallel()

	var (
		ctx    context.Context
		bucket string = "test-bucket"
		region string = "test-region"
		key    string = "test/keyOne"
	)

	t.Run("Get object, fail", func(t *testing.T) {
		mockClient := &MockS3APIMethods{

			GetObjectFunction: func(
				_ctx context.Context,
				_params *s3.GetObjectInput,
				_optFns ...func(*s3.Options),
			) (*s3.GetObjectOutput, error) {
				assert.Equal(t, ctx, _ctx)
				assert.Equal(t, bucket, *_params.Bucket)
				assert.Equal(t, key, *_params.Key)
				return nil, errForTesting
			},
		}

		storage := &S3ClientImpl{
			config: S3StorageConfig{Bucket: bucket, Region: region},
			client: mockClient,
		}

		actualOutput, actualError := storage.GetObject(ctx, key)

		assert.Error(t, actualError)
		assert.Nil(t, actualOutput)
	})
	t.Run("Get object, success", func(t *testing.T) {
		mockClient := &MockS3APIMethods{

			GetObjectFunction: func(
				_ctx context.Context,
				_params *s3.GetObjectInput,
				_optFns ...func(*s3.Options),
			) (*s3.GetObjectOutput, error) {
				assert.Equal(t, ctx, _ctx)
				assert.Equal(t, bucket, *_params.Bucket)
				assert.Equal(t, key, *_params.Key)
				return &s3.GetObjectOutput{}, nil
			},
		}

		storage := &S3ClientImpl{
			config: S3StorageConfig{Bucket: bucket, Region: region},
			client: mockClient,
		}

		actualOutput, actualError := storage.GetObject(ctx, key)

		assert.NoError(t, actualError)
		assert.NotNil(t, actualOutput)
	})
}
