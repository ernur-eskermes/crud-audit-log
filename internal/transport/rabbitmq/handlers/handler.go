package amqpHandler

import (
	"context"

	"github.com/ernur-eskermes/crud-audit-log/pkg/logging"

	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"
)

type Handler struct {
	Audit *Audit
}

type AuditService interface {
	Create(ctx context.Context, item *audit.LogItem) error
}

type Deps struct {
	AuditService AuditService
	Logger       *logging.Logger
}

func New(deps Deps) *Handler {
	return &Handler{
		Audit: NewAudit(deps.AuditService, deps.Logger),
	}
}
