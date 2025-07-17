.PHONY: all build build-cli build-webhook-server clean test install

# Build variables
GO := go
GOFLAGS := -v
LDFLAGS := -s -w

# Binary names
CLI_BINARY := commitlint
WEBHOOK_BINARY := commitlint-webhook-server

# Build directories
BUILD_DIR := build
CLI_DIR := cmd/cli
WEBHOOK_DIR := cmd/webhook-server

all: build

build: build-cli build-webhook-server

build-cli:
	@echo "Building CLI application..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(CLI_BINARY) ./$(CLI_DIR)
	@echo "CLI built: $(BUILD_DIR)/$(CLI_BINARY)"

build-webhook-server:
	@echo "Building webhook server..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(WEBHOOK_BINARY) ./$(WEBHOOK_DIR)
	@echo "Webhook server built: $(BUILD_DIR)/$(WEBHOOK_BINARY)"

install: install-cli

install-cli: build-cli
	@echo "Installing CLI..."
	@cp $(BUILD_DIR)/$(CLI_BINARY) $(GOPATH)/bin/
	@echo "CLI installed to $(GOPATH)/bin/$(CLI_BINARY)"

install-webhook-server: build-webhook-server
	@echo "Installing webhook server..."
	@cp $(BUILD_DIR)/$(WEBHOOK_BINARY) $(GOPATH)/bin/
	@echo "Webhook server installed to $(GOPATH)/bin/$(WEBHOOK_BINARY)"

test:
	@echo "Running tests..."
	$(GO) test ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Docker builds
docker-build-cli:
	docker build -f Dockerfile.cli -t commitlint:latest .

docker-build-webhook-server:
	docker build -f Dockerfile.webhook -t commitlint-webhook-server:latest .