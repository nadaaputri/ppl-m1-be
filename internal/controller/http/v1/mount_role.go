package v1

import (
	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/controller/http/v1/role_route"
	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/role_usecase"
	"ppl-m1-be/internal/usecase/utils"
)

func mountRole(api *gin.RouterGroup, store *sqlc.Queries, tokens *utils.TokenMaker) {
	uc := role_usecase.New(store)
	role_route.Mount(api, uc, tokens, store)
}
