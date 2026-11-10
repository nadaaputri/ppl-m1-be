package v1

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/controller/http/v1/division_route"
	"ppl-m1-be/internal/usecase/division_usecase"
	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func mountDivision(api *gin.RouterGroup, db *sql.DB, tokens *utils.TokenMaker) {
	uc := division_usecase.New(sqlc.New(db))
	division_route.Mount(api, uc, tokens)
}