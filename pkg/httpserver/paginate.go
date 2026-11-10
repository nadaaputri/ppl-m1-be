package httpserver

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const maxPerPage = 100

type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func ParsePage(c *gin.Context) (int, int, bool) {
	page, errPage := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, errPer := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if errPage != nil || errPer != nil || page < 1 || perPage < 1 || perPage > maxPerPage {
		Fail(c, http.StatusBadRequest, CodeInvalidParameter, "Parameter tidak valid")
		return 0, 0, false
	}
	return page, perPage, true
}

func SuccessList(c *gin.Context, data any, page, perPage, total int) {
	meta := Meta{Page: page, PerPage: perPage, Total: total, TotalPages: (total + perPage - 1) / perPage}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data, "meta": meta})
}