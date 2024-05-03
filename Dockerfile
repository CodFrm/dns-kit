FROM golang:1.22-alpine as build

WORKDIR /app

COPY . .

RUN go build -o dns-kit ./cmd/main/main.go

