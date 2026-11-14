package division_usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ppl-m1-be/internal/usecase/repository/sqlc"
)

var ErrCodeTaken = errors.New("kode divisi global sudah dipakai")

type divisionStore interface {
	CountDivisions(context.Context, sql.NullBool) (int64, error)
	ListDivisions(ctx context.Context, arg sqlc.ListDivisionsParams) ([]sqlc.Division, error)
	CreateDivision(ctx context.Context, arg sqlc.CreateDivisionParams) (sqlc.Division, error)
}

type Division struct {
	ID        int32
	PartnerID *int32
	Code      string
	Name      string
	IsActive  bool
	CreatedAt time.Time
}

type ListInput struct {
	IsActive *bool
	Page     int
	PerPage  int
}

type ListResult struct {
	Items []Division
	Total int
}

type CreateInput struct {
	Code string
	Name string
}
