package audit_usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"strings"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

type Event struct {
	Type       string
	UserID     int32
	PartnerID  *int32
	ResourceID *int32
	Resource   string
	ActorName  string
	Email      string
	IP         string
	OldValue   any
	NewValue   any
	Detail     string
}

type Logger interface {
	Log(context.Context, Event)
}

type Store interface {
	CreateAuditLog(context.Context, sqlc.CreateAuditLogParams) error
	CountAuditLogs(context.Context, sqlc.CountAuditLogsParams) (int64, error)
	ListAuditLogs(context.Context, sqlc.ListAuditLogsParams) ([]sqlc.ListAuditLogsRow, error)
}

type Usecase struct {
	store Store
}

func New(store Store) *Usecase {
	return &Usecase{store: store}
}

func (u *Usecase) Log(ctx context.Context, e Event) {
	oldValue, err := auditJSON(e.OldValue)
	if err != nil {
		slog.Error("gagal membaca data audit lama", "error", err)
		return
	}
	newValue, err := auditJSON(e.NewValue)
	if err != nil {
		slog.Error("gagal membaca data audit baru", "error", err)
		return
	}
	if !newValue.Valid && e.Detail != "" {
		newValue, err = auditJSON(map[string]string{"detail": e.Detail})
		if err != nil {
			slog.Error("gagal membaca detail audit", "error", err)
			return
		}
	}
	if e.Resource == "" {
		e.Resource = resourceForAction(e.Type)
	}
	if err := u.store.CreateAuditLog(ctx, auditParams(e, oldValue, newValue)); err != nil {
		slog.Error("gagal menyimpan audit log", "action", e.Type, "error", err)
	}
}

func auditJSON(value any) (sql.NullString, error) {
	if value == nil {
		return sql.NullString{}, nil
	}
	raw, err := json.Marshal(value)
	return sql.NullString{String: string(raw), Valid: err == nil}, err
}

func auditParams(e Event, oldValue, newValue sql.NullString) sqlc.CreateAuditLogParams {
	return sqlc.CreateAuditLogParams{
		PartnerID: utils.NullInt32(e.PartnerID),
		UserID:    nullableID(e.UserID),
		Action:    e.Type,
		Resource:  e.Resource,
		ResourceID: utils.NullInt32(e.ResourceID),
		OldValue: oldValue, NewValue: newValue,
		IpAddress: utils.NullText(e.IP),
		ActorName: utils.NullText(e.ActorName),
		ActorEmail: utils.NullText(e.Email),
	}
}

func nullableID(id int32) sql.NullInt32 {
	return sql.NullInt32{Int32: id, Valid: id > 0}
}

func resourceForAction(action string) string {
	switch {
	case strings.HasPrefix(action, "LOGIN_"):
		return "auth"
	case strings.HasPrefix(action, "PASSWORD_"),
		strings.HasPrefix(action, "USER_"),
		strings.HasPrefix(action, "CUSTOMER_"):
		return "users"
	case strings.HasPrefix(action, "ROLE_"):
		return "roles"
	case strings.HasPrefix(action, "DIVISION_"):
		return "divisions"
	default:
		return "system"
	}
}