package main

// ExampleConfig 示例配置，展示如何配置不同的网站
func ExampleConfig() {
	// 示例1：凤展超市
	// url := "https://www.fengzhansy.com/wchyzyg/wap.shtml?method=ztmodel&ztid=gfl00%E7%93%9C%E6%9E%9C%E8%8A%B1%E8%8F%9C%E7%B1%BB"
	// cookie := "shdzarea=%E6%96%87%E5%8D%8E%E8%B7%AF; scsmdid=012; shdzmdname=%E5%87%A4%E5%B1%95%E8%B6%85%E5%B8%82%E6%96%87%E5%8D%8E%E8%B7%AF%E5%BA%97"

	// 示例2：其他超市网站
	// url := "https://other-supermarket.com/product/123"
	// cookie := "session=abc123; user=test"

	// 示例3：电商平台
	// url := "https://ecommerce.com/product/detail/456"
	// cookie := "auth_token=xyz789; user_id=12345"

	// 如何获取Cookie：
	// 1. 在浏览器中登录目标网站
	// 2. 按F12打开开发者工具
	// 3. 切换到Network标签页
	// 4. 刷新页面，找到请求头中的Cookie字段
	// 5. 复制Cookie值到程序中
}

// 常见网站配置示例
var CommonConfigs = map[string]struct {
	URL    string
	Cookie string
	Notes  string
}{
	"凤展超市": {
		URL:    "https://www.fengzhansy.com/wchyzyg/wap.shtml?method=ztmodel&ztid=gfl00%E7%93%9C%E6%9E%9C%E8%8A%B1%E8%8F%9C%E7%B1%BB",
		Cookie: "shdzarea=%E6%96%87%E5%8D%8E%E8%B7%AF; scsmdid=012; shdzmdname=%E5%87%A4%E5%B1%95%E8%B6%85%E5%B8%82%E6%96%87%E5%8D%8E%E8%B7%AF%E5%BA%97",
		Notes:  "瓜果花菜类商品页面",
	},
	"通用模板": {
		URL:    "https://example.com/product",
		Cookie: "session=abc123; user=test",
		Notes:  "请替换为实际网站",
	},
}
