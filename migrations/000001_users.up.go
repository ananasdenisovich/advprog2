package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Up is executed when the migration is applied
func Up_XXXXXXXXXX_create_users(ctx context.Context, client *mongo.Client) error {
	// Access the furnitureShopDB database
	database := client.Database("furnitureShopDB")

	// Create the users collection
	usersCollection := database.Collection("users")

	// Create indexes, set validations, or perform other necessary operations
	// ...

	// Example: Create an index on the "email" field
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}},
		// Other index options can be set as needed
	}

	_, err := usersCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create index on users collection: %w", err)
	}

	return nil
}
