//go:build !windows
// +build !windows

package cache

import (
	"context"

	"golang.org/x/sys/unix"
)

func (f *FilesystemCache) PingContext(ctx context.Context) error {
	return unix.Access(f.dir, unix.W_OK)
}
