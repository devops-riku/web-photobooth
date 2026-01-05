# =========================
# Build stage
# =========================
FROM golang:1.24-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# =========================
# Runtime stage
# =========================
FROM alpine:latest

WORKDIR /app

# 🔑 REQUIRED for JWT generation
RUN apk --no-cache add ca-certificates openssl

COPY --from=builder /app/main ./main
COPY .env.template ./.env.template
COPY entrypoint.sh ./entrypoint.sh

RUN sed -i 's/\r$//' ./entrypoint.sh && chmod +x ./entrypoint.sh

EXPOSE 3000

ENTRYPOINT ["./entrypoint.sh"]
CMD ["./main"]
