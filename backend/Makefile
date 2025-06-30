# Go 参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 二进制文件名
BINARY_NAME=drone-scanner
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe

# 构建目录
BUILD_DIR=../dist/backend

.PHONY: all build clean test deps run dev docker

all: test build

build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v ./main.go

build-linux:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v ./main.go

build-windows:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_WINDOWS) -v ./main.go

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

deps:
	$(GOMOD) download
	$(GOMOD) tidy

run:
	$(GOCMD) run main.go

dev:
	air -c .air.toml

docker:
	docker build -t drone-scanner-backend .

# 安装开发工具
install-tools:
	$(GOGET) github.com/cosmtrek/air@latest
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 代码检查
lint:
	golangci-lint run

# 格式化代码
fmt:
	$(GOCMD) fmt ./...

# 生成 protobuf 文件
protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/*.proto

# 数据库迁移
migrate-up:
	$(GOCMD) run cmd/migrate/main.go up

migrate-down:
	$(GOCMD) run cmd/migrate/main.go down

# 生成模拟数据
seed:
	$(GOCMD) run cmd/seed/main.go

# 性能测试
bench:
	$(GOTEST) -bench=. -benchmem ./...

# 安全检查
security:
	gosec ./...

# 依赖检查
deps-check:
	$(GOCMD) list -u -m all

# 更新依赖
deps-update:
	$(GOGET) -u ./...
	$(GOMOD) tidy
