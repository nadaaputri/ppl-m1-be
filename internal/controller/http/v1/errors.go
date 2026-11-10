package v1

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"ppl-m1-be/pkg/httpserver"
)

func recoverPanic(c *gin.Context, err any) {
	slog.Error("panic tertangkap", "error", err, "path", c.Request.URL.Path)
	httpserver.Fail(c, http.StatusInternalServerError, httpserver.CodeInternal, "Terjadi kesalahan pada server")
	c.Abort()
}

func notFound(c *gin.Context) {
	httpserver.Fail(c, http.StatusNotFound, httpserver.CodeNotFound, "Endpoint tidak ditemukan")
}

func methodNotAllowed(c *gin.Context) {
	httpserver.Fail(c, http.StatusMethodNotAllowed, httpserver.CodeMethodNotAllowed, "Method tidak diizinkan untuk endpoint ini")
}
