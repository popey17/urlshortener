package main

import (
	"log"
	"net/http"

	"github.com/popey17/urlshortener/internal/api"
	"github.com/popey17/urlshortener/internal/config"
	"github.com/popey17/urlshortener/internal/store"
)

func main() {
	cfg, _ := config.Load()

	port := cfg.Addr
	// fmt.Println(cfg)

	mux := http.NewServeMux()

	s, err := store.OpenDB(cfg.DbPath)
	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()

	h := &api.Handler{
		Store:   s,
		BaseUrl: cfg.BaseUrl,
	}

	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("POST /api/v1/urls", h.CreateUrl)
	mux.HandleFunc("GET /{code}", h.Redirect)

	log.Fatal(http.ListenAndServe(port, mux))
}
