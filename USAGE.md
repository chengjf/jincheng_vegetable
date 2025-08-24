# 使用说明

## 快速开始

### 1. 运行程序

```bash
# 正常模式 - 访问网页获取商品信息
./vegetable-price

# 测试模式 - 测试程序功能
./vegetable-price test
```

### 2. 配置你的目标网站

编辑 `main.go` 文件中的以下部分：

```go
url := "https://your-target-website.com/product" // 替换为实际URL
cookie := "your-cookie-string" // 替换为实际cookie
```

## 如何获取Cookie

### 方法1：浏览器开发者工具

1. 在浏览器中登录目标网站
2. 按 `F12` 打开开发者工具
3. 切换到 `Network` 标签页
4. 刷新页面
5. 找到任意一个请求，查看请求头中的 `Cookie` 字段
6. 复制整个Cookie值

### 方法2：浏览器扩展

安装类似 "Cookie Editor" 的浏览器扩展，可以更方便地查看和复制Cookie。

## 支持的网站类型

### 1. 超市网站
- 凤展超市（已配置）
- 其他连锁超市

### 2. 电商平台
- 淘宝、天猫
- 京东
- 拼多多

### 3. 生鲜平台
- 盒马鲜生
- 每日优鲜
- 叮咚买菜

## 自定义HTML解析规则

如果程序无法正确提取商品信息，你可能需要调整HTML解析规则。

### 1. 商品名称提取

在 `parseHTML` 函数中修改：

```go
// 提取商品名称
product.Name = extractText(doc, "h1", "class", "product-title")
if product.Name == "" {
    product.Name = extractText(doc, "h1", "", "")
}
// 可以添加更多选择器
if product.Name == "" {
    product.Name = extractText(doc, "div", "class", "product-name")
}
```

### 2. 价格提取

```go
// 提取价格
priceText := extractText(doc, "span", "class", "price")
if priceText == "" {
    priceText = extractText(doc, "span", "", "")
}
// 添加更多选择器
if priceText == "" {
    priceText = extractText(doc, "div", "class", "current-price")
}
```

### 3. 规格提取

```go
// 提取规格
product.Spec = extractText(doc, "span", "class", "spec")
if product.Spec == "" {
    product.Spec = extractText(doc, "div", "class", "product-spec")
}
// 添加更多选择器
if product.Spec == "" {
    product.Spec = extractText(doc, "td", "", "")
}
```

## 重量单位支持

程序自动识别以下重量单位并转换为斤：

| 单位 | 转换比例 | 示例 |
|------|----------|------|
| 斤   | 1.0      | 2斤 = 2斤 |
| 公斤 | 2.0      | 1公斤 = 2斤 |
| 千克 | 2.0      | 0.5千克 = 1斤 |
| 克   | 0.002    | 500克 = 1斤 |
| 两   | 0.1      | 10两 = 1斤 |
| 磅   | 0.907    | 1磅 ≈ 0.907斤 |

## 错误处理

### 常见错误及解决方案

1. **网络连接失败**
   - 检查网络连接
   - 确认URL是否正确
   - 检查防火墙设置

2. **认证失败**
   - 更新Cookie值
   - 确认Cookie是否过期
   - 检查是否需要重新登录

3. **解析失败**
   - 运行测试模式验证功能
   - 检查网页结构是否变化
   - 调整HTML解析规则

4. **价格计算错误**
   - 检查规格格式是否正确
   - 确认重量单位是否支持
   - 验证价格数据是否完整

## 高级功能

### 1. 批量处理

可以修改程序支持批量URL处理：

```go
urls := []string{
    "https://site1.com/product1",
    "https://site1.com/product2",
    "https://site2.com/product1",
}

for _, url := range urls {
    product, err := fetchProductInfo(url, cookie)
    if err != nil {
        continue
    }
    // 处理商品信息
}
```

### 2. 数据导出

可以将结果导出为CSV或JSON格式：

```go
// 导出为JSON
jsonData, _ := json.Marshal(product)
os.WriteFile("product.json", jsonData, 0644)

// 导出为CSV
csvData := fmt.Sprintf("%s,%.2f,%s,%.2f\n", 
    product.Name, product.Price, product.Spec, product.PricePerJin)
os.WriteFile("products.csv", []byte(csvData), 0644)
```

### 3. 定时任务

可以设置定时运行：

```go
import "time"

func main() {
    ticker := time.NewTicker(1 * time.Hour)
    for range ticker.C {
        // 执行价格查询
        fetchProductInfo(url, cookie)
    }
}
```

## 注意事项

1. **合法使用**: 遵守网站的robots.txt和使用条款
2. **频率控制**: 避免过于频繁的请求
3. **数据准确性**: 定期验证解析结果的准确性
4. **Cookie管理**: 及时更新过期的Cookie
5. **错误监控**: 监控程序运行状态和错误日志

## 技术支持

如果遇到问题，可以：

1. 运行测试模式检查功能
2. 查看错误信息和日志
3. 检查网络连接和认证状态
4. 调整HTML解析规则
5. 参考示例配置文件
