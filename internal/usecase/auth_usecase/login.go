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
			Type: "LOGIN_FAILED", Email: email,
			IP: in.IP, Detail: reason,
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
	if user.ID == 0 {
		utils.CheckPassword(u.dummyHash, password)
		return "email tidak terdaftar"
	}
	if utils.CheckPassword(user.PasswordHash, password) != nil {
		return "sandi salah"
	}
	return accountFailureReason(user)
}

func accountFailureReason(user sqlc.GetUserByEmailRow) string {
	switch {
	case !user.IsActive:
		return "akun nonaktif"
	case !user.IsSuperAdmin && !user.RoleID.Valid:
		return "peran tidak tersedia"
	case user.RoleCode.Valid && !user.RoleActive:
		return "peran nonaktif"
	case user.DivisionCode.Valid && !user.DivisionActive:
		return "divisi nonaktif"
	}
	return ""
}

func (u *Usecase) issueToken(ctx context.Context, user sqlc.GetUserByEmailRow, ip string) (LoginResult, error) {
	permissions, err := u.permissionsForUser(ctx, user.ID)
	if err != nil {
		return LoginResult{}, err
	}
	token, err := u.tokens.Generate(toClaims(user, permissions))
	if err != nil {
		return LoginResult{}, err
	}
	u.auditLoginSuccess(ctx, user, ip)
	return LoginResult{AccessToken: token, ExpiresIn: int(u.tokens.TTL().Seconds())}, nil
}

func (u *Usecase) permissionsForUser(ctx context.Context, userID int32) ([]string, error) {
	permissions, err := u.store.GetUserPermissionCodes(ctx, userID)
	if permissions == nil {
		permissions = []string{}
	}
	return permissions, err
}

func (u *Usecase) auditLoginSuccess(ctx context.Context, user sqlc.GetUserByEmailRow, ip string) {
	u.audit.Log(ctx, audit_usecase.Event{
		Type: "LOGIN_SUCCESS", UserID: user.ID, Email: user.Email,
		ActorName: user.Name, IP: ip,
	})
}
