.PHONY: help infra-up infra-down proto-gen migrate-up migrate-down run-admin run-merchant run-client run-consumer run-websocket run-all build lint test

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ==================== Infrastructure ====================

infra-up: ## Start infrastructure (PostgreSQL, Redis, Kafka, MongoDB, ELK)
	docker compose -f docker-compose.infra.yml up -d

infra-down: ## Stop infrastructure
	docker compose -f docker-compose.infra.yml down

infra-reset: ## Stop infrastructure and remove volumes
	docker compose -f docker-compose.infra.yml down -v

infra-logs: ## Show infrastructure logs
	docker compose -f docker-compose.infra.yml logs -f

# ==================== Proto Generation ====================

proto-gen: ## Generate gRPC code from proto files
	./scripts/proto-gen.sh

proto-install: ## Install protoc plugins
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# ==================== Database Migrations ====================

MIGRATE_URL ?= postgres://goodie:goodie_secret@localhost:5432/goodie?sslmode=disable

migrate-up: ## Run all migrations
	migrate -database "$(MIGRATE_URL)" -path services/admin/migrations up
	migrate -database "$(MIGRATE_URL)" -path services/merchant/migrations up
	migrate -database "$(MIGRATE_URL)" -path services/client/migrations up

migrate-down: ## Rollback all migrations
	migrate -database "$(MIGRATE_URL)" -path services/client/migrations down -all
	migrate -database "$(MIGRATE_URL)" -path services/merchant/migrations down -all
	migrate -database "$(MIGRATE_URL)" -path services/admin/migrations down -all

migrate-install: ## Install golang-migrate CLI
	brew install golang-migrate

# ==================== Run Services ====================

run-admin: ## Run admin service
	cd services/admin && go run cmd/main.go

run-merchant: ## Run merchant service
	cd services/merchant && go run cmd/main.go

run-client: ## Run client service
	cd services/client && go run cmd/main.go

run-consumer: ## Run consumer service
	cd services/consumer && go run cmd/main.go

run-websocket: ## Run websocket service
	cd services/websocket && go run cmd/main.go

# ==================== Build ====================

build: ## Build all services
	cd services/admin && go build -o ../../bin/admin cmd/main.go
	cd services/merchant && go build -o ../../bin/merchant cmd/main.go
	cd services/client && go build -o ../../bin/client cmd/main.go
	cd services/consumer && go build -o ../../bin/consumer cmd/main.go
	cd services/websocket && go build -o ../../bin/websocket cmd/main.go

# ==================== Quality ====================

lint: ## Run linter
	golangci-lint run ./...

test: ## Run all tests
	cd pkg && go test ./...
	cd services/admin && go test ./...
	cd services/merchant && go test ./...
	cd services/client && go test ./...
	cd services/consumer && go test ./...
	cd services/websocket && go test ./...

tidy: ## Tidy all go modules
	cd pkg && go mod tidy
	cd services/admin && go mod tidy
	cd services/merchant && go mod tidy
	cd services/client && go mod tidy
	cd services/consumer && go mod tidy
	cd services/websocket && go mod tidy

# ==================== Docker ====================

docker-build: ## Build Docker images for all services
	docker build -t goodie-admin -f services/admin/Dockerfile .
	docker build -t goodie-merchant -f services/merchant/Dockerfile .
	docker build -t goodie-client -f services/client/Dockerfile .
	docker build -t goodie-consumer -f services/consumer/Dockerfile .
	docker build -t goodie-websocket -f services/websocket/Dockerfile .

# ==================== Seed ====================

seed: ## Seed initial data
	./scripts/seed.sh
