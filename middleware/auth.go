package middleware

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
	"ppl-m1-be/pkg/httpserver"
)

const claimsKey = "claims"
const RolePartner = "partner"

type SessionStore interface {
	GetUserSession(context.Context, int32) (sqlc.GetUserSessionRow, error)
	GetUserPermissionCodes(context.Context, int32) ([]string, error)
}

func Authenticate(tokens *utils.TokenMaker) gin.HandlerFunc {
	return authenticate(tokens, nil)
}

func AuthenticateCurrent(tokens *utils.TokenMaker, store SessionStore) gin.HandlerFunc {
	return authenticate(tokens, store)
}

func authenticate(tokens *utils.TokenMaker, store SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := tokens.Parse(bearerToken(c.GetHeader("Authorization")))
		if err != nil {
			unauthenticated(c)
			return
		}
		if store != nil && !refreshSession(c, claims, store) {
			return
		}
		c.Set(claimsKey, claims)
		c.Next()
	}
}

func refreshSession(c *gin.Context, claims *utils.Claims, store SessionStore) bool {
	session, ok := loadSession(c, claims.UserID, store)
	if !ok {
		return false
	}
	permissions, err := store.GetUserPermissionCodes(c.Request.Context(), claims.UserID)
	if err != nil {
		serviceUnavailable(c)
		return false
	}
	applySessionClaims(claims, session, permissions)
	return true
}

func loadSession(
	c *gin.Context,
	userID int32,
	store SessionStore,
) (sqlc.GetUserSessionRow, bool) {
	session, err := store.GetUserSession(c.Request.Context(), userID)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && !session.IsActive) {
		unauthenticated(c)
		return session, false
	}
	if err != nil {
		serviceUnavailable(c)
		return session, false
	}
	if !activeReferences(session) {
		forbidden(c)
		return session, false
	}
	return session, true
}

func activeReferences(s sqlc.GetUserSessionRow) bool {
	if s.IsSuperAdmin {
		return true
	}
	return s.RoleID.Valid && s.RoleActive && (!s.DivisionID.Valid || s.DivisionActive)
}

func applySessionClaims(c *utils.Claims, s sqlc.GetUserSessionRow, permissions []string) {
	c.Name = s.Name
	c.Role = s.RoleCode.String
	c.Division = utils.StringPtr(s.DivisionCode)
	c.PartnerID = utils.Int32Ptr(s.PartnerID)
	c.UnitUsahaID = utils.Int32Ptr(s.UnitUsahaID)
	c.IsSuperAdmin = s.IsSuperAdmin
	c.Permissions = permissions
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaims(c)
		if claims == nil {
			unauthenticated(c)
			return
		}
		if !claims.IsSuperAdmin && !slices.Contains(roles, claims.Role) {
			forbidden(c)
			return
		}
		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaims(c)
		if claims == nil {
			unauthenticated(c)
			return
		}
		if !claims.IsSuperAdmin && !slices.Contains(claims.Permissions, permission) {
			forbidden(c)
			return
		}
		c.Next()
	}
}

func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaims(c)
		if claims == nil {
			unauthenticated(c)
			return
		}
		if !claims.IsSuperAdmin {
			forbidden(c)
			return
		}
		c.Next()
	}
}

func RequireDivision(divisions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaims(c)
		if claims == nil {
			unauthenticated(c)
			return
		}
		if !claims.IsSuperAdmin && (claims.Division == nil || !slices.Contains(divisions, *claims.Division)) {
			forbidden(c)
			return
		}
		c.Next()
	}
}

func RequireRequestedDivision(param string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaims(c)
		if claims == nil {
			unauthenticated(c)
			return
		}
		if claims.IsSuperAdmin {
			c.Next()
			return
		}
		if !matchesRequestedDivision(c, claims.Division, param) {
			forbidden(c)
			return
		}
		c.Next()
	}
}

func matchesRequestedDivision(c *gin.Context, division *string, param string) bool {
	target := c.Param(param)
	if target == "" {
		target = c.Query(param)
	}
	return division != nil && target != "" && strings.EqualFold(*division, target)
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

func serviceUnavailable(c *gin.Context) {
	httpserver.Fail(c, http.StatusServiceUnavailable, httpserver.CodeDatabaseUnavailable,
		"Database tidak dapat diakses")
	c.Abort()
}
