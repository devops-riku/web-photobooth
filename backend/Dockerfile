# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for fetching dependencies
RUN apk add --no-cache git

# Copy go module files first to leverage Docker cache
COPY go.mod go.sum ./

# Install all backend requirements/dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage - Production
FROM alpine:latest

WORKDIR /app

# Install CA certificates for external API calls (DigitalOcean Spaces, etc.)
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]
