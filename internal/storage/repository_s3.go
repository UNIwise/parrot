package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

type S3StorageConfig struct {
	Region string `mapstructure:"region" validate:"required"`
	Bucket string `mapstructure:"bucket" validate:"required"`
}

type S3ClientImpl struct {
	config S3StorageConfig
	client *s3.Client
}

func NewS3Client(ctx context.Context, conf S3StorageConfig) (*S3ClientImpl, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(conf.Region))
	if err != nil {
		return nil, errors.Wrap(err, "failed to load configuration")
	}

	client := s3.NewFromConfig(cfg)

	return &S3ClientImpl{
		config: conf,
		client: client,
	}, nil
}

func (s *S3ClientImpl) ListObjectsV2(ctx context.Context) (*s3.ListObjectsV2Output, error) {
	
	response, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &s.config.Bucket,
	})

	if err != nil {
		return &s3.ListObjectsV2Output{}, errors.Wrap(err, "failed to get data")
	}

	return response, nil
}
