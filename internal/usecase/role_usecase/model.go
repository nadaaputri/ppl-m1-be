package role_usecase

import (
	"context"
	"errors"
	"time"

	"ppl-m1-be/internal/usecase/repository/sqlc"
)

var (
	ErrCodeTaken    = errors.New("kode peran sudah dipakai di partner ini")
	ErrCodeReserved = errors.New("kode peran dicadangkan sistem")
)

var reservedCodes = []string{"administrator", "partner", "manajer", "customer"}

type roleStore interface {
	CountRoles(ctx context.Context, arg sqlc.CountRolesParams) (int64, error)
	ListRoles(ctx context.Context, arg sqlc.ListRolesParams) ([]sqlc.Role, error)
	CreateRole(ctx context.Context, arg sqlc.CreateRoleParams) (sqlc.Role, error)
}

type Role struct {
	ID          int32
	PartnerID   *int32
	Code        string
	Name        string
	Description *string
	IsActive    bool
	CreatedAt   time.Time
}

type ListInput struct {
	PartnerID *int32
	IsActive  *bool
	Page      int
	PerPage   int
}

type ListResult struct {
	Items []Role
	Total int
}

type CreateInput struct {
	PartnerID   int32
	Code        string
	Name        string
	Description string
}
