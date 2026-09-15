package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/Arjun8242/ticket-system.git/store"
	"github.com/Arjun8242/ticket-system.git/utils"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}