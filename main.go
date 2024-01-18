package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	mongoURI       = "mongodb://localhost:27017"
	databaseName   = "furnitureShopDB"
	collectionName = "users"
)

var client *mongo.Client
var database *mongo.Database

type Furniture struct {
	ID          int     `json:"id" bson:"_id,omitempty"`
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`
	Price       float64 `json:"price" bson:"price"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Version   int                `bson:"version"`
}

var inventory = []Furniture{
	{ID: 1, Name: "Chair", Description: "Comfortable chair", Price: 49.99},
	{ID: 2, Name: "Table", Description: "Sturdy table", Price: 99.99},
}

func init() {

	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(mongoURI))
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

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Error pinging MongoDB:", err)
		return
	}

	fmt.Println("Connected to MongoDB successfully!")

	database = client.Database(databaseName)
}

func registerUser(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the user's password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = string(hashedPassword)

	// Set creation and update times
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Insert the user into the database
	usersCollection := client.Database(databaseName).Collection(collectionName)
	result, err := usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "userID": result.InsertedID})
}

func loginUser(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user by email
	usersCollection := client.Database(databaseName).Collection(collectionName)
	var user User
	err := usersCollection.FindOne(context.TODO(), bson.M{"email": loginRequest.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare the stored hashed password with the input password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// TODO: Implement token generation and response

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func handleGetFurniture(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(inventory)
}

func handlePostOrder(w http.ResponseWriter, r *http.Request) {
	var order map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		response := map[string]string{"status": "400", "message": "Invalid JSON-message"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Printf("Received order data: %+v\n", order)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "200", "message": "Order received successfully"}
	json.NewEncoder(w).Encode(response)
}

func handleHTML(w http.ResponseWriter, r *http.Request) {
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

	_, err := usersCollection.UpdateMany(
		context.TODO(),
		bson.D{},
		bson.M{"$set": bson.M{"age": 0}},
	)

	return err
}

func main() {
	r := gin.Default()

	r.POST("/register", registerUser)
	r.POST("/login", loginUser)

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

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Error pinging MongoDB:", err)
		return
	}

	fmt.Println("Connected to MongoDB successfully!")

	defer client.Disconnect(ctx)

	database := client.Database(databaseName)

	if err := createUsersCollection(); err != nil {
		fmt.Println("Error creating users collection:", err)
		return
	}

	if err := addAgeField(); err != nil {
		fmt.Println("Error adding age field:", err)
		return
	}
	exampleUser := User{
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	usersCollection := database.Collection(collectionName)
	insertResult, err := usersCollection.InsertOne(ctx, exampleUser)
	if err != nil {
		fmt.Println("Error inserting user:", err)
		return
	}

	fmt.Println("Inserted user with ID:", insertResult.InsertedID)

	http.Handle("/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/getFurniture", handleGetFurniture)
	http.HandleFunc("/submitOrder", handlePostOrder)

	// routes and handlers for CRUD operations
	http.HandleFunc("/createUser", createUser)
	http.HandleFunc("/getUser", getUserByID)
	http.HandleFunc("/updateUser", updateUser)
	http.HandleFunc("/deleteUser", deleteUser)
	http.HandleFunc("/getAllUsers", getAllUsers)

	fmt.Println("Server is running on :8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

// CRUD
func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = newUser.CreatedAt

	usersCollection := database.Collection(collectionName)
	insertResult, err := usersCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(insertResult)
}
func getUserByID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	objID, _ := primitive.ObjectIDFromHex(userID)

	var user User
	usersCollection := database.Collection(collectionName)
	err := usersCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
func updateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	objID, _ := primitive.ObjectIDFromHex(userID)

	var updateData struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersCollection := database.Collection(collectionName)
	_, err = usersCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"name": updateData.Name, "updated_at": time.Now()}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	objID, _ := primitive.ObjectIDFromHex(userID)

	usersCollection := database.Collection(collectionName)
	_, err := usersCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	usersCollection := database.Collection(collectionName)
	cursor, err := usersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}
