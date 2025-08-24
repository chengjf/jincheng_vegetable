# 晋城蔬菜价格查询工具

这是一个Go语言编写的网页爬虫工具，用于获取商品信息并计算每斤价格。

## 功能特性

- 🔍 支持带Cookie认证的HTTP GET请求
- 📊 自动解析HTML页面中的商品信息
- 💰 智能提取价格和规格信息
- ⚖️ 自动计算每斤价格
- 🏷️ 支持多种重量单位（斤、公斤、千克、克、两、磅）
- 📦 批量解析多个商品（支持`index_picAD`容器结构）
- 🐛 内置调试模式，帮助分析HTML结构
- ☁️ 支持阿里云函数计算部署
- 🌐 HTTP API接口，支持多种调用模式

## 使用方法

### 1. 配置URL和Cookie

创建 `config.json` 配置文件：

```json
{
  "url": "https://your-target-website.com/product",
  "cookie": "your-cookie-string",
  "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
  "timeout": 30,
  "retry_count": 3
}
```

或者复制示例配置文件：

```bash
cp config.example.json config.json
# 然后编辑 config.json 文件
```

### 2. 运行程序

```bash
# 正常模式 - 获取所有商品信息
./vegetable-price

# 测试模式 - 测试程序功能
./vegetable-price test

# 调试模式 - 调试HTML解析
./vegetable-price debug
```

或者使用Makefile：

```bash
make run     # 正常运行
make test    # 测试模式
make debug   # 调试模式
```

## 程序结构

- `main.go` - 主程序文件，包含所有核心功能
- `config.go` - 配置管理模块
- `config.json` - 主配置文件（包含URL、Cookie等）
- `config.example.json` - 配置文件示例
- `template.yml` - 阿里云函数部署配置
- `deploy.sh` - 阿里云函数部署脚本
- `test_function.go` - 阿里云函数测试客户端
- `ALIYUN_FC.md` - 阿里云函数部署指南
- `go.mod` - Go模块依赖文件

## 核心功能说明

### 1. 网页请求
- 使用Go标准库的`net/http`包
- 支持自定义Cookie和User-Agent
- 自动处理HTTP响应

### 2. HTML解析
- 使用`golang.org/x/net/html`包解析HTML
- 支持多种CSS选择器策略
- 智能提取文本内容

### 3. 价格计算
- 自动识别重量单位
- 统一转换为斤为单位
- 计算每斤价格

## 支持的重量单位

| 单位 | 转换比例（到斤） |
|------|------------------|
| 斤   | 1.0             |
| 公斤 | 2.0             |
| 千克 | 2.0             |
| 克   | 0.002           |
| 两   | 0.1             |
| 磅   | 0.907           |

## 注意事项

1. **合法使用**: 请确保遵守目标网站的robots.txt和使用条款
2. **频率控制**: 建议添加适当的请求间隔，避免对服务器造成压力
3. **Cookie管理**: 定期更新Cookie以保持登录状态
4. **错误处理**: 程序包含完善的错误处理机制

## 配置文件说明

### config.json 配置项

| 配置项 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| `url` | string | 要访问的网页URL | 必需 |
| `cookie` | string | 认证Cookie字符串 | 必需 |
| `user_agent` | string | 浏览器User-Agent | Mozilla/5.0... |
| `timeout` | int | HTTP请求超时时间（秒） | 30 |
| `retry_count` | int | 请求失败重试次数 | 3 |

### 配置文件示例

```json
{
  "url": "https://www.fengzhansy.com/wchyzyg/wap.shtml?method=ztmodel&ztid=gfl00%E7%93%9C%E6%9E%9C%E8%8A%B1%E8%8F%9C%E7%B1%BB",
  "cookie": "shdzarea=%E6%96%87%E5%8D%8E%E8%B7%AF; scsmdid=012; shdzmdname=%E5%87%A4%E5%B1%95%E8%B6%85%E5%B8%82%E6%96%87%E5%8D%8E%E8%B7%AF%E5%BA%97",
  "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
  "timeout": 30,
  "retry_count": 3
}
```

## 部署方式

### 本地运行

```bash
make run     # 正常运行
make test    # 测试模式
make debug   # 调试模式
```

### 阿里云函数部署

```bash
make build-fc  # 构建阿里云函数
./deploy.sh    # 部署到阿里云
make test-fc   # 测试阿里云函数
```

详细部署说明请参考：[ALIYUN_FC.md](ALIYUN_FC.md)

## 依赖包

- `golang.org/x/net/html` - HTML解析
- Go标准库：`net/http`, `regexp`, `strconv`, `strings`, `encoding/json`

## 许可证

本项目仅供学习和研究使用。
