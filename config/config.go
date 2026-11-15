package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	minJWTSecretBytes = 32
	sharedJWTLifetime = 8 * time.Hour
)

type Config struct {
	AppName            string
	AppPort            string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	JWTSecret          string
	JWTExpiresIn       time.Duration
	CORSAllowedOrigins []string
	TrustedProxies     []string
}

func Load() (Config, error) {
	if err := loadDotEnv(); err != nil {
		return Config{}, err
	}
	cfg := readConfig()
	if err := validateDatabase(cfg); err != nil {
		return Config{}, err
	}
	if err := loadJWT(&cfg); err != nil {
		return Config{}, err
	}
	loadCORS(&cfg)
	if err := loadProxies(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func loadDotEnv() error {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("gagal membaca file .env: %w", err)
	}
	return nil
}

func readConfig() Config {
	cfg := Config{
		AppName: os.Getenv("APP_NAME"), AppPort: os.Getenv("APP_PORT"),
		DBHost: os.Getenv("DB_HOST"), DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"), DBPassword: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"), DBSSLMode: os.Getenv("DB_SSLMODE"),
	}
	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}
	if cfg.DBSSLMode == "" {
		cfg.DBSSLMode = "disable"
	}
	return cfg
}

func validateDatabase(cfg Config) error {
	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBName == "" {
		return fmt.Errorf("konfigurasi database belum lengkap")
	}
	return nil
}

func loadJWT(cfg *Config) error {
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if len(cfg.JWTSecret) < minJWTSecretBytes {
		return fmt.Errorf("JWT_SECRET minimal %d byte", minJWTSecretBytes)
	}

	ttl, err := time.ParseDuration(os.Getenv("JWT_EXPIRES_IN"))
	if err != nil || ttl != sharedJWTLifetime {
		return fmt.Errorf("JWT_EXPIRES_IN harus 8h sesuai Panduan API Bersama")
	}

	cfg.JWTExpiresIn = sharedJWTLifetime
	return nil
}

func loadCORS(cfg *Config) {
	raw := os.Getenv("CORS_ALLOWED_ORIGINS")
	if raw == "" {
		raw = "*"
	}
	for _, origin := range strings.Split(raw, ",") {
		cfg.CORSAllowedOrigins = append(cfg.CORSAllowedOrigins, strings.TrimSpace(origin))
	}
}

func loadProxies(cfg *Config) error {
	for _, raw := range strings.Split(os.Getenv("TRUSTED_PROXIES"), ",") {
		entry := strings.TrimSpace(raw)
		if entry == "" {
			continue
		}
		if _, _, err := net.ParseCIDR(entry); err != nil && net.ParseIP(entry) == nil {
			return fmt.Errorf("TRUSTED_PROXIES tidak valid: %q", entry)
		}
		cfg.TrustedProxies = append(cfg.TrustedProxies, entry)
	}
	return nil
}
