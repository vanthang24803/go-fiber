package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	Name      string    `json:"name" bson:"name"`
	Value     string    `json:"value" bson:"value"`
	ExpiredAt time.Time `json:"expired_at" bson:"expired_at"`
}

type User struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email"`
	Avatar       string             `json:"avatar" bson:"avatar"`
	HashPassword string             `json:"-" bson:"hash_password"`
	FirstName    string             `json:"first_name" bson:"first_name"`
	LastName     string             `json:"last_name" bson:"last_name"`
	Username     string             `json:"username" bson:"username"`
	Roles        []string           `json:"roles" bson:"roles"`
	Tokens       []Token            `json:"-" bson:"tokens,omitempty"`
}
