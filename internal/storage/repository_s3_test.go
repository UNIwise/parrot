package storage

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

var errForTesting = errors.New("this is an error for testing")

type MockS3APIMethods struct {
	DeleteObjectFunction            func(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	PutObjectFunction               func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObjectFunction               func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	ListObjectsV2Function           func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	DeleteObjectsFunction           func(ctx context.Context, params *s3.DeleteObjectsInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error)
	AbortMultipartUploadFunction    func(ctx context.Context, params *s3.AbortMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.AbortMultipartUploadOutput, error)
	CompleteMultipartUploadFunction func(ctx context.Context, params *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error)
	CreateMultipartUploadFunction   func(ctx context.Context, params *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error)
	UploadPartFunction              func(ctx context.Context, params *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error)
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

func (m *MockS3APIMethods) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	return m.ListObjectsV2Function(ctx, params, optFns...)
}

func (m *MockS3APIMethods) DeleteObjects(ctx context.Context, params *s3.DeleteObjectsInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {
	return m.DeleteObjectsFunction(ctx, params, optFns...)
}

func (m *MockS3APIMethods) AbortMultipartUpload(ctx context.Context, params *s3.AbortMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.AbortMultipartUploadOutput, error) {
	return m.AbortMultipartUploadFunction(ctx, params, optFns...)
}

func (m *MockS3APIMethods) CompleteMultipartUpload(ctx context.Context, params *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error) {
	return m.CompleteMultipartUploadFunction(ctx, params, optFns...)
}

func (m *MockS3APIMethods) CreateMultipartUpload(ctx context.Context, params *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
	return m.CreateMultipartUploadFunction(ctx, params, optFns...)
}

func (m *MockS3APIMethods) UploadPart(ctx context.Context, params *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error) {
	return m.UploadPartFunction(ctx, params, optFns...)
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
		ctx           context.Context
		bucket        string    = "test-bucket"
		region        string    = "test-region"
		key           string    = "test/keyOne"
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

		actualError := storage.PutObject(ctx, key, payloadReader, nil, "txt")

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

		actualError := storage.PutObject(ctx, key, payloadReader, nil, "txt")

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

func TestListObjects(t *testing.T) {
	t.Parallel()

	var (
		ctx        context.Context
		bucket     string = "test-bucket"
		region     string = "test-region"
		storageKey string = "test/key"
	)

	t.Run("List objects, fail", func(t *testing.T) {
		mockClient := &MockS3APIMethods{
			ListObjectsV2Function: func(
				_ctx context.Context,
				_params *s3.ListObjectsV2Input,
				_optFns ...func(*s3.Options),
			) (*s3.ListObjectsV2Output, error) {
				assert.Equal(t, ctx, _ctx)
				assert.Equal(t, bucket, *_params.Bucket)
				return nil, errForTesting
			},
		}

		storage := &S3ClientImpl{
			config: S3StorageConfig{Bucket: bucket, Region: region},
			client: mockClient,
		}

		actualOutput, actualError := storage.ListObjects(ctx, storageKey)

		assert.Error(t, actualError)
		assert.Nil(t, actualOutput)
	})
	t.Run("List objects, success", func(t *testing.T) {
		mockClient := &MockS3APIMethods{
			ListObjectsV2Function: func(
				_ctx context.Context,
				_params *s3.ListObjectsV2Input,
				_optFns ...func(*s3.Options),
			) (*s3.ListObjectsV2Output, error) {
				assert.Equal(t, ctx, _ctx)
				assert.Equal(t, bucket, *_params.Bucket)
				return &s3.ListObjectsV2Output{}, nil
			},
		}

		storage := &S3ClientImpl{
			config: S3StorageConfig{Bucket: bucket, Region: region},
			client: mockClient,
		}

		actualOutput, actualError := storage.ListObjects(ctx, storageKey)

		assert.NoError(t, actualError)
		assert.NotNil(t, actualOutput)
	})
}

func TestDeleteObjects(t *testing.T) {
	t.Parallel()

	var (
		ctx    context.Context
		bucket string = "test-bucket"
		region string = "test-region"
		key    string = "test/key"
	)

	t.Run("Delete objects, fail", func(t *testing.T) {
		mockClient := &MockS3APIMethods{
			ListObjectsV2Function: func(
				_ctx context.Context,
				_params *s3.ListObjectsV2Input,
				_optFns ...func(*s3.Options),
			) (*s3.ListObjectsV2Output, error) {
				assert.Equal(t, ctx, _ctx)
				assert.Equal(t, bucket, *_params.Bucket)
				return nil, errForTesting
			},
		}

		storage := &S3ClientImpl{
			config: S3StorageConfig{Bucket: bucket, Region: region},
			client: mockClient,
		}

		actualError := storage.DeleteObjects(ctx, key)

		assert.Error(t, actualError)
	})
	t.Run("Delete objects, success", func(t *testing.T) {
		mockClient := &MockS3APIMethods{
			ListObjectsV2Function: func(
				_ctx context.Context,
				_params *s3.ListObjectsV2Input,
				_optFns ...func(*s3.Options),
			) (*s3.ListObjectsV2Output, error) {
				assert.Equal(t, ctx, _ctx)
				assert.Equal(t, bucket, *_params.Bucket)
				return &s3.ListObjectsV2Output{}, nil
			},
			DeleteObjectsFunction: func(
				_ctx context.Context,
				_params *s3.DeleteObjectsInput,
				_optFns ...func(*s3.Options),
			) (*s3.DeleteObjectsOutput, error) {
				assert.Equal(t, ctx, _ctx)
				assert.Equal(t, bucket, *_params.Bucket)
				return &s3.DeleteObjectsOutput{}, nil
			},
		}

		storage := &S3ClientImpl{
			config: S3StorageConfig{Bucket: bucket, Region: region},
			client: mockClient,
		}

		actualError := storage.DeleteObjects(ctx, key)

		assert.NoError(t, actualError)
	})
}
