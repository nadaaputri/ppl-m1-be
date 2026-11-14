package audit_route

import (
	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/audit_usecase"
	"ppl-m1-be/internal/usecase/utils"
	"ppl-m1-be/middleware"
)

func Mount(rg *gin.RouterGroup, uc *audit_usecase.Usecase, tokens *utils.TokenMaker) {
	handler := &Handler{uc: uc}

	rg.GET(
		"/audit-logs",
		middleware.Authenticate(tokens),
		middleware.RequireRole(middleware.RolePartner),
		handler.ListPartner,
	)

	rg.GET(
		"/platform/audit-logs",
		middleware.Authenticate(tokens),
		middleware.RequireSuperAdmin(),
		handler.ListPlatform,
	)
}