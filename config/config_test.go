package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadJWTRequiresEightHours(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("s", 32))
	t.Setenv("JWT_EXPIRES_IN", "8h")
	cfg := Config{}
	if err := loadJWT(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.JWTExpiresIn != 8*time.Hour {
		t.Fatalf("expiration tidak sesuai: %v", cfg.JWTExpiresIn)
	}
}

func TestLoadJWTRejectsOtherLifetime(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("s", 32))
	t.Setenv("JWT_EXPIRES_IN", "1h")
	if err := loadJWT(&Config{}); err == nil {
		t.Fatal("durasi selain 8h seharusnya ditolak")
	}
}

func TestLoadJWTRejectsShortSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("s", 31))
	t.Setenv("JWT_EXPIRES_IN", "8h")
	if err := loadJWT(&Config{}); err == nil {
		t.Fatal("secret kurang dari 32 byte seharusnya ditolak")
	}
}
