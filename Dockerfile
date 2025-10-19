FROM golang:1.25.1-alpine AS builder

WORKDIR /app
RUN go.mod go.sum ./
RUN go mod download

COPY . .



RUN go build -o delayed_notifier ./cmd

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/delayed_notifier .
COPY config.yaml .
COPY .env .

EXPOSE 7540

CMD ["./delayed_notifier"]
