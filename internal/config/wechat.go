package config

import (
	"os"
)

// WeChatConfig 微信配置
type WeChatConfig struct {
	AppID     string
	AppSecret string
}

// GetWeChatConfig 获取微信配置
func GetWeChatConfig() *WeChatConfig {
	return &WeChatConfig{
		AppID:     "wxad9e66bf74484a17",
		AppSecret: "fc79b7de2d170c8adf4746089dc245b4",
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
