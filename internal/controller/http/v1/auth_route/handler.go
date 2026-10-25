package auth_route

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/auth_usecase"
	"ppl-m1-be/pkg/httpserver"
)

type Handler struct {
	uc *auth_usecase.Usecase
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Password) > maxPasswordBytes {
		httpserver.Fail(c, http.StatusBadRequest, "Data yang dikirim tidak valid")
		return
	}

	result, err := h.uc.Login(c.Request.Context(), toLoginInput(req, c.ClientIP()))
	if err != nil {
		h.writeError(c, err)
		return
	}
	httpserver.Success(c, http.StatusOK, "Login berhasil", toLoginResponse(result))
}

func (h *Handler) writeError(c *gin.Context, err error) {
	if errors.Is(err, auth_usecase.ErrInvalidCredentials) {
		httpserver.Fail(c, http.StatusUnauthorized, "Email atau kata sandi salah")
		return
	}
	slog.Error("login gagal", "error", err)
	httpserver.Fail(c, http.StatusInternalServerError, "Terjadi kesalahan pada server")
}