package models

import "time"

type Ticket struct {
	ID          string
	Title       string
	Description string
	Status      string // "open" | "in_progress" | "closed"
	OwnerID     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}