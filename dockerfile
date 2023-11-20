FROM golang:1.20 as builder

WORKDIR /app
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /bin/graphql

LABEL org.opencontainers.image.source="https://github.com/kilianp07/graphql"

