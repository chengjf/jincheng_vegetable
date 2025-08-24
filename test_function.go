package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// FunctionTestClient 函数测试客户端
type FunctionTestClient struct {
	FunctionURL string
	HTTPClient  *http.Client
}

// NewFunctionTestClient 创建新的测试客户端
func NewFunctionTestClient(functionURL string) *FunctionTestClient {
	return &FunctionTestClient{
		FunctionURL: functionURL,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// TestFunction 测试函数调用
func (c *FunctionTestClient) TestFunction(request AliyunFunctionRequest) error {
	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	// 发送请求
	resp, err := c.HTTPClient.Post(c.FunctionURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	var response AliyunFunctionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	// 显示结果
	fmt.Printf("状态码: %d\n", response.StatusCode)
	fmt.Printf("响应头: %+v\n", response.Headers)
	fmt.Printf("响应体: %s\n", response.Body)

	return nil
}

// TestAllModes 测试所有模式
func (c *FunctionTestClient) TestAllModes(url, cookie string) {
	fmt.Println("=== 测试所有模式 ===")

	// 测试正常模式
	fmt.Println("\n--- 测试正常模式 ---")
	err := c.TestFunction(AliyunFunctionRequest{
		URL:    url,
		Cookie: cookie,
		Mode:   "normal",
	})
	if err != nil {
		fmt.Printf("正常模式测试失败: %v\n", err)
	}

	// 测试调试模式
	fmt.Println("\n--- 测试调试模式 ---")
	err = c.TestFunction(AliyunFunctionRequest{
		URL:    url,
		Cookie: cookie,
		Mode:   "debug",
	})
	if err != nil {
		fmt.Printf("调试模式测试失败: %v\n", err)
	}

	// 测试测试模式
	fmt.Println("\n--- 测试测试模式 ---")
	err = c.TestFunction(AliyunFunctionRequest{
		URL:    url,
		Cookie: cookie,
		Mode:   "test",
	})
	if err != nil {
		fmt.Printf("测试模式测试失败: %v\n", err)
	}
}

// 如果直接运行此文件，执行测试
func testMain() {
	// 从配置文件读取URL和Cookie
	config, err := LoadConfig("config.json")
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		return
	}

	// 设置函数URL（需要替换为实际的函数URL）
	functionURL := "https://your-function-url.fc.aliyuncs.com/2016-08-15/proxy/vegetable-price-service/vegetable-price-function/"

	fmt.Println("=== 阿里云函数测试客户端 ===")
	fmt.Printf("函数URL: %s\n", functionURL)
	fmt.Printf("目标URL: %s\n", config.UrlFv)
	fmt.Println("")

	// 创建测试客户端
	client := NewFunctionTestClient(functionURL)

	// 测试所有模式
	client.TestAllModes(config.UrlFv, config.Cookie)
}
