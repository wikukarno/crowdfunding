package storage

import (
	"context"
	"io"
	"log/slog"

	"backend-crowdfunding/config"
)

// Uploader stores a file and returns a URL that can be used to fetch it back.
type Uploader interface {
	Upload(ctx context.Context, key string, body io.Reader, size int64, contentType string) (string, error)
}

// New returns an R2-backed uploader when credentials are configured, otherwise
// it falls back to local disk so development works without cloud storage.
func New(cfg *config.Config) (Uploader, error) {
	if cfg.Storage.Enabled() {
		slog.Info("storage: using Cloudflare R2", "bucket", cfg.Storage.Bucket)
		return newR2Uploader(cfg.Storage)
	}

	slog.Warn("storage: R2 not configured, falling back to local disk")
	return newLocalUploader("./images", "/images"), nil
}
