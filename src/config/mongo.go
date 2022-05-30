package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx context.Context

func ConnectMongo(ctx context.Context) *mongo.Client {
	credentials := options.Credential{Username: "root", Password: "root"}
	opt := options.Client().ApplyURI("mongodb://127.0.0.1:27017").SetAuth(credentials)

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Minute)

	client, err := mongo.Connect(ctx, opt)
	HandlerError(err, "Can't connect to Mongo")

	err = client.Ping(ctx, nil)
	HandlerError(err, "Can't ping Mongo client")

	return client
}

func CloseMongo(client *mongo.Client) {
	client.Disconnect(ctx)
}
