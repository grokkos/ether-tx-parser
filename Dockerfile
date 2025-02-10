FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/ether-tx-parser ./cmd/api

# Use a minimal alpine image for the final container
FROM alpine:3.18

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/ether-tx-parser .
# Copy config files
COPY --from=builder /app/config.yaml /app/

# Create directory for logs
RUN mkdir -p /app/logs

EXPOSE 8080

CMD ["./ether-tx-parser"]