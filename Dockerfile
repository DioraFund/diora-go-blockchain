# Multi-stage build for Go blockchain node
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o diora \
    ./core

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN addgroup -g 1000 -S diora && \
    adduser -u 1000 -S diora -G diora

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/diora /app/diora

# Create necessary directories
RUN mkdir -p /app/data /app/config /app/logs

# Change ownership
RUN chown -R diora:diora /app

# Switch to non-root user
USER diora

# Expose ports
EXPOSE 8545 8546

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8545/health || exit 1

# Default command
CMD ["./diora", "node", "start", "--config=/app/config/config.toml"]
