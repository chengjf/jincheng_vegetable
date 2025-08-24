package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config 配置结构体
type Config struct {
	UrlFv string `json:"url_fv"` // 瓜果花菜类
	UrlLv string `json:"url_lv"` // 叶菜类
	UrlRv string `json:"url_rv"` // 根茎类
	UrlM  string `json:"url_m"`  // 菌菇类
	UrlC  string `json:"url_c"`  // 调味菜

	Cookie     string `json:"cookie"`      // 认证cookie
	UserAgent  string `json:"user_agent"`  // 用户代理
	Timeout    int    `json:"timeout"`     // 超时时间（秒）
	RetryCount int    `json:"retry_count"` // 重试次数
}

// LoadConfig 从配置文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证必要的配置项
	if config.UrlFv == "" {
		return nil, fmt.Errorf("配置文件中缺少URL")
	}
	if config.Cookie == "" {
		return nil, fmt.Errorf("配置文件中缺少Cookie")
	}

	// 设置默认值
	if config.UserAgent == "" {
		config.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36"
	}
	if config.Timeout <= 0 {
		config.Timeout = 30
	}
	if config.RetryCount <= 0 {
		config.RetryCount = 3
	}

	return &config, nil
}

// GetConfig 获取配置信息（向后兼容）
func GetConfig() *Config {
	config, err := LoadConfig("config.json")
	if err != nil {
		// 如果配置文件不存在，返回默认配置
		return &Config{
			UrlFv:      "https://www.fengzhansy.com/wchyzyg/wap.shtml?method=ztmodel&ztid=gfl00%E7%93%9C%E6%9E%9C%E8%8A%B1%E8%8F%9C%E7%B1%BB",
			UrlLv:      "https://www.fengzhansy.com/wchyzyg/wap.shtml?method=ztmodel&ztid=gfl00%E5%8F%B6%E8%8F%9C%E7%B1%BB",
			UrlRv:      "https://www.fengzhansy.com/wchyzyg/wap.shtml?method=ztmodel&ztid=gfl00%E6%A0%B9%E8%8C%8E%E7%B1%BB",
			UrlM:       "https://www.fengzhansy.com/wchyzyg/wap.shtml?method=ztmodel&ztid=gfl00%E8%8F%8C%E8%8F%87%E7%B1%BB",
			UrlC:       "https://www.fengzhansy.com/wchyzyg/wap.shtml?method=ztmodel&ztid=gfl00%E8%B0%83%E5%91%B3%E8%8F%9C",
			Cookie:     "shdzarea=%E6%96%87%E5%8D%8E%E8%B7%AF; scsmdid=012; shdzmdname=%E5%87%A4%E5%B1%95%E8%B6%85%E5%B8%82%E6%96%87%E5%8D%8E%E8%B7%AF%E5%BA%97",
			UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
			Timeout:    30,
			RetryCount: 3,
		}
	}
	return config
}
