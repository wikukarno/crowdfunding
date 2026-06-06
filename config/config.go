package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort  string
	Database DatabaseConfig
	JWT      JWTConfig
	Midtrans MidtransConfig
	Storage  StorageConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey string
}

type MidtransConfig struct {
	ServerKey string
	ClientKey string
}

// StorageConfig points at a Cloudflare R2 bucket (S3-compatible). When the
// credentials are empty the app falls back to local disk so development still
// works without cloud storage.
type StorageConfig struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	PublicURL       string
}

// Enabled reports whether R2 credentials are present.
func (s StorageConfig) Enabled() bool {
	return s.AccountID != "" && s.AccessKeyID != "" && s.SecretAccessKey != "" && s.Bucket != ""
}

// Endpoint builds the R2 S3 API endpoint from the account id.
func (s StorageConfig) Endpoint() string {
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com", s.AccountID)
}

// Load builds the configuration from environment variables. A .env file is
// read when present so local development works without exporting anything by
// hand, but real environment variables win and are the only source in prod.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		// Cloud Run (and similar platforms) inject the port via PORT; fall back
		// to APP_PORT for local runs.
		AppPort: getEnv("PORT", getEnv("APP_PORT", "8080")),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "crowdfunding"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
		},
		Midtrans: MidtransConfig{
			ServerKey: os.Getenv("MIDTRANS_SERVER_KEY"),
			ClientKey: os.Getenv("MIDTRANS_CLIENT_KEY"),
		},
		Storage: StorageConfig{
			AccountID:       os.Getenv("R2_ACCOUNT_ID"),
			AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
			Bucket:          os.Getenv("R2_BUCKET"),
			PublicURL:       os.Getenv("R2_PUBLIC_URL"),
		},
	}

	if cfg.JWT.SecretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY is required")
	}

	return cfg, nil
}

// DSN returns the PostgreSQL connection string consumed by GORM.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
