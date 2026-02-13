# Build stage
FROM golang:1.25.6-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Install build dependencies
RUN apk update --no-cache
RUN apk add --no-cache git ca-certificates
RUN update-ca-certificates

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN go build -o server cmd/server/main.go

# Final stage
FROM alpine:latest

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Expose port
EXPOSE 8080

# Set default environment variables
ENV PORT=8080
ENV DATABASE_PATH=/data/trash_taste.db
COPY trash_taste.db /data/trash_taste.db

# Run the binary
CMD ["./server"]