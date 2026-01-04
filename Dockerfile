# syntax=docker/dockerfile:1.4

# ============================================
# Stage 1: Dependencies
# ============================================
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS deps

WORKDIR /build

# Copy dependency files first (better layer caching)
COPY go.mod go.sum ./

# Download dependencies with cache mount for faster rebuilds
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

# ============================================
# Stage 2: Build Stage
# ============================================
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

# Build arguments to support cross-compilation
ARG TARGETOS
ARG TARGETARCH
ARG VERSION

WORKDIR /build

# Copy go.mod/go.sum and pre-downloaded modules from deps stage
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY go.mod go.sum ./

# Copy source code
COPY . .

# Build with security flags and build cache mount
# CGO_ENABLED=0 - Static binary (no C dependencies)
# -ldflags="-w -s" - Strip debug info (smaller binary)
# -trimpath - Remove file system paths from binary
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-w -s -X main.version=${VERSION:-dev}" \
    -trimpath \
    -o api \
    ./cmd/api

# ============================================
# Stage 3: Runtime Stage (Distroless)
# ============================================
FROM gcr.io/distroless/static-debian12:nonroot

# Copy binary from builder
COPY --from=builder /build/api /app/api

# Metadata labels
LABEL org.opencontainers.image.source="https://github.com/amorags/golangtickets"
LABEL org.opencontainers.image.description="Go API for ticket booking system"
LABEL org.opencontainers.image.licenses="MIT"

# Use non-root user (distroless provides 'nonroot' user with UID 65532)
USER nonroot:nonroot

# Expose port (documentation only, doesn't actually open port)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/api", "healthcheck"]

# Run the binary
ENTRYPOINT ["/app/api"]