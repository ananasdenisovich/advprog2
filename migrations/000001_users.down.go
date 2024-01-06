package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// Down is executed when the migration is rolled back
func Down_XXXXXXXXXX_create_users(ctx context.Context, client *mongo.Client) error {
	// Access the furnitureShopDB database
	database := client.Database("furnitureShopDB")

	// Drop the users collection
	err := database.Collection("users").Drop(ctx)
	if err != nil {
		return fmt.Errorf("failed to drop users collection: %w", err)
	}

	return nil
}
