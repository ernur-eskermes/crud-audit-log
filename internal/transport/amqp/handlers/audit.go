package amqpHandler

import (
	"context"
	"encoding/json"
	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"
)

type Audit struct {
	service AuditService
}

func NewAudit(service AuditService) *Audit {
	return &Audit{
		service: service,
	}
}

func (h Audit) Log(ctx context.Context, message []byte) error {
	var inp audit.LogItem
	if err := json.Unmarshal(message, &inp); err != nil {
		return err
	}

	return h.service.Create(ctx, &inp)
}
