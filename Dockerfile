# ---- Build stage ----
    FROM golang:1.24-alpine AS builder

    # Install tools (e.g. git if needed)
    RUN apk update && apk add --no-cache git && rm -rf /var/cache/apk/*
    
    # Set working directory
    WORKDIR /app
    
    # Copy go.mod and go.sum first for better caching
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy the rest of the source code
    COPY . .
    
    # Build the Go app
    RUN go build -o cubiqapi main.go
    
    # ---- Run stage ----
    FROM alpine:latest
    
    # Set working directory
    WORKDIR /root/
    
    # Copy binary from build stage
    COPY --from=builder /app/cubiqapi .
    COPY --from=builder /app/.env .env
    COPY --from=builder /app/service-account.json service-account.json

    # Print contents of .env for debugging purposes
    RUN cat .env
    
    # Expose application port (optional)
    # Application runs on port 8080
    EXPOSE 8080
    
    # Run the binary
    CMD ["sh", "-c", "cat .env && ./cubiqapi"]
    