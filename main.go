package main

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/Arjun8242/ticket-system.git/handlers"
)

func main() {
	router := chi.NewRouter()
	router.Get("/health", handlers.HealthHandler)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	http.ListenAndServe(":8080", router)
}