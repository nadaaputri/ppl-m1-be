package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort   string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPassword string
	DBName    string
	DBSSLMode string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("gagal membaca file .env: %w", err)
	}

	cfg := Config{
		AppPort:    os.Getenv("APP_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
	}

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	if cfg.DBSSLMode == "" {
		cfg.DBSSLMode = "disable"
	}

	if cfg.DBHost == "" || cfg.DBPort == "" ||
		cfg.DBUser == "" || cfg.DBName == "" {
		return Config{}, fmt.Errorf("konfigurasi database belum lengkap")
	}

	return cfg, nil
}