package auth_usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"ppl-m1-be/internal/usecase/audit_usecase"
	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

type loginStoreStub struct {
	row   sqlc.GetUserByEmailRow
	err   error
	email string
}

func (s *loginStoreStub) GetUserByEmail(_ context.Context, email string) (sqlc.GetUserByEmailRow, error) {
	s.email = email
	return s.row, s.err
}

func (*loginStoreStub) GetAccountCredentials(context.Context, int32) (sqlc.GetAccountCredentialsRow, error) {
	return sqlc.GetAccountCredentialsRow{}, nil
}

func (*loginStoreStub) UpdatePassword(context.Context, sqlc.UpdatePasswordParams) (time.Time, error) {
	return time.Time{}, nil
}

type loginAuditStub struct{}

func (loginAuditStub) Log(context.Context, audit_usecase.Event) {}

func loginTestRow() sqlc.GetUserByEmailRow {
	return sqlc.GetUserByEmailRow{
		ID: 7, Email: "user@example.com", Name: "Rani", IsActive: true,
		PartnerID:    sql.NullInt32{Int32: 7, Valid: true},
		UnitUsahaID:  sql.NullInt32{Int32: 3, Valid: true},
		RoleCode:     sql.NullString{String: "staff", Valid: true},
		DivisionCode: sql.NullString{String: "sales", Valid: true},
	}
}

func newLoginTestUsecase(t *testing.T, row sqlc.GetUserByEmailRow) (*Usecase, *loginStoreStub) {
	t.Helper()
	hash, err := utils.HashPassword("Password123")
	if err != nil {
		t.Fatal(err)
	}
	row.PasswordHash = hash
	store := &loginStoreStub{row: row}
	tokens := utils.NewTokenMaker(strings.Repeat("s", 32), "araya", 8*time.Hour)
	return New(store, tokens, loginAuditStub{}), store
}

func parseLoginClaims(t *testing.T, uc *Usecase, result LoginResult) *utils.Claims {
	t.Helper()
	claims, err := uc.tokens.Parse(result.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	return claims
}

func TestLoginSuccess(t *testing.T) {
	uc, store := newLoginTestUsecase(t, loginTestRow())
	result, err := uc.Login(context.Background(), LoginInput{
		Email: " USER@EXAMPLE.COM ", Password: "Password123",
	})
	if err != nil {
		t.Fatal(err)
	}
	if store.email != "user@example.com" || result.ExpiresIn != 28800 {
		t.Fatalf("email atau expiration tidak sesuai: %+v", result)
	}
	claims := parseLoginClaims(t, uc, result)
	if claims.UserID != 7 || claims.Role != "staff" || claims.IsSuperAdmin {
		t.Fatalf("claims tidak sesuai: %+v", claims)
	}
	if claims.PartnerID == nil || *claims.PartnerID != 7 {
		t.Fatalf("partner_id tidak sesuai: %+v", claims.PartnerID)
	}
}

func TestLoginIncludesSuperAdminClaim(t *testing.T) {
	row := loginTestRow()
	row.IsSuperAdmin = true
	row.RoleCode = sql.NullString{String: "administrator", Valid: true}
	row.PartnerID = sql.NullInt32{}
	uc, _ := newLoginTestUsecase(t, row)
	result, err := uc.Login(context.Background(), LoginInput{
		Email: "user@example.com", Password: "Password123",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !parseLoginClaims(t, uc, result).IsSuperAdmin {
		t.Fatal("claim is_super_admin seharusnya true")
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	uc, _ := newLoginTestUsecase(t, loginTestRow())
	_, err := uc.Login(context.Background(), LoginInput{
		Email: "user@example.com", Password: "WrongPassword123",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("seharusnya ditolak: %v", err)
	}
}

func TestLoginRejectsInactiveAccount(t *testing.T) {
	row := loginTestRow()
	row.IsActive = false
	uc, _ := newLoginTestUsecase(t, row)
	_, err := uc.Login(context.Background(), LoginInput{
		Email: "user@example.com", Password: "Password123",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("seharusnya ditolak: %v", err)
	}
}

func TestLoginRejectsUnknownEmail(t *testing.T) {
	uc, _ := newLoginTestUsecase(t, sqlc.GetUserByEmailRow{})
	_, err := uc.Login(context.Background(), LoginInput{
		Email: "unknown@example.com", Password: "Password123",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("seharusnya ditolak: %v", err)
	}
}

func TestLoginIncludesContractClaims(t *testing.T) {
	uc, _ := newLoginTestUsecase(t, loginTestRow())
	result, err := uc.Login(context.Background(), LoginInput{
		Email: "user@example.com", Password: "Password123",
	})
	if err != nil {
		t.Fatal(err)
	}
	claims := parseLoginClaims(t, uc, result)
	if claims.Name != "Rani" || claims.Division == nil || *claims.Division != "sales" {
		t.Fatalf("identity claims tidak sesuai: %+v", claims)
	}
	if claims.UnitUsahaID == nil || *claims.UnitUsahaID != 3 {
		t.Fatalf("unit_usaha_id tidak sesuai: %+v", claims.UnitUsahaID)
	}
	if claims.Issuer != "araya" || claims.IssuedAt == nil || claims.ExpiresAt == nil {
		t.Fatalf("registered claims tidak lengkap: %+v", claims)
	}
	if claims.ExpiresAt.Sub(claims.IssuedAt.Time) != 8*time.Hour {
		t.Fatal("masa berlaku token bukan 8 jam")
	}
}
