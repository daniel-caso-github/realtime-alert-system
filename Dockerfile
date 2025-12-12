FROM golang:1.24-alpine AS builder

LABEL authors="danielcaso"

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files first (better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/api \
    ./cmd/api

# ============================================================================
# STAGE 2: Development (with hot-reload)
# ============================================================================
FROM golang:1.24-alpine AS development

# Install Air for hot-reload
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code (will be overwritten by volume mount)
COPY . .

# Expose port
EXPOSE 8080

# Run with Air for hot-reload
CMD ["air", "-c", ".air.toml"]

# ============================================================================
# STAGE 3: Production
# ============================================================================
FROM alpine:3.19 AS production

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/api .

# Copy config file
COPY config.yaml .

# Change ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./api"]
