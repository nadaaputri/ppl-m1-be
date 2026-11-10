package middleware

import (
	"ppl-m1-be/pkg/httpserver"

	"github.com/gin-gonic/gin"
)

func PartnerScope(c *gin.Context) (*int32, bool) {
	claims := GetClaims(c)
	if claims.IsSuperAdmin {
		return nil, true
	}
	if claims.PartnerID == nil {
		forbidden(c)
		return nil, false
	}
	return claims.PartnerID, true
}

func PartnerFilter(c *gin.Context) (*int32, bool) {
	scope, ok := PartnerScope(c)
	if !ok || scope != nil {
		return scope, ok
	}
	return httpserver.QueryInt32(c, "partner_id")
}

func TargetPartnerID(c *gin.Context, bodyID *int32) (int32, bool) {
	scope, ok := PartnerScope(c)
	if !ok {
		return 0, false
	}
	if scope != nil {
		return *scope, true
	}
	if bodyID == nil || *bodyID < 1 {
		httpserver.FailFields(c, httpserver.MsgValidation, map[string]string{"partner_id": "Isian ini wajib diisi"})
		return 0, false
	}
	return *bodyID, true
}
