# Docker & Pipeline Security Analysis

## Overview

This project uses a comprehensive Docker-based CI/CD pipeline with GitHub Actions for automated building, testing, and security scanning. The pipeline builds multi-architecture container images and pushes them to GitHub Container Registry (ghcr.io).

**Pipeline Workflow**:
- Automated builds on push to `main` or `develop` branches
- Multi-stage Docker builds for minimal image size
- Security scanning with Gosec (SAST) and Trivy (vulnerability scanning)
- Multi-architecture support (linux/amd64, linux/arm64)
- Automated image pushing to GitHub Container Registry

---

## How to Access Built Images

### Registry Information

**Registry**: GitHub Container Registry (ghcr.io)

**Image Locations**:
- API: `ghcr.io/amorags/golangtickets/api`
- Web: `ghcr.io/amorags/golangtickets/web`

### Access Methods

**1. GitHub Packages UI**:
https://github.com/amorags/golangtickets/packages

**2. Pull Commands**:
```bash
# Pull latest images
docker pull ghcr.io/amorags/golangtickets/api:latest
docker pull ghcr.io/amorags/golangtickets/web:latest

# Pull specific version by SHA
docker pull ghcr.io/amorags/golangtickets/api:main-2cfd7e3
docker pull ghcr.io/amorags/golangtickets/web:main-2cfd7e3

# Pull by branch
docker pull ghcr.io/amorags/golangtickets/api:develop
```

**3. Deploy Production**:
```bash
# Pull latest images and deploy
make prod-pull
make prod-deploy

# Or manually
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
```

**Note**: Images are built and pushed only on push events (not on pull requests). Pull requests build images for testing but do not push to the registry.

---

## Current Security Strengths

The project demonstrates several strong security practices:

### Docker Image Security
- **Multi-stage builds**: Separates build and runtime environments, reducing attack surface
- **Distroless base image** (API): Uses `gcr.io/distroless/static-debian12:nonroot` - no shell, no package manager
- **Alpine base** (Web): Minimal footprint with regular security updates
- **Non-root execution**: All containers run as unprivileged users (UID 65532 for API, UID 1001 for Web)
- **Static binary compilation**: Go API compiled with `CGO_ENABLED=0` for maximum portability
- **Security compiler flags**: Uses `-ldflags="-w -s"` and `-trimpath` to reduce information leakage

### Container Runtime Security
- **Capability dropping**: All Linux capabilities dropped with `cap_drop: ALL`, minimal capabilities added back
- **No-new-privileges**: Prevents privilege escalation attacks
- **Resource limits**: CPU, memory, and PID limits prevent resource exhaustion
- **Ulimits configured**: File descriptor and process limits enforced
- **Health checks**: Automated container health monitoring

### CI/CD Security
- **Automated security scanning**:
  - **Gosec** (SAST): Static analysis for Go code vulnerabilities
  - **Trivy filesystem**: Scans source code and dependencies for vulnerabilities
  - **Trivy image**: Post-build image scanning for CVE detection
- **SARIF upload**: Security scan results uploaded to GitHub Security tab for tracking
- **Build caching**: GitHub Actions cache reduces build time and bandwidth
- **Dependency verification**: Frozen lockfiles (`go.sum`, `bun.lockb`) ensure reproducible builds

---

## Security Gaps

### CRITICAL Severity

#### 1. No Image Signing
**Issue**: Container images are not cryptographically signed.

**Risk**: Without signing, there's no way to verify that an image pulled from the registry hasn't been tampered with. An attacker who compromises the registry or performs a man-in-the-middle attack could inject malicious images.

**Missing**:
- Cosign signatures (CNCF standard)
- Docker Content Trust (DCT) with Notary
- Signature verification in deployment pipelines

**Impact**: Cannot guarantee image authenticity, integrity, or provenance. Violates supply chain security best practices.

---

#### 2. Secret Management in Plaintext
**Issue**: Sensitive credentials stored in plaintext `.env` files.

**Risk**:
- JWT_SECRET exposed in `.env.example` (even if placeholder)
- Database credentials in plaintext
- No rotation mechanism
- Risk of accidental commit to version control
- Secrets visible in container environment variables

**Missing**:
- Secrets manager integration (GitHub Secrets, HashiCorp Vault, AWS Secrets Manager)
- Encrypted secret storage
- Secret rotation policies
- Runtime secret injection

**Impact**: Credential compromise leads to authentication bypass, unauthorized database access, and potential data breach.

