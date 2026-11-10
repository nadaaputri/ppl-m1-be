package middleware

import "github.com/gin-gonic/gin"

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