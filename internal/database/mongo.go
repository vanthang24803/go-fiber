package database

import (
	"context"

	"github.com/vanthang24803/go-api/infra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func ConnectDB(connection string) (db *mongo.Database, err error) {
	clientOptions := options.Client().ApplyURI(connection)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		infra.Msg.Error("Error while connecting to MongoDB:", err)
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		infra.Msg.Error("MongoDB ping failed:", err)
		return nil, err
	}

	db = client.Database("db")
	if db == nil {
		infra.Msg.Error("Database instance is nil after connection")
	}

	infra.Msg.Info("Connected to MongoDB successfully!")
	return db, nil
}

func GetDatabase() *mongo.Database {
	return db
}
