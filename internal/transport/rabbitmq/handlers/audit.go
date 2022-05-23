package amqpHandler

import (
	"context"
	"encoding/json"

	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"
	"github.com/ernur-eskermes/crud-audit-log/pkg/logging"
	"github.com/streadway/amqp"
)

type Audit struct {
	service AuditService
	logger  *logging.Logger
}

func NewAudit(service AuditService, logger *logging.Logger) *Audit {
	return &Audit{
		service: service,
		logger:  logger,
	}
}

func (h *Audit) Log(ctx context.Context, d amqp.Delivery) {
	l := &audit.LogItem{}

	if err := json.Unmarshal(d.Body, &l); err != nil {
		h.logger.Errorf("%s: %s", "Failed to unmarshel LogItem", err)

		return
	}

	if err := h.service.Create(ctx, l); err != nil {
		h.logger.Errorf("%s: %s", "Failed to create Log", err)

		return
	}

	if err := d.Ack(false); err != nil {
		h.logger.Errorf("%s: %s", "Failed to ack", err)

		return
	}
}
