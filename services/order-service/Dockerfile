FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ecommerce-product-service .


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ecommerce-product-service .

EXPOSE 3001

CMD ["./ecommerce-product-service"]