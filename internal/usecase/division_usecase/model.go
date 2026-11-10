package division_usecase

import (
	"context"
	"errors"
	"time"

	"ppl-m1-be/internal/usecase/repository/sqlc"
)

var ErrCodeTaken = errors.New("kode divisi sudah dipakai di partner ini")

type divisionStore interface {
	CountDivisions(ctx context.Context, arg sqlc.CountDivisionsParams) (int64, error)
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
	PartnerID *int32
	IsActive  *bool
	Page      int
	PerPage   int
}

type ListResult struct {
	Items []Division
	Total int
}

type CreateInput struct {
	PartnerID int32
	Code      string
	Name      string
}
