package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Quicksand06/loyalty/cmd/api-loyalty/internal/domain"
	"github.com/Quicksand06/loyalty/cmd/api-loyalty/internal/store"
)

type createCustomerRequest struct {
	CustomerID     string               `json:"customerId"`
	IdentifierType domain.IdentifierType `json:"identifierType"`
}

func CreateCustomer(cs store.CustomerStore) http.HandlerFunc {
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

		if !req.IdentifierType.IsValid() {
			http.Error(w, "invalid identifierType", http.StatusBadRequest)
			return
		}

		c := domain.Customer{
			CustomerID:     req.CustomerID,
			IdentifierType: req.IdentifierType,
		}

		if err := cs.CreateCustomer(r.Context(), c); err != nil {
			log.Println("failed to insert customer:", err)
			http.Error(w, "failed to store customer", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
