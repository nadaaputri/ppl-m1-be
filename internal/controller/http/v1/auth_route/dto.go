package auth_route

import "ppl-m1-be/internal/usecase/auth_usecase"

const maxPasswordBytes = 72

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type divisionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type userResponse struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Email    string            `json:"email"`
	Role     string            `json:"role"`
	Division *divisionResponse `json:"division"`
}

type loginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int          `json:"expires_in"`
	User        userResponse `json:"user"`
}

func toLoginInput(req loginRequest, ip string) auth_usecase.LoginInput {
	return auth_usecase.LoginInput{Email: req.Email, Password: req.Password, IP: ip}
}

func toLoginResponse(r auth_usecase.LoginResult) loginResponse {
	return loginResponse{
		AccessToken: r.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   r.ExpiresIn,
		User:        toUserResponse(r.User),
	}
}

func toUserResponse(u auth_usecase.UserInfo) userResponse {
	return userResponse{
		ID: u.ID, Name: u.Name, Email: u.Email, Role: u.Role,
		Division: toDivisionResponse(u),
	}
}

func toDivisionResponse(u auth_usecase.UserInfo) *divisionResponse {
	if u.DivisionID == "" {
		return nil
	}
	return &divisionResponse{ID: u.DivisionID, Name: u.DivisionName}
}