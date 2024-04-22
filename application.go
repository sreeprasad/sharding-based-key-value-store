package main

import (
	"airline-checkin-system/toy_store"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var toyStore *toy_store.ToyStore

func main() {

	db1ConnStr := "host=localhost port=6432 user=user4 dbname=mydatabase4_1 password=password4 sslmode=disable"
	db1, err := sql.Open("postgres", db1ConnStr)

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	db2ConnStr := "host=localhost port=6432 user=user4 dbname=mydatabase4_2 password=password4 sslmode=disable"
	db2, err := sql.Open("postgres", db2ConnStr)

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	toyStore = toy_store.NewToyStore(db1, db2)

	defer db1.Close()
	defer db2.Close()

	http.HandleFunc("/set", handleSet)
	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/delete", handleDelete)

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

type SetRequest struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func handleSet(w http.ResponseWriter, r *http.Request) {

	var req SetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Key == "" || req.Value == "" || req.ExpiredAt.IsZero() {
		http.Error(w, "Missing key, value, or expiredAt", http.StatusBadRequest)
		return
	}

	_, err := toyStore.Set(req.Key, req.Value, req.ExpiredAt)
	if err != nil {
		log.Printf("Failed to insert or update record: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func handleGet(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	toyStoreRecord, err := toyStore.Get(key)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Value: %s, ExpiredAt: %v", toyStoreRecord.Value, toyStoreRecord.ExpiredAt)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	_, err := toyStore.Delete(key)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Key:%s deleted", key)
}
