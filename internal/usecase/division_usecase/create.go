package division_usecase

import (
	"context"
	"strings"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func (u *Usecase) Create(ctx context.Context, in CreateInput) (Division, error) {
    row, err := u.store.CreateDivision(ctx, sqlc.CreateDivisionParams{
        Code: in.Code, Name: strings.TrimSpace(in.Name),
    })
    if utils.IsUniqueViolation(err) {
        return Division{}, ErrCodeTaken
    }
    if err != nil {
        return Division{}, err
    }
    return toDivision(row), nil
}