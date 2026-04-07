APP_NAME=cinema-api

# ======================
# DEV
# ======================

dev:
	@echo "Starting DEV environment..."
	docker compose -f docker-compose.dev.yml up -d
	air

dev-down:
	@echo "Stopping DEV environment..."
	docker compose -f docker-compose.dev.yml down

dev-reset:
	@echo "Reset DEV (DB included)..."
	docker compose -f docker-compose.dev.yml down -v
	docker compose -f docker-compose.dev.yml up -d

# ======================
# PROD
# ======================

prod:
	@echo "Starting PROD..."
	docker compose up --build -d

prod-down:
	@echo "Stopping PROD..."
	docker compose down

# ======================
# UTIL
# ======================

logs:
	docker compose logs -f

ps:
	docker compose ps

swagger:
	@echo "Generating Swagger docs..."
	swag init -g cmd/main.go --parseDependency --parseInternal