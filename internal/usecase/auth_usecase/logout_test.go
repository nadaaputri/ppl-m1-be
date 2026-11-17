package auth_usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"ppl-m1-be/internal/usecase/audit_usecase"
	"ppl-m1-be/internal/usecase/utils"
)

type auditRecorder struct{ events []audit_usecase.Event }

func (r *auditRecorder) Log(_ context.Context, e audit_usecase.Event) {
	r.events = append(r.events, e)
}

func newLogoutUsecase(store *loginStoreStub, audit *auditRecorder) *Usecase {
	tokens := utils.NewTokenMaker(strings.Repeat("s", 32), "araya", 8*time.Hour)
	return New(store, tokens, audit)
}

func logoutInput() LogoutInput {
	return LogoutInput{
		UserID: 7, Name: "Rani", TokenID: "jti-1",
		ExpiresAt: time.Now().Add(time.Hour), IP: "127.0.0.1",
	}
}

func TestLogoutRevokesTokenAndWritesAudit(t *testing.T) {
	store, audit := &loginStoreStub{}, &auditRecorder{}
	if err := newLogoutUsecase(store, audit).Logout(context.Background(), logoutInput()); err != nil {
		t.Fatal(err)
	}
	if store.revoked.Jti != "jti-1" || store.revoked.UserID != 7 {
		t.Fatalf("token tidak dicabut: %+v", store.revoked)
	}
	if len(audit.events) != 1 || audit.events[0].Type != "LOGOUT" || audit.events[0].IP != "127.0.0.1" {
		t.Fatalf("audit logout tidak sesuai: %+v", audit.events)
	}
}

func TestLogoutSkipsAuditWhenRevokeFails(t *testing.T) {
	store := &loginStoreStub{revokeErr: errors.New("database down")}
	audit := &auditRecorder{}
	if err := newLogoutUsecase(store, audit).Logout(context.Background(), logoutInput()); err == nil {
		t.Fatal("galat database seharusnya dikembalikan")
	}
	if len(audit.events) != 0 {
		t.Fatal("audit tidak boleh tercatat bila pencabutan gagal")
	}
}
