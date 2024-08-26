package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ServiceImpl struct {
	storage Storage
}

type Service interface {
	PutObject(ctx context.Context, key string, payloadReader io.Reader, mimeType string) error
	DeleteObject(ctx context.Context, key string) error
	GetObject(ctx context.Context, key string) (*s3.GetObjectOutput, error)
}

func NewService(ctx context.Context, storage Storage) *ServiceImpl {
	return &ServiceImpl{
		storage: storage,
	}
}

func (s *ServiceImpl) PutObject(ctx context.Context, key string, payloadReader io.Reader, mimeType string) error {
	return s.storage.PutObject(ctx, key, payloadReader, mimeType)
}

func (s *ServiceImpl) DeleteObject(ctx context.Context, key string) error {
	return s.storage.DeleteObject(ctx, key)
}

func (s *ServiceImpl) GetObject(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	return s.storage.GetObject(ctx, key)
}