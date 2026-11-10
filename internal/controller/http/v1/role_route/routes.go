package role_route

import (
	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/role_usecase"
	"ppl-m1-be/internal/usecase/utils"
	"ppl-m1-be/middleware"
)

func Mount(rg *gin.RouterGroup, uc *role_usecase.Usecase, tokens *utils.TokenMaker) {
	handler := &Handler{uc: uc}
	roles := rg.Group("/roles", middleware.Authenticate(tokens), middleware.RequireRole(middleware.RolePartner))
	roles.GET("", handler.List)
	roles.POST("", handler.Create)
}
