FROM golang:1.22.2-alpine as builder

WORKDIR /app

COPY . /app


RUN CGO_ENABLE=0 go build -o authservice ./cmd/api




FROM alpine:latest

WORKDIR /app

COPY --from=builder app/authservice /app

ENTRYPOINT ["/app/authservice"]