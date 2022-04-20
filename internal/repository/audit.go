package repository

import (
	"context"

	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type Audit struct {
	db *mongo.Collection
}

func NewAudit(db *mongo.Database) *Audit {
	return &Audit{
		db: db.Collection("logs"),
	}
}

func (r *Audit) Insert(ctx context.Context, item audit.LogItem) error {
	_, err := r.db.InsertOne(ctx, item)

	return err
}
