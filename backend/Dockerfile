# ---- Build Stage ----
FROM golang:1.25.1-alpine AS builder

# Enable Go modules and set working directory
WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN go build -o app ./cmd/main

# ---- Run Stage ----
FROM alpine:latest

# Copy binary from builder
COPY --from=builder /app/app /app

# Expose the port
EXPOSE 3000

# Run the binary
ENTRYPOINT ["/app"]
