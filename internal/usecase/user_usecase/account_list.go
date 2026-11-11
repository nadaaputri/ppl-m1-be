package user_usecase

import (
	"context"
	"strings"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func (u *Usecase) ListAccounts(ctx context.Context, in AccountListInput) (AccountListResult, error) {
	total, err := u.store.CountAccounts(ctx, toCountParams(in))
	if err != nil {
		return AccountListResult{}, err
	}
	rows, err := u.store.ListAccounts(ctx, toListParams(in))
	if err != nil {
		return AccountListResult{}, err
	}
	items := make([]Account, 0, len(rows))
	for _, row := range rows {
		items = append(items, toListedAccount(row))
	}
	return AccountListResult{Items: items, Total: int(total)}, nil
}

func toCountParams(in AccountListInput) sqlc.CountAccountsParams {
	return sqlc.CountAccountsParams{
		PartnerID:    utils.NullInt32(in.PartnerID),
		Search:       utils.NullText(utils.EscapeLike(strings.TrimSpace(in.Search))),
		RoleCode:     utils.NullText(in.Role),
		DivisionCode: utils.NullText(in.Division),
		UnitUsahaID:  utils.NullInt32(in.UnitUsahaID),
		IsActive:     utils.NullBool(in.IsActive),
	}
}

func toListParams(in AccountListInput) sqlc.ListAccountsParams {
	return sqlc.ListAccountsParams{
		PartnerID:    utils.NullInt32(in.PartnerID),
		Search:       utils.NullText(utils.EscapeLike(strings.TrimSpace(in.Search))),
		RoleCode:     utils.NullText(in.Role),
		DivisionCode: utils.NullText(in.Division),
		UnitUsahaID:  utils.NullInt32(in.UnitUsahaID),
		IsActive:     utils.NullBool(in.IsActive),
		Sort:         in.Sort,
		PageLimit:    int32(in.PerPage),
		PageOffset:   int32((in.Page - 1) * in.PerPage),
	}
}

func toListedAccount(r sqlc.ListAccountsRow) Account {
	return Account{
		ID: r.ID, Name: r.Name, Email: r.Email,
		Role: utils.StringPtr(r.RoleCode), Division: utils.StringPtr(r.DivisionCode),
		PartnerID: utils.Int32Ptr(r.PartnerID), UnitUsahaID: utils.Int32Ptr(r.UnitUsahaID),
		IsActive: r.IsActive, CreatedAt: r.CreatedAt,
	}
}
