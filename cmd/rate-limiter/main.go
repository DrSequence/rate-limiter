package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"rate-limiter/internal/app"
	"rate-limiter/internal/cache"
	"rate-limiter/internal/config"
	internal_http "rate-limiter/internal/http"
)

func main() {
	cfg := config.ReadFile()

	cacheRepo := cache.NewMemcachedRepository(cfg.Cache.Limit, cfg.Cache.Ttl, cfg.Cache.Server)
	rateLimiterService := app.NewRateLimiterService(cacheRepo)
	handler := internal_http.NewRateLimiterHandler(rateLimiterService)

	r := mux.NewRouter()
	r.HandleFunc("/verify", handler.Verify).Methods("POST")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(cfg.Server.Host+":"+cfg.Server.Port, r); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}

	// reset counters
	app.ResetCounters(cacheRepo)
}
