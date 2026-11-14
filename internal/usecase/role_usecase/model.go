package role_usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ppl-m1-be/internal/usecase/repository/sqlc"
)

var (
	ErrCodeTaken             = errors.New("kode peran global sudah dipakai")
	ErrCodeReserved          = errors.New("kode peran dicadangkan sistem")
	ErrRoleNotFound          = errors.New("peran tidak ditemukan")
	ErrPermissionNotFound    = errors.New("permission tidak terdaftar")
	ErrDuplicatePermission   = errors.New("permission duplikat")
	ErrInvalidPermissionCode = errors.New("kode permission tidak valid")
)

var reservedCodes = []string{"administrator", "partner", "manajer", "customer"}

type roleStore interface {
	CountRoles(context.Context, sql.NullBool) (int64, error)
	ListRoles(context.Context, sqlc.ListRolesParams) ([]sqlc.Role, error)
	CreateRole(context.Context, sqlc.CreateRoleParams) (sqlc.CreateRoleRow, error)
	CountPermissions(context.Context) (int64, error)
	ListPermissions(context.Context, sqlc.ListPermissionsParams) ([]sqlc.Permission, error)
	CountPermissionsByCodes(context.Context, string) (int64, error)
	GetGlobalRoleByID(context.Context, int32) (sqlc.Role, error)
	GetRolePermissionCodes(context.Context, int32) ([]string, error)
	SetRolePermissions(context.Context, sqlc.SetRolePermissionsParams) error
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
	IsActive *bool
	Page     int
	PerPage  int
}

type ListResult struct {
	Items []Role
	Total int
}

type CreateInput struct {
	Code            string
	Name            string
	Description     string
	PermissionCodes []string
}

type Permission struct {
	Code        string
	Module      string
	Description string
}

type PermissionListInput struct {
	Page    int
	PerPage int
}

type PermissionListResult struct {
	Items []Permission
	Total int
}
