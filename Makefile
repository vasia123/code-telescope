.PHONY: build clean test lint run help install-deps

# Основные переменные
BINARY_NAME=code-telescope
MAIN_PATH=./cmd/codetelescope
BUILD_DIR=./build
CONFIG_DIR=./configs
VERSION=$(shell git describe --tags --always --abbrev=0 || echo "dev")
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Цель по умолчанию
.DEFAULT_GOAL := help

# Справка
help:
	@echo "Доступные команды:"
	@echo "  make build      - Сборка бинарного файла"
	@echo "  make run        - Запуск программы с конфигурацией по умолчанию"
	@echo "  make test       - Запуск тестов"
	@echo "  make lint       - Проверка кода линтером"
	@echo "  make clean      - Очистка артефактов сборки"
	@echo "  make install-deps - Установка зависимостей"

# Сборка бинарного файла
build:
	@echo "Сборка $(BINARY_NAME) версии $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Готово! Бинарный файл создан: $(BUILD_DIR)/$(BINARY_NAME)"

# Запуск программы
run: build
	@echo "Запуск $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME) --config $(CONFIG_DIR)/default.yaml

# Запуск тестов
test:
	@echo "Запуск тестов..."
	go test -v ./...

# Проверка кода линтером
lint:
	@echo "Проверка кода линтером..."
	golangci-lint run ./...

# Очистка артефактов сборки
clean:
	@echo "Очистка артефактов сборки..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Очистка завершена."

# Установка зависимостей
install-deps:
	@echo "Установка зависимостей..."
	@go mod tidy
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Установка завершена." 