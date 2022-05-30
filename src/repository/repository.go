package repository

import (
	"context"
	"file_processor/src/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx context.Context
var client *mongo.Client

func CollectionMongo() *mongo.Collection {
	client = config.ConnectMongo(ctx)
	return client.Database("file").Collection("fields")
}

func Find(filter bson.D, collection *mongo.Collection) []map[string]string {
	cursor, err := collection.Find(ctx, filter)
	config.HandlerError(err, "Can't get documents")

	var response []map[string]string
	err = cursor.All(ctx, &response)
	config.HandlerError(err, "Can't convert documents")

	return response
}

func InsertOne(doc interface{}, collection *mongo.Collection) bool {
	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		config.HandlerError(err, "Can't insert docs")
		return false
	}
	return true
}

func CloseCollection() {
	config.CloseMongo(client)
}
