//go:generate mockgen --source=repository.go -destination=repository_mock.go -package=storage
package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3API interface {
	DeleteObject(
		ctx context.Context,
		params *s3.DeleteObjectInput,
		optFns ...func(*s3.Options),
	) (*s3.DeleteObjectOutput, error)

	PutObject(
		ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options),
	) (*s3.PutObjectOutput, error)

	GetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.Options),
	) (*s3.GetObjectOutput, error)

	ListObjectsV2(
		ctx context.Context,
		params *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options),
	) (*s3.ListObjectsV2Output, error)

	DeleteObjects(
		ctx context.Context,
		params *s3.DeleteObjectsInput,
		optFns ...func(*s3.Options),
	) (*s3.DeleteObjectsOutput, error)
}

type Storage interface {
	PutObject(ctx context.Context, key string, payloadReader io.Reader, metadata map[string]string, mimeType string) error
	DeleteObject(ctx context.Context, key string) error
	GetObject(ctx context.Context, key string) (*s3.GetObjectOutput, error)
	ListObjects(ctx context.Context, storageKey string) (*s3.ListObjectsV2Output, error)
	DeleteObjects(ctx context.Context, key string) error
}
