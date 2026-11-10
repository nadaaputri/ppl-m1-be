package role_route

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/role_usecase"
	"ppl-m1-be/middleware"
	"ppl-m1-be/pkg/httpserver"
)

type Handler struct {
	uc *role_usecase.Usecase
}

func (h *Handler) List(c *gin.Context) {
	in, ok := listInput(c)
	if !ok {
		return
	}
	result, err := h.uc.List(c.Request.Context(), in)
	if err != nil {
		h.writeError(c, err)
		return
	}
	httpserver.SuccessList(c, toRoleList(result.Items), in.Page, in.PerPage, result.Total)
}

func (h *Handler) Create(c *gin.Context) {
	req, ok := bindCreate(c)
	if !ok {
		return
	}
	partnerID, ok := middleware.TargetPartnerID(c, req.PartnerID)
	if !ok {
		return
	}
	role, err := h.uc.Create(c.Request.Context(), role_usecase.CreateInput{
		PartnerID: partnerID, Code: req.Code, Name: req.Name, Description: req.Description,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	httpserver.Success(c, http.StatusCreated, toRoleResponse(role))
}

func listInput(c *gin.Context) (role_usecase.ListInput, bool) {
	page, perPage, ok := httpserver.ParsePage(c)
	if !ok {
		return role_usecase.ListInput{}, false
	}
	partnerID, ok := middleware.PartnerFilter(c)
	if !ok {
		return role_usecase.ListInput{}, false
	}
	active, ok := httpserver.QueryBool(c, "is_active")
	return role_usecase.ListInput{PartnerID: partnerID, IsActive: active, Page: page, PerPage: perPage}, ok
}

func bindCreate(c *gin.Context) (createRequest, bool) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpserver.BindError(c, err)
		return req, false
	}
	if fields := extraFieldErrors(req); len(fields) > 0 {
		httpserver.FailFields(c, httpserver.MsgValidation, fields)
		return req, false
	}
	return req, true
}

func (h *Handler) writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, role_usecase.ErrCodeReserved):
		httpserver.FailFields(c, httpserver.MsgValidation, map[string]string{"code": "Kode ini dicadangkan sistem"})
	case errors.Is(err, role_usecase.ErrCodeTaken):
		httpserver.Fail(c, http.StatusConflict, httpserver.CodeRoleCodeDuplicate, "Kode peran sudah dipakai")
	default:
		slog.Error("operasi peran gagal", "error", err)
		httpserver.Fail(c, http.StatusInternalServerError, httpserver.CodeInternal, "Terjadi kesalahan pada server")
	}
}
