# CI/CD Setup Guide

## Overview
This project now includes a complete CI/CD pipeline using GitHub Actions that:
- Runs tests for both API and Web applications
- Builds Docker images
- Pushes images to GitHub Container Registry (GHCR)
- Deploys to production automatically on main branch updates

## What's Been Added

### 1. Docker Configuration
- **[web/Dockerfile](web/Dockerfile)** - Multi-stage build for Nuxt web application
- **[docker-compose.yml](docker-compose.yml)** - Updated to include web service for local development
- **[docker-compose.prod.yml](docker-compose.prod.yml)** - Production configuration using registry images

### 2. GitHub Actions Workflow
- **[.github/workflows/ci-cd.yml](.github/workflows/ci-cd.yml)** - Complete CI/CD pipeline

The workflow includes 4 jobs:
1. **test-api** - Runs Go tests with coverage
2. **test-web** - Runs npm tests and builds
3. **build-and-push** - Builds and pushes Docker images to GHCR
4. **deploy** - Deploys to production server via SSH

## Setup Steps

### Step 1: Add Tests to Your Project

#### Go API Tests
Create test files for your Go code:

```bash
# Example: internal/handlers/event_handler_test.go
touch internal/handlers/event_handler_test.go
touch internal/handlers/booking_handler_test.go
touch internal/repository/event_repository_test.go
touch internal/repository/booking_repository_test.go
```

Basic test structure:
```go
package handlers

import "testing"

func TestEventHandler(t *testing.T) {
    // Your test code
    t.Run("should create event", func(t *testing.T) {
        // Test implementation
    })
}
```

#### Web Application Tests
Update [web/package.json](web/package.json) to add test script:

```json
{
  "scripts": {
    "test": "bun test"
  }
}
```

Or for a placeholder until you add actual tests:

```json
{
  "scripts": {
    "test": "echo 'No tests yet' && exit 0"
  }
}
```

Later, add actual tests using Bun's built-in test runner or Vitest.

### Step 2: Configure GitHub Repository Secrets

Go to your GitHub repository → Settings → Secrets and variables → Actions

Add the following secrets:

#### For Deployment (Required only if deploying)
- **DEPLOY_HOST** - Your production server IP or hostname
- **DEPLOY_USER** - SSH username for deployment
- **DEPLOY_SSH_KEY** - Private SSH key for authentication
- **DEPLOY_PATH** - Path on server where project is located (e.g., `/home/user/app`)

#### For Database (Production)
- **PROD_DB_PASSWORD** - Production database password
- **PROD_JWT_SECRET** - Production JWT secret (min 32 characters)

### Step 3: Enable GitHub Container Registry

1. The workflow uses GitHub Container Registry (GHCR) automatically
2. Images will be pushed to: `ghcr.io/YOUR_USERNAME/golang_test/api` and `ghcr.io/YOUR_USERNAME/golang_test/web`
3. No additional setup needed - GHCR is free for public repositories

### Step 4: Update Production docker-compose

Update the image names in [docker-compose.prod.yml](docker-compose.prod.yml):

```yaml
api:
  image: ghcr.io/YOUR_USERNAME/golang_test/api:latest

web:
  image: ghcr.io/YOUR_USERNAME/golang_test/web:latest
```

Replace `YOUR_USERNAME` with your actual GitHub username.

### Step 5: Set Up Production Server

On your production server:

```bash
# 1. Install Docker and Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# 2. Create deployment directory
mkdir -p /home/user/app
cd /home/user/app

# 3. Copy docker-compose.prod.yml and .env files
# (You can do this via git clone or scp)

# 4. Create .env file with production secrets
cat > .env << EOF
DB_USER=user
DB_PASSWORD=your-secure-password
DB_NAME=ticket_db
JWT_SECRET=your-super-secret-jwt-key-min-32-chars
API_IMAGE=ghcr.io/YOUR_USERNAME/golang_test/api:latest
WEB_IMAGE=ghcr.io/YOUR_USERNAME/golang_test/web:latest
API_BASE_URL=https://your-domain.com/api
EOF

# 5. Set up SSH key for GitHub Actions
# Add the public key to ~/.ssh/authorized_keys
```

