# Docker registry
REGISTRY ?= ikhsan123
TAG ?= latest

# Daftar service yang akan di-build dan push
SERVICES = bot-tele bot-wa core-api gateway

# Default ENV mode
ENV ?= development
APP_ENV := $(ENV)

# Docker Compose file
COMPOSE_FILE=compose.yml
COMPOSE_DEV=docker-compose.override.yml
COMPOSE_PROD=docker-compose.prod.yml

build-images:
	@echo "🔧 Building Docker images for $(ENV) environment..."
	@for service in $(SERVICES); do \
		echo "🔨 Building $$service..."; \
		docker build -t $(REGISTRY)/$$service:$(TAG) ./$$service; \
	done

push-images:
	@echo "🚀 Pushing Docker images to registry $(REGISTRY)..."
	@for service in $(SERVICES); do \
		echo "📤 Pushing $(REGISTRY)/$$service:$(TAG)..."; \
		docker push $(REGISTRY)/$$service:$(TAG); \
	done

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

up:
ifeq ($(ENV),production)
	@echo "📦 Removing old images to ensure fresh pull..."
	@for service in $(SERVICES); do \
		docker image rm $(REGISTRY)/$$service:$(TAG) || true; \
	done
	@echo "📦 Bringing up production stack..."
	docker compose -f $(COMPOSE_FILE) -f $(COMPOSE_PROD) up -d --pull always
else
	@echo "📦 Bringing up development stack..."
	docker compose -f $(COMPOSE_FILE) -f $(COMPOSE_DEV) up -d --build
endif

down:
ifeq ($(ENV),production)
	@echo "🛑 Tearing down production stack..."
	docker compose -f $(COMPOSE_FILE) -f $(COMPOSE_PROD) down
else
	@echo "🛑 Tearing down development stack..."
	docker compose -f $(COMPOSE_FILE) -f $(COMPOSE_DEV) down
endif
