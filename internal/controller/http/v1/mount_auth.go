package v1

import (
	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/controller/http/v1/auth_route"
	"ppl-m1-be/internal/usecase/audit_usecase"
	"ppl-m1-be/internal/usecase/auth_usecase"
	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func mountAuth(
	api *gin.RouterGroup,
	store *sqlc.Queries,
	tokens *utils.TokenMaker,
	audit audit_usecase.Logger,
) {
	uc := auth_usecase.New(store, tokens, audit)
	auth_route.Mount(api, uc, tokens, store)
}
