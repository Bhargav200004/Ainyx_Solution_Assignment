# ----- Build Stage -----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache Go modules.
COPY go.mod go.sum ./
RUN go mod download

# Copy source code.
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server/

# ----- Runtime Stage -----
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder.
COPY --from=builder /app/server .

# Copy migration files.
COPY --from=builder /app/db/migrations/ ./db/migrations/

EXPOSE 9090

ENTRYPOINT ["./server"]
