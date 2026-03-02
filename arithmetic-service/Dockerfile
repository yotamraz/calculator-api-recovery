# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod first for better layer caching
COPY go.mod ./

# Copy source code
COPY *.go ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o arithmetic-service .

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/arithmetic-service .

# Default port (can be overridden via ARITHMETIC_SERVICE_PORT env var)
ENV ARITHMETIC_SERVICE_PORT=8080
EXPOSE 8080

ENTRYPOINT ["./arithmetic-service"]
