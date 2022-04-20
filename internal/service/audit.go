package service

import (
	"context"

	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"
)

type Audit struct {
	repo AuditStorage
}

func NewAudit(repo AuditStorage) *Audit {
	return &Audit{
		repo: repo,
	}
}

func (s *Audit) Create(ctx context.Context, item *audit.LogItem) error {
	return s.repo.Insert(ctx, item)
}

func (s *Audit) GetAll(ctx context.Context) ([]audit.LogItem, error) {
	return s.repo.GetAll(ctx)
}
