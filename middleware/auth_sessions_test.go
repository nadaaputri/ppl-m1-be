package middleware

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

type sessionStoreStub struct {
	session     sqlc.GetUserSessionRow
	permissions []string
	err         error
	revoked     bool
	revokedErr  error
}

func (s sessionStoreStub) IsTokenRevoked(context.Context, string) (bool, error) {
	return s.revoked, s.revokedErr
}

func sessionTestTokens() *utils.TokenMaker {
	return utils.NewTokenMaker(strings.Repeat("s", 32), "araya", 8*time.Hour)
}

func runSessionRequest(t *testing.T, store SessionStore) (int, *utils.Claims) {
	t.Helper()
	raw, err := sessionTestTokens().Generate(utils.Claims{UserID: 7})
	if err != nil {
		t.Fatal(err)
	}
	return serveSessionRequest(store, raw)
}

func serveSessionRequest(store SessionStore, raw string) (int, *utils.Claims) {
	var claims *utils.Claims
	router := gin.New()
	router.GET("/", AuthenticateCurrent(sessionTestTokens(), store), func(c *gin.Context) {
		claims = GetClaims(c)
		c.Status(http.StatusNoContent)
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+raw)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res.Code, claims
}

func (s sessionStoreStub) GetUserSession(
	context.Context, int32,
) (sqlc.GetUserSessionRow, error) {
	return s.session, s.err
}

func (s sessionStoreStub) GetUserPermissionCodes(
	context.Context, int32,
) ([]string, error) {
	return s.permissions, s.err
}

func TestCurrentSessionRejectsRevokedToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := activeSessionStore()
	store.revoked = true
	status, _ := runSessionRequest(t, store)
	if status != http.StatusUnauthorized {
		t.Fatalf("status seharusnya 401, didapat %d", status)
	}
}

func TestCurrentSessionRevocationCheckFailsClosed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := activeSessionStore()
	store.revokedErr = errors.New("database down")
	status, _ := runSessionRequest(t, store)
	if status != http.StatusServiceUnavailable {
		t.Fatalf("status seharusnya 503, didapat %d", status)
	}
}

func TestCurrentSessionRejectsTokenWithoutID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	claims := utils.Claims{UserID: 7, RegisteredClaims: jwt.RegisteredClaims{
		Issuer: "araya", ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}}
	raw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(strings.Repeat("s", 32)))
	if err != nil {
		t.Fatal(err)
	}
	if status, _ := serveSessionRequest(activeSessionStore(), raw); status != http.StatusUnauthorized {
		t.Fatalf("status seharusnya 401, didapat %d", status)
	}
}

func TestCurrentSessionRejectsInactiveAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := sessionStoreStub{
		session: sqlc.GetUserSessionRow{ID: 7, IsActive: false},
	}
	status, _ := runSessionRequest(t, store)
	if status != http.StatusUnauthorized {
		t.Fatalf("status seharusnya 401, didapat %d", status)
	}
}

func TestCurrentSessionRefreshesClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := activeSessionStore()
	status, claims := runSessionRequest(t, store)
	if status != http.StatusNoContent {
		t.Fatalf("status seharusnya 204, didapat %d", status)
	}
	if claims.Division == nil || *claims.Division != "finance" {
		t.Fatal("division belum diperbarui")
	}
	if len(claims.Permissions) != 1 || claims.Permissions[0] != "m1.audit.read" {
		t.Fatal("permission belum diperbarui")
	}
}

func activeSessionStore() sessionStoreStub {
	return sessionStoreStub{
		session: sqlc.GetUserSessionRow{
			ID:             7,
			Name:           "Nama Baru",
			IsActive:       true,
			RoleID:         sql.NullInt32{Int32: 2, Valid: true},
			RoleCode:       sql.NullString{String: "staff", Valid: true},
			RoleActive:     true,
			DivisionID:     sql.NullInt32{Int32: 3, Valid: true},
			DivisionCode:   sql.NullString{String: "finance", Valid: true},
			DivisionActive: true,
		},
		permissions: []string{"m1.audit.read"},
	}
}

func TestCurrentSessionRejectsInactiveDivision(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := activeSessionStore()
	store.session.DivisionActive = false
	status, _ := runSessionRequest(t, store)
	if status != http.StatusForbidden {
		t.Fatalf("status seharusnya 403, didapat %d", status)
	}
}

func TestCurrentSessionFailsClosed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := sessionStoreStub{err: errors.New("database down")}
	status, _ := runSessionRequest(t, store)
	if status != http.StatusServiceUnavailable {
		t.Fatalf("status seharusnya 503, didapat %d", status)
	}
}
