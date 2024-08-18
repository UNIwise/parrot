package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewS3Client(t *testing.T) {
	t.Parallel()

	client, err := NewS3Client(context.Background(), S3StorageConfig{})

	assert.NoError(t, err)
	assert.NotNil(t, client.client)
	assert.NotNil(t, client.config)
}