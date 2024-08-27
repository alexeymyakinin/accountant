FROM golang:1.23 AS build

WORKDIR /usr/app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY cmd cmd
COPY internal internal
COPY migrations migrations
COPY test test

RUN go build -a -o build/app cmd/main.go
