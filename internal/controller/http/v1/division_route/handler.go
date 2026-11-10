package division_route

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/division_usecase"
	"ppl-m1-be/middleware"
	"ppl-m1-be/pkg/httpserver"
)

type Handler struct {
	uc *division_usecase.Usecase
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
	httpserver.SuccessList(c, toDivisionList(result.Items), in.Page, in.PerPage, result.Total)
}

func (h *Handler) Create(c *gin.Context) {
	req, ok := bindCreate(c)
	if !ok {
		return
	}
	partnerID, ok := createPartnerID(c, req)
	if !ok {
		return
	}
	div, err := h.uc.Create(c.Request.Context(), division_usecase.CreateInput{
		PartnerID: partnerID, Code: req.Code, Name: req.Name,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	httpserver.Success(c, http.StatusCreated, toDivisionResponse(div))
}

func listInput(c *gin.Context) (division_usecase.ListInput, bool) {
	page, perPage, ok := httpserver.ParsePage(c)
	if !ok {
		return division_usecase.ListInput{}, false
	}
	partnerID, ok := partnerFilter(c)
	if !ok {
		return division_usecase.ListInput{}, false
	}
	active, ok := httpserver.QueryBool(c, "is_active")
	return division_usecase.ListInput{PartnerID: partnerID, IsActive: active, Page: page, PerPage: perPage}, ok
}

func partnerFilter(c *gin.Context) (*int32, bool) {
	scope, ok := middleware.PartnerScope(c)
	if !ok || scope != nil {
		return scope, ok
	}
	return httpserver.QueryInt32(c, "partner_id")
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

func createPartnerID(c *gin.Context, req createRequest) (int32, bool) {
	scope, ok := middleware.PartnerScope(c)
	if !ok {
		return 0, false
	}
	if scope != nil {
		return *scope, true
	}
	if req.PartnerID == nil || *req.PartnerID < 1 {
		httpserver.FailFields(c, httpserver.MsgValidation, map[string]string{"partner_id": "Isian ini wajib diisi"})
		return 0, false
	}
	return *req.PartnerID, true
}

func (h *Handler) writeError(c *gin.Context, err error) {
	if errors.Is(err, division_usecase.ErrCodeTaken) {
		httpserver.Fail(c, http.StatusConflict, httpserver.CodeDivisionCodeDuplicate, "Kode divisi sudah dipakai")
		return
	}
	slog.Error("operasi divisi gagal", "error", err)
	httpserver.Fail(c, http.StatusInternalServerError, httpserver.CodeInternal, "Terjadi kesalahan pada server")
}