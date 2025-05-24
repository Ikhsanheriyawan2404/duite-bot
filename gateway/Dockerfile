# syntax=docker/dockerfile:1

# 1. Builder Stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Install git only if needed for modules
RUN apk add --no-cache git ca-certificates tzdata \
 && update-ca-certificates 2>/dev/null || true \
 && addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy go.mod & go.sum dulu, download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy kode sumber, build binary statically
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway main.go

# 2. Runtime Stage
FROM scratch
# Optional: gunakan alpine minimal jika butuh sh, ca-certificates
# FROM alpine:latest

# Copy user dari builder agar non-root
COPY --from=builder /etc/passwd /etc/passwd
USER appuser

# Copy binary dari builder
COPY --from=builder /app/gateway /gateway
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Jalankan binary
ENTRYPOINT ["/gateway"]
