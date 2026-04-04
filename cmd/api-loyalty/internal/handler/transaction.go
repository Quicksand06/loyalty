package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Quicksand06/loyalty/cmd/api-loyalty/internal/domain"
	"github.com/Quicksand06/loyalty/cmd/api-loyalty/internal/store"
)

// EventPublisher publishes domain events to a message broker.
type EventPublisher interface {
	Publish(ctx context.Context, key string, value []byte) error
}

type createTrxRequest struct {
	TransactionID string    `json:"transactionId"`
	Amount        float64   `json:"amount"`
	StoreID       string    `json:"storeId"`
	Timestamp     time.Time `json:"timestamp"`
	CustomerID    string    `json:"customerId"`
}

func CreateTransaction(ts store.TransactionStore, pub EventPublisher) http.HandlerFunc {
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

		t := domain.Transaction{
			TransactionID: req.TransactionID,
			Amount:        req.Amount,
			StoreID:       req.StoreID,
			Timestamp:     req.Timestamp,
			CustomerID:    req.CustomerID,
		}

		if err := ts.CreateTransaction(r.Context(), t); err != nil {
			switch {
			case errors.Is(err, store.ErrCustomerNotFound):
				http.Error(w, "customer not found", http.StatusConflict)
			case errors.Is(err, store.ErrCustomerInactive):
				http.Error(w, "customer is inactive", http.StatusConflict)
			default:
				log.Println("failed to insert transaction:", err)
				http.Error(w, "failed to store transaction", http.StatusInternalServerError)
			}
			return
		}

		data, err := json.Marshal(t)
		if err != nil {
			log.Println("failed to marshal transaction event:", err)
		} else if err := pub.Publish(r.Context(), t.TransactionID, data); err != nil {
			log.Println("failed to publish transaction event:", err)
		}

		w.WriteHeader(http.StatusCreated)
	}
}
