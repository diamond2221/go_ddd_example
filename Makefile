# Makefile for Recommendation Service
#
# 常用命令：
# make help     - 显示帮助信息
# make gen      - 生成 Kitex 代码
# make build    - 编译服务
# make run      - 运行服务
# make test     - 运行测试
# make clean    - 清理构建产物

.PHONY: help gen build run test clean docker

# 默认目标
.DEFAULT_GOAL := help

# 服务名称
SERVICE_NAME := recommendation-service

# 帮助信息
help: ## 显示帮助信息
	@echo "Recommendation Service - Makefile Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""

# 生成代码
gen: ## 生成 Kitex 代码
	@echo "Generating Kitex code..."
	@bash script/bootstrap.sh

# 安装依赖
deps: ## 安装依赖
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# 编译
build: ## 编译服务
	@echo "Building $(SERVICE_NAME)..."
	@bash build.sh

# 运行
run: ## 运行服务
	@echo "Running $(SERVICE_NAME)..."
	@go run main.go

# 测试
test: ## 运行测试
	@echo "Running tests..."
	@go test ./... -v -cover

# 测试覆盖率
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# 单元测试
test-unit: ## 运行单元测试
	@echo "Running unit tests..."
	@go test ./tests/unit/... -v

# 集成测试
test-integration: ## 运行集成测试
	@echo "Running integration tests..."
	@go test ./tests/integration/... -v

# 代码检查
lint: ## 运行代码检查
	@echo "Running linter..."
	@golangci-lint run

# 格式化代码
fmt: ## 格式化代码
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

# 清理
clean: ## 清理构建产物
	@echo "Cleaning..."
	@rm -f $(SERVICE_NAME)
	@rm -rf output/
	@rm -f coverage.out coverage.html

# Docker 构建
docker-build: ## 构建 Docker 镜像
	@echo "Building Docker image..."
	@docker build -t $(SERVICE_NAME):latest .

# Docker 运行
docker-run: ## 运行 Docker 容器
	@echo "Running Docker container..."
	@docker run -p 8888:8888 $(SERVICE_NAME):latest

# 安装工具
install-tools: ## 安装开发工具
	@echo "Installing development tools..."
	@go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
	@go install github.com/cloudwego/thriftgo@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest

# 初始化项目
init: install-tools deps gen ## 初始化项目（安装工具、依赖、生成代码）
	@echo "Project initialized successfully!"

# 开发模式（热重载）
dev: ## 开发模式（需要安装 air）
	@echo "Starting development mode..."
	@air

# 查看依赖
deps-list: ## 查看依赖列表
	@go list -m all

# 更新依赖
deps-update: ## 更新依赖
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# 检查依赖安全性
deps-check: ## 检查依赖安全性
	@echo "Checking dependencies..."
	@go list -json -m all | nancy sleuth
