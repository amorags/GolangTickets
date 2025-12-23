# 1. Use --platform=$BUILDPLATFORM to force the build stages to run on the 
#    native GitHub Runner (Fast AMD64), regardless of the target output.
FROM --platform=$BUILDPLATFORM oven/bun:1.2.0-alpine AS base

# Install dependencies (Runs fast on AMD64)
FROM base AS deps
WORKDIR /app
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

# Build the application (Runs fast on AMD64)
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN bun run build

# 2. Production image 
#    Remove the --platform flag here so this stage uses the TARGET arch (e.g. ARM64)
FROM oven/bun:1.2.0-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nuxt

# 3. Copy the built artifacts from the fast AMD64 builder to this ARM64 runner
COPY --from=builder --chown=nuxt:nodejs /app/.output /app/.output

# Metadata labels
LABEL org.opencontainers.image.source="https://github.com/amorags/golangtickets"
LABEL org.opencontainers.image.description="Nuxt web frontend for ticket booking system"
LABEL org.opencontainers.image.licenses="MIT"

USER nuxt

EXPOSE 3000

ENV PORT=3000
ENV HOST=0.0.0.0

CMD ["bun", "run", ".output/server/index.mjs"]