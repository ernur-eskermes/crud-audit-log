package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ernur-eskermes/crud-audit-log/docs"
	grpcHandler "github.com/ernur-eskermes/crud-audit-log/internal/transport/grpc/handlers"
	"github.com/ernur-eskermes/crud-audit-log/internal/transport/rest"
	restHandler "github.com/ernur-eskermes/crud-audit-log/internal/transport/rest/handlers"

	"github.com/ernur-eskermes/crud-audit-log/internal/config"
	"github.com/ernur-eskermes/crud-audit-log/internal/service"
	"github.com/ernur-eskermes/crud-audit-log/internal/storage"
	"github.com/ernur-eskermes/crud-audit-log/internal/transport/grpc"
	"github.com/ernur-eskermes/crud-audit-log/pkg/database/mongodb"
	"github.com/ernur-eskermes/crud-audit-log/pkg/logging"
)

// @title Audit Service
// @version 1.0
// @description REST API for audit log service

// @host localhost:9001
// @BasePath /api/

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

// Run initializes whole application.
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

	storages := storage.New(db)
	services := service.New(service.Deps{
		AuditStorage: storages.Audit,
	})

	grpcHandlers := grpcHandler.New(grpcHandler.Deps{
		AuditService: services.Audit,
	})
	grpcSrv := grpc.New(grpcHandlers)

	restHandlers := restHandler.New(restHandler.Deps{
		AuditService: services.Audit,
		Logger:       logger,
	})
	restSrv := rest.NewServer(cfg, restHandlers)

	go func() {
		logger.Info("Starting gRPC server")

		if err = grpcSrv.ListenAndServe(cfg.GRPC.Port); err != nil {
			logger.Error("gRPC ListenAndServer error", err)
		}
	}()

	go func() {
		logger.Info("Starting HTTP server")

		if err = restSrv.ListenAndServe(cfg.HTTP.Port); err != nil {
			logger.Error("HTTP ListenAndServer error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Info("Shutting down server")
	grpcSrv.Stop()

	if err = restSrv.Stop(); err != nil {
		logger.Errorf("failed to stop http server: %v", err)
	}

	if err = mongoClient.Disconnect(context.Background()); err != nil {
		logger.Error(err)
	}
}
