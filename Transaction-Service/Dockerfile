# Build stage - Updated to Go 1.23
FROM golang:1.24.4-alpine AS builder

# Install git and ca-certificates for module downloads
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Set Go module proxy and download with retries
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org
ENV CGO_ENABLED=0
ENV GOOS=linux

# Download dependencies with retry logic
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application
RUN go build -ldflags="-w -s" -o account-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS and wget for health check
RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/account-service .

# Copy any necessary config files
COPY --from=builder /app/.env* ./

# Expose port
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Run the binary
CMD ["./account-service"]