// Package config centralizes all environment-driven configuration.
//
// NOTE: The original Spring Boot app read these values from
// src/main/resources/application.properties using ${ENV_VAR} placeholders:
//
//	spring.datasource.url=jdbc:postgresql://aws-1-ap-northeast-1.pooler.supabase.com:5432/postgres
//	spring.datasource.username=postgres.ouyymrwygelugaptdgmh
//	spring.datasource.password=${DB_PASS}
//	app.jwt.secret=${JWT_SECRET_64}
//	app.jwt.expiration=604800000   <-- NOTE: this property is defined but never
//	                                    actually read anywhere in the Kotlin code;
//	                                    JWTService hardcodes 15 minutes (access)
//	                                    and 30 days (refresh) instead. We preserve
//	                                    that exact (dead-config) behavior below.
//	gemini.api-key=${GEMINI_API_KEY}
//	gemini.url=https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent
//	server.port=${PORT:8080}
//	server.address=0.0.0.0
//
// We keep the same defaults and the same set of required environment
// variables so a deployment can be migrated by copying over the same
// env vars (DB_PASS, JWT_SECRET_64, GEMINI_API_KEY, PORT).
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	// Set DATABASE_URL to fully override the DSN built from the fields above.
	DatabaseURL string

	// JWT
	// Base64-encoded HMAC-SHA256 secret, identical format to the Kotlin
	// JWTService which does: Keys.hmacShaKeyFor(Base64.getDecoder().decode(secret))
	JWTSecretBase64 string

	// Gemini
	GeminiAPIKey string
	GeminiURL    string

	// Server
	Port string
}

// Load reads configuration from environment variables, applying the same
// defaults as application.properties.
func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		DBHost:      getEnv("DB_HOST", "aws-1-ap-northeast-1.pooler.supabase.com"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBName:      getEnv("DB_NAME", "postgres"),
		DBUser:      getEnv("DB_USER", "postgres.ouyymrwygelugaptdgmh"),
		DBPassword:  os.Getenv("DB_PASS"),
		DatabaseURL: os.Getenv("DATABASE_URL"),

		JWTSecretBase64: os.Getenv("JWT_SECRET_64"),

		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		GeminiURL:    getEnv("GEMINI_URL", "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent"),

		Port: getEnv("PORT", "8080"),
	}

	if cfg.JWTSecretBase64 == "" {
		return nil, fmt.Errorf("JWT_SECRET_64 environment variable is required")
	}
	if cfg.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}
	if cfg.DatabaseURL == "" && cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASS (or DATABASE_URL) environment variable is required")
	}

	return cfg, nil
}

// DSN builds a libpq-style connection string equivalent to the JDBC URL
// spring.datasource.url used in application.properties.
func (c *Config) DSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
