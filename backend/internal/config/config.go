package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	AllowOrigin    string
	UploadDir      string
	MaxUploadMB    int64
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8081"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://teststore:teststore@localhost:5432/teststore?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-me"),
		AllowOrigin: getEnv("ALLOW_ORIGIN", "http://localhost:5173"),
		UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadMB: getEnvInt("MAX_UPLOAD_MB", 5),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fallback
	}
	return n
}
