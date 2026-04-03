package main

import (
	"log"
	"net/http"

	"github.com/Quicksand06/loyalty/cmd/api-loyalty/internal/config"
	"github.com/Quicksand06/loyalty/cmd/api-loyalty/internal/handler"
	"github.com/Quicksand06/loyalty/cmd/api-loyalty/internal/migrations"
	"github.com/Quicksand06/loyalty/cmd/api-loyalty/internal/store/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := postgres.Migrate(db, migrations.FS); err != nil {
		log.Fatal(err)
	}

	s := postgres.NewStore(db)
	router := handler.NewRouter(s, s)

	log.Println("starting server on", cfg.HTTPAddr)
	log.Fatal((&http.Server{Addr: cfg.HTTPAddr, Handler: router}).ListenAndServe())
}
