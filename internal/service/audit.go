package service

import (
	"context"

	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"
)

type Repository interface {
	Insert(ctx context.Context, item audit.LogItem) error
}

type Audit struct {
	repo Repository
}

func NewAudit(repo Repository) *Audit {
	return &Audit{
		repo: repo,
	}
}

func (s *Audit) Create(ctx context.Context, req *audit.LogRequest) error {
	return s.repo.Insert(ctx, audit.LogItem{
		Action:    req.GetAction().String(),
		Entity:    req.GetEntity().String(),
		EntityID:  req.GetEntityId(),
		Timestamp: req.GetTimestamp().AsTime(),
	})
}
