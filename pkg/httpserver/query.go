package httpserver

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryInt32(c *gin.Context, key string) (*int32, bool) {
	raw := c.Query(key)
	if raw == "" {
		return nil, true
	}
	n, err := strconv.ParseInt(raw, 10, 32)
	if err != nil || n < 1 {
		Fail(c, http.StatusBadRequest, CodeInvalidParameter, "Parameter tidak valid")
		return nil, false
	}
	value := int32(n)
	return &value, true
}

func QueryBool(c *gin.Context, key string) (*bool, bool) {
	raw := c.Query(key)
	if raw == "" {
		return nil, true
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		Fail(c, http.StatusBadRequest, CodeInvalidParameter, "Parameter tidak valid")
		return nil, false
	}
	return &value, true
}
