package v1

import (
	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/controller/http/v1/audit_route"
	"ppl-m1-be/internal/usecase/audit_usecase"
	"ppl-m1-be/internal/usecase/utils"
)

func mountAudit(
	api *gin.RouterGroup,
	uc *audit_usecase.Usecase,
	tokens *utils.TokenMaker,
) {
	audit_route.Mount(api, uc, tokens)
}