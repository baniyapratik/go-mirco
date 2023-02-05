# base go image
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o broker-service ./

RUN chmod +x /app/broker-service

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/broker-service /app

CMD ["/app/broker-service"]