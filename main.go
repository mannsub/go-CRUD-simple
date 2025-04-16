package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type Item struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

var (
	items  = make(map[int]Item)
	mu     sync.Mutex
	nextID = 1
)

func main() {
	http.HandleFunc("/items", handleItems)
	http.HandleFunc("/items/", handleItem)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func handleItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getAllItems(w, r)
	} else if r.Method == http.MethodPost {
		createItem(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/items/"):])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		getItem(w, r, id)
	} else if r.Method == http.MethodPut {
		updateItem(w, r, id)
	} else if r.Method == http.MethodDelete {
		deleteItem(w, r, id)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllItems(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var list []Item
	for _, item := range items {
		list = append(list, item)
	}

	json.NewEncoder(w).Encode(list)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	item.ID = nextID
	items[nextID] = item
	nextID++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func getItem(w http.ResponseWriter, r *http.Request, id int) {
	mu.Lock()
	defer mu.Unlock()

	item, exists := items[id]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request, id int) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := items[id]; !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	item.ID = id
	items[id] = item
	json.NewEncoder(w).Encode(item)
}

func deleteItem(w http.ResponseWriter, r *http.Request, id int) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := items[id]; !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	delete(items, id)
	w.WriteHeader(http.StatusNoContent)
}
