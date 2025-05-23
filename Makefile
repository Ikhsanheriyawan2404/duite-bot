# Default ENV mode
ENV ?= development
APP_ENV := $(ENV)

# Jalankan bot-wa dengan mode env yang sesuai
run-bot-wa:
	@echo "Running bot-wa in $(ENV) mode..."
ifeq ($(ENV),production)
	@cd bot-wa && APP_ENV=$(APP_ENV) node index.js --env-file .env.production
else
	@cd bot-wa && APP_ENV=$(APP_ENV) node index.js
endif

# Jalankan bot-tele dengan mode env yang sesuai
run-bot-tele:
	@echo "Running bot-tele in $(ENV) mode..."
ifeq ($(ENV),production)
	@cd bot-tele && APP_ENV=$(APP_ENV) go run cmd/main.go --env-file .env.production
else
	@cd bot-tele && APP_ENV=$(APP_ENV) go run cmd/main.go
endif

# Jalankan core-api dengan mode env yang sesuai
run-core-api:
	@echo "Running core-api in $(ENV) mode..."
ifeq ($(ENV),production)
	@cd core-api && APP_ENV=$(APP_ENV) go run cmd/main.go --env-file .env.production
else
	@cd core-api && APP_ENV=$(APP_ENV) go run cmd/main.go
endif

# Docker Compose untuk production
up:
	@echo "📦 Bringing up production stack..."
	docker compose up -d --build

down:
	@echo "🛑 Tearing down production stack..."
	docker compose down
