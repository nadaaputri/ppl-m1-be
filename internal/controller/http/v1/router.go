package v1

import (
	"database/sql"
	"net/http"

	"ppl-m1-be/config"
	apidocs "ppl-m1-be/docs"
	"ppl-m1-be/internal/usecase/audit_usecase"
	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
	"ppl-m1-be/middleware"
	"ppl-m1-be/pkg/httpserver"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB, cfg config.Config) *gin.Engine {
	router := gin.New()
	configureRouter(router, cfg)
	mountSystemRoutes(router, db)
	mountModuleRoutes(router, db, cfg)
	return router
}

func configureRouter(router *gin.Engine, cfg config.Config) {
	if err := router.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		panic(err)
	}
	router.Use(gin.Logger(), gin.CustomRecovery(recoverPanic), middleware.CORS(cfg.CORSAllowedOrigins))
	router.HandleMethodNotAllowed = true
	router.NoRoute(notFound)
	router.NoMethod(methodNotAllowed)
}

func mountSystemRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/openapi.yaml", openAPISpec)
	router.GET("/api/m1/v1/docs", docsHandler)
	router.GET("/health", healthCheck)
	router.GET("/ready", readinessCheck(db))
}

func openAPISpec(c *gin.Context) {
	c.Data(http.StatusOK, "application/yaml; charset=utf-8", apidocs.OpenAPI)
}

func healthCheck(c *gin.Context) {
	httpserver.Success(c, http.StatusOK, gin.H{"status": "ok"})
}

func readinessCheck(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.PingContext(c.Request.Context()); err != nil {
			httpserver.Fail(c, http.StatusServiceUnavailable, httpserver.CodeDatabaseUnavailable, "Database tidak dapat diakses")
			return
		}
		httpserver.Success(c, http.StatusOK, gin.H{"status": "ready", "database": "connected"})
	}
}

func mountModuleRoutes(router *gin.Engine, db *sql.DB, cfg config.Config) {
	tokens := utils.NewTokenMaker(cfg.JWTSecret, cfg.AppName, cfg.JWTExpiresIn)
	api := router.Group("/api/m1/v1")
	audit := audit_usecase.New(sqlc.New(db))

	mountAuth(api, db, tokens, audit)
	mountDivision(api, db, tokens)
	mountRole(api, db, tokens)
	mountUser(api, db, tokens, audit)
	mountAudit(api, audit, tokens)
}