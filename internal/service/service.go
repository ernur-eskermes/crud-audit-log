package service

import (
	"context"

	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"
)

type Service struct {
	Audit *Audit
}

type AuditStorage interface {
	Insert(ctx context.Context, item *audit.LogItem) error
	GetAll(ctx context.Context) ([]audit.LogItem, error)
}

type Deps struct {
	AuditStorage AuditStorage
}

func New(deps Deps) *Service {
	return &Service{
		Audit: NewAudit(deps.AuditStorage),
	}
}
