package queries

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MeService struct {
	db *mongo.Database
}

func NewMeService(db *mongo.Database) *MeService {
	return &MeService{db: db}
}
