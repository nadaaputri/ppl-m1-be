package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ppl-m1-be/internal/usecase/utils"

	"github.com/gin-gonic/gin"
)

func authTestRouter(tokens *utils.TokenMaker) *gin.Engine {
	router := gin.New()
	router.GET("/private", Authenticate(tokens), func(c *gin.Context) {
		if GetClaims(c) == nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	})
	return router
}

func TestAuthenticateRejectsMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tokens := utils.NewTokenMaker(strings.Repeat("s", 32), "araya", 8*time.Hour)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	authTestRouter(tokens).ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status seharusnya 401, didapat %d", res.Code)
	}
}

func TestAuthenticateRejectsInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tokens := utils.NewTokenMaker(strings.Repeat("s", 32), "araya", 8*time.Hour)
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer token-tidak-valid")
	res := httptest.NewRecorder()
	authTestRouter(tokens).ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status seharusnya 401, didapat %d", res.Code)
	}
}

func TestAuthenticateAcceptsValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tokens := utils.NewTokenMaker(strings.Repeat("s", 32), "araya", 8*time.Hour)
	raw, err := tokens.Generate(utils.Claims{UserID: 1, Name: "Rani", Role: "staff"})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer "+raw)
	res := httptest.NewRecorder()
	authTestRouter(tokens).ServeHTTP(res, req)
	if res.Code != http.StatusNoContent {
		t.Fatalf("status seharusnya 204, didapat %d", res.Code)
	}
}

func superAdminTestRouter(isSuperAdmin bool) *gin.Engine {
    router := gin.New()
    router.GET("/admin", func(c *gin.Context) {
        c.Set(claimsKey, &utils.Claims{IsSuperAdmin: isSuperAdmin})
        c.Next()
    }, RequireSuperAdmin(), func(c *gin.Context) {
        c.Status(http.StatusNoContent)
    })
    return router
}

func TestRequireSuperAdminRejectsPartner(t *testing.T) {
    gin.SetMode(gin.TestMode)
    res := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin", nil)
    superAdminTestRouter(false).ServeHTTP(res, req)
    if res.Code != http.StatusForbidden {
        t.Fatalf("status seharusnya 403, didapat %d", res.Code)
    }
}

func TestRequireSuperAdminAllowsSuperAdmin(t *testing.T) {
    gin.SetMode(gin.TestMode)
    res := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/admin", nil)
    superAdminTestRouter(true).ServeHTTP(res, req)
    if res.Code != http.StatusNoContent {
        t.Fatalf("status seharusnya 204, didapat %d", res.Code)
    }
}