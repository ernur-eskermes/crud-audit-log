package grpcHandler

import (
	"context"

	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"
)

type Audit struct {
	service AuditService
	audit.UnimplementedAuditServiceServer
}

func NewAudit(service AuditService) *Audit {
	return &Audit{
		service: service,
	}
}

func (h *Audit) Log(ctx context.Context, req *audit.LogRequest) (*audit.Empty, error) {
	err := h.service.Create(ctx, &audit.LogItem{
		Action:    req.GetAction().String(),
		Entity:    req.GetEntity().String(),
		EntityID:  req.GetEntityId(),
		Timestamp: req.GetTimestamp().AsTime(),
	})

	return &audit.Empty{}, err
}
