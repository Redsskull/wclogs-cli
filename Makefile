# Warcraft Logs CLI Makefile

.PHONY: build clean install uninstall test help dev container-build container-run container-shell

# Binary name
BINARY_NAME=wclogs

# Build directory
BUILD_DIR=build

# Container settings
CONTAINER_NAME=wclogs-cli
IMAGE_NAME=wclogs-cli
CONTAINER_TOOL := $(shell which podman 2>/dev/null || which docker 2>/dev/null || echo "container-tool-not-found")

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# Default target
all: build

# Build the binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v

# Build for current directory (quick build)
build-local:
	@echo "🔨 Building $(BINARY_NAME) in current directory..."
	$(GOBUILD) -o $(BINARY_NAME) -v

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)

# Install globally (Unix/Linux/macOS)
install: build
	@echo "📦 Installing $(BINARY_NAME) to /usr/local/bin..."
	@if [ "$$(id -u)" -ne 0 ]; then \
		echo "⚠️  Installing to /usr/local/bin requires sudo"; \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/; \
	else \
		cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/; \
	fi
	@echo "✅ $(BINARY_NAME) installed successfully!"
	@echo "📝 You can now run '$(BINARY_NAME)' from anywhere"

# Uninstall from system
uninstall:
	@echo "🗑️  Removing $(BINARY_NAME) from /usr/local/bin..."
	@if [ "$$(id -u)" -ne 0 ]; then \
		echo "⚠️  Removing from /usr/local/bin requires sudo"; \
		sudo rm -f /usr/local/bin/$(BINARY_NAME); \
	else \
		rm -f /usr/local/bin/$(BINARY_NAME); \
	fi
	@echo "✅ $(BINARY_NAME) uninstalled successfully!"

# Development mode - build and run with config
dev: build-local
	@echo "🚀 Starting development mode..."
	@echo "📝 Run './$(BINARY_NAME) config' to set up credentials first"
	@echo "💡 Try: './$(BINARY_NAME) --help'"
	@echo "💡 Example: './$(BINARY_NAME) damage ABC123 last'"

# Quick setup for new users
setup: deps build
	@echo "🎯 Setup complete! Next steps:"
	@echo "1️⃣  Get API credentials from: https://www.warcraftlogs.com/api/clients"
	@echo "2️⃣  Run: ./$(BUILD_DIR)/$(BINARY_NAME) config"
	@echo "3️⃣  Start analyzing: ./$(BUILD_DIR)/$(BINARY_NAME) --help"
	@echo ""
	@echo "💡 To install globally: make install"

# Build for multiple platforms
build-all: clean
	@echo "🌍 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)

	# Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64

	# macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64

	# Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe

	@echo "✅ Built for all platforms in $(BUILD_DIR)/"

# Container targets
container-build:
	@if [ "$(CONTAINER_TOOL)" = "container-tool-not-found" ]; then \
		echo "❌ Neither podman nor docker found. Please install one of them."; \
		exit 1; \
	fi
	@echo "🐳 Building container image with $(notdir $(CONTAINER_TOOL))..."
	$(CONTAINER_TOOL) build -t $(IMAGE_NAME) .

container-run: container-build
	@echo "🚀 Running $(BINARY_NAME) in container with $(notdir $(CONTAINER_TOOL))..."
	@echo "💡 Your ~/.wclogs.yaml will be mounted for config (if it exists)"
	@if [ -f "$(HOME)/.wclogs.yaml" ]; then \
		$(CONTAINER_TOOL) run --rm -it \
			-v $(HOME)/.wclogs.yaml:/home/wclogs/.wclogs.yaml:ro \
			$(IMAGE_NAME) $(ARGS); \
	else \
		$(CONTAINER_TOOL) run --rm -it $(IMAGE_NAME) $(ARGS); \
	fi

container-shell: container-build
	@echo "🐚 Opening shell in container with $(notdir $(CONTAINER_TOOL))..."
	@if [ -f "$(HOME)/.wclogs.yaml" ]; then \
		$(CONTAINER_TOOL) run --rm -it \
			-v $(HOME)/.wclogs.yaml:/home/wclogs/.wclogs.yaml:ro \
			--entrypoint /bin/sh \
			$(IMAGE_NAME); \
	else \
		$(CONTAINER_TOOL) run --rm -it \
			--entrypoint /bin/sh \
			$(IMAGE_NAME); \
	fi

container-clean:
	@if [ "$(CONTAINER_TOOL)" = "container-tool-not-found" ]; then \
		echo "❌ Neither podman nor docker found."; \
		exit 1; \
	fi
	@echo "🧹 Cleaning container images with $(notdir $(CONTAINER_TOOL))..."
	$(CONTAINER_TOOL) rmi $(IMAGE_NAME) 2>/dev/null || true

# Show help
help:
	@echo "🗡️  Warcraft Logs CLI - Build Commands"
	@echo ""
	@echo "📋 Native Build Targets:"
	@echo "   build        Build binary to $(BUILD_DIR)/"
	@echo "   build-local  Build binary to current directory"
	@echo "   build-all    Build for multiple platforms"
	@echo "   install      Install binary globally (requires sudo)"
	@echo "   uninstall    Remove binary from system"
	@echo "   deps         Install Go dependencies"
	@echo "   test         Run all tests"
	@echo "   clean        Remove build artifacts"
	@echo "   dev          Quick development setup"
	@echo "   setup        Complete setup for new users"
	@echo ""
	@echo "🐳 Container Targets (using $(notdir $(CONTAINER_TOOL))):"
	@echo "   container-build    Build container image"
	@echo "   container-run      Run wclogs in container"
	@echo "   container-shell    Open shell in container"
	@echo "   container-clean    Remove container images"
	@echo ""
	@echo "🚀 Quick start:"
	@echo "   make setup           # Native setup"
	@echo "   make container-run   # Container setup"
	@echo ""
	@echo "📝 Container Examples:"
	@echo "   make container-run ARGS='config'"
	@echo "   make container-run ARGS='config'"
	@echo "   make container-run ARGS='damage ABC123 5'"
	@echo "   make container-run ARGS='damage ABC123 last'"
