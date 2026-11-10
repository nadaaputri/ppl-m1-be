package audit_usecase

import (
	"context"
	"log/slog"
)

type Event struct {
	Type   string
	UserID int32
	Email  string
	IP     string
	Detail string
}

type Logger interface {
	Log(ctx context.Context, e Event)
}

type stdLogger struct{}

func NewStdLogger() Logger {
	return stdLogger{}
}

func (stdLogger) Log(ctx context.Context, e Event) {
	slog.InfoContext(ctx, "audit",
		"type", e.Type, "user_id", e.UserID,
		"email", e.Email, "ip", e.IP, "detail", e.Detail,
	)
}