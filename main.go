package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Furniture struct to represent the JSON data for a furniture item
type Furniture struct {
	ID          int     `json:"id" bson:"_id,omitempty"`
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`
	Price       float64 `json:"price" bson:"price"`
}

// User represents a user document in the MongoDB users collection
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

var inventory = []Furniture{
	{1, "Comfortable Sofa", "A stylish and comfortable sofa for your living room.", 499.99},
	{2, "Elegant Dining Table", "A beautiful dining table for your family gatherings.", 299.99},
	{3, "Modern Coffee Table", "A sleek and modern coffee table for your lounge.", 129.99},
}

// MongoDB configuration
const (
	mongoURI       = "mongodb://localhost:27017"
	databaseName   = "furnitureShopDB"
	collectionName = "users"
)

func handleGetFurniture(w http.ResponseWriter, r *http.Request) {
	// Return the list of furniture items as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(inventory)
}

func handlePostOrder(w http.ResponseWriter, r *http.Request) {
	// Parse JSON data from the request body
	var order map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		// Return a JSON response in case of decoding error
		response := map[string]string{"status": "400", "message": "Invalid JSON-message"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Do something with the received order data
	fmt.Printf("Received order data: %+v\n", order)

	// Send a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "200", "message": "Order received successfully"}
	json.NewEncoder(w).Encode(response)
}

func handleHTML(w http.ResponseWriter, r *http.Request) {
	// Serve the HTML file
	http.ServeFile(w, r, "index.html")
}

func main() {

	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Println("Error creating MongoDB client:", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	defer client.Disconnect(ctx)

	// Access the furnitureShopDB database
	database := client.Database(databaseName)

	// Create an example user
	exampleUser := User{
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert the user into the users collection
	usersCollection := database.Collection(collectionName)
	insertResult, err := usersCollection.InsertOne(ctx, exampleUser)
	if err != nil {
		fmt.Println("Error inserting user:", err)
		return
	}

	fmt.Println("Inserted user with ID:", insertResult.InsertedID)

	http.Handle("/", http.FileServer(http.Dir(".")))

	// Define the routes and handlers
	http.HandleFunc("/getFurniture", handleGetFurniture)
	http.HandleFunc("/submitOrder", handlePostOrder)

	// Start the server on port 8080
	fmt.Println("Server is running on :8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
