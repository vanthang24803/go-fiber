package queries

import (
	"context"

	"github.com/vanthang24803/go-api/internal/database"
	"github.com/vanthang24803/go-api/internal/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var db *mongo.Database

func init() {
	// Initialize the MongoDB database instance
	db = database.GetDatabase() // assuming this returns *mongo.Database
}

// Register a new user in the database
func Register(user *schema.User) (*schema.User, error) {
	// Hash the password using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(user.HashPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Set the hashed password
	user.HashPassword = string(hash)

	// Get the "users" collection
	usersCollection := db.Collection("users")

	// Insert the user into the collection
	result, err := usersCollection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	// Set the user ID from the inserted result
	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, nil
}
