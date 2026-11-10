package user_route

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/user_usecase"
	"ppl-m1-be/middleware"
	"ppl-m1-be/pkg/httpserver"
)

func (h *Handler) CreateUser(c *gin.Context) {
	req, ok := bindCreateUser(c)
	if !ok {
		return
	}
	partnerID, ok := middleware.TargetPartnerID(c, req.PartnerID)
	if !ok {
		return
	}
	in := toAccountInput(req, partnerID, middleware.GetClaims(c).UserID)
	account, err := h.uc.CreateAccount(c.Request.Context(), in)
	if err != nil {
		h.writeAccountError(c, err)
		return
	}
	httpserver.Success(c, http.StatusCreated, toUserResponse(account))
}

func bindCreateUser(c *gin.Context) (createUserRequest, bool) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpserver.BindError(c, err)
		return req, false
	}
	if fields := accountFieldErrors(req); len(fields) > 0 {
		httpserver.FailFields(c, httpserver.MsgValidation, fields)
		return req, false
	}
	return req, true
}

func (h *Handler) writeAccountError(c *gin.Context, err error) {
	var verr *user_usecase.ValidationError
	switch {
	case errors.As(err, &verr):
		httpserver.FailFields(c, httpserver.MsgValidation, verr.Fields)
	case errors.Is(err, user_usecase.ErrEmailTaken):
		httpserver.Fail(c, http.StatusConflict, httpserver.CodeEmailDuplicate, "Email sudah terdaftar")
	default:
		slog.Error("tambah akun gagal", "error", err)
		httpserver.Fail(c, http.StatusInternalServerError, httpserver.CodeInternal, "Terjadi kesalahan pada server")
	}
}
