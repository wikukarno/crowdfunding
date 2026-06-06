package database

import (
	"fmt"
	"log/slog"
	"time"

	"backend-crowdfunding/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	maxAttempts   = 15
	retryInterval = 2 * time.Second
)

// NewConnection opens the database and retries for a short while. A fresh MySQL
// container is not ready the instant it starts, so without this the app would
// crash on the first boot of a `docker compose up`.
func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
		if err == nil {
			return db, nil
		}

		lastErr = err
		slog.Warn("database not ready, retrying",
			"attempt", attempt,
			"max_attempts", maxAttempts,
			"error", err,
		)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("connect database after %d attempts: %w", maxAttempts, lastErr)
}