### Step 6: Configure Nginx (Optional)

For production, set up Nginx as reverse proxy:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Usage

### Local Development

```bash
# Start all services (db, api, web)
docker compose up

# Build and start
docker compose up --build

# Stop services
docker compose down
```

Access:
- Web: http://localhost:3000
- API: http://localhost:8080
- Database: localhost:5432

### Production Deployment

Deployment happens automatically when you push to the `main` branch:

```bash
git add .
git commit -m "Your changes"
git push origin main
```

The GitHub Actions workflow will:
1. Run all tests
2. Build Docker images
3. Push to GitHub Container Registry
4. SSH into your server and deploy

### Manual Deployment

If you need to deploy manually:

```bash
# On production server
cd /home/user/app

# Pull latest images
docker compose -f docker-compose.prod.yml pull

# Restart services
docker compose -f docker-compose.prod.yml up -d

# View logs
docker compose -f docker-compose.prod.yml logs -f
```

## Workflow Triggers

The CI/CD pipeline runs on:
- **Pull Requests** to `main` - Runs tests only
- **Push to `main`** - Runs tests, builds images, and deploys
- **Push to `develop`** - Runs tests and builds images (no deploy)

## Monitoring Deployments

1. Go to your GitHub repository → Actions
2. Click on the latest workflow run
3. View logs for each job
4. Check for any errors

## Troubleshooting

### Tests Failing
- Check the test logs in GitHub Actions
- Run tests locally: `go test ./...` or `cd web && npm test`

### Build Failing
- Verify Dockerfiles are correct
- Check that all dependencies are properly defined
- Test build locally: `docker compose build`

### Deployment Failing
- Verify SSH connection: `ssh DEPLOY_USER@DEPLOY_HOST`
- Check server disk space: `df -h`
- Verify secrets are set correctly in GitHub
- Check server logs: `docker compose logs`

### Images Not Pulling
- Ensure GHCR permissions are set correctly
- Make packages public in GitHub repository settings
- Verify image names match in docker-compose.prod.yml

## Next Steps

1. **Add actual tests** - Replace placeholder tests with real test coverage
2. **Set up monitoring** - Add application monitoring (e.g., Prometheus, Grafana)
3. **Add health checks** - Implement `/health` endpoints
4. **Set up SSL** - Use Let's Encrypt with Certbot
5. **Database backups** - Implement automated backup strategy
6. **Staging environment** - Create a staging branch/environment
7. **Rollback strategy** - Tag images for easy rollbacks

## Advanced Configuration

### Using Different Registries

To use Docker Hub instead of GHCR, update [.github/workflows/ci-cd.yml](.github/workflows/ci-cd.yml):

```yaml
env:
  REGISTRY: docker.io
  API_IMAGE_NAME: your-username/api
  WEB_IMAGE_NAME: your-username/web
```

### Adding More Environments

Create separate workflow files:
- `.github/workflows/staging.yml` - For staging deployments
- `.github/workflows/production.yml` - For production

### Database Migrations

Add migration step to deployment:

```yaml
- name: Run migrations
  run: |
    docker compose exec -T api ./migrate up
```

## Security Best Practices

1. Never commit `.env` files - they're in `.gitignore`
2. Use strong passwords for production
3. Rotate JWT secrets regularly
4. Keep Docker images updated
5. Use secrets for sensitive data in GitHub Actions
6. Enable branch protection rules
7. Require PR reviews before merging to main

## Cost Considerations

- **GitHub Actions**: 2,000 minutes/month free for public repos
- **GHCR Storage**: 500MB free, then paid
- **Server**: You'll need a VPS (DigitalOcean, AWS, etc.)

## Support

If you encounter issues:
1. Check GitHub Actions logs
2. Review Docker logs: `docker compose logs`
3. Verify all secrets are configured
4. Test locally before pushing
