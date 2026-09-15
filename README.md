# Ticket System Backend (Go)

A lightweight RESTful API for ticket management built with Go, Chi router, JWT authentication, and an in-memory thread-safe store.

## Prerequisites
- Go 1.22+
- Docker (optional, for containerized run)

## Environment Variables
Copy `.env.example` to `.env`:
```env
PORT=8080
JWT_SECRET=your-secret-here
```

## Run Locally
```bash
go run main.go
```

## Run Unit & Integration Tests
```bash
go test -v ./...
```

## Run with Docker
```bash
docker build -t ticket-system .
docker run -p 8080:8080 -e JWT_SECRET=mynameisarjunjaiswal ticket-system
```

## API Endpoints
- `GET /health` — Health check endpoint (Public)
- `POST /auth/register` — Register a new user (Public)
- `POST /auth/login` — Login and receive JWT token (Public)
- `POST /tickets` — Create a new ticket (Protected)
- `GET /tickets` — List owner's tickets (Protected)
- `GET /tickets/{id}` — Get ticket details by ID (Protected)
- `PATCH /tickets/{id}/status` — Update ticket status (`open` -> `in_progress` -> `closed`) (Protected)

## Key Design Assumptions
- **404 for Non-Owned Tickets**: Accessing or updating a ticket owned by another user returns `404 Not Found` (rather than `403 Forbidden`) to avoid leaking ticket existence.
- **Terminal Closed State**: Once a ticket's status is set to `closed`, further status updates are rejected with `400 Bad Request`.
- **In-Memory Store**: Data is stored in-memory using a thread-safe mutex store. Data resets when the application stops.
