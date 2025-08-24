# 发布说明

## 版本信息

- **版本**: 1.0.0
- **构建时间**: $(date)
- **Go版本**: 1.24.4
- **架构**: AMD64

## 可执行文件

### Linux AMD64
- **文件**: `vegetable-price-linux-amd64`
- **大小**: ~6.2MB
- **格式**: ELF 64-bit LSB executable
- **部署包**: `vegetable-price-linux-amd64.tar.gz`

### Windows AMD64
- **文件**: `vegetable-price-windows-amd64.exe`
- **大小**: ~6.4MB
- **格式**: PE32+ executable (GUI) Intel 80386
- **部署包**: `vegetable-price-windows-amd64.zip`

### macOS AMD64
- **文件**: `vegetable-price-darwin-amd64`
- **大小**: ~6.4MB
- **格式**: Mach-O 64-bit executable x86_64
- **部署包**: `vegetable-price-darwin-amd64.tar.gz`

### 当前平台
- **文件**: `vegetable-price`
- **大小**: ~6.0MB
- **部署包**: `vegetable-price-current.tar.gz`

## 部署说明

### Linux 用户
```bash
# 下载并解压
wget https://github.com/your-repo/releases/download/v1.0.0/vegetable-price-linux-amd64.tar.gz
tar -xzf vegetable-price-linux-amd64.tar.gz

# 运行程序
./vegetable-price-linux-amd64
```

### Windows 用户
```bash
# 下载并解压
# 使用Windows资源管理器解压 vegetable-price-windows-amd64.zip

# 运行程序
vegetable-price-windows-amd64.exe
```

### macOS 用户
```bash
# 下载并解压
curl -LO https://github.com/your-repo/releases/download/v1.0.0/vegetable-price-darwin-amd64.tar.gz
tar -xzf vegetable-price-darwin-amd64.tar.gz

# 运行程序
./vegetable-price-darwin-amd64
```

## 功能特性

- ✅ 支持带Cookie认证的HTTP GET请求
- ✅ 自动解析HTML页面中的商品信息
- ✅ 智能提取价格和规格信息
- ✅ 自动计算每斤价格
- ✅ 支持多种重量单位（斤、公斤、千克、克、两、磅）
- ✅ 批量解析多个商品（支持`index_picAD`容器结构）
- ✅ 内置调试模式，帮助分析HTML结构
- ✅ 支持阿里云函数计算部署
- ✅ HTTP API接口，支持多种调用模式

## 系统要求

### 最低要求
- **操作系统**: Linux 2.6+, Windows 7+, macOS 10.12+
- **架构**: x86_64 (AMD64)
- **内存**: 64MB RAM
- **磁盘**: 10MB 可用空间

### 推荐配置
- **操作系统**: Linux 4.x+, Windows 10+, macOS 11+
- **架构**: x86_64 (AMD64)
- **内存**: 256MB RAM
- **磁盘**: 50MB 可用空间
- **网络**: 稳定的互联网连接

## 配置文件

### 必需配置
- `config.json` - 包含URL和Cookie的配置文件

### 配置示例
```json
{
  "url": "https://example.com/product",
  "cookie": "your-cookie-string",
  "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
  "timeout": 30,
  "retry_count": 3
}
```

## 使用方法

### 命令行模式
```bash
# 正常运行
./vegetable-price

# 测试模式
./vegetable-price test

# 调试模式
./vegetable-price debug
```

### 阿里云函数模式
```bash
# 构建云函数
make build-fc

# 部署到阿里云
./deploy.sh
```

## 故障排除

### 常见问题

1. **权限错误**
   ```bash
   chmod +x vegetable-price-linux-amd64
   ```

2. **配置文件缺失**
   ```bash
   cp config.example.json config.json
   # 编辑配置文件填入真实信息
   ```

3. **网络连接问题**
   - 检查防火墙设置
   - 验证网络连接
   - 确认目标网站可访问

### 日志和调试

- 使用 `debug` 模式获取详细信息
- 检查网络请求和响应
- 验证HTML解析结果

## 更新日志

### v1.0.0 (2024-08-24)
- 🎉 首次发布
- ✨ 支持多种重量单位计算
- ✨ 智能识别包装商品和重量商品
- ✨ 批量解析多个商品
- ✨ 阿里云函数计算支持
- ✨ 跨平台AMD64可执行文件
- 📚 完整的文档和示例

## 技术支持

- **文档**: [README.md](README.md)
- **配置说明**: [CONFIG.md](CONFIG.md)
- **使用指南**: [USAGE.md](USAGE.md)
- **阿里云部署**: [ALIYUN_FC.md](ALIYUN_FC.md)
- **问题反馈**: [GitHub Issues](https://github.com/your-repo/issues)

## 许可证

本项目仅供学习和研究使用。请遵守相关网站的使用条款和robots.txt规定。


