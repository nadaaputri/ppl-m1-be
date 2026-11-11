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
