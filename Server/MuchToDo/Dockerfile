# ---------- BUILD STAGE ----------
# ─────────────────────────────────────────────
# Stage 1: Build
# ─────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

# Install git and certs (needed for go module downloads)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o muchtodo-api \
    ./cmd/api

# ─────────────────────────────────────────────
# Stage 2: Final minimal image
# ─────────────────────────────────────────────
FROM alpine:3.19

# Security: add ca-certs and timezone data
RUN apk --no-cache add ca-certificates tzdata wget && \
    update-ca-certificates

# Create non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy only the compiled binary from builder
COPY --from=builder /app/muchtodo-api .

# Set ownership to non-root user
RUN chown -R appuser:appgroup /app

USER appuser

# Expose application port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=15s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./muchtodo-api"]
