package role_usecase

import (
	"context"
	"database/sql"
	"slices"
	"strings"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func (u *Usecase) Create(ctx context.Context, in CreateInput) (Role, error) {
	if slices.Contains(reservedCodes, in.Code) {
		return Role{}, ErrCodeReserved
	}
	row, err := u.store.CreateRole(ctx, toCreateParams(in))
	if utils.IsUniqueViolation(err) {
		return Role{}, ErrCodeTaken
	}
	if err != nil {
		return Role{}, err
	}
	return toRole(row), nil
}

func toCreateParams(in CreateInput) sqlc.CreateRoleParams {
	description := strings.TrimSpace(in.Description)
	return sqlc.CreateRoleParams{
		PartnerID:   sql.NullInt32{Int32: in.PartnerID, Valid: true},
		Code:        in.Code,
		Name:        strings.TrimSpace(in.Name),
		Description: sql.NullString{String: description, Valid: description != ""},
	}
}
