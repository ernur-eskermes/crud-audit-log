package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ernur-eskermes/crud-audit-log/internal/config"
	"github.com/ernur-eskermes/crud-audit-log/internal/repository"
	"github.com/ernur-eskermes/crud-audit-log/internal/service"
	"github.com/ernur-eskermes/crud-audit-log/internal/transport/grpc"
	"github.com/ernur-eskermes/crud-audit-log/pkg/database/mongodb"
	"github.com/ernur-eskermes/crud-audit-log/pkg/logging"
)

func main() {
	logger := logging.GetLogger()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal(err)
	}

	mongoClient, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		logger.Fatal(err)

		return
	}

	db := mongoClient.Database(cfg.Mongo.Database)

	auditRepo := repository.NewAudit(db)
	auditService := service.NewAudit(auditRepo)

	auditSrv := grpc.NewAuditServer(auditService)
	grpcSrv := grpc.New(auditSrv)

	go func() {
		logger.Info("Starting gRPC server")

		if err = grpcSrv.ListenAndServe(cfg.Server.Port); err != nil {
			logger.Error("gRPC ListenAndServer error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Info("Shutting down server")
	grpcSrv.Stop()

	if err = mongoClient.Disconnect(context.Background()); err != nil {
		logger.Error(err)
	}
}
