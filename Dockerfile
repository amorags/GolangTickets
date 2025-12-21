# syntax=docker/dockerfile:1.4

# ============================================ 
# Stage 1: Build Stage
# ============================================ 
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user for runtime
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

WORKDIR /build

# Copy dependency files first (layer caching)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build with security flags
# CGO_ENABLED=0 - Static binary (no C dependencies)
# -ldflags="-w -s" - Strip debug info (smaller binary)
# -trimpath - Remove file system paths from binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-w -s -X main.version=${VERSION:-dev}" \
    -trimpath \
    -o api \
    ./cmd/api

# ============================================ 
# Stage 2: Runtime Stage (Distroless)
# ============================================ 
FROM gcr.io/distroless/static-debian12:nonroot

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data (if needed)
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary from builder
COPY --from=builder /build/api /app/api

# Use non-root user (distroless provides 'nonroot' user with UID 65532)
USER nonroot:nonroot

# Expose port (documentation only, doesn't actually open port)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/api", "healthcheck"]

# Run the binary
ENTRYPOINT ["/app/api"]