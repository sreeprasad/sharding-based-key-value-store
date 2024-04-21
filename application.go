package main

import (
	"airline-checkin-system/toy_store"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var toyStore *toy_store.ToyStore

func main() {

	connStr := "host=localhost port=6432 user=user4 dbname=mydatabase4 password=password4 sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	toyStore = toy_store.NewToyStore(db)

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/set", handleSet)
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
