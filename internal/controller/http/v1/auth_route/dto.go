package auth_route

import "ppl-m1-be/internal/usecase/auth_usecase"

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