package user_usecase

import (
	"context"
	"errors"
	"time"

	"ppl-m1-be/internal/usecase/repository/sqlc"
)

const customerRoleName = "customer"
const partnerRoleCode = "partner"

var (
	ErrEmailTaken   = errors.New("email sudah terdaftar")
	ErrRoleNotFound = errors.New("role customer belum tersedia di tabel roles")
)

type userStore interface {
	GetRoleIDByName(ctx context.Context, roleName string) (int32, error)
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.CreateUserRow, error)
	GetGlobalRoleIDByCode(ctx context.Context, code string) (int32, error)
	GetGlobalDivisionIDByCode(ctx context.Context, code string) (int32, error)
	CreateAccount(ctx context.Context, arg sqlc.CreateAccountParams) (sqlc.User, error)
	CountAccounts(ctx context.Context, arg sqlc.CountAccountsParams) (int64, error)
	ListAccounts(ctx context.Context, arg sqlc.ListAccountsParams) ([]sqlc.ListAccountsRow, error)
	GetAccountForStatus(ctx context.Context, arg sqlc.GetAccountForStatusParams) (sqlc.GetAccountForStatusRow, error)
	SetAccountStatus(ctx context.Context, arg sqlc.SetAccountStatusParams) (sqlc.SetAccountStatusRow, error)
}

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "data tidak valid"
}

type AccountInput struct {
	PartnerID   int32
	CreatedBy   int32
	Name        string
	Email       string
	Password    string
	Role        string
	Division    string
	UnitUsahaID *int32
}

type Account struct {
	ID          int32
	Name        string
	Email       string
	Role        *string
	Division    *string
	PartnerID   *int32
	UnitUsahaID *int32
	IsActive    bool
	CreatedAt   time.Time
}

type AccountListInput struct {
	PartnerID   *int32
	Search      string
	Role        string
	Division    string
	UnitUsahaID *int32
	IsActive    *bool
	Sort        string
	Page        int
	PerPage     int
}

type AccountListResult struct {
	Items []Account
	Total int
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
	IP       string
}

type RegisterResult struct {
	UserID int32
	Name   string
	Email  string
	Role   string
}
