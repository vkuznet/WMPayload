// mongo.go
package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoClient creates a new MongoDB client.
func NewMongoClient(uri string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	return client, err
}

// InsertDocuments inserts multiple documents into a MongoDB collection.
func InsertDocuments(client *mongo.Client, dbName, collName string, docs []map[string]interface{}) error {
	collection := client.Database(dbName).Collection(collName)
	var interfaceDocs []interface{}
	for _, doc := range docs {
		interfaceDocs = append(interfaceDocs, doc)
	}
	_, err := collection.InsertMany(context.Background(), interfaceDocs)
	return err
}

// SearchDocuments retrieves documents with pagination and filtering.
func SearchDocuments(client *mongo.Client, dbName, collName string, filter map[string]interface{}, skip, limit int64) ([]map[string]interface{}, error) {
	collection := client.Database(dbName).Collection(collName)
	findOptions := options.Find().SetSkip(skip).SetLimit(limit)
	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []map[string]interface{}
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
