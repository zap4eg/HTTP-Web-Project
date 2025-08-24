package config

import (
	"context"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Uri  string
	Name string
}

func SetupMongoDataBase(ctx context.Context, cancel context.CancelFunc) (*mongo.Database, error) {
	config := &MongoConfig{}

	err := viper.UnmarshalKey("mongo.database", config)
	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	defer cancel()

	return client.Database(config.Name), nil
}
