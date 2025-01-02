package functions

import (
	"context"

	"github.com/vanthang24803/go-api/internal/schema"
	"github.com/vanthang24803/go-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MeService struct {
	db *mongo.Database
}

func NewMeService(db *mongo.Database) *MeService {
	return &MeService{db: db}
}

func (s *MeService) GetProfileHandler(payload *utils.JwtPayload) (*schema.User, *utils.AppError) {

	var user schema.User

	err := s.db.Collection("users").FindOne(context.Background(), bson.M{"_id": payload.Sub}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.NewAppError(404, "User not found")
		}
		return nil, utils.NewAppError(500, err.Error())
	}

	return &user, nil
}
