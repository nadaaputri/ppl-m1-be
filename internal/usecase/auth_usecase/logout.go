package auth_usecase

import (
	"context"
	"time"

	"ppl-m1-be/internal/usecase/audit_usecase"
	"ppl-m1-be/internal/usecase/repository/sqlc"
)

type LogoutInput struct {
	UserID    int32
	Name      string
	PartnerID *int32
	TokenID   string
	ExpiresAt time.Time
	IP        string
}

func (u *Usecase) Logout(ctx context.Context, in LogoutInput) error {
	err := u.store.RevokeToken(ctx, sqlc.RevokeTokenParams{
		Jti: in.TokenID, UserID: in.UserID, ExpiresAt: in.ExpiresAt,
	})
	if err != nil {
		return err
	}
	u.audit.Log(ctx, audit_usecase.Event{
		Type: "LOGOUT", UserID: in.UserID, PartnerID: in.PartnerID,
		ActorName: in.Name, Email: u.actorEmail(ctx, in.UserID), IP: in.IP,
	})
	return nil
}

func (u *Usecase) actorEmail(ctx context.Context, id int32) string {
	acc, err := u.store.GetAccountCredentials(ctx, id)
	if err != nil {
		return ""
	}
	return acc.Email
}
