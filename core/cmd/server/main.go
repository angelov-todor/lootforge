package main

import (
	"log"
	"net/http"
	"os"

	"github.com/angelov-todor/lootforge/core/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handlers.HealthCheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("LootForge API starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
