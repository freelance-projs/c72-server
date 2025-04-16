FROM golang:1.23.3 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o main ./cmd/api

FROM alpine
# Install tzdata to configure the time zone
RUN apk add --no-cache tzdata
# Set the time zone to Ho Chi Minh (Viet Nam)
RUN ln -fs /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime \
    && echo "Asia/Ho_Chi_Minh" > /etc/timezone

WORKDIR /

COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/account_credentials.json .

EXPOSE 8080

ENTRYPOINT ["/main"]
