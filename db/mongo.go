package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client            *mongo.Client
	db                *mongo.Database
	targetCollection  *mongo.Collection
	backupsCollection *mongo.Collection
)

func Connect() error {
	uri := os.Getenv("MONGO_HOST")
	if uri == "" {
		return fmt.Errorf("MONGO_HOST env var not set")
	}

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	// load stuff
	db = client.Database("backups")

	targetCollection = db.Collection("targets")
	backupsCollection = db.Collection("backups")

	return nil
}

func Disconnect() error {
	return client.Disconnect(context.TODO())
}
