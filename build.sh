#!/bin/bash

# 跨平台构建脚本
# 使用方法: ./build.sh

set -e

echo "=== 跨平台构建脚本 ==="

# 创建输出目录
mkdir -p dist

echo "🔨 构建 Linux AMD64 版本..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/vegetable-price-linux-amd64

echo "🔨 构建 Windows AMD64 版本..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/vegetable-price-windows-amd64.exe

echo "🔨 构建 macOS AMD64 版本..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/vegetable-price-darwin-amd64

echo "🔨 构建当前平台版本..."
go build -ldflags="-s -w" -o dist/vegetable-price

echo "📦 创建部署包..."
cd dist

# 复制配置文件
cp ../config.json . 2>/dev/null || echo "⚠️  config.json 不存在，使用示例配置"
cp ../config.example.json . 2>/dev/null || echo "⚠️  config.example.json 不存在"
cp ../product_list.html . 2>/dev/null || echo "⚠️  product_list.html 不存在"

# Linux包
tar -czf vegetable-price-linux-amd64.tar.gz vegetable-price-linux-amd64 config.json config.example.json product_list.html  2>/dev/null || tar -czf vegetable-price-linux-amd64.tar.gz vegetable-price-linux-amd64

# Windows包
zip -r vegetable-price-windows-amd64.zip vegetable-price-windows-amd64.exe config.json config.example.json product_list.html  2>/dev/null || zip -r vegetable-price-windows-amd64.zip vegetable-price-windows-amd64.exe

# macOS包
tar -czf vegetable-price-darwin-amd64.tar.gz vegetable-price-darwin-amd64 config.json config.example.json product_list.html  2>/dev/null || tar -czf vegetable-price-darwin-amd64.tar.gz vegetable-price-darwin-amd64

# 当前平台包
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    tar -czf vegetable-price-current.tar.gz vegetable-price config.json config.example.json 2>/dev/null || tar -czf vegetable-price-current.tar.gz vegetable-price
elif [[ "$OSTYPE" == "darwin"* ]]; then
    tar -czf vegetable-price-current.tar.gz vegetable-price config.json config.example.json 2>/dev/null || tar -czf vegetable-price-current.tar.gz vegetable-price
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
    zip -r vegetable-price-current.zip vegetable-price config.json config.example.json 2>/dev/null || zip -r vegetable-price-current.zip vegetable-price
fi

cd ..

echo "✅ 构建完成！"
echo ""
echo "📁 输出目录: dist/"
echo "📦 生成的文件:"
ls -la dist/
echo ""
echo "🚀 部署说明:"
echo "   Linux:   tar -xzf vegetable-price-linux-amd64.tar.gz"
echo "   Windows: unzip vegetable-price-windows-amd64.zip"
echo "   macOS:   tar -xzf vegetable-price-darwin-amd64.tar.gz"
