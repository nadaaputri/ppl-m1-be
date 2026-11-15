package user_usecase

import (
	"testing"
)

func TestToAccountParamsAllowsGlobalAccount(t *testing.T) {
	params := toAccountParams(AccountInput{
		Name:      "Pengguna Global",
		Email:     "global@example.com",
		Role:      "administrator",
		Division:  "pusat",
		CreatedBy: 1,
	}, "hashed-password", 2, 3)

	if params.PartnerID.Valid {
		t.Fatal("partner_id harus NULL untuk akun global")
	}
}

func TestToAccountParamsKeepsPartnerID(t *testing.T) {
	partnerID := int32(9)
	params := toAccountParams(AccountInput{
		PartnerID: &partnerID,
		CreatedBy: 1,
	}, "hashed-password", 2, 3)

	if !params.PartnerID.Valid || params.PartnerID.Int32 != partnerID {
		t.Fatalf("partner_id tidak tersimpan: %+v", params.PartnerID)
	}
}
