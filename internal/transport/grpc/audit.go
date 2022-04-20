package grpc

import (
	"context"

	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"
)

type AuditService interface {
	Create(ctx context.Context, req *audit.LogRequest) error
}

type AuditServer struct {
	service AuditService
	audit.UnimplementedAuditServiceServer
}

func NewAuditServer(service AuditService) *AuditServer {
	return &AuditServer{
		service: service,
	}
}

func (h *AuditServer) Log(ctx context.Context, req *audit.LogRequest) (*audit.Empty, error) {
	err := h.service.Create(ctx, req)

	return &audit.Empty{}, err
}
