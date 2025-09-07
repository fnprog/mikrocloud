# Multi-stage build for optimal image size
FROM node:20-alpine AS frontend-builder

# Set working directory for frontend
WORKDIR /app/web

# Copy frontend package files
COPY web/package.json web/pnpm-lock.yaml ./

# Install pnpm and dependencies
RUN npm install -g pnpm && pnpm install

# Copy frontend source
COPY web/ .

# Build frontend
RUN pnpm build

# Backend builder stage
FROM golang:1.22-alpine AS backend-builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend from frontend-builder
COPY --from=frontend-builder /app/web/dist ./web/dist

# Build the application with embedded frontend
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mikrocloud ./main.go

# Final stage - minimal image
FROM alpine:3.18

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata && \
    adduser -D -s /bin/sh mikrocloud

WORKDIR /app

# Copy binary from backend-builder stage
COPY --from=backend-builder /app/mikrocloud .
COPY --from=backend-builder /app/mikrocloud.toml .

# Copy migrations
COPY --from=backend-builder /app/migrations ./migrations

# Change ownership
RUN chown -R mikrocloud:mikrocloud /app

# Switch to non-root user
USER mikrocloud

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# Default command
CMD ["./mikrocloud", "serve"]
