package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Arjun8242/ticket-system.git/middleware"
	"github.com/Arjun8242/ticket-system.git/models"
	"github.com/Arjun8242/ticket-system.git/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TicketHandler struct {
	Store *store.Store
}

func NewTicketHandler(s *store.Store) *TicketHandler {
	return &TicketHandler{Store: s}
}

type CreateTicketRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

func (h *TicketHandler) CreateTicketHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateTicketRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	now := time.Now()
	ticket := &models.Ticket{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Status:      "open",
		OwnerID:     userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = h.Store.CreateTicket(ticket)
	if err != nil {
		http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ticket)
}

func (h *TicketHandler) GetTicketsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tickets := h.Store.GetTicketsByOwner(userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tickets)
}

func (h *TicketHandler) GetTicketByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ticketID := chi.URLParam(r, "id")
	if ticketID == "" {
		http.Error(w, "Ticket ID is required", http.StatusBadRequest)
		return
	}

	ticket, exists := h.Store.GetTicketByID(ticketID)
	if !exists {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	if ticket.OwnerID != userID {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ticket)
}

func (h *TicketHandler) UpdateTicketStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ticketID := chi.URLParam(r, "id")
	if ticketID == "" {
		http.Error(w, "Ticket ID is required", http.StatusBadRequest)
		return
	}

	ticket, exists := h.Store.GetTicketByID(ticketID)
	if !exists {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	if ticket.OwnerID != userID {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	var req UpdateStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status != "open" {
		if req.Status != "in_progress" {
			if req.Status != "closed" {
				http.Error(w, "Invalid status value", http.StatusBadRequest)
				return
			}
		}
	}

	if ticket.Status == "closed" {
		http.Error(w, "Cannot update status of a closed ticket", http.StatusBadRequest)
		return
	}

	ticket.Status = req.Status
	ticket.UpdatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ticket)
}