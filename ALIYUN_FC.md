# 阿里云函数部署指南

## 概述

本项目已适配阿里云函数计算（Function Compute），可以在阿里云上以无服务器方式运行蔬菜价格查询服务。

## 架构特点

- **无服务器**: 按需运行，按实际使用量计费
- **自动扩缩容**: 根据请求量自动调整实例数量
- **HTTP触发器**: 支持HTTP/HTTPS访问
- **CORS支持**: 支持跨域请求
- **多种运行模式**: normal、debug、test三种模式

## 部署前准备

### 1. 安装阿里云CLI

```bash
# macOS
brew install aliyun-cli

# Linux
curl -o aliyun-cli-linux-latest-amd64.tgz https://aliyuncli.alicdn.com/aliyun-cli-linux-latest-amd64.tgz
tar xzvf aliyun-cli-linux-latest-amd64.tgz
sudo mv aliyun /usr/local/bin/

# Windows
# 下载 https://aliyuncli.alicdn.com/aliyun-cli-windows-latest-amd64.zip
```

### 2. 配置阿里云账号

```bash
aliyun configure
# 输入AccessKey ID、AccessKey Secret、默认地域等
```

### 3. 安装Fun工具

```bash
npm install -g @alicloud/fun
```

## 部署步骤

### 1. 自动部署（推荐）

```bash
# 给部署脚本添加执行权限
chmod +x deploy.sh

# 运行部署脚本
./deploy.sh
```

### 2. 手动部署

```bash
# 构建程序
go build -o main

# 创建部署包
zip -r function.zip main config.json config.go

# 部署函数
aliyun fun deploy --template-file template.yml
```

## 配置说明

### template.yml 配置

```yaml
Resources:
  vegetable-price-service:          # 服务名称
    Type: 'Aliyun::Serverless::Service'
    Properties:
      Description: '晋城蔬菜价格查询服务'
      
  vegetable-price-function:         # 函数名称
    Type: 'Aliyun::Serverless::Function'
    Properties:
      Handler: 'main.HandleRequest'  # 入口函数
      Runtime: 'go1.x'               # 运行时
      Timeout: 60                    # 超时时间（秒）
      MemorySize: 512                # 内存大小（MB）
      
      Events:
        httpTrigger:                 # HTTP触发器
          Type: 'HTTP'
          Properties:
            AuthType: 'ANONYMOUS'    # 匿名访问
            Methods: ['GET', 'POST', 'OPTIONS']
            Cors:                     # 跨域配置
              AllowOrigin: '*'
              AllowMethods: 'GET,POST,OPTIONS'
```

## 调用方式

### HTTP请求格式

```bash
POST https://your-function-url.fc.aliyuncs.com/2016-08-15/proxy/vegetable-price-service/vegetable-price-function/
Content-Type: application/json

{
  "url": "https://www.fengzhansy.com/wchyzyg/wap.shtml?method=ztmodel&ztid=gfl00%E7%93%9C%E6%9E%9C%E8%8A%B1%E8%8F%9C%E7%B1%BB",
  "cookie": "your-cookie-string",
  "mode": "normal"
}
```

### 请求参数

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `url` | string | ✅ | 要访问的网页URL |
| `cookie` | string | ✅ | 认证Cookie |
| `mode` | string | ❌ | 运行模式：normal/debug/test |

### 响应格式

```json
{
  "statusCode": 200,
  "headers": {
    "Content-Type": "application/json; charset=utf-8",
    "Access-Control-Allow-Origin": "*"
  },
  "body": "{\"mode\":\"normal\",\"status\":\"success\",\"total_products\":33,\"products\":[...]}"
}
```

## 运行模式

### 1. Normal模式（默认）

获取商品信息并返回JSON格式结果：

```json
{
  "mode": "normal",
  "url": "https://example.com",
  "message": "商品信息获取完成",
  "status": "success",
  "total_products": 33,
  "products": [
    {
      "name": "线椒（盒）",
      "price": 5.8,
      "spec": "1盒",
      "price_per_jin": 5.8,
      "is_packaged": true,
      "unit": "元/盒"
    }
  ]
}
```

