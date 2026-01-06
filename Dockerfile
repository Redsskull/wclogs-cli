# Multi-stage Dockerfile for Warcraft Logs CLI
FROM golang:1.25-alpine AS builder

# Install git for go modules
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wclogs .

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS calls to WCL API
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN adduser -D -s /bin/sh wclogs

# Set working directory
WORKDIR /home/wclogs

# Copy binary from builder stage
COPY --from=builder /app/wclogs .

# Change ownership to non-root user
RUN chown wclogs:wclogs wclogs

# Switch to non-root user
USER wclogs

# Create directory for config
RUN mkdir -p /home/wclogs/.config

# Expose the binary
ENTRYPOINT ["./wclogs"]

# Default command shows help
CMD ["--help"]