---

#### 3. Mutable "latest" Tag Usage
**Issue**: Production deployment uses `latest` tag, which is mutable and unpredictable.

**Risk**:
- Unpredictable deployments - "latest" can change at any time
- Rollback difficulty - cannot guarantee which version is deployed
- No audit trail of what's running in production
- Cache invalidation issues across nodes

**Missing**:
- Immutable SHA-based tags for production (e.g., `main-2cfd7e3`)
- Semantic versioning (e.g., `v1.2.3`)
- Tag enforcement policies

**Impact**: Production deployments are non-reproducible and unreliable. Violates infrastructure-as-code principles.

---

### HIGH Severity

#### 4. No SBOM (Software Bill of Materials) Generation
**Issue**: No inventory of software components and dependencies in container images.

**Risk**:
- Cannot track which dependencies are in production
- Slow response to vulnerability disclosures (e.g., Log4Shell-type events)
- No visibility into transitive dependencies
- Compliance violations (Executive Order 14028 for government contractors)

**Missing**:
- SBOM generation during build (SPDX or CycloneDX format)
- SBOM attestation and signing
- SBOM distribution with images
- Automated SBOM scanning for vulnerabilities

**Impact**: Blind to supply chain risks. Cannot quickly identify affected systems when vulnerabilities are disclosed.

---

#### 5. Network Isolation Gaps
**Issue**: All services run on the default Docker bridge network with no segmentation.

**Risk**:
- Web frontend can directly connect to database (violates defense-in-depth)
- No east-west traffic filtering between containers
- Compromised container can pivot to attack others
- No microsegmentation

**Missing**:
- Custom Docker networks with service isolation
- Network policies restricting inter-service communication
- Database network isolation (db should only accept connections from api)
- Ingress/egress filtering

**Impact**: Lateral movement is trivial if any container is compromised. A web application vulnerability could lead to direct database compromise.

---

#### 6. Read-Only Filesystem Not Enforced
**Issue**: Containers run with writable root filesystems.

**Risk**:
- Malware can persist on filesystem
- Attackers can modify binaries or inject backdoors
- Log tampering possible
- Runtime modifications undetected

**Missing**:
- `read_only: true` in docker-compose files
- `tmpfs` mounts for writable directories (/tmp, /var/run)
- Immutable infrastructure enforcement

**Impact**: Container compromise can lead to persistent infections. Violates immutable infrastructure principles.

---

#### 7. No Automated Dependency Scanning
**Issue**: Missing automated tools for dependency vulnerability management.

**Risk**:
- Outdated dependencies with known CVEs
- No alerts for new vulnerabilities
- Manual dependency updates only
- Slow patching cycle

**Missing**:
- Dependabot for automated pull requests
- Renovate or similar dependency update automation
- Continuous vulnerability monitoring
- SLA for security patch deployment

**Impact**: Known vulnerabilities remain unpatched. Increases window of exploitation.

---

#### 8. Supply Chain Security (SLSA) Missing
**Issue**: No SLSA (Supply-chain Levels for Software Artifacts) attestations or provenance tracking.

**Risk**:
- Cannot verify build integrity
- Build process could be compromised without detection
- No tamper evidence for built artifacts
- Cannot prove "what was built from what source"

**Missing**:
- SLSA provenance generation
- Build attestations (who built, when, from what commit)
- Provenance verification in deployment
- Build environment isolation and verification

**Impact**: No guarantee that deployed images were built from trusted source code. Build pipeline compromise could inject backdoors undetected.

---

### MEDIUM Severity

#### 9. Security Context (MAC) Not Configured
**Issue**: No Mandatory Access Control (MAC) profiles applied to containers.

**Risk**:
- Containers rely solely on DAC (Discretionary Access Control)
- No kernel-level confinement beyond capabilities
- Limited defense against kernel exploits
- No default-deny security policy

**Missing**:
- AppArmor profiles
- SELinux contexts
- Custom security profiles for container workloads

**Impact**: Reduced defense-in-depth. Kernel vulnerabilities or capability bypasses have fewer mitigations.

---

#### 10. Log Sanitization Absent
**Issue**: No filtering of sensitive data from container logs.

**Risk**:
- Secrets, tokens, or PII could leak into logs
- Logs may be stored unencrypted
- Logs accessible to broad set of users
- Compliance violations (GDPR, PCI-DSS)

