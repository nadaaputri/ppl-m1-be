package audit_route

import (
	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/audit_usecase"
	"ppl-m1-be/internal/usecase/utils"
	"ppl-m1-be/middleware"
)

func Mount(
	rg *gin.RouterGroup,
	uc *audit_usecase.Usecase,
	tokens *utils.TokenMaker,
	store middleware.SessionStore,
) {
	handler := &Handler{uc: uc}
	auth := middleware.AuthenticateCurrent(tokens, store)

	rg.GET("/audit-logs", auth,
		middleware.RequireRole(middleware.RolePartner),
		middleware.RequirePermission("m1.audit.read"),
		handler.ListPartner,
	)

	rg.GET("/platform/audit-logs", auth,
		middleware.RequireSuperAdmin(),
		handler.ListPlatform,
	)
}
