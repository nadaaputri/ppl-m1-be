package auth_route

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/auth_usecase"
	"ppl-m1-be/internal/usecase/utils"
	"ppl-m1-be/middleware"
	"ppl-m1-be/pkg/httpserver"
)

func (h *Handler) ChangePassword(c *gin.Context) {
	req, ok := bindChangePassword(c)
	if !ok {
		return
	}
	in := auth_usecase.ChangePasswordInput{
		UserID: middleware.GetClaims(c).UserID, OldPassword: req.OldPassword,
		NewPassword: req.NewPassword, IP: c.ClientIP(),
	}
	changedAt, err := h.uc.ChangePassword(c.Request.Context(), in)
	if err != nil {
		h.writePasswordError(c, err)
		return
	}
	httpserver.Success(c, http.StatusOK, changePasswordResponse{ChangedAt: changedAt.In(utils.WIB)})
}

func bindChangePassword(c *gin.Context) (changePasswordRequest, bool) {
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpserver.BindError(c, err)
		return req, false
	}
	fields := changePasswordFieldErrors(req)
	if len(fields) > 0 {
		httpserver.FailFields(c, httpserver.MsgValidation, fields)
		return req, false
	}
	return req, true
}

func changePasswordFieldErrors(req changePasswordRequest) map[string]string {
	fields := map[string]string{}
	if req.OldPassword == "" {
		fields["old_password"] = "Isian ini wajib diisi"
	}
	if req.NewPassword == "" {
		fields["new_password"] = "Isian ini wajib diisi"
	}
	if req.ConfirmPassword == "" {
		fields["confirm_password"] = "Isian ini wajib diisi"
	} else if req.ConfirmPassword != req.NewPassword {
		fields["confirm_password"] = "Konfirmasi kata sandi baru tidak cocok"
	}
	return fields
}

func (h *Handler) writePasswordError(c *gin.Context, err error) {
	var newErr *auth_usecase.NewPasswordError
	switch {
	case errors.As(err, &newErr):
		httpserver.FailFields(c, httpserver.MsgValidation, map[string]string{"new_password": newErr.Message})
	case errors.Is(err, auth_usecase.ErrOldPasswordInvalid):
		httpserver.FailFieldsCode(c, httpserver.CodeOldPasswordInvalid, "Kata sandi lama salah",
			map[string]string{"old_password": "Kata sandi lama salah"})
	case errors.Is(err, auth_usecase.ErrAccountUnavailable):
		httpserver.Fail(c, http.StatusUnauthorized, httpserver.CodeUnauthenticated, "Sesi tidak valid atau telah berakhir, silakan login kembali")
	default:
		slog.Error("ubah kata sandi gagal", "error", err)
		httpserver.Fail(c, http.StatusInternalServerError, httpserver.CodeInternal, "Terjadi kesalahan pada server")
	}
}
