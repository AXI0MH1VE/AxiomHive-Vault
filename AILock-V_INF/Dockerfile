# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate protobuf code (if needed)
# RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/auth-service
RUN CGO_ENABLED=0 GOOS=linux go build -o rbac-service ./cmd/rbac-service
RUN CGO_ENABLED=0 GOOS=linux go build -o detenforce-proxy ./cmd/detenforce-proxy

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/auth-service .
COPY --from=builder /app/rbac-service .
COPY --from=builder /app/detenforce-proxy .
# Copy config if needed
# COPY config/ ./config/
