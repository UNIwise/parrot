package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ServiceImpl struct {
	storage Storage
}

type Service interface {
	ListObjectsV2(ctx context.Context)  (*s3.ListObjectsV2Output, error)
}

func NewService(ctx context.Context, storage Storage) *ServiceImpl {
	return &ServiceImpl{
		storage: storage,
	}
}

// TODO: Add methods to implement