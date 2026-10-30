package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/utils"
	"ppl-m1-be/pkg/httpserver"
)

const claimsKey = "claims"

func Authenticate(tokens *utils.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := tokens.Parse(bearerToken(c.GetHeader("Authorization")))
		if err != nil {
			unauthenticated(c)
			return
		}
		c.Set(claimsKey, claims)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaims(c)
		if claims == nil {
			unauthenticated(c)
			return
		}
		if !slices.Contains(roles, claims.Role) {
			forbidden(c)
			return
		}
		c.Next()
	}
}

func GetClaims(c *gin.Context) *utils.Claims {
	value, ok := c.Get(claimsKey)
	if !ok {
		return nil
	}
	claims, _ := value.(*utils.Claims)
	return claims
}

func bearerToken(header string) string {
	scheme, token, found := strings.Cut(header, " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return ""
	}
	return strings.TrimSpace(token)
}

func unauthenticated(c *gin.Context) {
	httpserver.Fail(c, http.StatusUnauthorized, httpserver.CodeUnauthenticated,
		"Sesi tidak valid atau telah berakhir, silakan login kembali")
	c.Abort()
}

func forbidden(c *gin.Context) {
	httpserver.Fail(c, http.StatusForbidden, httpserver.CodeForbidden,
		"Anda tidak memiliki akses ke fitur ini")
	c.Abort()
}