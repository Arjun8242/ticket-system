package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arjun8242/ticket-system.git/handlers"
	"github.com/Arjun8242/ticket-system.git/middleware"
	"github.com/Arjun8242/ticket-system.git/store"
	"github.com/go-chi/chi/v5"
)

func setupTestRouter() http.Handler {
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

	return router
}

func TestHealthHandler(t *testing.T) {
	router := setupTestRouter()

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestAllRoutesWorkflow(t *testing.T) {
	router := setupTestRouter()

	regBodyA := map[string]string{
		"username": "userA",
		"email":    "userA@example.com",
		"password": "password123",
	}
	jsonBodyA, err := json.Marshal(regBodyA)
	if err != nil {
		t.Fatal(err)
	}

	reqRegA, err := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBodyA))
	if err != nil {
		t.Fatal(err)
	}
	reqRegA.Header.Set("Content-Type", "application/json")

	recRegA := httptest.NewRecorder()
	router.ServeHTTP(recRegA, reqRegA)

	if recRegA.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recRegA.Code)
	}

	regBodyB := map[string]string{
		"username": "userB",
		"email":    "userB@example.com",
		"password": "password123",
	}
	jsonBodyB, err := json.Marshal(regBodyB)
	if err != nil {
		t.Fatal(err)
	}

	reqRegB, err := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBodyB))
	if err != nil {
		t.Fatal(err)
	}
	reqRegB.Header.Set("Content-Type", "application/json")

	recRegB := httptest.NewRecorder()
	router.ServeHTTP(recRegB, reqRegB)

	if recRegB.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recRegB.Code)
	}

	loginBodyA := map[string]string{
		"email":    "userA@example.com",
		"password": "password123",
	}
	jsonLoginA, err := json.Marshal(loginBodyA)
	if err != nil {
		t.Fatal(err)
	}

	reqLoginA, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonLoginA))
	if err != nil {
		t.Fatal(err)
	}
	reqLoginA.Header.Set("Content-Type", "application/json")

	recLoginA := httptest.NewRecorder()
	router.ServeHTTP(recLoginA, reqLoginA)

	if recLoginA.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recLoginA.Code)
	}

	var loginRespA handlers.LoginResponse
	err = json.NewDecoder(recLoginA.Body).Decode(&loginRespA)
	if err != nil {
		t.Fatal(err)
	}
	tokenA := loginRespA.Token

	loginBodyB := map[string]string{
		"email":    "userB@example.com",
		"password": "password123",
	}
	jsonLoginB, err := json.Marshal(loginBodyB)
	if err != nil {
		t.Fatal(err)
	}

	reqLoginB, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonLoginB))
	if err != nil {
		t.Fatal(err)
	}
	reqLoginB.Header.Set("Content-Type", "application/json")

	recLoginB := httptest.NewRecorder()
	router.ServeHTTP(recLoginB, reqLoginB)

	if recLoginB.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recLoginB.Code)
	}

	var loginRespB handlers.LoginResponse
	err = json.NewDecoder(recLoginB.Body).Decode(&loginRespB)
	if err != nil {
		t.Fatal(err)
	}
	tokenB := loginRespB.Token

	ticketBody := map[string]string{
		"title":       "Database issue",
		"description": "Connection pool timeout",
	}
	jsonTicket, err := json.Marshal(ticketBody)
	if err != nil {
		t.Fatal(err)
	}

	reqCreateTicket, err := http.NewRequest("POST", "/tickets", bytes.NewBuffer(jsonTicket))
	if err != nil {
		t.Fatal(err)
	}
	reqCreateTicket.Header.Set("Content-Type", "application/json")
	reqCreateTicket.Header.Set("Authorization", "Bearer "+tokenA)

	recCreateTicket := httptest.NewRecorder()
	router.ServeHTTP(recCreateTicket, reqCreateTicket)

	if recCreateTicket.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recCreateTicket.Code)
	}

	var createdTicket struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	err = json.NewDecoder(recCreateTicket.Body).Decode(&createdTicket)
	if err != nil {
		t.Fatal(err)
	}
	ticketID := createdTicket.ID

	reqGetTicketsA, err := http.NewRequest("GET", "/tickets", nil)
	if err != nil {
		t.Fatal(err)
	}
	reqGetTicketsA.Header.Set("Authorization", "Bearer "+tokenA)

	recGetTicketsA := httptest.NewRecorder()
	router.ServeHTTP(recGetTicketsA, reqGetTicketsA)

	if recGetTicketsA.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recGetTicketsA.Code)
	}

	reqGetTicketsB, err := http.NewRequest("GET", "/tickets", nil)
	if err != nil {
		t.Fatal(err)
	}
	reqGetTicketsB.Header.Set("Authorization", "Bearer "+tokenB)

	recGetTicketsB := httptest.NewRecorder()
	router.ServeHTTP(recGetTicketsB, reqGetTicketsB)

	if recGetTicketsB.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recGetTicketsB.Code)
	}

	reqGetByIDB, err := http.NewRequest("GET", "/tickets/"+ticketID, nil)
	if err != nil {
		t.Fatal(err)
	}
	reqGetByIDB.Header.Set("Authorization", "Bearer "+tokenB)

	recGetByIDB := httptest.NewRecorder()
	router.ServeHTTP(recGetByIDB, reqGetByIDB)

	if recGetByIDB.Code != http.StatusNotFound {
		t.Fatalf("expected status %d for unauthorized ticket view, got %d", http.StatusNotFound, recGetByIDB.Code)
	}

	patchBody1 := map[string]string{"status": "in_progress"}
	jsonPatch1, err := json.Marshal(patchBody1)
	if err != nil {
		t.Fatal(err)
	}

	reqPatch1, err := http.NewRequest("PATCH", "/tickets/"+ticketID+"/status", bytes.NewBuffer(jsonPatch1))
	if err != nil {
		t.Fatal(err)
	}
	reqPatch1.Header.Set("Content-Type", "application/json")
	reqPatch1.Header.Set("Authorization", "Bearer "+tokenA)

	recPatch1 := httptest.NewRecorder()
	router.ServeHTTP(recPatch1, reqPatch1)

	if recPatch1.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recPatch1.Code)
	}

	patchBody2 := map[string]string{"status": "closed"}
	jsonPatch2, err := json.Marshal(patchBody2)
	if err != nil {
		t.Fatal(err)
	}

	reqPatch2, err := http.NewRequest("PATCH", "/tickets/"+ticketID+"/status", bytes.NewBuffer(jsonPatch2))
	if err != nil {
		t.Fatal(err)
	}
	reqPatch2.Header.Set("Content-Type", "application/json")
	reqPatch2.Header.Set("Authorization", "Bearer "+tokenA)

	recPatch2 := httptest.NewRecorder()
	router.ServeHTTP(recPatch2, reqPatch2)

	if recPatch2.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recPatch2.Code)
	}

	patchBody3 := map[string]string{"status": "open"}
	jsonPatch3, err := json.Marshal(patchBody3)
	if err != nil {
		t.Fatal(err)
	}

	reqPatch3, err := http.NewRequest("PATCH", "/tickets/"+ticketID+"/status", bytes.NewBuffer(jsonPatch3))
	if err != nil {
		t.Fatal(err)
	}
	reqPatch3.Header.Set("Content-Type", "application/json")
	reqPatch3.Header.Set("Authorization", "Bearer "+tokenA)

	recPatch3 := httptest.NewRecorder()
	router.ServeHTTP(recPatch3, reqPatch3)

	if recPatch3.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for reopening closed ticket, got %d", http.StatusBadRequest, recPatch3.Code)
	}
}

func TestAuthMiddleware(t *testing.T) {
	router := setupTestRouter()

	reqNoToken, err := http.NewRequest("GET", "/tickets", nil)
	if err != nil {
		t.Fatal(err)
	}

	recNoToken := httptest.NewRecorder()
	router.ServeHTTP(recNoToken, reqNoToken)

	if recNoToken.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recNoToken.Code)
	}

	reqBadToken, err := http.NewRequest("GET", "/tickets", nil)
	if err != nil {
		t.Fatal(err)
	}
	reqBadToken.Header.Set("Authorization", "Bearer invalidtoken")

	recBadToken := httptest.NewRecorder()
	router.ServeHTTP(recBadToken, reqBadToken)

	if recBadToken.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recBadToken.Code)
	}
}
