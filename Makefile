SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -c
.DEFAULT_GOAL := help
.DELETE_ON_ERROR:

SERVICE_NAME := notify-svc
BIN_DIR := bin
BINARY := $(BIN_DIR)/$(SERVICE_NAME)

SERVER_CMD := ./cmd/svc-starter
MIGRATOR_CMD := ./cmd/migrator

.PHONY: help build migrate run run-only \
		clean clean-bin print-config 

help: ## Показать список доступных команд
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

migrate: ## Применить миграции
	@echo "⚙️  Применение миграций $(SERVICE_NAME)..."
	@go run $(MIGRATOR_CMD)

build: ## Собрать бинарник сервиса
	@echo "🔨  Сборка $(SERVICE_NAME)..."
	@mkdir -p "$(BIN_DIR)"
	@go build -o "$(BINARY)" $(SERVER_CMD)
	@echo "✅  Собран $(BINARY)"

run: migrate build ## Применить миграции, собрать бинарник и запустить сервис
	@echo "🚀  Запуск $(SERVICE_NAME)..."
	@exec "$(BINARY)"

run-only: ## Запустить сервис без миграций и сборки бинарника
	@test -f "$(BINARY)" || { echo "❌  $(BINARY) не найден. Выполните make build"; exit 1; }
	@echo "🚀  Запуск $(SERVICE_NAME)..."
	@exec "$(BINARY)"

clean: clean-bin ## Очистить локальные артефакты
clean-bin: ## Удалить собранные бинарники
	@rm -rf "$(BIN_DIR)"
	@echo "🧹  Директория $(BIN_DIR) удалена"

print-config: ## Показать текущую конфигурацию
	@echo "SERVICE_NAME        = $(SERVICE_NAME)"
	@echo "BINARY              = $(BINARY)"