FROM golang:1.23-alpine AS builder

COPY . /github.com/codeboris/music-lib/
WORKDIR /github.com/codeboris/music-lib/

RUN go mod download
RUN go build -o ./bin/main cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/codeboris/music-lib/bin/main .
COPY --from=builder /github.com/codeboris/music-lib/migrations migrations/

EXPOSE 8000

CMD ["./main"]
