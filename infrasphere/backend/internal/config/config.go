package config

import (
	"os"
	"strings"
)

type Config struct {
	Env                string
	Port               string
	DatabaseURL        string
	RedisURL           string
	JWTSecret          string
	EncryptionKey      string
	CORSAllowedOrigins []string
	LogLevel           string
}

func Load() Config {
	return Config{
		Env:                getenv("APP_ENV", "development"),
		Port:               getenv("APP_PORT", "8080"),
		DatabaseURL:        getenv("DATABASE_URL", "postgres://infrasphere:infrasphere@localhost:5432/infrasphere?sslmode=disable"),
		RedisURL:           getenv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:          getenv("JWT_SECRET", "change-me-dev-only"),
		EncryptionKey:      getenv("ENCRYPTION_KEY", "change-me-32-byte-key-dev-only"),
		CORSAllowedOrigins: split(getenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")),
		LogLevel:           getenv("LOG_LEVEL", "debug"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func split(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return out
}

