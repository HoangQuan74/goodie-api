FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy shared packages
COPY pkg/ ./pkg/

# Copy service
COPY services/admin/ ./services/admin/

# Build
WORKDIR /app/services/admin
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/admin cmd/main.go

# Runtime
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/bin/admin .
COPY services/admin/migrations ./migrations/

EXPOSE 8081 9081

CMD ["./admin"]
