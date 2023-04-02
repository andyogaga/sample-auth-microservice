FROM golang:1.18 AS base

WORKDIR /app

COPY . .

EXPOSE 3001

CMD go run main.go
