# Build stage
FROM golang:1.19.2-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o aggregator main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/aggregator .

EXPOSE 8080
CMD [ "/app/aggregator" ]
