FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]