package main

import (
	"fmt"
	"testing"
)

// TestParseWeightFromSpec 测试重量解析功能
func TestParseWeightFromSpec(t *testing.T) {
	testCases := []struct {
		spec     string
		expected float64
	}{
		{"500g", 1.0}, // 500克 = 1斤
		{"1kg", 2.0},  // 1公斤 = 2斤
		{"2斤", 2.0},   // 2斤 = 2斤
		{"500克", 1.0}, // 500克 = 1斤
		{"1千克", 2.0},  // 1千克 = 2斤
		{"10两", 1.0},  // 10两 = 1斤
		{"1磅", 0.907}, // 1磅 ≈ 0.907斤
		{"", 0},       // 空字符串
		{"无规格", 0},    // 无数字
	}

	for _, tc := range testCases {
		result := parseWeightFromSpec(tc.spec)
		if result != tc.expected {
			t.Errorf("parseWeightFromSpec(%q) = %f, 期望 %f", tc.spec, result, tc.expected)
		}
	}
}

// TestCleanPriceText 测试价格文本清理功能
func TestCleanPriceText(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"¥12.50", "12.50"},
		{"$15.99", "15.99"},
		{"价格：8.80元", "8.80"},
		{"特价 5.5", "5.5"},
		{"免费", "0"},
		{"", "0"},
	}

	for _, tc := range testCases {
		result := cleanPriceText(tc.input)
		if result != tc.expected {
			t.Errorf("cleanPriceText(%q) = %s, 期望 %s", tc.input, result, tc.expected)
		}
	}
}

// TestCalculatePricePerJin 测试每斤价格计算功能
func TestCalculatePricePerJin(t *testing.T) {
	testCases := []struct {
		price    float64
		spec     string
		expected float64
	}{
		{10.0, "1斤", 10.0},   // 10元/斤
		{20.0, "2斤", 10.0},   // 20元/2斤 = 10元/斤
		{15.0, "500g", 15.0}, // 15元/500g = 15元/斤
		{0, "1斤", 0},         // 价格为0
		{10.0, "", 0},        // 无规格
	}

	for _, tc := range testCases {
		result := calculatePricePerJin(tc.price, tc.spec)
		if result != tc.expected {
			t.Errorf("calculatePricePerJin(%.2f, %q) = %.2f, 期望 %.2f", tc.price, tc.spec, result, tc.expected)
		}
	}
}

// 运行测试的函数（非标准测试）
func RunTests() {
	fmt.Println("=== 运行功能测试 ===")

	// 测试重量解析
	fmt.Println("\n1. 测试重量解析:")
	testSpecs := []string{"500g", "1kg", "2斤", "500克", "1千克", "10两", "1磅"}
	for _, spec := range testSpecs {
		weight := parseWeightFromSpec(spec)
		fmt.Printf("   %s -> %.3f斤\n", spec, weight)
	}

	// 测试价格清理
	fmt.Println("\n2. 测试价格清理:")
	testPrices := []string{"¥12.50", "$15.99", "价格：8.80元", "特价 5.5"}
	for _, price := range testPrices {
		cleaned := cleanPriceText(price)
		fmt.Printf("   %s -> %s\n", price, cleaned)
	}

	// 测试每斤价格计算
	fmt.Println("\n3. 测试每斤价格计算:")
	testCases := []struct {
		price float64
		spec  string
	}{
		{10.0, "1斤"},
		{20.0, "2斤"},
		{15.0, "500g"},
		{25.0, "1kg"},
	}

	for _, tc := range testCases {
		pricePerJin := calculatePricePerJin(tc.price, tc.spec)
		fmt.Printf("   价格%.2f元，规格%s -> %.2f元/斤\n", tc.price, tc.spec, pricePerJin)
	}

	fmt.Println("\n=== 测试完成 ===")
}