**Missing**:
- Log filtering/redaction middleware
- Structured logging with sensitive field exclusion
- Log encryption at rest
- Log access controls

**Impact**: Credential or data exposure through log aggregation systems. Audit trail could reveal sensitive information.

---

#### 11. Image Provenance Verification Missing
**Issue**: Cannot cryptographically verify the source and build context of images.

**Risk**:
- Cannot prove image was built from a specific commit
- No verification of build environment integrity
- Cannot detect tampering during build process
- Trust relies solely on registry access controls

**Missing**:
- In-toto attestations
- Build provenance metadata
- Verification in deployment pipeline
- Chain of custody documentation

**Impact**: Trust is based on registry security alone. Compromised build systems or insider threats could inject malicious code.

---

### LOW Severity

#### 12. Temporary Filesystem (tmpfs) Not Configured
**Issue**: No tmpfs mounts for ephemeral storage in containers.

**Risk**:
- Minor performance degradation for temporary files
- Temporary data persists on disk longer than necessary
- Slightly larger attack surface for forensic analysis

**Missing**:
- tmpfs mounts for /tmp, /var/run, /var/cache
- Size limits on temporary filesystems

**Impact**: Minimal. Primarily a defense-in-depth measure when combined with read-only root filesystem.

---

#### 13. Multi-Stage Build Secrets Exposure
**Issue**: Build-time secrets (if needed) not mounted securely.

**Risk**:
- Secrets could leak into intermediate build layers
- Docker history could expose sensitive build arguments
- Build cache could retain credentials

**Missing**:
- BuildKit `--mount=type=secret` for build secrets
- Build-time secret rotation
- Scan for secrets in image layers

**Impact**: Currently low (no build secrets in use), but architectural gap if future builds require credentials (private packages, API keys for builds).

---

## Quick Reference

### Deployment Commands

```bash
# Development
make dev                    # Build and start dev environment
make up                     # Start services (detached)
make down                   # Stop all services
make logs                   # View logs

# Production
make prod-pull              # Pull latest images from ghcr.io
make prod-deploy            # Pull and deploy production
make prod-logs              # View production logs

# Testing & Security
make test                   # Run all tests
make test-coverage          # Generate coverage reports
```

### Useful Links

- **GitHub Packages**: https://github.com/amorags/golangtickets/packages
- **GitHub Actions**: Check `.github/workflows/ci-cd.yml` for pipeline configuration
- **Security Scanning**: GitHub Security tab for SARIF results
- **Docker Best Practices**: [Docker Security Docs](https://docs.docker.com/engine/security/)
- **SLSA Framework**: https://slsa.dev/
- **Sigstore/Cosign**: https://docs.sigstore.dev/

### Image Tagging Strategy

Current pipeline creates these tags:
- `latest` - Latest build from default branch (main)
- `main` / `develop` - Branch-based tags
- `main-<sha>` - Immutable SHA-based tags (e.g., `main-2cfd7e3`)

**Recommendation**: Use SHA-based tags for production deployments for reproducibility.

---

## Presentation Notes (15-Minute Exam)

### 7-Minute Presentation Outline

1. **Overview** (1 min): Docker setup, CI/CD with GitHub Actions, GHCR registry
2. **Security Strengths** (2 min): Highlight distroless images, multi-stage builds, automated scanning, capability dropping
3. **Critical Gaps** (3 min): Focus on image signing, secret management, mutable tags
4. **High Severity Issues** (1 min): Briefly mention SBOM, network isolation, read-only filesystem

### 7-Minute Q&A Preparation

**Likely Questions**:
- **"Why no image signing?"** → Cost/complexity vs risk. Plan to implement Cosign with SLSA provenance.
- **"How do you handle secrets now?"** → .env files (acknowledge weakness). Migration path: GitHub Secrets → Vault.
- **"Latest tag concerns?"** → Development convenience. Production should use SHA tags for immutability.
- **"Network isolation - why important?"** → Defense-in-depth. Prevents lateral movement post-compromise.
- **"What's the priority fix?"** → Secret management (highest impact), then image signing (supply chain).

**Strengths to Emphasize**:
- Already using distroless (advanced practice)
- Multi-stage builds show security awareness
- Automated scanning integrated into CI/CD
- Non-root execution enforced
- Resource limits prevent DoS

**Trade-offs to Acknowledge**:
- Some gaps (SBOM, SLSA) require significant tooling investment
- Balance between security and development velocity
- Prioritized practical security over theoretical perfection
