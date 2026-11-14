package division_usecase

import (
	"context"
	"strings"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func (u *Usecase) Create(ctx context.Context, in CreateInput) (Division, error) {
    row, err := u.store.CreateDivision(ctx, sqlc.CreateDivisionParams{
        Code: in.Code,
        Name: strings.TrimSpace(in.Name),
        ActorID: in.ActorID,
        IpAddress: utils.NullText(in.IP),
    })
    if utils.IsUniqueViolation(err) {
        return Division{}, ErrCodeTaken
    }
    if err != nil {
        return Division{}, err
    }
    return toCreatedDivision(row), nil
}

func toCreatedDivision(d sqlc.CreateDivisionRow) Division {
	return Division{
		ID: d.ID,
		PartnerID: utils.Int32Ptr(d.PartnerID),
		Code: d.Code,
		Name: d.Name,
		IsActive: d.IsActive,
		CreatedAt: d.CreatedAt,
	}
}