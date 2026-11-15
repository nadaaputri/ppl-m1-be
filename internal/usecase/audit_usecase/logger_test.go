package audit_usecase

import (
	"context"
	"strings"
	"testing"

	"ppl-m1-be/internal/usecase/repository/sqlc"
)

type auditStoreStub struct {
	created sqlc.CreateAuditLogParams
}

func (s *auditStoreStub) CreateAuditLog(
	_ context.Context, arg sqlc.CreateAuditLogParams,
) error {
	s.created = arg
	return nil
}

func (*auditStoreStub) CountAuditLogs(
	context.Context, sqlc.CountAuditLogsParams,
) (int64, error) {
	return 0, nil
}

func (*auditStoreStub) ListAuditLogs(
	context.Context, sqlc.ListAuditLogsParams,
) ([]sqlc.ListAuditLogsRow, error) {
	return nil, nil
}

func TestLogPersistsAuditFields(t *testing.T) {
	store := &auditStoreStub{}
	uc := New(store)
	partnerID, resourceID := int32(5), int32(12)
	uc.Log(context.Background(), Event{
		Type: "USER_CREATED", UserID: 7, PartnerID: &partnerID,
		Resource: "users", ResourceID: &resourceID,
		ActorName: "Administrator", IP: "127.0.0.1",
		NewValue: map[string]string{"name": "Rani"},
	})
	if store.created.Action != "USER_CREATED" ||
		store.created.ResourceID.Int32 != resourceID {
		t.Fatalf("data audit tidak sesuai: %+v", store.created)
	}
	if !store.created.IpAddress.Valid || store.created.ActorName.String != "Administrator" {
		t.Fatal("identitas pelaku atau IP tidak tercatat")
	}
	if strings.Contains(store.created.NewValue.String, "password") {
		t.Fatal("password tidak boleh masuk audit")
	}
}