# AILock Dockerfile - Phase 2
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.work ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ailock-api ./cmd/ailock-api

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/ailock-api .

# Copy configuration files
COPY --from=builder /app/prom_metrics.yml .

# Expose ports
EXPOSE 8080 9090

# Run the application
CMD ["./ailock-api"]
