# base image 

FROM golang:1.22.2-alpine as builder

WORKDIR /app

COPY . /app


RUN CGO_ENABLE=0 go build -o brokerApp ./cmd/api

# RUN chmod +x /app/brokerApp



FROM alpine:latest

WORKDIR /app

COPY --from=builder app/brokerApp /app

ENTRYPOINT ["/app/brokerApp"]




