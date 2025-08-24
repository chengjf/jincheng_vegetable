#!/bin/bash

# 阿里云函数部署脚本
# 使用方法: ./deploy.sh

set -e

echo "=== 阿里云函数部署脚本 ==="

# 检查是否安装了阿里云CLI
if ! command -v aliyun &> /dev/null; then
    echo "❌ 未安装阿里云CLI，请先安装："
    echo "   https://help.aliyun.com/document_detail/121541.html"
    exit 1
fi

# 检查是否已配置阿里云账号
if ! aliyun sts GetCallerIdentity &> /dev/null; then
    echo "❌ 未配置阿里云账号，请先配置："
    echo "   aliyun configure"
    exit 1
fi

echo "✅ 阿里云CLI已配置"

# 构建程序
echo "🔨 构建程序..."
go build -o main

# 创建部署包
echo "📦 创建部署包..."
zip -r function.zip main config.json config.go

# 部署函数
echo "🚀 部署函数到阿里云..."
aliyun fun deploy --template-file template.yml

echo "✅ 部署完成！"
echo ""
echo "📋 部署信息："
echo "   函数名称: vegetable-price-function"
echo "   服务名称: vegetable-price-service"
echo "   运行时: go1.x"
echo "   超时时间: 60秒"
echo "   内存: 512MB"
echo ""
echo "🌐 访问方式："
echo "   1. 通过阿里云控制台获取函数URL"
echo "   2. 使用HTTP触发器调用"
echo ""
echo "📖 调用示例："
echo "   POST /"
echo "   Content-Type: application/json"
echo "   {"
echo "     \"url\": \"https://example.com\","
echo "     \"cookie\": \"your-cookie\","
echo "     \"mode\": \"normal\""
echo "   }"
