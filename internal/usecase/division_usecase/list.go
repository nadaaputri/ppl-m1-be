package division_usecase

import (
	"context"
	"database/sql"

	"ppl-m1-be/internal/usecase/repository/sqlc"
)

func (u *Usecase) List(ctx context.Context, in ListInput) (ListResult, error) {
	partner, active := toNullInt32(in.PartnerID), toNullBool(in.IsActive)
	total, err := u.store.CountDivisions(ctx, sqlc.CountDivisionsParams{PartnerID: partner, IsActive: active})
	if err != nil {
		return ListResult{}, err
	}
	rows, err := u.store.ListDivisions(ctx, sqlc.ListDivisionsParams{
		PartnerID: partner, IsActive: active,
		PageLimit: int32(in.PerPage), PageOffset: int32((in.Page - 1) * in.PerPage),
	})
	if err != nil {
		return ListResult{}, err
	}
	items := make([]Division, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDivision(row))
	}
	return ListResult{Items: items, Total: int(total)}, nil
}

func toNullInt32(v *int32) sql.NullInt32 {
	if v == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *v, Valid: true}
}

func toNullBool(v *bool) sql.NullBool {
	if v == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *v, Valid: true}
}

func toDivision(d sqlc.Division) Division {
	var partnerID *int32
	if d.PartnerID.Valid {
		partnerID = &d.PartnerID.Int32
	}
	return Division{
		ID: d.ID, PartnerID: partnerID, Code: d.Code, Name: d.Name,
		IsActive: d.IsActive, CreatedAt: d.CreatedAt,
	}
}