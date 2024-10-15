package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"
)

const batchSize = 1000 // Max objects to delete at once, restricted by AWS

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
func (c *S3ClientImpl) PutObject(ctx context.Context, key string, payloadReader io.Reader, metadata map[string]string, mimeType string) error {
	uploader := manager.NewUploader(c.client)

	if _, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Body:        payloadReader,
		Bucket:      aws.String(c.config.Bucket),
		Key:         aws.String(key),
		Metadata:    metadata,
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

func (s *S3ClientImpl) DeleteObjects(ctx context.Context, key string) error {
	output, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.config.Bucket),
		Prefix: &key,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to list objects")
	}

	keys := make([]string, len(output.Contents))
	for i, object := range output.Contents {
		keys[i] = aws.ToString(object.Key)
	}

	objects := make([]types.ObjectIdentifier, len(keys))
	for i, key := range keys {
		objects[i] = types.ObjectIdentifier{
			Key: aws.String(key),
		}
	}

	for i := 0; i < len(objects); i += batchSize {
		j := i + batchSize
		if j > len(objects) {
			j = len(objects)
		}

		if _, err := s.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(s.config.Bucket),
			Delete: &types.Delete{
				Objects: objects[i:j],
				Quiet:   aws.Bool(true),
			},
		}); err != nil {
			return errors.Wrap(err, "Failed to delete objects")
		}
	}

	return nil
}

// GetObject gets a single S3 object
func (s *S3ClientImpl) GetObject(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	object, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get object from S3")
	}

	return object, err
}

// ListObjects lists all objects in an S3 bucket
func (s *S3ClientImpl) ListObjects(ctx context.Context, storageKey string) (*s3.ListObjectsV2Output, error) {
	output, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.config.Bucket),
		Prefix:    &storageKey,
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to list objects in S3")
	}

	return output, nil
}

func (s *S3ClientImpl) AbortMultipartUpload(ctx context.Context, uploadID string) error {
	_, err := s.client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(s.config.Bucket),
		Key:      aws.String(uploadID),
		UploadId: aws.String(uploadID),
	})
	if err != nil {
		return errors.Wrap(err, "Failed to abort multipart upload")
	}

	return nil
}

func (s *S3ClientImpl) CompleteMultipartUpload(ctx context.Context, input *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error) {
	output, err := s.client.CompleteMultipartUpload(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to complete multipart upload")
	}

	return output, nil
}

func (s *S3ClientImpl) CreateMultipartUpload(ctx context.Context, input *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
	output, err := s.client.CreateMultipartUpload(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create multipart upload")
	}

	return output, nil
}

func (s *S3ClientImpl) UploadPart(ctx context.Context, input *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
	output, err := s.client.UploadPart(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to upload part")
	}

	return output, nil
}
