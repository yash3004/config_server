package users

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitMongoDB initializes a MongoDB connection
func InitMongoDB(ctx context.Context, uri string) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database("config_server"), nil
}

// CloseMongoDB closes the MongoDB connection
func CloseMongoDB(ctx context.Context, client *mongo.Client) error {
	return client.Disconnect(ctx)
}