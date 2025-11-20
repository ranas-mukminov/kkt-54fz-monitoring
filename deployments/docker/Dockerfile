# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kkt-monitor ./cmd/kkt-monitor

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/kkt-monitor .

# Copy configuration
COPY configs/ ./configs/

# Create directory for logs
RUN mkdir -p /var/log/kkt

# Expose metrics port
EXPOSE 9090

# Run as non-root user
RUN adduser -D -u 1000 kktmon
USER kktmon

ENTRYPOINT ["/app/kkt-monitor"]
CMD ["--config", "/app/configs/config.yaml"]
