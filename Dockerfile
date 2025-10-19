FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY . .


COPY .env config.yaml ./
RUN go build -o delayed_notifier ./cmd

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/delayed_notifier .
COPY .env .
COPY config.yaml .

EXPOSE 7540

CMD ["./delayed_notifier"]
