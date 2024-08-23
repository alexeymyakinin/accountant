FROM golang:1.23 AS build

WORKDIR /usr/src

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd cmd
COPY internal internal
COPY migrations migrations

RUN go build -a -o app cmd/main.go
