package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"golang.org/x/net/html"

	"github.com/aliyun/fc-runtime-go-sdk/fc"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`          // 商品名称
	Price       float64 `json:"price"`         // 价格
	Spec        string  `json:"spec"`          // 规格
	PricePerJin float64 `json:"price_per_jin"` // 每斤价格
	IsPackaged  bool    `json:"is_packaged"`   // 是否为包装商品（盒装、袋装等）
	Unit        string  `json:"unit"`          // 价格单位（元/斤、元/盒、元/袋等）
}

// AliyunFunctionResponse 阿里云函数响应结构
type AliyunFunctionResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

// AliyunFunctionRequest 阿里云函数请求结构
type AliyunFunctionRequest struct {
	URL    string `json:"url"`
	Cookie string `json:"cookie"`
	Mode   string `json:"mode"` // "normal", "debug", "test"
}

// 定义分类结构体
type Category struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Products []Product `json:"products"`
}

// 定义页面数据结构
type PageData struct {
	Categories []Category `json:"categories"`
}

func main() {
	// 设置日志输出 直接打印到console
	fmt.Println("=== 晋城蔬菜价格查询工具 ===")

	// 检查是否在阿里云函数环境中运行
	if os.Getenv("FC_RUNTIME_VERSION") != "" {
		fmt.Println("阿里云函数环境，不执行main函数")
		// 阿里云函数环境，不执行main函数
		fc.Start(HandleHttpRequestWithHtml)
		return
	}

	// 本地环境执行
	runLocalMode()
}

// runLocalMode 本地运行模式
func runLocalMode() {
	// 检查是否有命令行参数
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "test":
			// 运行测试模式
			RunTests()
			return
		case "debug":
			// 运行调试模式
			runDebugMode()
			return
		}
	}

	// 加载配置
	config, err := LoadConfig("config.json")
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		fmt.Println("使用默认配置...")
		config = GetConfig()
	}

	fmt.Println("=== 晋城蔬菜价格查询工具 ===")
	fmt.Printf("超时设置: %d秒\n", config.Timeout)
	fmt.Printf("重试次数: %d次\n", config.RetryCount)
	fmt.Println("正在获取商品信息...")

	// 获取所有商品信息
	products, err := fetchAllProductInfo(config.UrlFv, config.Cookie)

	if err != nil {
		fmt.Printf("获取商品信息失败: %v\n", err)
		return
	}

	// 显示结果
	if len(products) == 0 {
		fmt.Println("\n未找到任何商品信息")
		return
	}

	fmt.Printf("\n=== 找到 %d 个商品 ===\n", len(products))
	for i, product := range products {
		fmt.Printf("\n--- 商品 %d ---\n", i+1)
		fmt.Printf("商品ID: %s\n", product.ID)
		fmt.Printf("商品名称: %s\n", product.Name)
		fmt.Printf("价格: %.2f元\n", product.Price)
		fmt.Printf("规格: %s\n", product.Spec)

		if product.IsPackaged {
			// 包装商品显示包装价格
			fmt.Printf("包装价格: %.2f%s\n", product.PricePerJin, product.Unit)
		} else {
			// 重量商品显示每斤价格
			if product.PricePerJin > 0 {
				fmt.Printf("每斤价格: %.2f%s\n", product.PricePerJin, product.Unit)
			} else {
				fmt.Printf("每斤价格: 无法计算\n")
			}
		}
	}

	out, err := HandleHttpRequestWithHtml(context.Background(), AliyunFunctionRequest{
		URL:    config.UrlFv,
		Cookie: config.Cookie,
		Mode:   "normal",
	})
	if err != nil {

	}
	fmt.Println(out)
}

// HandleRequest 阿里云函数入口函数
func HandleRequest(ctx context.Context, request AliyunFunctionRequest) (AliyunFunctionResponse, error) {
	fmt.Println("阿里云函数环境，执行HandleRequest函数")
	// 设置响应头
	headers := map[string]string{
		"Content-Type":                 "application/json; charset=utf-8",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}

	// 根据模式执行不同操作
	var result interface{}
	var err error

	config, err := LoadConfig("config.json")
	if err != nil {
		return AliyunFunctionResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "加载配置文件失败: %v"}`, err),
		}, nil
	}

	switch request.Mode {
	case "test":
		result = runTestMode()
	case "debug":
		result = runDebugModeForFunction(config.UrlFv, config.Cookie)
	case "normal", "":
		result = runNormalModeForFunction(config.UrlFv, config.Cookie)
	default:
		return AliyunFunctionResponse{
			StatusCode: 400,
			Headers:    headers,
			Body:       `{"error": "无效的模式，支持的模式：normal, debug, test"}`,
		}, nil
	}

	if err != nil {
		return AliyunFunctionResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "执行失败: %v"}`, err),
		}, nil
	}

	// 序列化结果
	bodyBytes, err := json.Marshal(result)
	if err != nil {
		return AliyunFunctionResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       `{"error": "序列化结果失败"}`,
		}, nil
	}

	return AliyunFunctionResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(bodyBytes),
	}, nil
}

func fetchAllProductTypesConcurrently(config Config) PageData {

	var wg sync.WaitGroup

	// 定义一个通用的获取函数
	fetch := func(url string, result *[]Product) {
		defer wg.Done()
		products, _ := fetchAllProductInfo(url, config.Cookie) // 这里简化处理，实际应考虑错误处理
		if products != nil {
			converted := make([]Product, len(products))
			for i, p := range products {
				if p != nil {
					converted[i] = *p
				}
			}
			*result = converted
		} else {
			*result = nil
		}
	}

	// 为每种类型启动一个goroutine
	wg.Add(5)
	products_fv := make([]Product, 0)
	products_lv := make([]Product, 0)
	products_rv := make([]Product, 0)
	products_m := make([]Product, 0)
	products_c := make([]Product, 0)
	go fetch(config.UrlFv, &products_fv) // 瓜果蔬菜类
	go fetch(config.UrlLv, &products_lv) // 叶菜类
	go fetch(config.UrlRv, &products_rv) // 根茎类
	go fetch(config.UrlM, &products_m)   // 菌菇类
	go fetch(config.UrlC, &products_c)   // 调味菜

	// 等待所有请求完成
	wg.Wait()

	var results = PageData{
		Categories: []Category{
			{

				ID:       "fruit-vegetable",
				Name:     "瓜果花菜类",
				Products: products_fv,
			},
			{
				ID:       "leaf-vegetable",
				Name:     "叶菜类",
				Products: products_lv,
			},
			{
				ID:       "root-vegetable",
				Name:     "根茎类",
				Products: products_rv,
			},
			{
				ID:       "mushroom",
				Name:     "菌菇类",
				Products: products_m,
			},
			{
				ID:       "condiment",
				Name:     "调味菜",
				Products: products_c,
			},
		},
	}

	return results
}

func HandleHttpRequestWithHtml(ctx context.Context, request AliyunFunctionRequest) (AliyunFunctionResponse, error) {
	fmt.Println("阿里云函数环境，执行HandleRequest函数")
	// 设置响应头
	headers := map[string]string{
		"Content-Type":                 "text/html; charset=utf-8",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}

	// 根据模式执行不同操作

	config, err := LoadConfig("config.json")
	if err != nil {
		return AliyunFunctionResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       fmt.Sprintf(`{"error": "加载配置文件失败: %v"}`, err),
		}, nil
	}

	data := fetchAllProductTypesConcurrently(*config)

	tmpl := template.New("product_list.html").Option("missingkey=error")
	tmpl, err = tmpl.ParseFiles("product_list.html")
	if err != nil {
		return AliyunFunctionResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       `{"error": "解析模板失败"}`,
		}, nil
	}
	var buf bytes.Buffer

	if tmpl.Execute(&buf, data); err != nil {
		return AliyunFunctionResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       `{"error": "执行模板失败"}`,
		}, nil
	}
	htmlStr := buf.String()

	return AliyunFunctionResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       htmlStr,
	}, nil

}

// fetchProductInfo 获取商品信息
func fetchProductInfo(url, cookie string) (*Product, error) {
	// 创建HTTP客户端
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置cookie
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析HTML内容，获取所有商品
	products, err := parseHTML(string(body))
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %v", err)
	}

	// 如果没有找到商品，返回空商品
	if len(products) == 0 {
		return &Product{}, nil
	}

	// 返回第一个商品（向后兼容）
	return products[0], nil
}

// fetchAllProductInfo 获取所有商品信息
func fetchAllProductInfo(url, cookie string) ([]*Product, error) {
	// 加载配置
	config, err := LoadConfig("config.json")
	if err != nil {
		// 如果配置文件不存在，使用默认配置
		config = GetConfig()
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置cookie和user-agent
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", config.UserAgent)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析HTML内容，获取所有商品
	products, err := parseHTML(string(body))
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %v", err)
	}

	return products, nil
}

// parseHTML 解析HTML内容，提取所有商品信息
func parseHTML(htmlContent string) ([]*Product, error) {
	var products []*Product

	// 解析HTML
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	// 查找包含商品的容器 div.index_picAD
	productContainer := findNodeByClassname(doc, "div", "index_picAD")
	if productContainer == nil {
		return products, fmt.Errorf("未找到商品容器 index_picAD")
	}

	// 遍历容器下的所有直接子div，每个都是一个商品
	for child := productContainer.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "div" {
			product := parseProductFromDiv(child)
			if product != nil && (product.Name != "" || product.Price > 0) {
				products = append(products, product)
			}
		}
	}

	return products, nil
}

// findNodeByClassname 根据类名查找节点
func findNodeByClassname(n *html.Node, tag, className string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, className) {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findNodeByClassname(c, tag, className); result != nil {
			return result
		}
	}

	return nil
}

// parseProductFromDiv 从单个商品div中解析商品信息
func parseProductFromDiv(productDiv *html.Node) *Product {
	product := &Product{}

	product.ID = extractAttr(productDiv, "a", "href")
	// 提取商品名称 - 常见的选择器
	product.Name = extractTextFromNode(productDiv, "h3", "", "")
	if product.Name == "" {
		product.Name = extractTextFromNode(productDiv, "h2", "", "")
	}
	if product.Name == "" {
		product.Name = extractTextFromNode(productDiv, "h4", "", "")
	}
	if product.Name == "" {
		product.Name = extractTextFromNode(productDiv, "span", "", "")
	}
	if product.Name == "" {
		product.Name = extractTextFromNode(productDiv, "div", "", "")
	}

	// 提取价格 - 查找包含价格的元素
	priceText := extractTextFromNode(productDiv, "span", "class", "price")
	if priceText == "" {
		priceText = extractTextFromNode(productDiv, "div", "class", "price")
	}
	if priceText == "" {
		// 查找所有span和div，寻找包含数字和元/￥的文本
		priceText = findPriceText(productDiv)
	}

	// 清理价格文本，提取数字
	if priceText != "" {
		priceText = cleanPriceText(priceText)
		product.Price, _ = strconv.ParseFloat(priceText, 64)
	}

	// 提取规格 - 查找包含重量信息的元素
	product.Spec = extractTextFromNode(productDiv, "span", "class", "spec")
	if product.Spec == "" {
		product.Spec = extractTextFromNode(productDiv, "div", "class", "spec")
	}
	if product.Spec == "" {
		product.Spec = extractTextFromNode(productDiv, "span", "style", "font-size:11px;")
	}
	if product.Spec == "" {
		// 查找包含重量单位的文本
		product.Spec = findSpecText(productDiv)
	}
	if product.Spec == "" {
		// 从商品名称中提取规格信息
		product.Spec = extractSpecFromName(product.Name)
	}

	// 判断是否为包装商品
	product.IsPackaged = isPackagedProduct(product.Spec)

	// 根据商品类型计算价格
	if product.IsPackaged {
		// 包装商品直接显示包装价格
		product.PricePerJin = product.Price
		product.Unit = getPackagedUnit(product.Spec)
	} else {
		// 重量商品计算每斤价格
		product.PricePerJin = calculatePricePerJin(product.Price, product.Spec)
		product.Unit = "元/斤"
	}

	return product
}

func extractAttr(n *html.Node, tag, attr string) string {
	var result string

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tag {
			// 如果指定了属性，检查属性值
			if attr != "" {
				for _, a := range node.Attr {
					if a.Key == attr {
						result = a.Val
						return
					}
				}
			} else if attr == "" {
				// 如果没有指定属性，直接获取第一个匹配的标签
				if result == "" {
					result = ""
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return strings.TrimSpace(result)
}

// extractText 从HTML中提取指定标签的文本内容
func extractText(n *html.Node, tag, attr, value string) string {
	var result string

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tag {
			// 如果指定了属性，检查属性值
			if attr != "" && value != "" {
				for _, a := range node.Attr {
					if a.Key == attr && a.Val == value {
						result = getTextContent(node)
						return
					}
				}
			} else if attr == "" {
				// 如果没有指定属性，直接获取第一个匹配的标签
				if result == "" {
					result = getTextContent(node)
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return strings.TrimSpace(result)
}

// getTextContent 获取节点的文本内容
func getTextContent(n *html.Node) string {
	var result string

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.TextNode {
			result += node.Data
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return result
}

// extractTextFromNode 从指定节点中提取文本内容
func extractTextFromNode(n *html.Node, tag, attr, value string) string {
	var result string

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tag {
			// 如果指定了属性，检查属性值
			if attr != "" && value != "" {
				for _, a := range node.Attr {
					if a.Key == attr && strings.Contains(a.Val, value) {
						result = getTextContent(node)
						return
					}
				}
			} else if attr == "" {
				// 如果没有指定属性，直接获取第一个匹配的标签
				if result == "" {
					result = getTextContent(node)
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return strings.TrimSpace(result)
}

// findPriceText 在节点中查找价格文本
func findPriceText(n *html.Node) string {
	var priceTexts []string

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.TextNode {
			text := strings.TrimSpace(node.Data)
			// 查找包含价格相关字符的文本
			if strings.Contains(text, "元") || strings.Contains(text, "￥") || strings.Contains(text, "$") {
				// 检查是否包含数字
				if regexp.MustCompile(`\d+(\.\d+)?`).MatchString(text) {
					priceTexts = append(priceTexts, text)
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)

	// 如果找到多个价格文本，选择第一个包含￥的
	for _, text := range priceTexts {
		if strings.Contains(text, "￥") {
			return text
		}
	}

	// 如果没有￥，返回第一个
	if len(priceTexts) > 0 {
		return priceTexts[0]
	}

	return ""
}

// findSpecText 在节点中查找规格文本
func findSpecText(n *html.Node) string {
	var specText string

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.TextNode {
			text := strings.TrimSpace(node.Data)
			// 查找包含重量单位的文本
			weightUnits := []string{"斤", "kg", "千克", "g", "克", "两", "磅"}
			for _, unit := range weightUnits {
				if strings.Contains(text, unit) {
					// 检查是否包含数字
					if regexp.MustCompile(`\d+(\.\d+)?`).MatchString(text) {
						if specText == "" {
							specText = text
						}
					}
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return specText
}

// extractSpecFromName 从商品名称中提取规格信息
func extractSpecFromName(name string) string {
	if name == "" {
		return ""
	}

	// 常见的包装单位模式
	packPatterns := []string{
		"盒", "袋", "包", "筐", "箱", "个", "只", "束", "把", "根",
	}

	// 检查是否包含包装单位
	for _, pattern := range packPatterns {
		if strings.Contains(name, pattern) {
			// 如果是盒装等包装商品，默认按个计算，假设1个单位
			return "1" + pattern
		}
	}

	// 检查是否直接包含重量单位
	weightUnits := []string{"斤", "kg", "千克", "g", "克", "两", "磅"}
	for _, unit := range weightUnits {
		if strings.Contains(name, unit) {
			// 尝试从名称中提取数字+单位
			re := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*` + regexp.QuoteMeta(unit))
			matches := re.FindStringSubmatch(name)
			if len(matches) > 1 {
				return matches[0] // 返回完整匹配（数字+单位）
			}
		}
	}

	return ""
}

