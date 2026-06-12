package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort  string
	Database DatabaseConfig
	JWT      JWTConfig
	Midtrans MidtransConfig
	Storage  StorageConfig
	// SupportEmail is surfaced to the web app so visitors can request access to
	// features that are gated behind demo mode (e.g. live payments).
	SupportEmail string
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
	// Enabled is the explicit MIDTRANS_ENABLED toggle. It only turns the live
	// checkout on; whether it's actually usable also depends on the keys below
	// (see Available).
	Enabled bool
}

// Available reports whether the live Midtrans checkout should be attempted. It
// requires the explicit MIDTRANS_ENABLED flag AND both keys, so flipping the
// flag on without credentials can't ship a broken checkout — the app stays in
// demo mode instead.
func (m MidtransConfig) Available() bool {
	return m.Enabled && m.ServerKey != "" && m.ClientKey != ""
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
			Enabled:   getBoolEnv("MIDTRANS_ENABLED", false),
		},
		Storage: StorageConfig{
			AccountID:       os.Getenv("R2_ACCOUNT_ID"),
			AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
			Bucket:          os.Getenv("R2_BUCKET"),
			PublicURL:       os.Getenv("R2_PUBLIC_URL"),
		},
		SupportEmail: getEnv("SUPPORT_EMAIL", "hi@wikukarno.dev"),
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

// getBoolEnv reads a boolean-ish env var. "1", "true", "yes", "on" (any case)
// count as true; anything else falls back to the default when unset, or false
// when set to something unrecognised.
func getBoolEnv(key string, fallback bool) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if v == "" {
		return fallback
	}
	switch v {
	case "1", "true", "yes", "on", "y", "t":
		return true
	default:
		return false
	}
}
