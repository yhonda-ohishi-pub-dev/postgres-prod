package config

import (
	"os"
	"strconv"
)

// Default configuration (override via environment variables)
const (
	DefaultRegion = "asia-northeast1"
)

type Config struct {
	// Cloud SQL connection
	ProjectID          string
	Region             string
	InstanceName       string
	DatabaseName       string
	DatabaseUser       string
	DatabasePassword   string // For local development only
	DatabasePort       int    // For local proxy port
	InstanceConnection string // format: project:region:instance

	// Server
	Port string

	// Frontend
	FrontendURL string // OAuth callback redirect URL
}

func Load() *Config {
	cfg := &Config{
		ProjectID:        getEnv("GCP_PROJECT_ID", ""),
		Region:           getEnv("GCP_REGION", DefaultRegion),
		InstanceName:     getEnv("CLOUDSQL_INSTANCE_NAME", ""),
		DatabaseName:     getEnv("DB_NAME", "postgres"),
		DatabaseUser:     getEnv("DB_USER", ""),
		DatabasePassword: getEnv("DB_PASSWORD", ""),
		DatabasePort:     getEnvInt("DB_PORT", 5432),
		Port:             getEnv("PORT", "8080"),
		FrontendURL:      getEnv("FRONTEND_URL", ""),
	}

	// Build instance connection string
	if cfg.ProjectID != "" && cfg.Region != "" && cfg.InstanceName != "" {
		cfg.InstanceConnection = cfg.ProjectID + ":" + cfg.Region + ":" + cfg.InstanceName
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
