# Stage 1: Build
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build aplikasi
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o backend-sarpras ./cmd/server

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates untuk HTTPS
RUN apk --no-cache add ca-certificates libc6-compat

# Copy binary dari builder
COPY --from=builder /app/backend-sarpras .

# Expose port
EXPOSE 8000

# Run aplikasi
CMD ["./backend-sarpras"]
