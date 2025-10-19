FROM golang:1.25.1-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download

RUN go build -o delayed_notifier ./cmd/app

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/delayed_notifier .
COPY config.yaml .
COPY .env .

EXPOSE 7540

CMD ["./delayed_notifier"]
