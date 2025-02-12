APP_NAME := kubectl-alias
OUTPUT_NAME := kubectl-alias  # Final binary name inside ZIP
VERSION  := $(shell cat VERSION)
BUILD_DIR := dist
LDFLAGS := -X main.Version=$(VERSION) -s -w

PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

all: clean deps test build package

deps:  ## Install dependencies
	@echo "Downloading dependencies..."
	go mod tidy

test:  ## Run tests
	@echo "Running tests..."
	go test ./... -v

coverage:  ## Run tests with coverage and show result
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

build: ## Build binaries for all platforms
	@echo "Building binaries..."
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; ARCH=$${platform#*/}; \
		EXT=""; \
		if [ "$${OS}" = "windows" ]; then EXT=".exe"; fi; \
		GOOS=$${OS} GOARCH=$${ARCH} \
		go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(OUTPUT_NAME)$${EXT} main.go; \
		mv $(BUILD_DIR)/$(OUTPUT_NAME)$${EXT} $(BUILD_DIR)/$(OUTPUT_NAME)-$${OS}-$${ARCH}$${EXT}; \
	done

package: build ## Package binaries into zip/tar.gz
	@echo "Packaging binaries..."
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; ARCH=$${platform#*/}; \
		EXT=""; \
		if [ "$${OS}" = "windows" ]; then EXT=".exe"; fi; \
		FILENAME="$(OUTPUT_NAME)-$(VERSION)-$${OS}-$${ARCH}"; \
		mkdir -p $(BUILD_DIR)/$${FILENAME}; \
		cp $(BUILD_DIR)/$(OUTPUT_NAME)-$${OS}-$${ARCH}$${EXT} $(BUILD_DIR)/$${FILENAME}/$(OUTPUT_NAME)$${EXT}; \
		if [ "$${OS}" = "windows" ]; then \
			zip -j $(BUILD_DIR)/$${FILENAME}.zip $(BUILD_DIR)/$${FILENAME}/*; \
		else \
			tar -czf $(BUILD_DIR)/$${FILENAME}.tar.gz -C $(BUILD_DIR)/$${FILENAME} .; \
		fi; \
		rm -rf $(BUILD_DIR)/$${FILENAME}; \
	done

clean: ## Clean build artifacts
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR) coverage.out coverage.html

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
