# Simple Next.js Docker build
FROM node:18-alpine AS base
RUN apk add --no-cache libc6-compat wget

# Install dependencies
FROM base AS deps
WORKDIR /app

# Copy root package files (required for workspace resolution)
COPY package.json package-lock.json ./
COPY ui/portfolio/package.json ./ui/portfolio/

# Install all dependencies for the workspace
RUN npm ci --workspace=ui/portfolio

# Build stage
FROM base AS builder
WORKDIR /app

# Copy installed dependencies
COPY --from=deps /app/node_modules ./node_modules

# Copy source files
COPY ui/portfolio/ ./ui/portfolio/

# Build the app
WORKDIR /app/ui/portfolio
RUN npm run build:docker

# Production stage
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copy everything needed to run Next.js
COPY --from=builder /app/ui/portfolio/package.json ./package.json
COPY --from=builder /app/ui/portfolio/.next ./.next
COPY --from=builder /app/ui/portfolio/public ./public
COPY --from=builder /app/node_modules ./node_modules

USER nextjs
EXPOSE 3000
CMD ["npm", "start"]
