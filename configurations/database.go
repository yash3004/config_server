package configurations

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB(ctx context.Context, uri string) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database("config_server"), nil
}

func CloseMongoDB(ctx context.Context, client *mongo.Client) error {
	return client.Disconnect(ctx)
}

type ConfigFile struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	Filename  string    `bson:"filename"`
	FileType  int       `bson:"file_type"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}