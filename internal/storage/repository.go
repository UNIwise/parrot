//go:generate mockgen --source=repository.go -destination=repository_mock.go -package=storage
package storage

import (
	"context"
	
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage interface {
	ListObjectsV2(ctx context.Context)  (*s3.ListObjectsV2Output, error)
}
