package auth_usecase

import (
	"context"
	"time"

	"ppl-m1-be/internal/usecase/repository/sqlc"
)

type LogoutInput struct {
	UserID    int32
	TokenID   string
	ExpiresAt time.Time
}

func (u *Usecase) Logout(ctx context.Context, in LogoutInput) error {
	return u.store.RevokeToken(ctx, sqlc.RevokeTokenParams{
		Jti: in.TokenID, UserID: in.UserID, ExpiresAt: in.ExpiresAt,
	})
}