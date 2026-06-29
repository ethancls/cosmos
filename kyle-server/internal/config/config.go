package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	GuacdHost   string
	GuacdPort   string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		Port:        envOrDefault("PORT", "8080"),
		DatabaseURL: envOrDefault("DATABASE_URL", "postgres://kyle:kyle@localhost:5432/kyle?sslmode=disable"),
		GuacdHost:   envOrDefault("GUACD_HOST", "localhost"),
		GuacdPort:   envOrDefault("GUACD_PORT", "4822"),
		JWTSecret:   envOrDefault("JWT_SECRET", "change-me-in-production"),
	}
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
