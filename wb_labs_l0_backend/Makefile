COMPOSE_FILE=./build/docker-compose.yaml
PROJECT_NAME=wb_labs_l0_backend
FRONTEND_DIR=./frontend

.PHONY: dev build up down restart logs ps clean prune deps backend-deps frontend-deps

# Основная команда для запуска dev окружения
dev: deps
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) up --build

# Установка всех зависимостей (Go + React)
deps: backend-deps frontend-deps

# Go-зависимости
backend-deps:
	@echo "📦 Installing Go dependencies..."
	go mod tidy
	go mod download

# React-зависимости
frontend-deps:
	@echo "📦 Installing frontend dependencies..."
	cd $(FRONTEND_DIR) && npm install

# Сборка образов без запуска
build:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) build

# Запуск контейнеров в фоне
up:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) up -d

# Остановка контейнеров
down:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) down

# Перезапуск (с пересборкой)
restart:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) down && \
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) up --build -d

# Логи всех контейнеров
logs:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) logs -f

# Список контейнеров
ps:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) ps

# Полная очистка окружения (контейнеры, volume'ы, сети)
clean:
	docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) down -v --remove-orphans

# Глобальная очистка Docker (удаляет все образы, кэш и volume'ы)
prune:
	docker system prune -af --volumes
