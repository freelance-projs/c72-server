FROM golang:1.23.3 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o main ./cmd/api

FROM alpine

WORKDIR /

COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["/main"]
