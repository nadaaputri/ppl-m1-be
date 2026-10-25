package auth_usecase

import (
	"context"
	"errors"

	"ppl-m1-be/internal/usecase/repository/sqlc"
)

var ErrInvalidCredentials = errors.New("email atau kata sandi salah")

type userStore interface {
	GetUserByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error)
}

type LoginInput struct {
	Email    string
	Password string
	IP       string
}

type UserInfo struct {
	ID           string
	Name         string
	Email        string
	Role         string
	DivisionID   string
	DivisionName string
}

type LoginResult struct {
	AccessToken string
	ExpiresIn   int
	User        UserInfo
}