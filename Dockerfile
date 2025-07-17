# Multi-stage build for commitlint webhook server
FROM golang:1.23-alpine AS builder

# Install dependencies
RUN apk add --no-cache git make ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' -u 10001 appuser

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the webhook server with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty) -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -a -installsuffix cgo \
    -o commitlint-webhook-server \
    ./cmd/webhook-server

# Final stage - minimal image
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /build/commitlint-webhook-server /commitlint-webhook-server

# Copy default config
COPY --from=builder /build/webhook-server.yml.example /etc/commitlint/webhook-server.yml.example

# Use non-root user
USER appuser

# Expose default port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/commitlint-webhook-server", "-health"]

ENTRYPOINT ["/commitlint-webhook-server"]
CMD ["-c", "/etc/commitlint/webhook-server.yml"]