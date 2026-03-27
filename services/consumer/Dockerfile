FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY pkg/ ./pkg/
COPY services/consumer/ ./services/consumer/

WORKDIR /app/services/consumer
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/consumer cmd/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/bin/consumer .

CMD ["./consumer"]
