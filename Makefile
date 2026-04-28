.PHONY: help build clean build-all

APP_NAME := certd-cli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
OUTPUT_DIR := dist

LDFLAGS := -s -w \
	-X certd-cli/cmd.Version=$(VERSION) \
	-X certd-cli/cmd.GitCommit=$(GIT_COMMIT) \
	-X certd-cli/cmd.BuildTime=$(BUILD_TIME)

help:
	@echo "可用命令:"
	@echo "  make build       - 构建当前平台"
	@echo "  make build-all   - 构建所有平台"
	@echo "  make clean       - 清理构建产物"
	@echo "  make version     - 显示版本信息"
	@echo "  make info        - 显示当前平台信息"

build-all: clean
	@echo "========================================"
	@echo "构建 $(APP_NAME)"
	@echo "版本: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"
	@echo "========================================"
	@mkdir -p $(OUTPUT_DIR)
	@$(MAKE) build-darwin-amd64
	@$(MAKE) build-darwin-arm64
	@$(MAKE) build-windows-amd64
	@$(MAKE) build-windows-arm64
	@$(MAKE) build-linux-amd64
	@$(MAKE) build-linux-arm64
	@echo ""
	@echo "========================================"
	@echo "所有平台构建完成!"
	@echo "========================================"
	@ls -lh $(OUTPUT_DIR)/

build:
	@echo "正在构建当前平台 ($(shell go env GOOS)/$(shell go env GOARCH))..."
	@mkdir -p $(OUTPUT_DIR)
	@if [ "$(shell go env GOOS)" = "windows" ]; then \
		CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o "$(OUTPUT_DIR)/$(APP_NAME)-$(shell go env GOOS)-$(shell go env GOARCH).exe" .; \
	else \
		CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o "$(OUTPUT_DIR)/$(APP_NAME)-$(shell go env GOOS)-$(shell go env GOARCH)" .; \
		chmod +x "$(OUTPUT_DIR)/$(APP_NAME)-$(shell go env GOOS)-$(shell go env GOARCH)"; \
	fi
	@echo "✓ 构建成功: $(OUTPUT_DIR)/$(APP_NAME)-$(shell go env GOOS)-$(shell go env GOARCH)"


build-darwin-amd64:
	@echo "正在构建 darwin/amd64..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o "$(OUTPUT_DIR)/$(APP_NAME)-darwin-amd64" .
	@chmod +x "$(OUTPUT_DIR)/$(APP_NAME)-darwin-amd64"
	@echo "✓ darwin/amd64 构建成功"

build-darwin-arm64:
	@echo "正在构建 darwin/arm64..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o "$(OUTPUT_DIR)/$(APP_NAME)-darwin-arm64" .
	@chmod +x "$(OUTPUT_DIR)/$(APP_NAME)-darwin-arm64"
	@echo "✓ darwin/arm64 构建成功"

build-windows-amd64:
	@echo "正在构建 windows/amd64..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o "$(OUTPUT_DIR)/$(APP_NAME)-windows-amd64.exe" .
	@echo "✓ windows/amd64 构建成功"

build-windows-arm64:
	@echo "正在构建 windows/arm64..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o "$(OUTPUT_DIR)/$(APP_NAME)-windows-arm64.exe" .
	@echo "✓ windows/arm64 构建成功"

build-linux-amd64:
	@echo "正在构建 linux/amd64..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o "$(OUTPUT_DIR)/$(APP_NAME)-linux-amd64" .
	@chmod +x "$(OUTPUT_DIR)/$(APP_NAME)-linux-amd64"
	@echo "✓ linux/amd64 构建成功"

build-linux-arm64:
	@echo "正在构建 linux/arm64..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o "$(OUTPUT_DIR)/$(APP_NAME)-linux-arm64" .
	@chmod +x "$(OUTPUT_DIR)/$(APP_NAME)-linux-arm64"
	@echo "✓ linux/arm64 构建成功"

clean:
	@echo "清理构建产物..."
	@rm -rf $(OUTPUT_DIR)
	@echo "✓ 清理完成"

version:
	@echo "版本: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"

info:
	@echo "当前平台信息:"
	@echo "  GOOS:   $(shell go env GOOS)"
	@echo "  GOARCH: $(shell go env GOARCH)"
	@echo "  GOPATH: $(shell go env GOPATH)"
	@echo "  GOROOT: $(shell go env GOROOT)"
