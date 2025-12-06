FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o keep-alive-service main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/keep-alive-service .

RUN apk --no-cache add ca-certificates

CMD ["./keep-alive-service"]