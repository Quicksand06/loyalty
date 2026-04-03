package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	db, err := sql.Open("pgx", "postgres://postgres:Pa55word!@localhost:5432/loyalty")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			transaction_id TEXT PRIMARY KEY,
			amount         NUMERIC(19, 4) NOT NULL,
			store_id       TEXT NOT NULL,
			timestamp      TIMESTAMPTZ NOT NULL,
			customer_id    TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal("failed to create transactions table:", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /transactions", handleCreateTrx(db))
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("starting server on :8080")
	log.Fatal(server.ListenAndServe())
}

type createTrxRequest struct {
	TransactionID string    `json:"transactionId"`
	Amount        float64   `json:"amount"`
	StoreID       string    `json:"storeId"`
	Timestamp     time.Time `json:"timestamp"`
	CustomerID    string    `json:"customerId"`
}

func handleCreateTrx(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createTrxRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.TransactionID == "" || req.StoreID == "" || req.CustomerID == "" {
			http.Error(w, "transactionId, storeId and customerId are required", http.StatusBadRequest)
			return
		}

		_, err := db.ExecContext(r.Context(), `
			INSERT INTO transactions (transaction_id, amount, store_id, timestamp, customer_id)
			VALUES ($1, $2, $3, $4, $5)
		`, req.TransactionID, req.Amount, req.StoreID, req.Timestamp.UTC(), req.CustomerID)
		if err != nil {
			log.Println("failed to insert transaction:", err)
			http.Error(w, "failed to store transaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
