package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func generateRandomDoc() ([]byte, string, error) {
	doc := map[string]interface{}{
		"id":    uuid.NewString(),
		"name":  randomString(10),
		"email": randomString(10) + "@example.com",
		"age":   rand.Intn(50) + 18,
	}
	var records []map[string]any
	records = append(records, doc)

	//     docBytes, err := json.Marshal(doc)
	docBytes, err := json.Marshal(records)
	if err != nil {
		return nil, "", err
	}
	return docBytes, doc["id"].(string), nil
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// getAllIDs retrieves all document IDs from a specified MongoDB collection.
func getAllIDs(collectionName, attr string) ([]string, error) {
	var ids []string
	collection := mongoDB.Collection(collectionName)

	// Find all documents, but only retrieve the _id field.
	cursor, err := collection.Find(context.Background(), bson.M{}, options.Find().SetProjection(bson.M{attr: 1}))
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor and extract the _id field from each document.
	for cursor.Next(context.Background()) {
		var result struct {
			ID interface{} `bson:"id"`
		}
		if err := cursor.Decode(&result); err != nil {
			log.Printf("failed to decode document: %v", err)
			continue
		}

		// Convert ObjectID or other ID type to string
		idStr := fmt.Sprintf("%v", result.ID)
		ids = append(ids, idStr)
	}

	// Check if any error occurred during the iteration
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return ids, nil
}
