//go:build windows
// +build windows

package cache

import "context"

func (f *FilesystemCache) PingContext(ctx context.Context) error {
	return nil
}
