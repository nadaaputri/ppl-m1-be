package auth_usecase

import (
	"context"
	"errors"
	"time"

	"ppl-m1-be/internal/usecase/repository/sqlc"
)

var ErrInvalidCredentials = errors.New("email atau kata sandi salah")

type userStore interface {
	GetUserByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error)
	GetAccountCredentials(ctx context.Context, id int32) (sqlc.GetAccountCredentialsRow, error)
	UpdatePassword(ctx context.Context, arg sqlc.UpdatePasswordParams) (time.Time, error)
}

type LoginInput struct {
	Email    string
	Password string
	IP       string
}

type LoginResult struct {
	AccessToken string
	ExpiresIn   int
}
