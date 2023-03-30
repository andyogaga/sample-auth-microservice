FROM golang:1.18 AS base
WORKDIR /app
COPY . .
EXPOSE 3001
RUN go mod download
CMD go run main.go