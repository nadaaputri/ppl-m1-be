package auth_usecase

import (
	"database/sql"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func toClaims(u sqlc.GetUserByEmailRow) utils.Claims {
	return utils.Claims{
		UserID:       u.ID,
		Name:         u.Name,
		Role:         u.RoleCode.String,
		Division:     stringOrNil(u.DivisionCode),
		PartnerID:    intOrNil(u.PartnerID),
		UnitUsahaID:  intOrNil(u.UnitUsahaID),
		IsSuperAdmin: u.IsSuperAdmin,
	}
}

func stringOrNil(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

func intOrNil(n sql.NullInt32) *int32 {
	if !n.Valid {
		return nil
	}
	return &n.Int32
}
