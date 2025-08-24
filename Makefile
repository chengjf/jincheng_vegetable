# 晋城蔬菜价格查询工具 Makefile

.PHONY: build run test clean help

# 默认目标
all: build

# 构建程序
build:
	@echo "构建程序..."
	go build -o vegetable-price
	@echo "构建完成: vegetable-price"

# 构建阿里云函数
build-fc:
	@echo "构建阿里云函数..."
	go build -o main
	@echo "构建完成: main (阿里云函数)"

# 跨平台构建
build-all: build-fc
	@echo "跨平台构建..."
	./build.sh
	@echo "跨平台构建完成"

# 运行程序
run: build
	@echo "运行程序..."
	./vegetable-price

# 运行测试模式
test: build
	@echo "运行测试模式..."
	./vegetable-price test

# 运行调试模式
debug: build
	@echo "运行调试模式..."
	./vegetable-price debug

# 构建并测试阿里云函数
test-fc: build-fc
	@echo "测试阿里云函数..."
	go run test_function.go

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -f vegetable-price
	@echo "清理完成"

# 安装依赖
deps:
	@echo "安装依赖..."
	go mod tidy
	go mod download

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run

# 运行Go测试
go-test:
	@echo "运行Go测试..."
	go test -v ./...

# 显示帮助信息
help:
	@echo "晋城蔬菜价格查询工具 - Makefile 帮助"
	@echo ""
	@echo "可用命令:"
	@echo "  build    构建程序 (默认)"
	@echo "  run      构建并运行程序"
	@echo "  test     构建并运行测试模式"
	@echo "  debug    构建并运行调试模式"
	@echo "  build-fc 构建阿里云函数"
	@echo "  test-fc  构建并测试阿里云函数"
	@echo "  build-all跨平台构建所有版本"
	@echo "  clean    清理构建文件"
	@echo "  deps     安装依赖"
	@echo "  fmt      格式化代码"
	@echo "  lint     代码检查"
	@echo "  go-test  运行Go测试"
	@echo "  help     显示此帮助信息"
	@echo ""
	@echo "示例:"
	@echo "  make build    # 构建程序"
	@echo "  make test     # 运行测试模式"
	@echo "  make debug    # 运行调试模式"
	@echo "  make build-fc # 构建阿里云函数"
	@echo "  make test-fc  # 测试阿里云函数"
	@echo "  make build-all# 跨平台构建"
	@echo "  make clean    # 清理文件"
