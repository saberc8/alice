# Alice Go Backend Makefile

.PHONY: build run clean test fmt vet deps help docker lint security coverage init-data build-init

# 默认目标
.DEFAULT_GOAL := help

# 项目配置
APP_NAME := alice
BUILD_DIR := ./bin
CONFIG_FILE := config.yaml
DOCKER_IMAGE := alice
DOCKER_TAG := latest

# 构建变量
VERSION?=$(shell git describe --tags --always --dirty)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go 相关配置
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt

# LDFLAGS for build
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}"
GOVET := $(GOCMD) vet

## help: 显示帮助信息
help:
	@echo "Alice Go Backend Project (with RBAC System)"
	@echo ""
	@echo "Available commands:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## deps: 下载依赖
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## fmt: 格式化代码
fmt:
	$(GOFMT) ./...

## vet: 静态分析
vet:
	$(GOVET) ./...

## lint: 运行代码检查
lint:
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin)
	golangci-lint run

## security: 运行安全扫描
security:
	@echo "Running security scan..."
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
	gosec ./...

## test: 运行测试
test:
	$(GOTEST) -v -race ./...

## test-coverage: 运行测试并生成覆盖率报告
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## coverage: 显示测试覆盖率
coverage: test-coverage
	@echo "Coverage report generated: coverage.html"
	$(GOCMD) tool cover -func=coverage.out | tail -1

## build: 编译项目
build: deps fmt vet
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./main.go

## build-linux: 交叉编译Linux版本
build-linux: deps fmt vet
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux ./main.go

## build-local: 本地开发编译
build-local: clean
	@echo "Building $(APP_NAME) for local development..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) -v

## run: 运行项目
run: build
	./$(BUILD_DIR)/$(APP_NAME)

## build-init: 编译数据初始化脚本
build-init: deps fmt vet
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/init cmd/init/main.go

## init-data: 运行数据初始化脚本（初始化RBAC数据）
init-data: build-init
	@echo "Initializing RBAC data..."
	./$(BUILD_DIR)/init
	@echo "RBAC data initialization completed!"

## rbac-setup: 完整的RBAC系统设置
rbac-setup: build build-init init-data
	@echo "RBAC system setup completed!"

## dev: 开发模式运行（使用go run）
dev:
	$(GOCMD) run main.go

## docker-build: 构建Docker镜像
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

## docker-run: 运行Docker容器
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8090:8090 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

## db-setup: 设置开发数据库
db-setup:
	@echo "Setting up development database..."
	docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Database is ready!"

## setup: 初始项目设置
setup: deps db-setup
	@echo "Creating config file..."
	@cp config.yaml.example config.yaml 2>/dev/null || echo "config.yaml already exists"
	@echo "Project setup complete!"

## check: 运行所有检查
check: fmt vet lint security test
	@echo "All checks passed!"

## ci: 本地运行CI流程
ci: deps check build
	@echo "CI pipeline completed successfully!"

## clean: 清理构建文件
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## docker-build: 构建Docker镜像
docker-build:
	docker build -t $(APP_NAME):latest .

## docker-run: 运行Docker容器
docker-run:
	docker run -p 8090:8090 --env-file .env $(APP_NAME):latest

## db-setup: 设置PostgreSQL数据库（需要Docker）
db-setup:
	docker run --name alice-postgres -e POSTGRES_DB=alice -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:13

## db-start: 启动数据库容器
db-start:
	docker start alice-postgres

## db-stop: 停止数据库容器
db-stop:
	docker stop alice-postgres

## db-remove: 删除数据库容器
db-remove:
	docker rm alice-postgres

## install: 安装到系统
install: build
	cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/

## check: 检查代码质量
check: fmt vet test

## all: 执行完整的构建和测试流程
all: clean deps fmt vet test build
