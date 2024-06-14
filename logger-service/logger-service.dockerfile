FROM golang:1.22.2-alpine as builder

WORKDIR /app

COPY . ./

RUN CGO_ENABLE=0 go build -o loggerservice ./cmd/api




FROM alpine:latest

WORKDIR /app

COPY --from=builder app/loggerservice /app

ENTRYPOINT ["/app/loggerservice"]