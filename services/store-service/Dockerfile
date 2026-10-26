FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ecommerce-store-service .


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ecommerce-store-service .

EXPOSE 3001

CMD ["./ecommerce-store-service"]