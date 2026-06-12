package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// localUploader writes files under a base directory and serves them through a
// URL prefix. It exists as a development fallback when R2 is not configured.
type localUploader struct {
	baseDir   string
	urlPrefix string
}

func newLocalUploader(baseDir, urlPrefix string) *localUploader {
	return &localUploader{baseDir: baseDir, urlPrefix: urlPrefix}
}

func (u *localUploader) Upload(_ context.Context, key string, body io.Reader, _ int64, _ string) (string, error) {
	dst := filepath.Join(u.baseDir, key)
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return "", err
	}

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, body); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", u.urlPrefix, key), nil
}
