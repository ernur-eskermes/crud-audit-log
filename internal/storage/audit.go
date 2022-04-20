package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type Audit struct {
	db *mongo.Collection
}

func NewAudit(db *mongo.Collection) *Audit {
	return &Audit{
		db: db,
	}
}

func (r *Audit) Insert(ctx context.Context, item *audit.LogItem) error {
	_, err := r.db.InsertOne(ctx, item)

	return err
}

func (r *Audit) GetAll(ctx context.Context) ([]audit.LogItem, error) {
	var items []audit.LogItem

	cur, err := r.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cur.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}
