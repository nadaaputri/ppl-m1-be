package division_route

import (
	"regexp"
	"strings"
	"time"

	"ppl-m1-be/internal/usecase/division_usecase"
)

var (
	codePattern = regexp.MustCompile(`^[a-z0-9]+([_-][a-z0-9]+)*$`)
	wib         = time.FixedZone("WIB", 7*60*60)
)

type createRequest struct {
    Code string `json:"code" binding:"required,max=50"`
    Name string `json:"name" binding:"required,max=100"`
}

type divisionResponse struct {
	ID        int32     `json:"id"`
	PartnerID *int32    `json:"partner_id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func extraFieldErrors(req createRequest) map[string]string {
	fields := map[string]string{}
	if !codePattern.MatchString(req.Code) {
		fields["code"] = "Hanya huruf kecil, angka, tanda hubung, dan garis bawah"
	}
	if strings.TrimSpace(req.Name) == "" {
		fields["name"] = "Isian ini wajib diisi"
	}
	return fields
}

func toDivisionResponse(d division_usecase.Division) divisionResponse {
	return divisionResponse{
		ID: d.ID, PartnerID: d.PartnerID, Code: d.Code, Name: d.Name,
		IsActive: d.IsActive, CreatedAt: d.CreatedAt.In(wib),
	}
}

func toDivisionList(items []division_usecase.Division) []divisionResponse {
	out := make([]divisionResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toDivisionResponse(item))
	}
	return out
}
