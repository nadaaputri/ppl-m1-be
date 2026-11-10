package role_route

import (
	"strings"
	"time"

	"ppl-m1-be/internal/usecase/role_usecase"
	"ppl-m1-be/internal/usecase/utils"
)

type createRequest struct {
	PartnerID   *int32 `json:"partner_id"`
	Code        string `json:"code" binding:"required,max=50"`
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description"`
}

type roleResponse struct {
	ID          int32     `json:"id"`
	PartnerID   *int32    `json:"partner_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

func extraFieldErrors(req createRequest) map[string]string {
	fields := map[string]string{}
	if !utils.CodePattern.MatchString(req.Code) {
		fields["code"] = "Hanya huruf kecil, angka, tanda hubung, dan garis bawah"
	}
	if strings.TrimSpace(req.Name) == "" {
		fields["name"] = "Isian ini wajib diisi"
	}
	return fields
}

func toRoleResponse(r role_usecase.Role) roleResponse {
	return roleResponse{
		ID: r.ID, PartnerID: r.PartnerID, Code: r.Code, Name: r.Name,
		Description: r.Description, IsActive: r.IsActive, CreatedAt: r.CreatedAt.In(utils.WIB),
	}
}

func toRoleList(items []role_usecase.Role) []roleResponse {
	out := make([]roleResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toRoleResponse(item))
	}
	return out
}
