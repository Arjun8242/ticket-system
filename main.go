package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/Arjun8242/ticket-system.git/handlers"
	"github.com/Arjun8242/ticket-system.git/middleware"
	"github.com/Arjun8242/ticket-system.git/store"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	}

	appStore := store.NewStore()
	authHandler := handlers.NewAuthHandler(appStore)
	ticketHandler := handlers.NewTicketHandler(appStore)

	router := chi.NewRouter()

	router.Get("/health", handlers.HealthHandler)
	router.Post("/auth/register", authHandler.RegisterHandler)
	router.Post("/auth/login", authHandler.LoginHandler)

	router.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)

		r.Post("/tickets", ticketHandler.CreateTicketHandler)
		r.Get("/tickets", ticketHandler.GetTicketsHandler)
		r.Get("/tickets/{id}", ticketHandler.GetTicketByIDHandler)
		r.Patch("/tickets/{id}/status", ticketHandler.UpdateTicketStatusHandler)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	http.ListenAndServe(port, router)
}