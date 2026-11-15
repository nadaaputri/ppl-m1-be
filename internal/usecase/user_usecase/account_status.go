package user_usecase

import (
	"context"
	"database/sql"
	"errors"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

var (
	ErrUserNotFound         = errors.New("pengguna tidak ditemukan")
	ErrNotAllowed           = errors.New("tidak berhak mengubah status akun ini")
	ErrCannotDeactivateSelf = errors.New("tidak dapat menonaktifkan akun sendiri")
	ErrAlreadyActive        = errors.New("akun sudah aktif")
	ErrAlreadyInactive      = errors.New("akun sudah nonaktif")
	ErrReferenceInactive    = errors.New("peran atau divisi akun nonaktif")
)

type StatusInput struct {
	ActorID      int32
	PartnerScope *int32
	TargetID     int32
	Active       bool
	IP           string
}

type AccountStatus struct {
	ID       int32
	IsActive bool
}

func (u *Usecase) SetAccountStatus(ctx context.Context, in StatusInput) (AccountStatus, error) {
	target, err := u.statusTarget(ctx, in)
	if err != nil {
		return AccountStatus{}, err
	}
	if err := checkStatusChange(in, target); err != nil {
		return AccountStatus{}, err
	}
	row, err := u.store.SetAccountStatus(ctx, sqlc.SetAccountStatusParams{
		ID:        in.TargetID,
		PartnerID: utils.NullInt32(in.PartnerScope),
		IsActive:  in.Active,
		ActorID:   in.ActorID,
		IpAddress: utils.NullText(in.IP),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return AccountStatus{}, statusConflict(in.Active)
	}
	return AccountStatus{ID: row.ID, IsActive: row.IsActive}, err
}

func (u *Usecase) statusTarget(ctx context.Context, in StatusInput) (sqlc.GetAccountForStatusRow, error) {
	row, err := u.store.GetAccountForStatus(ctx, sqlc.GetAccountForStatusParams{
		ID: in.TargetID, PartnerID: utils.NullInt32(in.PartnerScope),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return row, ErrUserNotFound
	}
	return row, err
}

func checkStatusChange(in StatusInput, t sqlc.GetAccountForStatusRow) error {
	switch {
	case !in.Active && in.TargetID == in.ActorID:
		return ErrCannotDeactivateSelf
	case in.PartnerScope != nil && t.RoleCode.String == partnerRoleCode:
		return ErrNotAllowed
	case t.IsActive == in.Active:
		return statusConflict(in.Active)
	case in.Active && (!t.RoleActive || !t.DivisionActive):
		return ErrReferenceInactive
	}
	return nil
}

func statusConflict(active bool) error {
	if active {
		return ErrAlreadyActive
	}
	return ErrAlreadyInactive
}
