package v1

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/controller/http/v1/user_route"
	"ppl-m1-be/internal/usecase/audit_usecase"
	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/user_usecase"
	"ppl-m1-be/internal/usecase/utils"
)

func mountUser(
	api *gin.RouterGroup,
	db *sql.DB,
	tokens *utils.TokenMaker,
	audit audit_usecase.Logger,
) {
	uc := user_usecase.New(sqlc.New(db), audit)
	user_route.Mount(api, uc, tokens)
}