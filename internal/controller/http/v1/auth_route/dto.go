package auth_route

import (
	"ppl-m1-be/internal/usecase/auth_usecase"
	"ppl-m1-be/internal/usecase/utils"
	"time"
)

const maxPasswordBytes = 72

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type meResponse struct {
	UserID       int32   `json:"user_id"`
	Name         string  `json:"name"`
	Role         string  `json:"role"`
	Division     *string `json:"division"`
	PartnerID    *int32  `json:"partner_id"`
	UnitUsahaID  *int32  `json:"unit_usaha_id"`
	IsSuperAdmin bool    `json:"is_super_admin"`
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type changePasswordResponse struct {
	ChangedAt time.Time `json:"changed_at"`
}

func toMeResponse(c *utils.Claims) meResponse {
	return meResponse{
		UserID:       c.UserID,
		Name:         c.Name,
		Role:         c.Role,
		Division:     c.Division,
		PartnerID:    c.PartnerID,
		UnitUsahaID:  c.UnitUsahaID,
		IsSuperAdmin: c.IsSuperAdmin,
	}
}

func toLoginInput(req loginRequest, ip string) auth_usecase.LoginInput {
	return auth_usecase.LoginInput{Email: req.Email, Password: req.Password, IP: ip}
}

func toLoginResponse(r auth_usecase.LoginResult) loginResponse {
	return loginResponse{
		AccessToken: r.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   r.ExpiresIn,
	}
}
