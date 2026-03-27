FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY pkg/ ./pkg/
COPY services/client/ ./services/client/

WORKDIR /app/services/client
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/client cmd/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/bin/client .
COPY services/client/migrations ./migrations/

EXPOSE 8083 9083

CMD ["./client"]
