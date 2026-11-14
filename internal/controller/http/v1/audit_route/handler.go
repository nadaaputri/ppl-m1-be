package audit_route

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/audit_usecase"
	"ppl-m1-be/middleware"
	"ppl-m1-be/pkg/httpserver"
)

type Handler struct {
	uc *audit_usecase.Usecase
}

func (h *Handler) ListPlatform(c *gin.Context) {
	h.list(c, nil)
}

func (h *Handler) ListPartner(c *gin.Context) {
	claims := middleware.GetClaims(c)
	if claims == nil || claims.PartnerID == nil {
		httpserver.Fail(c, http.StatusForbidden, httpserver.CodeForbidden,
			"Anda tidak memiliki akses ke fitur ini")
		return
	}
	h.list(c, claims.PartnerID)
}

func (h *Handler) list(c *gin.Context, partnerID *int32) {
	page, perPage, ok := httpserver.ParsePage(c)
	if !ok {
		return
	}
	result, err := h.uc.List(c.Request.Context(), audit_usecase.ListInput{
		PartnerID: partnerID, Page: page, PerPage: perPage,
	})
	if err != nil {
		slog.Error("gagal mengambil audit log", "error", err)
		httpserver.Fail(c, http.StatusInternalServerError, httpserver.CodeInternal,
			"Terjadi kesalahan pada server")
		return
	}
	httpserver.SuccessList(c, result.Items, page, perPage, result.Total)
}