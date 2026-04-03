package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	dsn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := sql.Open("pgx", dsn)
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

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS customers (
			customer_id      TEXT PRIMARY KEY,
			identifier_type  TEXT NOT NULL,
			customer_active  BOOLEAN NOT NULL DEFAULT TRUE
		)
	`)
	if err != nil {
		log.Fatal("failed to create customers table:", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /transactions", handleCreateTrx(db))
	mux.HandleFunc("POST /customers", handleCreateCustomer(db))
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("starting server on :8080")
	log.Fatal(server.ListenAndServe())
}

type createCustomerRequest struct {
	CustomerID     string         `json:"customerId"`
	IdentifierType IdentifierType `json:"identifierType"`
}

func handleCreateCustomer(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createCustomerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.CustomerID == "" {
			http.Error(w, "customerId is required", http.StatusBadRequest)
			return
		}

		if !req.IdentifierType.isValid() {
			http.Error(w, "invalid identifierType", http.StatusBadRequest)
			return
		}

		_, err := db.ExecContext(r.Context(), `
			INSERT INTO customers (customer_id, identifier_type)
			VALUES ($1, $2)
		`, req.CustomerID, req.IdentifierType)
		if err != nil {
			log.Println("failed to insert customer:", err)
			http.Error(w, "failed to store customer", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
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

type IdentifierType string

const (
	IdentifierTypeLoyaltyCard  IdentifierType = "loyalty_card"
	IdentifierTypeMembershipID IdentifierType = "membership_id"
	IdentifierTypeEmailAddress IdentifierType = "email_address"
)

func (e IdentifierType) isValid() bool {
	switch e {
	case IdentifierTypeLoyaltyCard, IdentifierTypeMembershipID, IdentifierTypeEmailAddress:
		return true
	}
	return false
}
