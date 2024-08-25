package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	client S3API
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

// PutObject creates an S3 object
func (c *S3ClientImpl) PutObject(ctx context.Context, key string, payloadReader io.Reader, mimeType string) error {
	if _, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Body:        payloadReader,
		Bucket:      aws.String(c.config.Bucket),
		Key:         aws.String(key),
		ContentType: &mimeType,
	}); err != nil {
		return errors.Wrap(err, "Failed to put object to S3")
	}

	return nil
}

// DeleteObject deletes a single S3 object
func (c *S3ClientImpl) DeleteObject(ctx context.Context, key string) error {
	if _, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.config.Bucket),
		Key:    aws.String(key),
	}); err != nil {
		return errors.Wrap(err, "Failed to delete object into S3")
	}

	return nil
}

// GetObject gets a single S3 object
func (c *S3ClientImpl) GetObject(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	doc, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.config.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get object from S3")
	}

	return doc, err
}
