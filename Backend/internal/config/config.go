package config

import "os"

type Config struct {
	DBPath        string
	HTTPPort      string
	JWTSecret     string
	AdminEmail    string
	AdminPassword string
}

func Load() Config {
	return Config{
		DBPath:        getEnv("PETMATCH_DB_PATH", "petmatch.db"),
		HTTPPort:      getEnv("PETMATCH_HTTP_PORT", "8084"),
		JWTSecret:     getEnv("PETMATCH_JWT_SECRET", "change-me"),
		AdminEmail:    getEnv("PETMATCH_ADMIN_EMAIL", "admin@petmatch.local"),
		AdminPassword: getEnv("PETMATCH_ADMIN_PASSWORD", "admin123"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