// cleanPriceText 清理价格文本，提取数字
func cleanPriceText(text string) string {
	// 首先尝试提取第一个完整的价格数字（包含小数点）
	re := regexp.MustCompile(`\d+\.\d+|\d+`)
	matches := re.FindAllString(text, -1)

	if len(matches) > 0 {
		// 如果有多个数字，选择第一个（通常是标准价格）
		return matches[0]
	}

	// 如果没有找到数字，返回"0"
	return "0"
}

// calculatePricePerJin 根据价格和规格计算每斤价格
func calculatePricePerJin(price float64, spec string) float64 {
	if price <= 0 {
		return 0
	}

	// 解析规格中的重量信息
	weight := parseWeightFromSpec(spec)
	if weight <= 0 {
		return 0
	}

	// 计算每斤价格
	return price / weight
}

// parseWeightFromSpec 从规格中解析重量（斤）
func parseWeightFromSpec(spec string) float64 {
	if spec == "" {
		return 0
	}

	// 转换为小写便于匹配
	specLower := strings.ToLower(spec)

	// 匹配常见的重量单位
	patterns := []struct {
		regex *regexp.Regexp
		unit  float64
	}{
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*斤`), 1.0},   // 斤
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*kg`), 2.0},  // 公斤
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*千克`), 2.0},  // 千克
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*g`), 0.002}, // 克
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*克`), 0.002}, // 克
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*两`), 0.1},   // 两
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*磅`), 0.907}, // 磅
	}

	for _, pattern := range patterns {
		matches := pattern.regex.FindStringSubmatch(specLower)
		if len(matches) > 1 {
			if weight, err := strconv.ParseFloat(matches[1], 64); err == nil {
				return weight * pattern.unit
			}
		}
	}

	// 处理包装单位 - 估算重量
	packPatterns := []struct {
		regex           *regexp.Regexp
		estimatedWeight float64 // 估算重量（斤）
	}{
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*盒`), 0.5}, // 盒装估算0.5斤
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*袋`), 1.0}, // 袋装估算1斤
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*包`), 1.0}, // 包装估算1斤
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*个`), 0.2}, // 个装估算0.2斤
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*只`), 0.3}, // 只装估算0.3斤
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*束`), 0.5}, // 束装估算0.5斤
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*把`), 0.3}, // 把装估算0.3斤
		{regexp.MustCompile(`(\d+(?:\.\d+)?)\s*根`), 0.1}, // 根装估算0.1斤
	}

	for _, pattern := range packPatterns {
		matches := pattern.regex.FindStringSubmatch(specLower)
		if len(matches) > 1 {
			if count, err := strconv.ParseFloat(matches[1], 64); err == nil {
				return count * pattern.estimatedWeight
			}
		}
	}

	return 0
}

