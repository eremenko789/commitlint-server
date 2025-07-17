# Переменные
BINARY_CLI = commitlint
BINARY_SERVER = commitlint-server
VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS = -ldflags "-X github.com/conventionalcommit/commitlint/internal.version=$(VERSION)"

.PHONY: all build build-cli build-server clean test help install-cli install-server

# По умолчанию собираем оба приложения
all: build

# Сборка обоих приложений
build: build-cli build-server

# Сборка CLI приложения
build-cli:
	@echo "Сборка CLI приложения..."
	go build $(LDFLAGS) -o $(BINARY_CLI) main.go

# Сборка веб-сервера
build-server:
	@echo "Сборка веб-сервера..."
	go build $(LDFLAGS) -o $(BINARY_SERVER) cmd/server/main.go

# Установка CLI приложения
install-cli:
	@echo "Установка CLI приложения..."
	go install $(LDFLAGS) .

# Установка веб-сервера
install-server:
	@echo "Установка веб-сервера..."
	go install $(LDFLAGS) ./cmd/server

# Запуск тестов
test:
	@echo "Запуск тестов..."
	go test -v ./...

# Очистка собранных файлов
clean:
	@echo "Очистка..."
	rm -f $(BINARY_CLI) $(BINARY_SERVER)

# Сборка для разных платформ
build-all-platforms: clean
	@echo "Сборка для всех платформ..."
	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_CLI)-linux-amd64 main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_SERVER)-linux-amd64 cmd/server/main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_CLI)-windows-amd64.exe main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_SERVER)-windows-amd64.exe cmd/server/main.go
	# MacOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_CLI)-darwin-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_SERVER)-darwin-amd64 cmd/server/main.go

# Показать справку
help:
	@echo "Доступные команды:"
	@echo "  build          - Собрать оба приложения"
	@echo "  build-cli      - Собрать только CLI приложение"
	@echo "  build-server   - Собрать только веб-сервер"
	@echo "  install-cli    - Установить CLI приложение"
	@echo "  install-server - Установить веб-сервер"
	@echo "  test           - Запустить тесты"
	@echo "  clean          - Удалить собранные файлы"
	@echo "  build-all-platforms - Собрать для всех платформ"
	@echo "  help           - Показать эту справку"