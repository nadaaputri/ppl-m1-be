package user_route

import (
	"strings"
	"time"

	"ppl-m1-be/internal/usecase/user_usecase"
	"ppl-m1-be/internal/usecase/utils"
)

type createUserRequest struct {
	PartnerID   *int32 `json:"partner_id"`
	Name        string `json:"name" binding:"required,max=150"`
	Email       string `json:"email" binding:"required,email,max=255"`
	Password    string `json:"password" binding:"required"`
	Role        string `json:"role" binding:"required"`
	Division    string `json:"division" binding:"required"`
	UnitUsahaID *int32 `json:"unit_usaha_id"`
}

type userResponse struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Role        *string   `json:"role"`
	Division    *string   `json:"division"`
	PartnerID   *int32    `json:"partner_id"`
	UnitUsahaID *int32    `json:"unit_usaha_id"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type statusRequest struct {
	IsActive *bool `json:"is_active"`
}

type statusResponse struct {
	ID       int32 `json:"id"`
	IsActive bool  `json:"is_active"`
}

func accountFieldErrors(req createUserRequest) map[string]string {
	fields := map[string]string{}
	if strings.TrimSpace(req.Name) == "" {
		fields["name"] = "Isian ini wajib diisi"
	}
	if msg := utils.ValidatePassword(req.Password, strings.TrimSpace(req.Email)); msg != "" {
		fields["password"] = msg
	}
	if req.UnitUsahaID != nil && *req.UnitUsahaID < 1 {
		fields["unit_usaha_id"] = "Nilai tidak valid"
	}
	return fields
}

func toAccountInput(
	req createUserRequest,
	partnerID, createdBy int32,
	ip string,
) user_usecase.AccountInput {
	return user_usecase.AccountInput{
		PartnerID: partnerID, CreatedBy: createdBy,
		Name: req.Name, Email: req.Email,
		Password: req.Password, Role: req.Role,
		Division: req.Division, UnitUsahaID: req.UnitUsahaID,
		IP: ip,
	}
}

func toUserResponse(a user_usecase.Account) userResponse {
	return userResponse{
		ID: a.ID, Name: a.Name, Email: a.Email, Role: a.Role, Division: a.Division,
		PartnerID: a.PartnerID, UnitUsahaID: a.UnitUsahaID, IsActive: a.IsActive,
		CreatedAt: a.CreatedAt.In(utils.WIB),
	}
}

func toUserList(items []user_usecase.Account) []userResponse {
	out := make([]userResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toUserResponse(item))
	}
	return out
}
