# DevOps & Docker Security Guide

**Context:** Secure containerized software delivery for isolated factory environments (Agramkow use case)
**Exam Focus:** Docker security best practices, CI/CD pipeline hardening, production deployment
**Date:** December 21, 2025

---

## Table of Contents

1. [Current State Analysis](#current-state-analysis)
2. [Docker Security Fundamentals](#docker-security-fundamentals)
3. [Secure Multi-Stage Builds](#secure-multi-stage-builds)
4. [Container Hardening](#container-hardening)
5. [Secrets Management](#secrets-management)
6. [Network Security](#network-security)
7. [CI/CD Pipeline Security](#cicd-pipeline-security)
8. [Image Scanning & Vulnerability Management](#image-scanning--vulnerability-management)
9. [Supply Chain Security](#supply-chain-security)
10. [Production Deployment Strategies](#production-deployment-strategies)
11. [Monitoring & Logging](#monitoring--logging)
12. [Agramkow-Specific Use Case](#agramkow-specific-use-case)
13. [Implementation Checklist](#implementation-checklist)

---

## Current State Analysis

### ✅ What You're Doing Right

1. **Multi-stage builds** - Nuxt Dockerfile uses proper build stages
2. **Non-root user** - Nuxt container runs as `nuxt` user (UID 1001)
3. **Alpine base images** - Smaller attack surface
4. **Health checks** - PostgreSQL has proper health monitoring
5. **Basic .dockerignore** - Web app excludes node_modules, .env

### ❌ Security Issues to Fix

#### Critical (Exam Killers):
1. **Hardcoded secrets** in docker-compose.yml (DB password, JWT secret)
2. **Root user** - Go API runs as root (no USER directive)
3. **No image scanning** - Vulnerable dependencies not detected
4. **Single-stage Go build** - Includes build tools in production image
5. **Missing .dockerignore** - Go app copies everything (including .git, .env)
6. **Exposed ports** - All services bind to 0.0.0.0 (accessible from host)
7. **No resource limits** - Containers can consume unlimited CPU/memory
8. **Latest tag** in production - docker-compose.prod.yml references :latest

#### High Priority:
9. **No read-only filesystem** - Containers can modify their filesystem
10. **No security profiles** - AppArmor/Seccomp not configured
11. **Privileged capabilities** - Running with default capabilities
12. **No network segmentation** - All containers on same network
13. **Missing image signing** - No verification of image authenticity
14. **No SBOM** - No Software Bill of Materials for dependency tracking

---

## Docker Security Fundamentals

### The 4 Pillars of Container Security

```
┌─────────────────────────────────────────────────────────┐
│                    Container Security                    │
├──────────────┬──────────────┬──────────────┬────────────┤
│   Image      │   Runtime    │   Network    │   Secrets  │
│   Security   │   Security   │   Security   │   Mgmt     │
└──────────────┴──────────────┴──────────────┴────────────┘
```

#### 1. Image Security
- Minimal base images (distroless, scratch, alpine)
- Multi-stage builds (no build tools in production)
- Regular base image updates
- Vulnerability scanning
- Image signing and verification

#### 2. Runtime Security
- Non-root users
- Read-only filesystem
- Resource limits (CPU, memory)
- Security profiles (AppArmor, Seccomp)
- Capability dropping

#### 3. Network Security
- Network segmentation (multiple Docker networks)
- Least privilege (only necessary ports exposed)
- TLS for inter-service communication
- Internal DNS resolution

#### 4. Secrets Management
- Environment variables from secrets (not hardcoded)
- Docker secrets or external vaults (HashiCorp Vault, AWS Secrets Manager)
- Secret rotation
- Encryption at rest and in transit

---

## Secure Multi-Stage Builds

### Current Go Dockerfile (INSECURE)

```dockerfile
# ❌ Problems:
# - Single stage (build tools in production)
# - Root user
# - Full Go toolchain in final image
# - No .dockerignore

FROM golang:1.25-alpine
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api
EXPOSE 8080
CMD ["./main"]
```

**Image size:** ~400MB
**Attack surface:** Full Go compiler, git, build tools
**User:** root (UID 0)

---

### Secure Multi-Stage Dockerfile (BEST PRACTICE)

```dockerfile
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
```

**Image size:** ~15MB (96% smaller!)
**Attack surface:** No shell, no package manager, no build tools
**User:** nonroot (UID 65532)

---

### Alternative: Scratch-Based (Even Smaller)

```dockerfile
FROM scratch

# Copy CA certs and timezone (if your app needs them)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary
COPY --from=builder /build/api /api

# Scratch doesn't have users, but you can still set UID
USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/api"]
```

**Image size:** ~8MB
**Attack surface:** Literally just your binary
**Limitation:** No shell (can't debug with `docker exec`)

---

### Secure Nuxt Dockerfile Improvements

Your current Nuxt Dockerfile is pretty good! Minor improvements:

```dockerfile
# syntax=docker/dockerfile:1.4

FROM oven/bun:1-alpine AS base

# ============================================
# Stage 1: Dependencies
# ============================================
FROM base AS deps
WORKDIR /app

# Install security updates
RUN apk upgrade --no-cache

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile --production=false

# ============================================
# Stage 2: Builder
# ============================================
FROM base AS builder
WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Build with production optimizations
ENV NODE_ENV=production
RUN bun run build

# ============================================
# Stage 3: Runtime
# ============================================
FROM base AS runner
WORKDIR /app

# Install security updates
RUN apk upgrade --no-cache

ENV NODE_ENV=production

# Create non-root user (already done, good!)
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nuxt

# Copy only production artifacts
COPY --from=builder --chown=nuxt:nodejs /app/.output /app/.output
COPY --from=builder --chown=nuxt:nodejs /app/package.json /app/package.json

# Switch to non-root user
USER nuxt

EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

ENV PORT=3000
ENV HOST=0.0.0.0

CMD ["bun", "run", ".output/server/index.mjs"]
```

---

## Container Hardening

### 1. Non-Root User (Critical)

**Why it matters:**
- Root in container = potential root on host (with kernel exploits)
- Agramkow factories: prevents compromised container from affecting host machine

**Implementation:**

```dockerfile
# Method 1: Create user in Dockerfile
RUN addgroup -g 1001 appuser && \
    adduser -u 1001 -D -G appuser appuser

USER appuser

# Method 2: Use numeric UID (works with distroless)
USER 1001:1001

# Method 3: Use distroless nonroot
FROM gcr.io/distroless/static-debian12:nonroot
# Automatically uses UID 65532
```

**Verification:**

```bash
# Check user inside container
docker exec -it <container> whoami
# Should NOT be "root"

# Check process owner on host
ps aux | grep api
# Should show UID 1001, not root
```

---

### 2. Read-Only Filesystem

**Why it matters:**
- Prevents malware from writing to disk
- Immutable infrastructure
- Factory environments: prevents tampering with software

**Implementation:**

```yaml
# docker-compose.yml
services:
  api:
    read_only: true
    tmpfs:
      - /tmp:noexec,nosuid,size=64m
      - /app/logs:noexec,nosuid,size=128m
```

**For apps that need write access:**

```yaml
services:
  api:
    read_only: true
    volumes:
      - ./logs:/app/logs:rw  # Only this directory is writable
    tmpfs:
      - /tmp:size=64m
```

**Go app adjustment:**

```go
// Use /tmp for temporary files (tmpfs mount)
tmpFile, err := os.CreateTemp("/tmp", "upload-*")

// Or use in-memory
var buf bytes.Buffer
```

---

### 3. Resource Limits (Prevent DoS)

**Why it matters:**
- Prevent container from consuming all host resources
- Factory QA machines: ensure stability when multiple containers run

**Implementation:**

```yaml
# docker-compose.yml
services:
  api:
    deploy:
      resources:
        limits:
          cpus: '1.0'        # Max 1 CPU core
          memory: 512M       # Max 512MB RAM
        reservations:
          cpus: '0.25'       # Guaranteed 0.25 cores
          memory: 128M       # Guaranteed 128MB RAM

    # Process limits
    ulimits:
      nproc: 512             # Max processes
      nofile:                # Max open files
        soft: 1024
        hard: 2048

    # Prevent fork bombs
    pids_limit: 200
```

**Testing resource limits:**

```bash
# Stress test CPU
docker run --rm --cpus="0.5" stress-ng --cpu 4 --timeout 30s

# Stress test memory
docker run --rm --memory="100m" stress-ng --vm 1 --vm-bytes 200M
# Should get OOM killed
```

---

### 4. Drop Capabilities (Principle of Least Privilege)

**Why it matters:**
- Linux capabilities give fine-grained permissions
- Default Docker has ~14 capabilities (including CAP_NET_RAW for packet sniffing!)
- Drop all, add only what you need

**Default capabilities Docker gives:**

```
CHOWN, DAC_OVERRIDE, FOWNER, FSETID, KILL, SETGID, SETUID, SETPCAP,
NET_BIND_SERVICE, NET_RAW, SYS_CHROOT, MKNOD, AUDIT_WRITE, SETFCAP
```

**Secure configuration:**

```yaml
# docker-compose.yml
services:
  api:
    cap_drop:
      - ALL                  # Drop all capabilities
    cap_add:
      - NET_BIND_SERVICE     # Only if binding to port <1024
      # Most apps need NOTHING!

    security_opt:
      - no-new-privileges:true  # Prevent privilege escalation
```

**For Go API that listens on 8080 (not privileged port):**

```yaml
services:
  api:
    cap_drop:
      - ALL                  # No capabilities needed!
    security_opt:
      - no-new-privileges:true
```

---

### 5. Seccomp Profile (Syscall Filtering)

**Why it matters:**
- Seccomp filters which system calls container can make
- Docker default allows 300+ syscalls
- Most apps use <100

**Custom Seccomp Profile:**

```json
// seccomp-profile.json
{
  "defaultAction": "SCMP_ACT_ERRNO",
  "architectures": [
    "SCMP_ARCH_X86_64",
    "SCMP_ARCH_X86",
    "SCMP_ARCH_X32"
  ],
  "syscalls": [
    {
      "names": [
        "accept4", "access", "arch_prctl", "bind", "brk",
        "clone", "close", "connect", "dup", "dup2",
        "epoll_create1", "epoll_ctl", "epoll_pwait", "epoll_wait",
        "execve", "exit", "exit_group", "fcntl", "fstat",
        "futex", "getcwd", "getdents64", "getpid", "getppid",
        "getrandom", "getsockname", "getsockopt", "getuid",
        "listen", "lseek", "madvise", "mmap", "mprotect",
        "munmap", "nanosleep", "open", "openat", "pipe2",
        "poll", "pread64", "read", "readlinkat", "recvfrom",
        "recvmsg", "rt_sigaction", "rt_sigprocmask", "rt_sigreturn",
        "sched_getaffinity", "sched_yield", "sendmsg", "sendto",
        "set_robust_list", "set_tid_address", "setgid", "setgroups",
        "setsockopt", "setuid", "sigaltstack", "socket", "stat",
        "tgkill", "uname", "write", "writev"
      ],
      "action": "SCMP_ACT_ALLOW"
    }
  ]
}
```

**Apply in docker-compose:**

```yaml
services:
  api:
    security_opt:
      - seccomp:./seccomp-profile.json
```

**Generate profile from running container:**

```bash
# Use docker-slim to auto-generate minimal profile
docker-slim build --http-probe your-api-image
```

---

### 6. AppArmor/SELinux Profiles

**AppArmor (Ubuntu/Debian):**

```yaml
services:
  api:
    security_opt:
      - apparmor=docker-default  # Or custom profile
```

**SELinux (RHEL/CentOS):**

```yaml
services:
  api:
    security_opt:
      - label=type:container_runtime_t
```

For Agramkow factories, check which security module the host OS uses:

```bash
# Check AppArmor
aa-status

# Check SELinux
sestatus
```

---

## Secrets Management

### ❌ NEVER Do This (Your Current Setup)

```yaml
# docker-compose.yml - INSECURE!
environment:
  DB_PASSWORD: password  # ❌ Hardcoded in file (committed to git!)
  JWT_SECRET: your-super-secret-jwt-key-change-this-in-production-min-32-chars  # ❌ In plain text
```

**Why it's bad:**
- Secrets in git history (even if you delete them later)
- Readable by anyone with access to docker-compose.yml
- Visible in `docker inspect`
- Logs may capture environment variables

---

### ✅ Method 1: Environment Files (Basic)

**Step 1: Create `.env` file (gitignored)**

```bash
# .env
DB_USER=ticket_user
DB_PASSWORD=XmK9$vP#2nQ@8wR
DB_NAME=ticket_db
JWT_SECRET=9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08
API_BASE_URL=http://localhost:8080
```

**Step 2: Update docker-compose.yml**

```yaml
services:
  db:
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}

  api:
    environment:
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: ${DB_NAME}
      JWT_SECRET: ${JWT_SECRET}
```

**Step 3: Add to .gitignore**

```gitignore
.env
.env.*
!.env.example
```

**Step 4: Create .env.example (committed)**

```bash
# .env.example
DB_USER=user
DB_PASSWORD=change_me
DB_NAME=ticket_db
JWT_SECRET=generate_with_openssl_rand_hex_32
API_BASE_URL=http://localhost:8080
```

---

### ✅ Method 2: Docker Secrets (Production)

**Why better:**
- Secrets encrypted at rest (Swarm mode)
- Never written to disk in plain text
- Mounted as in-memory filesystem at `/run/secrets/`
- Rotatable without rebuilding images

**Step 1: Create secrets**

```bash
# Generate strong secrets
openssl rand -hex 32 > db_password.txt
openssl rand -hex 32 > jwt_secret.txt

# Create Docker secrets
docker secret create db_password db_password.txt
docker secret create jwt_secret jwt_secret.txt

# Delete local files
shred -u db_password.txt jwt_secret.txt
```

**Step 2: Update docker-compose.yml for Swarm**

```yaml
version: '3.8'

services:
  api:
    image: your-api:latest
    secrets:
      - db_password
      - jwt_secret
    environment:
      # Read from /run/secrets/
      DB_PASSWORD_FILE: /run/secrets/db_password
      JWT_SECRET_FILE: /run/secrets/jwt_secret

secrets:
  db_password:
    external: true
  jwt_secret:
    external: true
```

**Step 3: Update Go app to read secrets**

```go
// internal/config/config.go

func loadSecret(envVar string) string {
    // Check if _FILE variant exists
    fileEnv := envVar + "_FILE"
    if filepath := os.Getenv(fileEnv); filepath != "" {
        data, err := os.ReadFile(filepath)
        if err != nil {
            log.Fatalf("Failed to read secret from %s: %v", filepath, err)
        }
        return strings.TrimSpace(string(data))
    }

    // Fallback to regular env var
    return os.Getenv(envVar)
}

func LoadConfig() *Config {
    return &Config{
        DBPassword: loadSecret("DB_PASSWORD"),
        JWTSecret:  loadSecret("JWT_SECRET"),
        // ...
    }
}
```

---

### ✅ Method 3: HashiCorp Vault (Enterprise)

**Best for:** Multi-environment deployments, dynamic secrets, audit logging

```go
// Example Vault integration
import (
    vault "github.com/hashicorp/vault/api"
)

func getVaultSecret(key string) string {
    client, err := vault.NewClient(&vault.Config{
        Address: os.Getenv("VAULT_ADDR"),
    })
    if err != nil {
        log.Fatal(err)
    }

    client.SetToken(os.Getenv("VAULT_TOKEN"))

    secret, err := client.Logical().Read(fmt.Sprintf("secret/data/%s", key))
    if err != nil {
        log.Fatal(err)
    }

    return secret.Data["data"].(map[string]interface{})["value"].(string)
}
```

---

### Secret Rotation Strategy

**Database password rotation:**

```bash
# 1. Create new secret
docker secret create db_password_v2 db_password_new.txt

# 2. Update service to use both secrets temporarily
docker service update --secret-rm db_password --secret-add db_password_v2 api

# 3. Update database to accept both passwords
ALTER USER ticket_user WITH PASSWORD 'new_password';

# 4. Deploy app update to use new secret
docker service update --env-add DB_PASSWORD_FILE=/run/secrets/db_password_v2 api

# 5. Remove old password from database
# (After verifying all services use new password)
```

---

## Network Security

### Current Setup (Insecure)

```yaml
# All services on default bridge network
# No isolation between services
services:
  db:
    ports:
      - "5432:5432"  # ❌ PostgreSQL exposed to host network!
  api:
    ports:
      - "8080:8080"
  web:
    ports:
      - "3000:3000"
```

**Problems:**
- Database accessible from host (should be internal-only)
- All containers can talk to each other
- No defense in depth

---

### Secure Network Architecture

```yaml
version: '3.8'

services:
  # ============================================
  # Frontend (Public Network)
  # ============================================
  web:
    image: your-web:latest
    networks:
      - frontend
    ports:
      - "3000:3000"  # Only web is exposed
    depends_on:
      - api

  # ============================================
  # API (Frontend + Backend Networks)
  # ============================================
  api:
    image: your-api:latest
    networks:
      - frontend   # Can receive requests from web
      - backend    # Can connect to database
    # ❌ No ports exposed to host!
    # Web communicates via Docker network
    depends_on:
      - db

  # ============================================
  # Database (Backend Network Only)
  # ============================================
  db:
    image: postgres:16-alpine
    networks:
      - backend    # Only accessible by API
    # ❌ No ports exposed to host!
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      POSTGRES_USER_FILE: /run/secrets/db_user
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
    secrets:
      - db_user
      - db_password

networks:
  frontend:
    driver: bridge
    driver_opts:
      com.docker.network.bridge.name: br-frontend

  backend:
    driver: bridge
    driver_opts:
      com.docker.network.bridge.name: br-backend
    internal: true  # No internet access (optional)

volumes:
  postgres_data:

secrets:
  db_user:
    external: true
  db_password:
    external: true
```

**Network diagram:**

```
Internet
   │
   ├─→ [web:3000] ────→ [api] ────→ [db]
   │    (frontend)      (frontend   (backend
   │                     + backend)  only)
   │
   └─→ [nginx:443]  (future)
        (reverse proxy)
```

---

### TLS Between Services (mTLS)

For sensitive data in Agramkow factories (QA test results, proprietary algorithms):

```yaml
# Use Traefik or Envoy for automatic mTLS
services:
  api:
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.api.rule=Host(`api.local`)"
      - "traefik.http.routers.api.tls=true"
      - "traefik.http.routers.api.tls.certresolver=myresolver"
```

Or manual with Go TLS:

```go
// Server (API)
cert, _ := tls.LoadX509KeyPair("server.crt", "server.key")
tlsConfig := &tls.Config{
    Certificates: []tls.Certificate{cert},
    MinVersion:   tls.VersionTLS13,
}
server := &http.Server{
    TLSConfig: tlsConfig,
}
server.ListenAndServeTLS("", "")

// Client (Web → API)
client := &http.Client{
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            MinVersion: tls.VersionTLS13,
        },
    },
}
```

---

## CI/CD Pipeline Security

### Current Pipeline Analysis

Your `.github/workflows/ci-cd.yml` has:
- ✅ Tests run on push/PR
- ✅ Go version pinned
- ❌ No security scanning
- ❌ No image building
- ❌ No vulnerability checks
- ❌ No SAST/DAST
- ❌ No secrets scanning

---

### Secure CI/CD Pipeline (Full Example)

```yaml
# .github/workflows/secure-ci-cd.yml
name: Secure CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]
  schedule:
    # Daily security scan
    - cron: '0 2 * * *'

env:
  REGISTRY: ghcr.io
  GO_VERSION: '1.25'
  DOCKER_BUILDKIT: 1

jobs:
  # ============================================
  # Job 1: Security Scanning (SAST)
  # ============================================
  security-scan:
    name: Security Scanning
    runs-on: ubuntu-latest
    permissions:
      contents: read
      security-events: write

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      # Go security scanning
      - name: Run Gosec (Go SAST)
        uses: securego/gosec@master
        with:
          args: '-fmt sarif -out gosec-results.sarif ./...'

      - name: Upload Gosec results to GitHub Security
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: gosec-results.sarif

      # Dependency vulnerability scanning
      - name: Run Trivy vulnerability scanner (filesystem)
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
          format: 'sarif'
          output: 'trivy-results.sarif'
          severity: 'CRITICAL,HIGH'

      - name: Upload Trivy results
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: trivy-results.sarif

      # Secret scanning
      - name: GitGuardian scan
        uses: GitGuardian/ggshield-action@v1
        env:
          GITHUB_PUSH_BEFORE_SHA: ${{ github.event.before }}
          GITHUB_PUSH_BASE_SHA: ${{ github.event.base }}
          GITHUB_DEFAULT_BRANCH: ${{ github.event.repository.default_branch }}
          GITGUARDIAN_API_KEY: ${{ secrets.GITGUARDIAN_API_KEY }}

      # Dockerfile linting
      - name: Lint Dockerfile with Hadolint
        uses: hadolint/hadolint-action@v3.1.0
        with:
          dockerfile: Dockerfile
          failure-threshold: warning

  # ============================================
  # Job 2: Test & Build
  # ============================================
  test-and-build:
    name: Test & Build
    runs-on: ubuntu-latest
    needs: security-scan

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: Download dependencies
        run: go mod download && go mod verify

      - name: Run tests with coverage
        run: |
          go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
          go tool cover -func=coverage.out

      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v4
        with:
          files: ./coverage.out
          flags: unittests

      - name: Build binary
        run: |
          CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
          go build -ldflags="-w -s" -trimpath -o api ./cmd/api

  # ============================================
  # Job 3: Build & Scan Docker Images
  # ============================================
  build-images:
    name: Build & Scan Images
    runs-on: ubuntu-latest
    needs: test-and-build
    permissions:
      contents: read
      packages: write
      security-events: write

    strategy:
      matrix:
        image: [api, web]

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ github.repository }}/${{ matrix.image }}
          tags: |
            type=ref,event=branch
            type=ref,event=pr
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=sha,prefix={{branch}}-

      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: ${{ matrix.image == 'web' && './web' || '.' }}
          file: ${{ matrix.image == 'web' && './web/Dockerfile' || './Dockerfile' }}
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
          platforms: linux/amd64,linux/arm64
          build-args: |
            VERSION=${{ github.sha }}
            BUILD_DATE=${{ github.event.head_commit.timestamp }}

      - name: Run Trivy scanner on image
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: ${{ env.REGISTRY }}/${{ github.repository }}/${{ matrix.image }}:${{ github.sha }}
          format: 'sarif'
          output: 'trivy-image-results.sarif'
          severity: 'CRITICAL,HIGH'

      - name: Upload Trivy image results
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: trivy-image-results.sarif

      # Sign image with Cosign (for supply chain security)
      - name: Install Cosign
        uses: sigstore/cosign-installer@v3

      - name: Sign container image
        run: |
          cosign sign --yes \
            ${{ env.REGISTRY }}/${{ github.repository }}/${{ matrix.image }}@${{ steps.build.outputs.digest }}

      # Generate SBOM
      - name: Generate SBOM with Syft
        uses: anchore/sbom-action@v0
        with:
          image: ${{ env.REGISTRY }}/${{ github.repository }}/${{ matrix.image }}:${{ github.sha }}
          format: spdx-json
          output-file: sbom-${{ matrix.image }}.spdx.json

      - name: Attach SBOM to image
        run: |
          cosign attach sbom --sbom sbom-${{ matrix.image }}.spdx.json \
            ${{ env.REGISTRY }}/${{ github.repository }}/${{ matrix.image }}@${{ steps.build.outputs.digest }}

      - name: Upload SBOM artifact
        uses: actions/upload-artifact@v4
        with:
          name: sbom-${{ matrix.image }}
          path: sbom-${{ matrix.image }}.spdx.json

  # ============================================
  # Job 4: Deploy to Staging
  # ============================================
  deploy-staging:
    name: Deploy to Staging
    runs-on: ubuntu-latest
    needs: build-images
    if: github.ref == 'refs/heads/develop'
    environment:
      name: staging
      url: https://staging.example.com

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Deploy to staging
        run: |
          # Example: SSH to server and pull new images
          ssh deploy@staging.example.com << 'EOF'
            cd /opt/app
            docker compose pull
            docker compose up -d
            docker compose ps
          EOF

  # ============================================
  # Job 5: Production Deployment (Manual Approval)
  # ============================================
  deploy-production:
    name: Deploy to Production
    runs-on: ubuntu-latest
    needs: build-images
    if: github.ref == 'refs/heads/main'
    environment:
      name: production
      url: https://app.example.com

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Verify image signatures
        run: |
          cosign verify \
            --certificate-identity-regexp="https://github.com/${{ github.repository }}" \
            --certificate-oidc-issuer=https://token.actions.githubusercontent.com \
            ${{ env.REGISTRY }}/${{ github.repository }}/api:${{ github.sha }}

      - name: Deploy to production
        run: |
          # Blue-green deployment example
          ssh deploy@prod.example.com << 'EOF'
            cd /opt/app

            # Pull new images
            docker compose -f docker-compose.prod.yml pull

            # Start new containers (blue-green)
            docker compose -f docker-compose.prod.yml up -d --scale api=2

            # Health check
            sleep 10
            curl -f http://localhost:8080/health || exit 1

            # Switch traffic (nginx/load balancer update)
            # ...

            # Remove old containers
            docker compose -f docker-compose.prod.yml up -d --scale api=1 --remove-orphans
          EOF
```

---

### Key CI/CD Security Features Explained

#### 1. Gosec - Go SAST
Finds security issues in Go code:
- SQL injection
- Command injection
- Path traversal
- Weak crypto
- Hardcoded secrets

```bash
# Run locally
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...
```

#### 2. Trivy - Vulnerability Scanner
Scans for:
- OS package vulnerabilities (Alpine APK, Debian APT)
- Go module vulnerabilities
- Misconfigurations in Dockerfiles
- Secrets in code

```bash
# Scan filesystem
trivy fs .

# Scan Docker image
trivy image your-api:latest

# Scan for HIGH/CRITICAL only
trivy image --severity HIGH,CRITICAL your-api:latest
```

#### 3. Hadolint - Dockerfile Linter
Checks Dockerfile best practices:
- DL3006: Always tag base images (not `:latest`)
- DL3008: Pin package versions
- DL4006: Set `SHELL` to pipefail mode
- DL3059: Multiple consecutive `RUN` commands (merge them)

```bash
# Run locally
docker run --rm -i hadolint/hadolint < Dockerfile
```

#### 4. Cosign - Image Signing
Cryptographically sign images to verify authenticity:

```bash
# Sign image
cosign sign --key cosign.key your-api:v1.0.0

# Verify image
cosign verify --key cosign.pub your-api:v1.0.0
```

**Why it matters:** Prevents pulling tampered images in factory environments.

#### 5. SBOM - Software Bill of Materials
Lists all dependencies (for supply chain attacks like Log4j):

```bash
# Generate SBOM
syft your-api:latest -o spdx-json > sbom.json

# Scan SBOM for vulnerabilities
grype sbom:./sbom.json
```

---

## Image Scanning & Vulnerability Management

### Daily Vulnerability Scanning

```yaml
# .github/workflows/daily-scan.yml
name: Daily Security Scan

on:
  schedule:
    - cron: '0 2 * * *'  # 2 AM daily

jobs:
  scan-production-images:
    runs-on: ubuntu-latest
    steps:
      - name: Scan API image
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: ghcr.io/${{ github.repository }}/api:latest
          format: 'table'
          exit-code: '1'  # Fail if vulnerabilities found
          severity: 'CRITICAL,HIGH'

      - name: Notify on Slack if vulnerabilities found
        if: failure()
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          text: '🚨 Critical vulnerabilities found in production images!'
          webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```

---

### Automated Dependency Updates

```yaml
# .github/dependabot.yml
version: 2
updates:
  # Go modules
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5
    reviewers:
      - "your-username"
    labels:
      - "dependencies"
      - "security"

  # Docker base images
  - package-ecosystem: "docker"
    directory: "/"
    schedule:
      interval: "weekly"

  # GitHub Actions
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
```

---

### Vulnerability Remediation Workflow

```
┌─────────────────────────────────────────────────────────┐
│ 1. Daily Trivy scan finds CVE-2024-1234 in alpine:3.19 │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│ 2. Automated PR from Dependabot: Update to alpine:3.20 │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│ 3. CI pipeline runs: tests pass, no new vulnerabilities│
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│ 4. Auto-merge PR (if tests pass)                       │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│ 5. Deploy to staging → smoke tests → production        │
└─────────────────────────────────────────────────────────┘
```

---

## Supply Chain Security

### Verify Base Image Integrity

```dockerfile
# Use SHA256 digest instead of tags
FROM golang:1.25-alpine@sha256:abc123...def456

# Verify checksum
FROM postgres:16-alpine@sha256:789ghi...jkl012
```

**How to get digest:**

```bash
docker pull golang:1.25-alpine
docker inspect --format='{{index .RepoDigests 0}}' golang:1.25-alpine
# Output: golang@sha256:abc123...
```

---

### Content Trust (Docker Notary)

```bash
# Enable Docker Content Trust
export DOCKER_CONTENT_TRUST=1

# Now pulls only verify signed images
docker pull your-api:v1.0.0
# Fails if image not signed
```

---

### Private Registry Mirror

For Agramkow factories (air-gapped environments):

```yaml
# docker-compose.yml with private registry
services:
  api:
    image: registry.agramkow.local/ticket-api:v1.0.0
    # Instead of public ghcr.io
```

**Setup private registry:**

```bash
# Run Docker Registry
docker run -d -p 5000:5000 \
  --restart=always \
  --name registry \
  -v /opt/registry:/var/lib/registry \
  registry:2

# Push images
docker tag your-api:latest localhost:5000/your-api:latest
docker push localhost:5000/your-api:latest

# Pull in factory
docker pull registry.agramkow.local:5000/your-api:latest
```

---

## Production Deployment Strategies

### 1. Blue-Green Deployment

```bash
#!/bin/bash
# deploy-blue-green.sh

# Pull new image (green)
docker pull your-api:v2.0.0

# Start green alongside blue
docker run -d --name api-green -p 8081:8080 your-api:v2.0.0

# Health check
curl -f http://localhost:8081/health || exit 1

# Update load balancer to point to green
nginx -s reload

# Wait for active connections to drain
sleep 30

# Stop blue
docker stop api-blue
docker rm api-blue

# Rename green to blue (for next deployment)
docker rename api-green api-blue
```

---

### 2. Rolling Update (Docker Swarm)

```yaml
# docker-compose.swarm.yml
version: '3.8'

services:
  api:
    image: your-api:latest
    deploy:
      replicas: 3
      update_config:
        parallelism: 1      # Update 1 at a time
        delay: 10s          # Wait 10s between updates
        failure_action: rollback
        monitor: 30s
        max_failure_ratio: 0.3
      rollback_config:
        parallelism: 1
        delay: 5s
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
```

```bash
# Deploy update
docker stack deploy -c docker-compose.swarm.yml ticket-app

# Rollback if issues
docker service rollback ticket-app_api
```

---

### 3. Canary Deployment

```yaml
# Deploy v2 to 10% of traffic
services:
  api-v1:
    image: your-api:v1.0.0
    deploy:
      replicas: 9
      labels:
        - "traefik.http.services.api.loadbalancer.weight=90"

  api-v2:
    image: your-api:v2.0.0
    deploy:
      replicas: 1
      labels:
        - "traefik.http.services.api.loadbalancer.weight=10"
```

Monitor error rates, latency. If good, gradually shift traffic to v2.

---

## Monitoring & Logging

### 1. Centralized Logging

```yaml
# docker-compose.yml with logging
services:
  api:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
        labels: "service,environment"
        tag: "{{.Name}}/{{.ID}}"
```

**Ship logs to ELK/Loki:**

```yaml
# Loki logging driver
services:
  api:
    logging:
      driver: loki
      options:
        loki-url: "http://loki:3100/loki/api/v1/push"
        loki-external-labels: "service=api,environment=production"
```

---

### 2. Security Monitoring (Falco)

Detect runtime threats:

```yaml
# Falco rules for container security
- rule: Unexpected outbound connection
  desc: Detect unexpected network connections
  condition: outbound and container.name = "api" and not allowed_destinations
  output: "Suspicious outbound connection (dest=%fd.sip.name port=%fd.sport)"
  priority: WARNING

- rule: Write below binary dir
  desc: Detect writes to /usr/bin (immutable!)
  condition: container.name = "api" and write and fd.name startswith /usr/bin
  output: "Modification to /usr/bin detected (file=%fd.name)"
  priority: ERROR
```

---

### 3. Metrics & Alerts (Prometheus)

```yaml
# Expose metrics from Go app
import "github.com/prometheus/client_golang/prometheus/promhttp"

http.Handle("/metrics", promhttp.Handler())
```

**Alert rules:**

```yaml
# prometheus-alerts.yml
groups:
  - name: containers
    rules:
      - alert: ContainerHighMemory
        expr: container_memory_usage_bytes{name="api"} > 450000000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "API container using >450MB RAM"

      - alert: ContainerRestarting
        expr: rate(container_restart_count[15m]) > 0
        labels:
          severity: critical
        annotations:
          summary: "Container restarting frequently"
```

---

## Agramkow-Specific Use Case

### Challenge: Isolated Factory Environments

**Requirements:**
1. Software runs on QA machines at production factories
2. Limited/no internet access
3. Must be tamper-proof (can't modify QA test logic)
4. Easy distribution to 100+ factory locations
5. Version control (rollback if QA machine issues)
6. Security compliance for factory network

---

### Proposed Architecture

```
┌──────────────────────────────────────────────────────────┐
│                  Agramkow Central HQ                     │
├──────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────┐    │
│  │  CI/CD Pipeline (GitHub Actions)                │    │
│  │  - Build Docker images                          │    │
│  │  - Security scan (Trivy, Gosec)                 │    │
│  │  - Sign images (Cosign)                         │    │
│  │  - Generate SBOM                                │    │
│  └─────────────────────────────────────────────────┘    │
│                         │                                │
│                         ▼                                │
│  ┌─────────────────────────────────────────────────┐    │
│  │  Private Container Registry                     │    │
│  │  registry.agramkow.com                          │    │
│  │  - Stores signed images                         │    │
│  │  - TLS + authentication required                │    │
│  └─────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────┘
                         │
                         │ (VPN or Direct Connection)
                         │
          ┌──────────────┴───────────────┐
          │                              │
          ▼                              ▼
┌───────────────────────┐      ┌───────────────────────┐
│  Factory #1 (Germany) │      │  Factory #2 (China)   │
├───────────────────────┤      ├───────────────────────┤
│ QA Machine            │      │ QA Machine            │
│ ┌──────────────────┐  │      │ ┌──────────────────┐  │
│ │ Docker Engine    │  │      │ │ Docker Engine    │  │
│ │                  │  │      │ │                  │  │
│ │ ┌──────────────┐ │  │      │ │ ┌──────────────┐ │  │
│ │ │ QA App v3.2  │ │  │      │ │ │ QA App v3.1  │ │  │
│ │ │ (container)  │ │  │      │ │ │ (container)  │ │  │
│ │ └──────────────┘ │  │      │ │ └──────────────┘ │  │
│ │                  │  │      │ │                  │  │
│ │ ┌──────────────┐ │  │      │ │ Read-only FS     │  │
│ │ │ PostgreSQL   │ │  │      │ │ Signed image     │  │
│ │ │ (test data)  │ │  │      │ │ No shell         │  │
│ │ └──────────────┘ │  │      │ └──────────────────┘  │
│ └──────────────────┘  │      └───────────────────────┘
│                       │
│ Fridge/Freezer ←──────┤ (Serial/USB connection)
│ (Test Subject)        │
└───────────────────────┘
```

---

### Deployment Workflow for Factories

#### Step 1: Build & Sign at HQ

```bash
# CI/CD builds secure image
docker build -t registry.agramkow.com/qa-software:3.2.0 \
  --build-arg VERSION=3.2.0 \
  --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
  -f Dockerfile.secure .

# Sign image
cosign sign --key cosign.key registry.agramkow.com/qa-software:3.2.0

# Generate SBOM
syft registry.agramkow.com/qa-software:3.2.0 -o spdx-json > sbom-3.2.0.json

# Push to private registry
docker push registry.agramkow.com/qa-software:3.2.0
```

---

#### Step 2: Factory Pulls Image (Scheduled or Manual)

```bash
# On QA machine at factory
#!/bin/bash
# /opt/agramkow/update.sh

# Verify image signature before pulling
cosign verify --key cosign.pub registry.agramkow.com/qa-software:3.2.0

# Pull new version
docker pull registry.agramkow.com/qa-software:3.2.0

# Stop old version
docker compose -f /opt/agramkow/docker-compose.yml down

# Start new version
docker compose -f /opt/agramkow/docker-compose.yml up -d

# Health check
sleep 10
curl -f http://localhost:8080/health || docker compose logs
```

---

#### Step 3: Factory docker-compose.yml

```yaml
# /opt/agramkow/docker-compose.yml
version: '3.8'

services:
  qa-app:
    image: registry.agramkow.com/qa-software:${VERSION:-3.2.0}
    container_name: agramkow-qa

    # Security hardening
    read_only: true
    cap_drop:
      - ALL
    security_opt:
      - no-new-privileges:true

    # Resource limits (QA machine specs)
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 2G
        reservations:
          cpus: '0.5'
          memory: 512M

    # Isolated network
    networks:
      - qa_network

    # USB/serial device access for fridge connection
    devices:
      - /dev/ttyUSB0:/dev/ttyUSB0  # Serial port

    # Writable volumes (test results, logs)
    tmpfs:
      - /tmp:size=100M
    volumes:
      - qa_results:/app/results:rw
      - qa_logs:/app/logs:rw
      - /etc/localtime:/etc/localtime:ro  # Match factory timezone

    # Secrets from file (not environment variables)
    secrets:
      - db_password
      - api_key

    environment:
      - TZ=Europe/Berlin  # Factory timezone
      - DB_HOST=db
      - DB_PASSWORD_FILE=/run/secrets/db_password
      - API_KEY_FILE=/run/secrets/api_key
      - FACTORY_ID=${FACTORY_ID}  # Unique per factory

    depends_on:
      db:
        condition: service_healthy

    restart: unless-stopped

    healthcheck:
      test: ["CMD", "/app/qa-app", "healthcheck"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s

  db:
    image: postgres:16-alpine@sha256:abc123...  # Pinned digest
    container_name: agramkow-db

    read_only: true
    cap_drop:
      - ALL
    cap_add:
      - CHOWN
      - DAC_OVERRIDE
      - FOWNER
      - SETGID
      - SETUID

    tmpfs:
      - /tmp:size=100M
      - /run/postgresql:size=10M

    volumes:
      - postgres_data:/var/lib/postgresql/data:rw

    networks:
      - qa_network

    secrets:
      - db_password

    environment:
      POSTGRES_USER: qauser
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
      POSTGRES_DB: qa_results

    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U qauser -d qa_results"]
      interval: 5s
      timeout: 3s
      retries: 5

networks:
  qa_network:
    driver: bridge
    internal: false  # Allow outbound for updates, or true for full isolation

volumes:
  postgres_data:
    driver: local
  qa_results:
    driver: local
  qa_logs:
    driver: local

secrets:
  db_password:
    file: /opt/agramkow/secrets/db_password.txt
  api_key:
    file: /opt/agramkow/secrets/api_key.txt
```

---

#### Step 4: Initial Factory Setup Script

```bash
#!/bin/bash
# setup-factory.sh - Run once on new QA machine

set -e

FACTORY_ID=$1
REGISTRY_URL="registry.agramkow.com"

if [ -z "$FACTORY_ID" ]; then
  echo "Usage: $0 <FACTORY_ID>"
  exit 1
fi

echo "Setting up Agramkow QA software for Factory: $FACTORY_ID"

# Create directories
mkdir -p /opt/agramkow/{secrets,config,backups}
cd /opt/agramkow

# Generate secrets
openssl rand -hex 32 > secrets/db_password.txt
chmod 600 secrets/db_password.txt

# Copy API key (provided by HQ)
echo "$API_KEY" > secrets/api_key.txt
chmod 600 secrets/api_key.txt

# Create .env file
cat > .env << EOF
VERSION=3.2.0
FACTORY_ID=$FACTORY_ID
EOF

# Download docker-compose.yml from HQ
curl -H "Authorization: Bearer $DEPLOY_TOKEN" \
  https://deploy.agramkow.com/docker-compose.yml \
  -o docker-compose.yml

# Download cosign public key for verification
curl https://deploy.agramkow.com/cosign.pub -o cosign.pub

# Login to private registry
docker login $REGISTRY_URL -u factory -p "$REGISTRY_PASSWORD"

# Pull initial image
cosign verify --key cosign.pub $REGISTRY_URL/qa-software:3.2.0
docker pull $REGISTRY_URL/qa-software:3.2.0

# Start services
docker compose up -d

# Wait for health checks
sleep 15

# Verify system
curl -f http://localhost:8080/health || {
  echo "Health check failed!"
  docker compose logs
  exit 1
}

echo "✅ Factory $FACTORY_ID setup complete!"
echo "QA software running on http://localhost:8080"

# Setup automatic updates (daily at 2 AM)
cat > /etc/cron.d/agramkow-update << 'CRON'
0 2 * * * root /opt/agramkow/update.sh >> /var/log/agramkow-update.log 2>&1
CRON
```

---

### Air-Gapped Factory Variant

For factories with **zero internet access**, ship updates on USB:

```bash
#!/bin/bash
# create-offline-package.sh (run at HQ)

VERSION="3.2.0"
OUTPUT_DIR="agramkow-update-$VERSION"

mkdir -p $OUTPUT_DIR

# Save Docker image to tarball
docker save registry.agramkow.com/qa-software:$VERSION -o $OUTPUT_DIR/qa-software.tar

# Copy docker-compose.yml
cp docker-compose.yml $OUTPUT_DIR/

# Copy signature and verification key
cosign sign-blob --key cosign.key $OUTPUT_DIR/qa-software.tar > $OUTPUT_DIR/qa-software.tar.sig
cp cosign.pub $OUTPUT_DIR/

# Copy install script
cat > $OUTPUT_DIR/install.sh << 'EOF'
#!/bin/bash
set -e

# Verify signature
cosign verify-blob --key cosign.pub --signature qa-software.tar.sig qa-software.tar

# Load image
docker load -i qa-software.tar

# Update version in .env
sed -i 's/VERSION=.*/VERSION=3.2.0/' /opt/agramkow/.env

# Restart services
cd /opt/agramkow
docker compose up -d

echo "✅ Update to 3.2.0 complete!"
EOF

chmod +x $OUTPUT_DIR/install.sh

# Create checksums
cd $OUTPUT_DIR
sha256sum * > SHA256SUMS

# Create signed archive
cd ..
tar czf agramkow-update-$VERSION.tar.gz $OUTPUT_DIR
gpg --sign --armor agramkow-update-$VERSION.tar.gz

echo "📦 Offline update package ready: agramkow-update-$VERSION.tar.gz"
echo "Copy to USB and run install.sh on factory QA machine"
```

---

### Security Benefits for Agramkow

| Security Feature | Benefit | Attack Prevented |
|-----------------|---------|------------------|
| **Signed images** | Verify authenticity | Tampered/malicious images |
| **Read-only filesystem** | Immutable code | Malware installation |
| **Non-root user** | Limited privileges | Container escape |
| **No shell in image** | Can't execute arbitrary commands | Reverse shell |
| **Network isolation** | Limited attack surface | Lateral movement |
| **Resource limits** | Stability | Resource exhaustion DoS |
| **Secrets from files** | Not in git/logs | Secret exposure |
| **SBOM tracking** | Know all dependencies | Supply chain attacks |
| **Vulnerability scanning** | Catch CVEs before deploy | Known exploits |
| **Versioned deployments** | Rollback capability | Bad updates |

---

## Implementation Checklist

### Phase 1: Immediate Wins (1 Week)

- [ ] Create `.dockerignore` for Go app
```gitignore
.git
.env
*.md
.github
.vscode
*.log
```

- [ ] Secure Go Dockerfile with multi-stage build
  - [ ] Use distroless base image
  - [ ] Run as non-root user
  - [ ] Static binary with CGO_ENABLED=0

- [ ] Move secrets to `.env` file
  - [ ] Add `.env` to `.gitignore`
  - [ ] Create `.env.example` template
  - [ ] Update docker-compose.yml to use `${VAR}` syntax

- [ ] Add resource limits to docker-compose.yml
  - [ ] CPU limits (1 core per service)
  - [ ] Memory limits (512MB API, 1GB DB)

- [ ] Implement read-only filesystem
  - [ ] `read_only: true` for all services
  - [ ] Add tmpfs mounts for /tmp

### Phase 2: Pipeline Security (1 Week)

- [ ] Add security scanning to CI/CD
  - [ ] Gosec for Go SAST
  - [ ] Trivy for vulnerability scanning
  - [ ] Hadolint for Dockerfile linting

- [ ] Build Docker images in CI
  - [ ] Multi-platform builds (amd64, arm64)
  - [ ] Tag with SHA and semver
  - [ ] Push to GitHub Container Registry

- [ ] Add test coverage reporting
  - [ ] Run tests with `-race` flag
  - [ ] Upload coverage to Codecov

- [ ] Implement branch protection
  - [ ] Require PR reviews
  - [ ] Require passing CI
  - [ ] Require up-to-date branches

### Phase 3: Advanced Security (2 Weeks)

- [ ] Image signing with Cosign
  - [ ] Generate signing key
  - [ ] Sign images in CI
  - [ ] Verify signatures before deployment

- [ ] SBOM generation
  - [ ] Use Syft to generate SBOMs
  - [ ] Attach to images
  - [ ] Store as artifacts

- [ ] Implement Docker secrets
  - [ ] Set up Docker Swarm (or use Compose v2 secrets)
  - [ ] Migrate from env vars to secrets
  - [ ] Update Go app to read from `/run/secrets/`

- [ ] Network segmentation
  - [ ] Create frontend/backend networks
  - [ ] Remove exposed DB port
  - [ ] Internal-only backend network

- [ ] Drop capabilities
  - [ ] `cap_drop: ALL` for all services
  - [ ] Test and add back only necessary caps

### Phase 4: Production Readiness (2 Weeks)

- [ ] Set up private container registry
  - [ ] Docker Registry or Harbor
  - [ ] TLS certificates
  - [ ] Authentication

- [ ] Implement deployment strategy
  - [ ] Blue-green or rolling updates
  - [ ] Automated rollback on failure
  - [ ] Health checks before traffic switch

- [ ] Monitoring & logging
  - [ ] Prometheus metrics
  - [ ] Centralized logging (Loki/ELK)
  - [ ] Alert rules for security events

- [ ] Runtime security
  - [ ] Falco for runtime threat detection
  - [ ] AppArmor/SELinux profiles
  - [ ] Seccomp profiles

- [ ] Disaster recovery
  - [ ] Database backups
  - [ ] Image backups
  - [ ] Restore testing

### Exam Prep Checklist

**Docker Security Concepts to Know:**

- [ ] Namespaces (PID, NET, IPC, MNT, UTS, USER)
- [ ] Cgroups (resource limits)
- [ ] Capabilities (Linux permissions model)
- [ ] Seccomp (syscall filtering)
- [ ] AppArmor/SELinux (MAC - Mandatory Access Control)
- [ ] Rootless Docker
- [ ] User namespaces
- [ ] Image signing & verification
- [ ] Supply chain security (SBOM, SLSA framework)
- [ ] Secrets management (Docker secrets, Vault)
- [ ] Network security (bridge, overlay, macvlan)
- [ ] Attack surface reduction (distroless, scratch)

**Hands-On Skills to Practice:**

- [ ] Build multi-stage Dockerfiles
- [ ] Scan images with Trivy/Grype/Snyk
- [ ] Sign images with Cosign
- [ ] Configure security profiles (AppArmor, Seccomp)
- [ ] Set up private registry with TLS
- [ ] Implement secrets rotation
- [ ] Configure network policies
- [ ] Debug container security (docker inspect, docker diff)
- [ ] Audit container runtime (docker events, falco)
- [ ] Perform vulnerability remediation

**Common Exam Questions:**

1. **Why run containers as non-root?**
   - Answer: Reduces blast radius of container escape; if attacker breaks out, they're unprivileged user on host, not root.

2. **What's the difference between CMD and ENTRYPOINT?**
   - CMD: Default args (can override with `docker run`)
   - ENTRYPOINT: Always runs (args append)
   - Best practice: Use both (`ENTRYPOINT ["app"]`, `CMD ["--help"]`)

3. **How do you secure secrets?**
   - Answer: Docker secrets (encrypted at rest, in-memory at runtime), never in Dockerfile/env vars, use external vaults (Vault, AWS Secrets Manager), rotate regularly.

4. **What's the attack surface of alpine vs distroless vs scratch?**
   - Alpine: ~5MB, has shell & package manager (attack surface: medium)
   - Distroless: ~2MB, no shell (attack surface: low)
   - Scratch: <1MB, literally empty (attack surface: minimal)

5. **How do you verify image hasn't been tampered with?**
   - Answer: Content Trust (Docker Notary), Cosign signatures, SHA256 digests, private registry with access controls.

6. **What's the principle of least privilege for containers?**
   - Answer: Drop all capabilities, add only necessary; read-only FS; non-root user; no privileged mode; network isolation; resource limits.

---

## Tools & Resources

### Essential Tools

```bash
# Install security tools
brew install cosign syft grype trivy hadolint
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

### Useful Commands Cheat Sheet

```bash
# ==== Image Analysis ====
# List image layers
docker history your-api:latest

# Inspect image config
docker inspect your-api:latest | jq '.[0].Config'

# Check image size
docker images your-api:latest --format "{{.Size}}"

# ==== Security Scanning ====
# Scan with Trivy
trivy image --severity HIGH,CRITICAL your-api:latest

# Scan Go code with Gosec
gosec ./...

# Lint Dockerfile
hadolint Dockerfile

# Generate SBOM
syft your-api:latest -o spdx-json

# Scan SBOM for vulnerabilities
grype sbom:./sbom.json

# ==== Runtime Inspection ====
# Check running user
docker exec -it container_name whoami

# Check capabilities
docker exec -it container_name cat /proc/1/status | grep Cap

# Check filesystem changes
docker diff container_name

# Check exposed ports
docker port container_name

# ==== Network Security ====
# Inspect networks
docker network inspect bridge

# Check iptables rules
sudo iptables -L DOCKER

# Monitor traffic (requires host access)
sudo tcpdump -i docker0

# ==== Secrets Management ====
# Create secret
echo "secret_value" | docker secret create my_secret -

# List secrets
docker secret ls

# Inspect secret (metadata only, not value!)
docker secret inspect my_secret

# ==== Image Signing ====
# Generate key pair
cosign generate-key-pair

# Sign image
cosign sign --key cosign.key your-api:latest

# Verify signature
cosign verify --key cosign.pub your-api:latest

# ==== Cleanup ====
# Remove unused images
docker image prune -a

# Remove stopped containers
docker container prune

# Full cleanup (BE CAREFUL!)
docker system prune -a --volumes
```

---

## Conclusion

### For Your Exam

**Focus on these topics:**
1. **Why** security practices matter (blast radius, defense in depth, least privilege)
2. **How** to implement them (Dockerfile best practices, runtime configs)
3. **Trade-offs** (security vs convenience, image size vs debugging)
4. **Agramkow use case** - explain how each practice applies to isolated factory deployments

### For Your Project

**Priority order:**
1. ✅ Secure Dockerfile (multi-stage, non-root, distroless)
2. ✅ Secrets management (.env → Docker secrets)
3. ✅ CI/CD security scanning (Trivy, Gosec)
4. ✅ Network segmentation
5. ✅ Resource limits & hardening

### Key Takeaway

Docker security is **layered defense**:
```
Image Security (build time)
    ↓
Runtime Security (container config)
    ↓
Network Security (isolation)
    ↓
Secrets Security (credential management)
    ↓
Monitoring & Response (detection)
```

Each layer reduces risk. No single layer is perfect, but together they create a robust security posture suitable for production deployments in sensitive environments like Agramkow's factory QA systems.

**Good luck on your exam!** 🎓

---

**Next Steps:**
1. Implement Phase 1 checklist this week
2. Create practice scenarios (simulate container escape, test rollbacks)
3. Document your setup for exam reference
4. Review OWASP Docker Security Cheat Sheet: https://cheatsheetseries.owasp.org/cheatsheets/Docker_Security_Cheat_Sheet.html
