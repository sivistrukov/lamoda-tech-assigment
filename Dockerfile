FROM golang:1.22.0-bookworm

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./bin/app ./cmd/main.go