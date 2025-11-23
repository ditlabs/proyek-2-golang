# Stage 1: Build
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build aplikasi dengan CGO disabled (untuk Railway Linux environment)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o backend-sarpras ./cmd/server

# Stage 2: Runtime
FROM alpine:3.19

WORKDIR /root/

# Install ca-certificates untuk HTTPS dan wget untuk health check
RUN apk --no-cache add ca-certificates wget

# Copy binary dari builder
COPY --from=builder /app/backend-sarpras .

# Expose port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8000/api/health || exit 1

# Run aplikasi
CMD ["./backend-sarpras"]
