FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY pkg/ ./pkg/
COPY services/merchant/ ./services/merchant/

WORKDIR /app/services/merchant
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/merchant cmd/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/bin/merchant .
COPY services/merchant/migrations ./migrations/

EXPOSE 8082 9082

CMD ["./merchant"]
