# Redis simple chat based on websockets

## Install

1. Install [docker](https://www.docker.com/) and docker compose
2. Install [websocat](https://github.com/vi/websocat) for testing

## Run

1. Up redis compose file in redis folder `docker compose up -d`
2. Run `go run main.go`
3. Connect from different terminal `websocat ws://0.0.0.0:8080/chat/{username}`
4. Text some messages