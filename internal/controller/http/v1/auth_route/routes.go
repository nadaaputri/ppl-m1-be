package auth_route

import (
	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/auth_usecase"
)

func Mount(rg *gin.RouterGroup, uc *auth_usecase.Usecase) {
	handler := &Handler{uc: uc}
	group := rg.Group("/auth")
	group.POST("/login", handler.Login)
}