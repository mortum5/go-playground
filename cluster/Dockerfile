FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY *.go .

RUN go build -o app *.go

FROM alpine:latest

RUN apk update && apk upgrade

RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]