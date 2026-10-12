# Multi-stage Dockerfile for Go Backend
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build server binary and reset tool binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o reset-db ./cmd/reset/main.go

FROM alpine:3.19

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata curl

ENV TZ=Asia/Jakarta

COPY --from=builder /app/server .
COPY --from=builder /app/reset-db .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./server"]
