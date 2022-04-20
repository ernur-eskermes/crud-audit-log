package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc/reflection"

	grpcHandler "github.com/ernur-eskermes/crud-audit-log/internal/transport/grpc/handlers"

	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"
	"google.golang.org/grpc"
)

type Server struct {
	grpcSrv      *grpc.Server
	auditHandler audit.AuditServiceServer
}

func New(handlers *grpcHandler.Handler) *Server {
	return &Server{
		grpcSrv:      grpc.NewServer(),
		auditHandler: handlers.Audit,
	}
}

func (s *Server) ListenAndServe(port int) error {
	addr := fmt.Sprintf(":%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	audit.RegisterAuditServiceServer(s.grpcSrv, s.auditHandler)
	reflection.Register(s.grpcSrv)

	return s.grpcSrv.Serve(lis)
}

func (s *Server) Stop() {
	s.grpcSrv.GracefulStop()
}
