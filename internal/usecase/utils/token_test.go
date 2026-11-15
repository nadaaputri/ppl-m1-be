package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func tokenTestMaker(issuer string, ttl time.Duration) *TokenMaker {
	return NewTokenMaker(strings.Repeat("s", 32), issuer, ttl)
}

func TestTokenGenerateAndParse(t *testing.T) {
	maker := tokenTestMaker("araya", 8*time.Hour)
	raw, err := maker.Generate(Claims{UserID: 7, Name: "Rani", Role: "staff"})
	if err != nil {
		t.Fatal(err)
	}
	got, err := maker.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.UserID != 7 || got.Role != "staff" {
		t.Fatalf("claims tidak sesuai: %+v", got)
	}
	if got.ExpiresAt.Sub(got.IssuedAt.Time) != 8*time.Hour {
		t.Fatal("masa berlaku token bukan 8 jam")
	}
}

func TestTokenRejectsExpired(t *testing.T) {
	maker := tokenTestMaker("araya", -time.Minute)
	raw, err := maker.Generate(Claims{UserID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := maker.Parse(raw); err == nil {
		t.Fatal("token kedaluwarsa seharusnya ditolak")
	}
}

func TestTokenRejectsWrongIssuer(t *testing.T) {
	raw, err := tokenTestMaker("araya", 8*time.Hour).Generate(Claims{UserID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tokenTestMaker("modul-lain", 8*time.Hour).Parse(raw); err == nil {
		t.Fatal("issuer berbeda seharusnya ditolak")
	}
}

func TestTokenRejectsWrongSecret(t *testing.T) {
	maker := NewTokenMaker(strings.Repeat("a", 32), "araya", 8*time.Hour)
	raw, err := maker.Generate(Claims{UserID: 1})
	if err != nil {
		t.Fatal(err)
	}
	other := NewTokenMaker(strings.Repeat("b", 32), "araya", 8*time.Hour)
	if _, err := other.Parse(raw); err == nil {
		t.Fatal("signature tidak valid seharusnya ditolak")
	}
}

func TestTokenRejectsWrongAlgorithm(t *testing.T) {
	claims := Claims{UserID: 1, RegisteredClaims: jwt.RegisteredClaims{
		Issuer: "araya", ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}}
	raw, err := jwt.NewWithClaims(jwt.SigningMethodHS384, claims).SignedString([]byte(strings.Repeat("s", 32)))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tokenTestMaker("araya", 8*time.Hour).Parse(raw); err == nil {
		t.Fatal("algoritma selain HS256 seharusnya ditolak")
	}
}
