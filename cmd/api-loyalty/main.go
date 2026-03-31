package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /transactions", handleCreateTrx)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("starting server on :8080")
	log.Fatal(server.ListenAndServe())
}

func handleCreateTrx(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Body)
	w.WriteHeader(http.StatusCreated)
}
