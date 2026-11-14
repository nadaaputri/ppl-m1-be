package user_route

import (
	"errors"
	"log/slog"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/user_usecase"
	"ppl-m1-be/middleware"
	"ppl-m1-be/pkg/httpserver"
)

var allowedSorts = []string{"created_at", "-created_at", "name", "-name", "email", "-email"}

var statusErrors = []struct {
	err     error
	status  int
	code    string
	message string
}{
	{user_usecase.ErrUserNotFound, http.StatusNotFound, httpserver.CodeUserNotFound, "Pengguna tidak ditemukan"},
	{user_usecase.ErrNotAllowed, http.StatusForbidden, httpserver.CodeForbidden, "Anda tidak berhak mengubah status akun ini"},
	{user_usecase.ErrCannotDeactivateSelf, http.StatusConflict, httpserver.CodeCannotDeactivateSelf, "Tidak dapat menonaktifkan akun sendiri"},
	{user_usecase.ErrAlreadyActive, http.StatusConflict, httpserver.CodeUserAlreadyActive, "Akun sudah aktif"},
	{user_usecase.ErrAlreadyInactive, http.StatusConflict, httpserver.CodeUserAlreadyInactive, "Akun sudah nonaktif"},
	{user_usecase.ErrReferenceInactive, http.StatusConflict, httpserver.CodeUserReferenceInactive, "Peran atau divisi akun nonaktif, ganti dulu sebelum mengaktifkan"},
}

func (h *Handler) CreateUser(c *gin.Context) {
	req, ok := bindCreateUser(c)
	if !ok {
		return
	}
	partnerID, ok := middleware.TargetPartnerID(c, req.PartnerID)
	if !ok {
		return
	}
	in := toAccountInput(
		req, partnerID, middleware.GetClaims(c).UserID, c.ClientIP(),
	)
	account, err := h.uc.CreateAccount(c.Request.Context(), in)
	if err != nil {
		h.writeAccountError(c, err)
		return
	}
	httpserver.Success(c, http.StatusCreated, toUserResponse(account))
}

func (h *Handler) GetUser(c *gin.Context) {
	id, ok := httpserver.ParamInt32(c, "id")
	if !ok {
		return
	}
	account, err := h.uc.GetAccount(c.Request.Context(), id)
	if err != nil {
		h.writeAccountError(c, err)
		return
	}
	httpserver.Success(c, http.StatusOK, toUserResponse(account))
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, ok := httpserver.ParamInt32(c, "id")
	if !ok {
		return
	}
	req, ok := bindUpdateUser(c)
	if !ok {
		return
	}
	in := toUpdateAccountInput(req, id, middleware.GetClaims(c).UserID, c.ClientIP())
	account, err := h.uc.UpdateAccount(c.Request.Context(), in)
	if err != nil {
		h.writeAccountError(c, err)
		return
	}
	httpserver.Success(c, http.StatusOK, toUserResponse(account))
}

func bindUpdateUser(c *gin.Context) (updateUserRequest, bool) {
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpserver.BindError(c, err)
		return req, false
	}
	if fields := updateFieldErrors(req); len(fields) > 0 {
		httpserver.FailFields(c, httpserver.MsgValidation, fields)
		return req, false
	}
	return req, true
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
	case errors.Is(err, user_usecase.ErrUserNotFound):
		httpserver.Fail(c, http.StatusNotFound, httpserver.CodeUserNotFound, "Pengguna tidak ditemukan")
	default:
		slog.Error("tambah akun gagal", "error", err)
		httpserver.Fail(c, http.StatusInternalServerError, httpserver.CodeInternal, "Terjadi kesalahan pada server")
	}
}

func (h *Handler) ListUsers(c *gin.Context) {
	in, ok := userListInput(c)
	if !ok {
		return
	}
	result, err := h.uc.ListAccounts(c.Request.Context(), in)
	if err != nil {
		h.writeAccountError(c, err)
		return
	}
	httpserver.SuccessList(c, toUserList(result.Items), in.Page, in.PerPage, result.Total)
}

