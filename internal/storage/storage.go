package storage

import "go.mongodb.org/mongo-driver/mongo"

type Storage struct {
	Audit *Audit
}

func New(db *mongo.Database) *Storage {
	return &Storage{
		Audit: NewAudit(db.Collection("storage")),
	}
}