### 2. Debug模式

调试网页解析过程：

```json
{
  "mode": "debug",
  "url": "https://example.com",
  "message": "调试模式执行完成",
  "html_length": 106103,
  "status": "success",
  "container_found": true,
  "product_count": 33
}
```

### 3. Test模式

测试程序功能：

```json
{
  "mode": "test",
  "message": "功能测试完成",
  "tests": {
    "weight_parsing": "通过",
    "price_cleaning": "通过",
    "price_calculation": "通过"
  }
}
```

## 测试部署

### 1. 使用测试客户端

```bash
# 修改test_function.go中的函数URL
# 然后运行测试
go run test_function.go
```

### 2. 使用curl测试

```bash
# 测试正常模式
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","cookie":"your-cookie","mode":"normal"}' \
  https://your-function-url.fc.aliyuncs.com/2016-08-15/proxy/vegetable-price-service/vegetable-price-function/

# 测试调试模式
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","cookie":"your-cookie","mode":"debug"}' \
  https://your-function-url.fc.aliyuncs.com/2016-08-15/proxy/vegetable-price-service/vegetable-price-function/
```

## 监控和日志

### 1. 阿里云控制台

- 函数计算控制台：https://fc.console.aliyun.com/
- 查看函数执行日志
- 监控函数调用次数和性能

### 2. 日志查看

```bash
# 查看函数日志
aliyun fun logs --tail
```

### 3. 性能监控

- 执行时间
- 内存使用
- 调用次数
- 错误率

## 成本优化

### 1. 内存配置

根据实际需求调整内存大小：

```yaml
MemorySize: 256  # 降低到256MB
```

### 2. 超时设置

合理设置超时时间：

```yaml
Timeout: 30  # 降低到30秒
```

### 3. 冷启动优化

- 使用预留实例
- 合理设置并发度
- 优化代码执行效率

## 故障排除

### 常见问题

1. **函数部署失败**
   - 检查阿里云CLI配置
   - 验证账号权限
   - 检查网络连接

2. **函数执行超时**
   - 增加超时时间
   - 优化网络请求
   - 检查目标网站响应

3. **内存不足**
   - 增加内存配置
   - 优化代码逻辑
   - 减少并发处理

4. **CORS错误**
   - 检查CORS配置
   - 验证请求头设置
   - 测试跨域访问

### 调试技巧

1. **使用Debug模式**
   - 获取详细的执行信息
   - 分析HTML解析过程
   - 定位问题所在

2. **查看函数日志**
   - 分析执行流程
   - 识别错误原因
   - 优化性能瓶颈

3. **本地测试**
   - 在本地环境验证
   - 模拟函数调用
   - 调试配置问题

## 扩展功能

### 1. 添加认证

```yaml
Events:
  httpTrigger:
    Properties:
      AuthType: 'FUNCTION'  # 函数认证
```

### 2. 配置环境变量

```yaml
EnvironmentVariables:
  DEBUG_MODE: 'true'
  LOG_LEVEL: 'info'
```

### 3. 集成其他服务

- 阿里云OSS：存储结果
- 阿里云RDS：数据库存储
- 阿里云SLS：日志服务

## 最佳实践

1. **安全性**
   - 不要在代码中硬编码敏感信息
   - 使用环境变量管理配置
   - 定期更新Cookie和认证信息

2. **性能优化**
   - 合理设置内存和超时
   - 优化网络请求
   - 使用连接池

3. **监控告警**
   - 设置错误率告警
   - 监控执行时间
   - 跟踪调用量变化

4. **版本管理**
   - 使用Git管理代码
   - 标记发布版本
   - 回滚机制

## 联系支持

- 阿里云函数计算文档：https://help.aliyun.com/product/50980.html
- 阿里云技术支持：https://help.aliyun.com/
- 项目Issues：https://github.com/your-repo/issues