func userListInput(c *gin.Context) (user_usecase.AccountListInput, bool) {
	page, perPage, ok := httpserver.ParsePage(c)
	if !ok {
		return user_usecase.AccountListInput{}, false
	}
	partnerID, ok := middleware.PartnerFilter(c)
	if !ok {
		return user_usecase.AccountListInput{}, false
	}
	in, ok := userFilters(c)
	in.Page, in.PerPage, in.PartnerID = page, perPage, partnerID
	return in, ok
}

func userFilters(c *gin.Context) (user_usecase.AccountListInput, bool) {
	unitUsahaID, ok := httpserver.QueryInt32(c, "unit_usaha_id")
	if !ok {
		return user_usecase.AccountListInput{}, false
	}
	active, ok := httpserver.QueryBool(c, "is_active")
	if !ok {
		return user_usecase.AccountListInput{}, false
	}
	sort, ok := userSort(c)
	return user_usecase.AccountListInput{
		Search: c.Query("q"), Role: c.Query("role"), Division: c.Query("division"),
		UnitUsahaID: unitUsahaID, IsActive: active, Sort: sort,
	}, ok
}

func userSort(c *gin.Context) (string, bool) {
	sort := c.DefaultQuery("sort", "-created_at")
	if !slices.Contains(allowedSorts, sort) {
		httpserver.Fail(c, http.StatusBadRequest, httpserver.CodeInvalidParameter, "Parameter tidak valid")
		return "", false
	}
	return sort, true
}

func (h *Handler) SetUserStatus(c *gin.Context) {
	in, ok := statusInput(c)
	if !ok {
		return
	}
	result, err := h.uc.SetAccountStatus(c.Request.Context(), in)
	if err != nil {
		h.writeStatusError(c, err)
		return
	}
	httpserver.Success(c, http.StatusOK, statusResponse{ID: result.ID, IsActive: result.IsActive})
}

func (h *Handler) DeactivateUser(c *gin.Context) {
	in, ok := deactivateInput(c)
	if !ok {
		return
	}
	result, err := h.uc.SetAccountStatus(c.Request.Context(), in)
	if err != nil {
		h.writeStatusError(c, err)
		return
	}
	httpserver.Success(c, http.StatusOK, statusResponse{
		ID: result.ID, IsActive: result.IsActive,
	})
}

func deactivateInput(c *gin.Context) (user_usecase.StatusInput, bool) {
	id, ok := httpserver.ParamInt32(c, "id")
	if !ok {
		return user_usecase.StatusInput{}, false
	}
	scope, ok := middleware.PartnerScope(c)
	if !ok {
		return user_usecase.StatusInput{}, false
	}
	return user_usecase.StatusInput{
		ActorID:      middleware.GetClaims(c).UserID,
		PartnerScope: scope,
		TargetID:     id,
		Active:       false,
		IP:           c.ClientIP(),
	}, true
}

func statusInput(c *gin.Context) (user_usecase.StatusInput, bool) {
	id, ok := httpserver.ParamInt32(c, "id")
	if !ok {
		return user_usecase.StatusInput{}, false
	}
	req, ok := bindStatus(c)
	if !ok {
		return user_usecase.StatusInput{}, false
	}
	scope, ok := middleware.PartnerScope(c)
	return user_usecase.StatusInput{
		ActorID: middleware.GetClaims(c).UserID, PartnerScope: scope, TargetID: id, Active: *req.IsActive, IP: c.ClientIP(),
	}, ok
}

func bindStatus(c *gin.Context) (statusRequest, bool) {
	var req statusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpserver.BindError(c, err)
		return req, false
	}
	if req.IsActive == nil {
		httpserver.FailFields(c, httpserver.MsgValidation, map[string]string{"is_active": "Isian ini wajib diisi"})
		return req, false
	}
	return req, true
}

func (h *Handler) writeStatusError(c *gin.Context, err error) {
	for _, e := range statusErrors {
		if errors.Is(err, e.err) {
			httpserver.Fail(c, e.status, e.code, e.message)
			return
		}
	}
	slog.Error("ubah status akun gagal", "error", err)
	httpserver.Fail(c, http.StatusInternalServerError, httpserver.CodeInternal, "Terjadi kesalahan pada server")
}
