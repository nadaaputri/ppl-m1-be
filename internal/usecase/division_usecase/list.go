package division_usecase

import (
	"context"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func (u *Usecase) List(ctx context.Context, in ListInput) (ListResult, error) {
	active := utils.NullBool(in.IsActive)
	total, err := u.store.CountDivisions(ctx, active)
	if err != nil {
		return ListResult{}, err
	}
	rows, err := u.store.ListDivisions(ctx, sqlc.ListDivisionsParams{
		IsActive: active, PageLimit: int32(in.PerPage),
		PageOffset: int32((in.Page - 1) * in.PerPage),
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