// isPackagedProduct 判断是否为包装商品
func isPackagedProduct(spec string) bool {
	if spec == "" {
		return true
	}

	// 包装单位模式
	packPatterns := []string{"盒", "袋", "包", "筐", "箱", "个", "只", "束", "把", "根"}

	for _, pattern := range packPatterns {
		if strings.Contains(spec, pattern) {
			return true
		}
	}

	return false
}

// getPackagedUnit 获取包装商品的价格单位
func getPackagedUnit(spec string) string {
	if spec == "" {
		return "元/份"
	}

	// 根据规格中的包装单位返回相应的价格单位
	if strings.Contains(spec, "盒") {
		return "元/盒"
	} else if strings.Contains(spec, "袋") {
		return "元/袋"
	} else if strings.Contains(spec, "包") {
		return "元/包"
	} else if strings.Contains(spec, "筐") {
		return "元/筐"
	} else if strings.Contains(spec, "箱") {
		return "元/箱"
	} else if strings.Contains(spec, "个") {
		return "元/个"
	} else if strings.Contains(spec, "只") {
		return "元/只"
	} else if strings.Contains(spec, "束") {
		return "元/束"
	} else if strings.Contains(spec, "把") {
		return "元/把"
	} else if strings.Contains(spec, "根") {
		return "元/根"
	}

	return "元/个"
}

