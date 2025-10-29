.PHONY: all up down down-v restart logs clean

all: up

up:
	@echo "Запуск стека + автопубликация заказа..."
	@cd build && docker compose up --build

up-d:
	@echo "Запуск в фоне..."
	@cd build && docker compose up -d --build

down:
	@cd build && docker compose down

down-v:
	@echo "Остановка + очистка..."
	@cd build && docker compose down -v
	@docker volume prune -f

restart: down-v up

logs:
	@cd build && docker compose logs -f

clean: down-v
	@docker system prune -f
	@docker volume prune -f