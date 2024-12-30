package database

import (
	"context"

	"github.com/vanthang24803/go-api/infra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func Connect(connection string) error {

	clientOptions := options.Client().ApplyURI(connection)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		return err
	}

	db = client.Database("db")

	infra.Msg.Info("Connected to MongoDB successfully!")

	return nil

}

func GetDatabase() *mongo.Database {
	return db
}
