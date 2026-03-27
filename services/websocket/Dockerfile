FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY pkg/ ./pkg/
COPY services/websocket/ ./services/websocket/

WORKDIR /app/services/websocket
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/websocket cmd/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/bin/websocket .

EXPOSE 8084

CMD ["./websocket"]
