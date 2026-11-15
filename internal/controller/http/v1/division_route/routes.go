package division_route

import (
	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/division_usecase"
	"ppl-m1-be/internal/usecase/utils"
	"ppl-m1-be/middleware"
)

func Mount(
	rg *gin.RouterGroup,
	uc *division_usecase.Usecase,
	tokens *utils.TokenMaker,
	store middleware.SessionStore,
) {
	handler := &Handler{uc: uc}
	divisions := rg.Group("/divisions",
		middleware.AuthenticateCurrent(tokens, store),
		middleware.RequireSuperAdmin(),
	)
	divisions.GET("", handler.List)
	divisions.POST("", handler.Create)
}
