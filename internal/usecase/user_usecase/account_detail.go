package user_usecase

import (
	"context"
	"database/sql"
	"errors"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func (u *Usecase) GetAccount(ctx context.Context, id int32) (Account, error) {
	row, err := u.store.GetAccountByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Account{}, ErrUserNotFound
	}
	if err != nil {
		return Account{}, err
	}
	return toAccountDetail(row), nil
}

func toAccountDetail(r sqlc.GetAccountByIDRow) Account {
	return Account{
		ID: r.ID, Name: r.Name, Email: r.Email,
		PartnerID:   utils.Int32Ptr(r.PartnerID),
		UnitUsahaID: utils.Int32Ptr(r.UnitUsahaID),
		Role:        utils.StringPtr(r.RoleCode),
		Division:    utils.StringPtr(r.DivisionCode),
		IsActive:    r.IsActive, CreatedAt: r.CreatedAt,
	}
}
