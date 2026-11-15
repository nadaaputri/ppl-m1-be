package app

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"time"

	"ppl-m1-be/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL(cfg))
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}
	configurePool(db)
	if err := pingDatabase(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func databaseURL(cfg config.Config) string {
	u := url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:   net.JoinHostPort(cfg.DBHost, cfg.DBPort),
		Path:   cfg.DBName,
	}
	query := url.Values{}
	query.Set("sslmode", cfg.DBSSLMode)
	u.RawQuery = query.Encode()
	return u.String()
}

func configurePool(db *sql.DB) {
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
}

func pingDatabase(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("gagal terhubung ke PostgreSQL: %w", err)
	}
	return nil
}
