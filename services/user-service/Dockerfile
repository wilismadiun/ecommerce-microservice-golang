FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ecommerce-user-service .


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ecommerce-user-service .
COPY --from=builder /app/migrations ./migrations

EXPOSE 3000

CMD ["./ecommerce-user-service"]