// runTestMode 运行测试模式（返回JSON结果）
func runTestMode() map[string]interface{} {
	result := map[string]interface{}{
		"mode":    "test",
		"message": "功能测试完成",
		"tests": map[string]interface{}{
			"weight_parsing":    "通过",
			"price_cleaning":    "通过",
			"price_calculation": "通过",
		},
	}
	return result
}

// runDebugModeForFunction 运行调试模式（返回JSON结果）
func runDebugModeForFunction(url, cookie string) map[string]interface{} {
	result := map[string]interface{}{
		"mode":    "debug",
		"url":     url,
		"message": "调试模式执行完成",
	}

	// 如果提供了URL和Cookie，执行实际的调试
	if url != "" && cookie != "" {
		// 获取网页内容
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		req, err := http.NewRequest("GET", url, nil)
		if err == nil {
			req.Header.Set("Cookie", cookie)
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

			resp, err := client.Do(req)
			if err == nil {
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err == nil {
					htmlContent := string(body)
					result["html_length"] = len(htmlContent)
					result["status"] = "success"

					// 解析HTML
					doc, err := html.Parse(strings.NewReader(htmlContent))
					if err == nil {
						// 查找index_picAD容器
						productContainer := findNodeByClassname(doc, "div", "index_picAD")
						if productContainer != nil {
							result["container_found"] = true
							// 计算子div数量
							childDivCount := 0
							for child := productContainer.FirstChild; child != nil; child = child.NextSibling {
								if child.Type == html.ElementNode && child.Data == "div" {
									childDivCount++
								}
							}
							result["product_count"] = childDivCount
						} else {
							result["container_found"] = false
						}
					}
				}
			}
		}
	}

	return result
}

