package main

import (
	"fmt"
	"net/http"

	"github.com/Arjun8242/ticket-system.git/handlers"
	"github.com/Arjun8242/ticket-system.git/store"
	"github.com/go-chi/chi/v5"
)

func main() {
	appStore := store.NewStore()
	authHandler := handlers.NewAuthHandler(appStore)

	router := chi.NewRouter()

	router.Get("/health", handlers.HealthHandler)
	router.Post("/auth/register", authHandler.RegisterHandler)
	router.Post("/auth/login", authHandler.LoginHandler)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	http.ListenAndServe(":8080", router)
}