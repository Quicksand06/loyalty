package handler

import (
	"net/http"

	"github.com/Quicksand06/loyalty/cmd/api-loyalty/internal/store"
)

func NewRouter(cs store.CustomerStore, ts store.TransactionStore, pub EventPublisher) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /customers", CreateCustomer(cs))
	mux.HandleFunc("POST /transactions", CreateTransaction(ts, pub))
	return mux
}
