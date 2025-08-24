# 配置文件说明

## 概述

程序现在使用JSON配置文件来管理所有设置，包括URL、Cookie、User-Agent等。这样可以更方便地管理不同环境的配置，而不需要修改代码。

## 配置文件结构

### 主配置文件：`config.json`

```json
{
  "url": "https://www.fengzhansy.com/wchyzyg/wap.shtml?method=ztmodel&ztid=gfl00%E7%93%9C%E6%9E%9C%E8%8A%B1%E8%8F%9C%E7%B1%BB",
  "cookie": "shdzarea=%E6%96%87%E5%8D%8E%E8%B7%AF; scsmdid=012; shdzmdname=%E5%87%A4%E5%B1%95%E8%B6%85%E5%B8%82%E6%96%87%E5%8D%8E%E8%B7%AF%E5%BA%97",
  "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
  "timeout": 30,
  "retry_count": 3
}
```

### 配置项说明

| 配置项 | 类型 | 必需 | 说明 | 默认值 |
|--------|------|------|------|--------|
| `url` | string | ✅ | 要访问的网页URL | 无 |
| `cookie` | string | ✅ | 认证Cookie字符串 | 无 |
| `user_agent` | string | ❌ | 浏览器User-Agent | Mozilla/5.0... |
| `timeout` | int | ❌ | HTTP请求超时时间（秒） | 30 |
| `retry_count` | int | ❌ | 请求失败重试次数 | 3 |

## 快速开始

### 1. 复制示例配置文件

```bash
cp config.example.json config.json
```

### 2. 编辑配置文件

```bash
# 使用文本编辑器编辑
nano config.json
# 或使用VS Code
code config.json
```

### 3. 修改配置项

```json
{
  "url": "https://your-target-website.com/product",
  "cookie": "your-cookie-string",
  "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
  "timeout": 60,
  "retry_count": 5
}
```

## 配置管理

### 多环境配置

你可以为不同环境创建不同的配置文件：

```bash
# 开发环境
cp config.json config.dev.json

# 生产环境
cp config.json config.prod.json

# 测试环境
cp config.json config.test.json
```

然后在代码中指定配置文件：

```go
config, err := LoadConfig("config.prod.json")
```

### 环境变量支持

如果需要，可以扩展配置系统支持环境变量：

```json
{
  "url": "${TARGET_URL}",
  "cookie": "${AUTH_COOKIE}",
  "timeout": "${TIMEOUT:-30}"
}
```

## 安全注意事项

### 1. Cookie安全

- 不要在代码仓库中提交包含真实Cookie的配置文件
- 将`config.json`添加到`.gitignore`
- 定期更新Cookie以保持登录状态

### 2. 配置文件权限

```bash
# 设置适当的文件权限
chmod 600 config.json
```

### 3. 敏感信息处理

对于生产环境，考虑使用环境变量或密钥管理服务：

```bash
export TARGET_URL="https://example.com"
export AUTH_COOKIE="your-cookie"
```

## 故障排除

### 常见问题

1. **配置文件不存在**
   ```
   加载配置文件失败: open config.json: no such file or directory
   使用默认配置...
   ```
   解决方案：复制示例配置文件

2. **JSON格式错误**
   ```
   解析配置文件失败: invalid character '}' looking for beginning of value
   ```
   解决方案：检查JSON语法，确保格式正确

3. **缺少必需配置项**
   ```
   配置文件中缺少URL
   ```
   解决方案：确保`url`和`cookie`字段存在且不为空

### 调试配置

使用调试模式查看当前配置：

```bash
./vegetable-price debug
```

输出会显示：
- 目标URL
- 超时设置
- 重试次数

## 配置示例

### 凤展超市配置

```json
{
  "url": "https://www.fengzhansy.com/wchyzyg/wap.shtml?method=ztmodel&ztid=gfl00%E7%93%9C%E6%9E%9C%E8%8A%B1%E8%8F%9C%E7%B1%BB",
  "cookie": "shdzarea=%E6%96%87%E5%8D%8E%E8%B7%AF; scsmdid=012; shdzmdname=%E5%87%A4%E5%B1%95%E8%B6%85%E5%B8%82%E6%96%87%E5%8D%8E%E8%B7%AF%E5%BA%97",
  "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
  "timeout": 30,
  "retry_count": 3
}
```

### 其他网站配置

```json
{
  "url": "https://other-supermarket.com/products",
  "cookie": "session=abc123; user=test",
  "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
  "timeout": 45,
  "retry_count": 5
}
```

## 最佳实践

1. **版本控制**: 将`config.example.json`提交到代码仓库，但不要提交`config.json`
2. **文档化**: 在团队中维护配置项的说明文档
3. **测试**: 在修改配置后，使用调试模式验证配置是否正确
4. **备份**: 定期备份重要的配置文件
5. **监控**: 监控配置文件的修改和访问

## 扩展配置

如果需要添加更多配置项，可以：

1. 在`Config`结构体中添加新字段
2. 在`LoadConfig`函数中设置默认值
3. 在相关代码中使用新配置项
4. 更新文档和示例配置文件