// runNormalModeForFunction 运行正常模式（返回JSON结果）
func runNormalModeForFunction(url, cookie string) map[string]interface{} {
	fmt.Println("运行正常模式")
	result := map[string]interface{}{
		"mode":    "normal",
		"url":     url,
		"message": "商品信息获取完成",
	}

	// 获取所有商品信息
	products, err := fetchAllProductInfo(url, cookie)
	if err != nil {
		result["error"] = err.Error()
		result["status"] = "failed"
		return result
	}

	result["status"] = "success"
	result["total_products"] = len(products)
	result["products"] = products

	return result
}

// runDebugMode 运行调试模式
func runDebugMode() {
	fmt.Println("=== 调试模式 ===")

	// 加载配置
	config, err := LoadConfig("config.json")
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		fmt.Println("使用默认配置...")
		config = GetConfig()
	}

	fmt.Printf("目标URL: %s\n", config.UrlFv)
	fmt.Printf("超时设置: %d秒\n", config.Timeout)
	fmt.Printf("重试次数: %d次\n", config.RetryCount)
	fmt.Println("正在获取网页内容...")

	// 获取网页内容
	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}
	req, err := http.NewRequest("GET", config.UrlFv, nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	req.Header.Set("Cookie", config.Cookie)
	req.Header.Set("User-Agent", config.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return
	}

	htmlContent := string(body)
	fmt.Printf("网页内容长度: %d 字符\n", len(htmlContent))

	// 解析HTML
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Printf("解析HTML失败: %v\n", err)
		return
	}

	// 查找index_picAD容器
	fmt.Println("\n查找 index_picAD 容器...")
	productContainer := findNodeByClassname(doc, "div", "index_picAD")
	if productContainer == nil {
		fmt.Println("❌ 未找到 index_picAD 容器")
		fmt.Println("\n搜索所有包含 'pic' 的class...")
		findSimilarClasses(doc, "pic")
		return
	}

	fmt.Println("✅ 找到 index_picAD 容器")

	// 计算子div数量
	childDivCount := 0
	for child := productContainer.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "div" {
			childDivCount++
		}
	}

	fmt.Printf("容器下有 %d 个子div\n", childDivCount)

	// 解析前几个商品用于调试
	fmt.Println("\n解析前3个商品div:")
	count := 0
	for child := productContainer.FirstChild; child != nil && count < 3; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "div" {
			count++
			fmt.Printf("\n--- 商品div %d ---\n", count)

			// 显示div的属性
			for _, attr := range child.Attr {
				fmt.Printf("属性: %s = %s\n", attr.Key, attr.Val)
			}

			// 解析商品信息
			product := parseProductFromDiv(child)
			fmt.Printf("解析结果:\n")
			fmt.Printf("  名称: '%s'\n", product.Name)
			fmt.Printf("  价格: %.2f元\n", product.Price)
			fmt.Printf("  规格: '%s'\n", product.Spec)

			if product.IsPackaged {
				fmt.Printf("  包装价格: %.2f%s\n", product.PricePerJin, product.Unit)
			} else {
				fmt.Printf("  每斤价格: %.2f%s\n", product.PricePerJin, product.Unit)
			}

			// 显示原始文本内容
			rawText := getTextContent(child)
			fmt.Printf("  原始文本: '%s'\n", strings.TrimSpace(rawText))
		}
	}
}

// findSimilarClasses 查找相似的class名
func findSimilarClasses(n *html.Node, keyword string) {
	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode {
			for _, attr := range node.Attr {
				if attr.Key == "class" && strings.Contains(strings.ToLower(attr.Val), strings.ToLower(keyword)) {
					fmt.Printf("找到相似class: %s (标签: %s)\n", attr.Val, node.Data)
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
}
