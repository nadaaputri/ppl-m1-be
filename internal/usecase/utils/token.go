package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenMaker struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewTokenMaker(secret, issuer string, ttl time.Duration) *TokenMaker {
	return &TokenMaker{secret: []byte(secret), issuer: issuer, ttl: ttl}
}

func (m *TokenMaker) TTL() time.Duration {
	return m.ttl
}

func (m *TokenMaker) Generate(userID string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    m.issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

func (m *TokenMaker) Parse(raw string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(raw, claims, m.keyFunc,
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(m.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}

func (m *TokenMaker) keyFunc(_ *jwt.Token) (any, error) {
	return m.secret, nil
}