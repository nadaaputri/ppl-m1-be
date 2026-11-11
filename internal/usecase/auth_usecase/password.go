package auth_usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ppl-m1-be/internal/usecase/audit_usecase"
	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

var (
	ErrOldPasswordInvalid = errors.New("kata sandi lama salah")
	ErrAccountUnavailable = errors.New("akun tidak tersedia")
)

type NewPasswordError struct {
	Message string
}

func (e *NewPasswordError) Error() string {
	return e.Message
}

type ChangePasswordInput struct {
	UserID      int32
	OldPassword string
	NewPassword string
	IP          string
}

func (u *Usecase) ChangePassword(ctx context.Context, in ChangePasswordInput) (time.Time, error) {
	acc, err := u.store.GetAccountCredentials(ctx, in.UserID)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && !acc.IsActive) {
		return time.Time{}, ErrAccountUnavailable
	}
	if err != nil {
		return time.Time{}, err
	}
	if err := checkPasswords(acc, in); err != nil {
		return time.Time{}, err
	}
	changedAt, err := u.saveNewPassword(ctx, acc.PasswordHash, in)
	if err == nil {
		u.audit.Log(ctx, audit_usecase.Event{Type: "PASSWORD_CHANGED", UserID: in.UserID, IP: in.IP})
	}
	return changedAt, err
}

func checkPasswords(acc sqlc.GetAccountCredentialsRow, in ChangePasswordInput) error {
	if len(in.OldPassword) > utils.MaxPasswordBytes || utils.CheckPassword(acc.PasswordHash, in.OldPassword) != nil {
		return ErrOldPasswordInvalid
	}
	if msg := utils.ValidatePassword(in.NewPassword, acc.Email); msg != "" {
		return &NewPasswordError{Message: msg}
	}
	if in.NewPassword == in.OldPassword {
		return &NewPasswordError{Message: "Tidak boleh sama dengan kata sandi lama"}
	}
	return nil
}

func (u *Usecase) saveNewPassword(ctx context.Context, oldHash string, in ChangePasswordInput) (time.Time, error) {
	newHash, err := utils.HashPassword(in.NewPassword)
	if err != nil {
		return time.Time{}, err
	}
	changedAt, err := u.store.UpdatePassword(ctx, sqlc.UpdatePasswordParams{
		ID: in.UserID, OldHash: oldHash, NewHash: newHash,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return time.Time{}, ErrOldPasswordInvalid
	}
	return changedAt, err
}
