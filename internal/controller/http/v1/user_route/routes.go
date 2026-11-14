package user_route

import (
	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/user_usecase"
	"ppl-m1-be/internal/usecase/utils"
	"ppl-m1-be/middleware"
)

func Mount(rg *gin.RouterGroup, uc *user_usecase.Usecase, tokens *utils.TokenMaker) {
	handler := &Handler{uc: uc}
	rg.POST("/auth/registrations", handler.Register)

	users := rg.Group("/users", middleware.Authenticate(tokens), middleware.RequireRole(middleware.RolePartner))
	users.POST("", handler.CreateUser)
	users.GET("", handler.ListUsers)
	users.PATCH("/:id/status", handler.SetUserStatus)
	users.DELETE("/:id", handler.DeactivateUser)
}
