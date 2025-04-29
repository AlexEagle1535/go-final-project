FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_app

FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /my_app /my_app

COPY web /web

COPY .env .env

COPY scheduler.db /app/scheduler.db

CMD ["/my_app"]