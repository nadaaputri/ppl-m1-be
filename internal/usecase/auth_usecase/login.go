package auth_usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"ppl-m1-be/internal/usecase/audit_usecase"
	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func (u *Usecase) Login(ctx context.Context, in LoginInput) (LoginResult, error) {
	email := normalizeEmail(in.Email)
	user, err := u.findUser(ctx, email)
	if err != nil {
		return LoginResult{}, err
	}

	if reason := u.failReason(user, in.Password); reason != "" {
		u.audit.Log(ctx, audit_usecase.Event{
			Type: "LOGIN_FAILED", Email: email, IP: in.IP, Detail: reason,
		})
		return LoginResult{}, ErrInvalidCredentials
	}
	return u.issueToken(ctx, user, in.IP)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (u *Usecase) findUser(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error) {
	user, err := u.store.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return sqlc.GetUserByEmailRow{}, nil
	}
	return user, err
}

func (u *Usecase) failReason(user sqlc.GetUserByEmailRow, password string) string {
	hash := user.PasswordHash
	if user.ID == 0 {
		hash = u.dummyHash
	}
	wrongPassword := utils.CheckPassword(hash, password) != nil

	switch {
	case user.ID == 0:
		return "email tidak terdaftar"
	case wrongPassword:
		return "sandi salah"
	case !user.IsActive:
		return "akun nonaktif"
	}
	return ""
}

func (u *Usecase) issueToken(ctx context.Context, user sqlc.GetUserByEmailRow, ip string) (LoginResult, error) {
	token, err := u.tokens.Generate(toClaims(user))
	if err != nil {
		return LoginResult{}, err
	}

	u.audit.Log(ctx, audit_usecase.Event{Type: "LOGIN_SUCCESS", UserID: user.ID, IP: ip})
	return LoginResult{
		AccessToken: token,
		ExpiresIn:   int(u.tokens.TTL().Seconds()),
	}, nil
}