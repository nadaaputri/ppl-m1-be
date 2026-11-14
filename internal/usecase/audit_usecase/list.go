package audit_usecase

import (
	"context"
	"encoding/json"
	"time"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

type ListInput struct {
	PartnerID *int32
	Page      int
	PerPage   int
}

type Entry struct {
	ID         int32           `json:"id"`
	PartnerID  *int32          `json:"partner_id"`
	UserID     *int32          `json:"user_id"`
	Action     string          `json:"action"`
	Resource   string          `json:"resource"`
	ResourceID *int32          `json:"resource_id"`
	OldValue   json.RawMessage `json:"old_value"`
	NewValue   json.RawMessage `json:"new_value"`
	IPAddress  *string         `json:"ip_address"`
	ActorName  *string         `json:"actor_name"`
	ActorEmail *string         `json:"actor_email"`
	CreatedAt  time.Time       `json:"created_at"`
}

type ListResult struct {
	Items []Entry
	Total int
}

func (u *Usecase) List(ctx context.Context, in ListInput) (ListResult, error) {
	includeAll := in.PartnerID == nil
	partnerID := utils.NullInt32(in.PartnerID)
	total, err := u.store.CountAuditLogs(ctx, sqlc.CountAuditLogsParams{
		IncludeAll: includeAll, PartnerID: partnerID,
	})
	if err != nil {
		return ListResult{}, err
	}
	rows, err := u.store.ListAuditLogs(ctx, sqlc.ListAuditLogsParams{
		IncludeAll: includeAll, PartnerID: partnerID,
		PageLimit: int32(in.PerPage), PageOffset: int32((in.Page - 1) * in.PerPage),
	})
	if err != nil {
		return ListResult{}, err
	}
	return ListResult{Items: toEntries(rows), Total: int(total)}, nil
}

func toEntries(rows []sqlc.ListAuditLogsRow) []Entry {
	items := make([]Entry, 0, len(rows))
	for _, row := range rows {
		items = append(items, toEntry(row))
	}
	return items
}

func toEntry(r sqlc.ListAuditLogsRow) Entry {
	return Entry{
		ID: r.ID, PartnerID: utils.Int32Ptr(r.PartnerID),
		UserID: utils.Int32Ptr(r.UserID), Action: r.Action,
		Resource: r.Resource, ResourceID: utils.Int32Ptr(r.ResourceID),
		OldValue: rawJSON(r.OldValue), NewValue: rawJSON(r.NewValue),
		IPAddress:  utils.StringPtr(r.IpAddress),
		ActorName:  utils.StringPtr(r.ActorName),
		ActorEmail: utils.StringPtr(r.ActorEmail),
		CreatedAt:  r.CreatedAt.In(utils.WIB),
	}
}

func rawJSON(value string) json.RawMessage {
	if value == "" {
		return nil
	}
	return json.RawMessage(value)
}
