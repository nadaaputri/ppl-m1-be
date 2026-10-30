package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID      int32   `json:"user_id"`
	Name        string  `json:"name"`
	Role        string  `json:"role"`
	Division    *string `json:"division"`
	UnitUsahaID *int32  `json:"unit_usaha_id"`
	jwt.RegisteredClaims
}

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

func (m *TokenMaker) Generate(c Claims) (string, error) {
	now := time.Now()
	c.Issuer = m.issuer
	c.IssuedAt = jwt.NewNumericDate(now)
	c.ExpiresAt = jwt.NewNumericDate(now.Add(m.ttl))
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(m.secret)
}

func (m *TokenMaker) Parse(raw string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(raw, claims, m.keyFunc,
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(m.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (m *TokenMaker) keyFunc(_ *jwt.Token) (any, error) {
	return m.secret, nil
}