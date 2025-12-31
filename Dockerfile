# STAGE 1: Build the binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files first to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build for ARM64 (Graviton)
# CGO_ENABLED=0 creates a static binary that works in Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main .

# STAGE 2: Run the binary
FROM alpine:latest
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]