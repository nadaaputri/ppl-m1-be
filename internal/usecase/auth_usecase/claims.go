package auth_usecase

import (
	"database/sql"
	"strings"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func toClaims(u sqlc.GetUserByEmailRow) utils.Claims {
	return utils.Claims{
		UserID:      u.ID,
		Name:        u.FullName,
		Role:        strings.ToLower(u.RoleName),
		Division:    lowerOrNil(u.DivisionName),
		UnitUsahaID: intOrNil(u.UnitUsahaID),
	}
}

func lowerOrNil(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	v := strings.ToLower(s.String)
	return &v
}

func intOrNil(n sql.NullInt32) *int32 {
	if !n.Valid {
		return nil
	}
	return &n.Int32
}