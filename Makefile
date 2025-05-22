# Default ENV mode
ENV ?= development
APP_ENV := $(ENV)

# Run local Go apps
run-bot-tele:
	@echo "Running bot-tele in $(ENV) mode..."
	@cd bot-tele && APP_ENV=$(APP_ENV) go run cmd/main.go

run-core-api:
	@echo "Running core-api in $(ENV) mode..."
	@cd core-api && APP_ENV=$(APP_ENV) go run cmd/main.go

# Build Docker (optional)
build-bot-tele:
	@docker build -t bot-tele ./bot-tele

build-core-api:
	@docker build -t core-api ./core-api

# Run Docker Compose
up:
	docker-compose up -d --build

down:
	docker-compose down
