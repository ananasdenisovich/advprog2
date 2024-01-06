package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI       = "mongodb://localhost:27017"
	databaseName   = "furnitureShopDB"
	collectionName = "users"
)

var client *mongo.Client
var database *mongo.Database

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
	Age       int                `bson:"age,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Version   int                `bson:"version"`
}

var inventory = []Furniture{
	{ID: 1, Name: "Chair", Description: "Comfortable chair", Price: 49.99},
	{ID: 2, Name: "Table", Description: "Sturdy table", Price: 99.99},
	// Add more furniture items as needed
}

func init() {
	// Initialize MongoDB client
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Println("Error creating MongoDB client:", err)
		return
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}

	// Check the success of the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Error pinging MongoDB:", err)
		return
	}

	// If the connection is successful, print a success message
	fmt.Println("Connected to MongoDB successfully!")

	// Access the furnitureShopDB database
	database = client.Database(databaseName)
}

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

func createUsersCollection() error {
	usersCollection := database.Collection(collectionName)

	_, err := usersCollection.InsertOne(context.TODO(), User{
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	})

	return err
}

func addAgeField() error {
	usersCollection := database.Collection(collectionName)

	// Update all existing documents to set the Age field to a default value
	_, err := usersCollection.UpdateMany(
		context.TODO(),
		bson.D{},
		bson.M{"$set": bson.M{"age": 0}},
	)

	return err
}

func main() {
	// Run migrations
	if err := createUsersCollection(); err != nil {
		fmt.Println("Error creating users collection:", err)
		return
	}

	if err := addAgeField(); err != nil {
		fmt.Println("Error adding age field:", err)
		return
	}

	// Your other application logic...

	// Serve static files (HTML, CSS, JS) from the root path
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Define the routes and handlers
	http.HandleFunc("/getFurniture", handleGetFurniture)
	http.HandleFunc("/submitOrder", handlePostOrder)

	// Start the server on port 8080
	fmt.Println("Server is running on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
