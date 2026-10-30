package auth_route

import (
	"ppl-m1-be/internal/usecase/auth_usecase"
	"ppl-m1-be/internal/usecase/utils"
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
	UserID      int32   `json:"user_id"`
	Name        string  `json:"name"`
	Role        string  `json:"role"`
	Division    *string `json:"division"`
	UnitUsahaID *int32  `json:"unit_usaha_id"`
}

func toMeResponse(c *utils.Claims) meResponse {
	return meResponse{
		UserID:      c.UserID,
		Name:        c.Name,
		Role:        c.Role,
		Division:    c.Division,
		UnitUsahaID: c.UnitUsahaID,
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