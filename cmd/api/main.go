package main

import (
	"log"
	"net/http"

	"github.com/popey17/urlshortener/internal/api"
	"github.com/popey17/urlshortener/internal/config"
	"github.com/popey17/urlshortener/internal/store"
)

func isAllowed(origin string, allowed []string) bool {
	for _, a := range allowed {
		if origin == a {
			return true
		}
	}
	return false
}

func withCORS(next http.Handler, allowed []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && isAllowed(origin, allowed) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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

	log.Fatal(http.ListenAndServe(port, withCORS(mux, cfg.CorsOrigins)))
}
