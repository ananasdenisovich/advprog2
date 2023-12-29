package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Furniture struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

var inventory = []Furniture{
	{1, "Comfortable Sofa", "A stylish and comfortable sofa for your living room.", 499.99},
	{2, "Elegant Dining Table", "A beautiful dining table for your family gatherings.", 299.99},
	{3, "Modern Coffee Table", "A sleek and modern coffee table for your lounge.", 129.99},
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
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received order data: %+v\n", order)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Order received successfully"}
	json.NewEncoder(w).Encode(response)
}

func main() {

	http.Handle("/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/getFurniture", handleGetFurniture)
	http.HandleFunc("/submitOrder", handlePostOrder)

	fmt.Println("Server is running on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
