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
	if err := u.validatePermissionCodes(ctx, in.PermissionCodes); err != nil {
		return Role{}, err
	}
	row, err := u.store.CreateRole(ctx, toCreateParams(in))
	if utils.IsUniqueViolation(err) {
		return Role{}, ErrCodeTaken
	}
	if err != nil {
		return Role{}, err
	}
	return toCreatedRole(row), nil
}

func toCreatedRole(row sqlc.CreateRoleRow) Role {
	return toRole(sqlc.Role(row))
}

func toCreateParams(in CreateInput) sqlc.CreateRoleParams {
	description := strings.TrimSpace(in.Description)
	return sqlc.CreateRoleParams{
		Code: in.Code,
		Name: strings.TrimSpace(in.Name),
		Description: sql.NullString{
			String: description, Valid: description != "",
		},
		PermissionCodes: strings.Join(in.PermissionCodes, ","),
	}
}
