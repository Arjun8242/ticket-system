FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ticket-system .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ticket-system .

EXPOSE 8080

CMD ["./ticket-system"]
