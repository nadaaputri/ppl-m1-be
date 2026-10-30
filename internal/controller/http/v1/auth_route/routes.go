package auth_route

import (
	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/auth_usecase"
	"ppl-m1-be/internal/usecase/utils"
	"ppl-m1-be/middleware"
)

func Mount(rg *gin.RouterGroup, uc *auth_usecase.Usecase, tokens *utils.TokenMaker) {
	handler := &Handler{uc: uc}
	protected := middleware.Authenticate(tokens)

	sessions := rg.Group("/auth/sessions")
	sessions.POST("", handler.Login)
	sessions.POST("/logout", protected, handler.Logout)

	rg.GET("/me", protected, handler.Me)
}