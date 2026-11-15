package role_usecase

import (
	"context"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func (u *Usecase) List(ctx context.Context, in ListInput) (ListResult, error) {
	active := utils.NullBool(in.IsActive)
	total, err := u.store.CountRoles(ctx, active)
	if err != nil {
		return ListResult{}, err
	}
	rows, err := u.store.ListRoles(ctx, sqlc.ListRolesParams{
		IsActive: active, PageLimit: int32(in.PerPage),
		PageOffset: int32((in.Page - 1) * in.PerPage),
	})
	if err != nil {
		return ListResult{}, err
	}
	items := make([]Role, 0, len(rows))
	for _, row := range rows {
		items = append(items, toRole(row))
	}
	return ListResult{Items: items, Total: int(total)}, nil
}

func toRole(r sqlc.Role) Role {
	out := Role{ID: r.ID, Code: r.Code, Name: r.Name, IsActive: r.IsActive, CreatedAt: r.CreatedAt}
	if r.PartnerID.Valid {
		out.PartnerID = &r.PartnerID.Int32
	}
	if r.Description.Valid {
		out.Description = &r.Description.String
	}
	return out
}
