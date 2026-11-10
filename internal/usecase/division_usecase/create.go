package division_usecase

import (
	"context"
	"database/sql"
	"strings"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func (u *Usecase) Create(ctx context.Context, in CreateInput) (Division, error) {
	row, err := u.store.CreateDivision(ctx, sqlc.CreateDivisionParams{
		PartnerID: sql.NullInt32{Int32: in.PartnerID, Valid: true},
		Code:      in.Code,
		Name:      strings.TrimSpace(in.Name),
	})
	if utils.IsUniqueViolation(err) {
		return Division{}, ErrCodeTaken
	}
	if err != nil {
		return Division{}, err
	}
	return toDivision(row), nil
}
