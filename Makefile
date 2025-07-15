# Makefile for bmstock

# 获取当前时间，格式：YYYYMMDDHHMM
NOW := $(shell date +"%Y%m%d%H%M")

# 获取当前时间（格式：YmdH）
TIMESTAMP := $(shell date "+%Y%m%d%H")
# 计算 MD5（Linux 使用 md5sum，macOS 使用 md5）
ifneq (,$(wildcard /etc/os-release))
    MD5HASH := $(shell echo -n $(TIMESTAMP) | md5sum | awk '{print $$1}')
else
    MD5HASH := $(shell echo -n $(TIMESTAMP) | md5)
endif

ENV ?= dev
ifeq ($(ENV),prod)
    UPLOAD_HOST := https://api.prod.com
else
    UPLOAD_HOST := https://api.dev.com
endif

# Variables
BINARY_NAME=bmstock
BIN_DIR=bin
SRC_DIR=cmd
# 强制立即求值，并确保 fallback 到 dev
VERSION := $(shell git describe --tags --always 2>/dev/null)
ifeq ($(VERSION),)
    VERSION = dev
endif
FULL_VERSION ?= $(VERSION)-$(NOW)

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: version
version:
	@echo Current version: $(FULL_VERSION)

# Build the binary
# 优先级: ldflags > tag > dev
# 方式1（ldflags）: make VERSION=v1.2.3 build
# 方式2（读取tag）: make build
# 方式2（默认值dev）: make build
.PHONY: clean build
build:
	@CGO_ENABLE=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME) -ldflags "-X 'main.version=$(FULL_VERSION)'" $(SRC_DIR)/main.go
	upx $(BIN_DIR)/$(BINARY_NAME)

# 交叉编译
.PHONY: clean build-cross
build-cross:
	@CGO_ENABLE=0 GO111MODULE=on GOOS=linux   GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-amd-linux   -ldflags "-X 'main.version=$(FULL_VERSION)'" $(SRC_DIR)/main.go
	@CGO_ENABLE=0 GO111MODULE=on GOOS=darwin  GOARCH=arm64 go build -o $(BIN_DIR)/$(BINARY_NAME)-arm-darwin  -ldflags "-X 'main.version=$(FULL_VERSION)'" $(SRC_DIR)/main.go
	@CGO_ENABLE=0 GO111MODULE=on GOOS=darwin  GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-amd-darwin  -ldflags "-X 'main.version=$(FULL_VERSION)'" $(SRC_DIR)/main.go
	@CGO_ENABLE=0 GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME).exe     	 -ldflags "-X 'main.version=$(FULL_VERSION)'" $(SRC_DIR)/main.go

# 上传编译包
# make deploy
# ENV=prod make deploy
.PHONY: deploy
upload:
	curl -Xi POST $(UPLOAD_HOST)/api/v1/public/deploy \
		-F "file=@./$(BIN_DIR)/$(BINARY_NAME)" \
		-H "Authorization: Bearer $(MD5HASH)"

# Run the application
.PHONY: run
run:
	@go run $(SRC_DIR)/main.go

# Clean build artifacts
.PHONY: clean
clean:
	@rm -rf $(BIN_DIR)/*
	@rm -rf logs/*

# Run tests
.PHONY: test
test:
	@go test ./... -v

# Generate Swagger docs
.PHONY: doc
doc:
	@swag init -g cmd/main.go -o doc

# Build and run migrations
.PHONY: migrate
migrate:
	@go run scripts/migrate.go