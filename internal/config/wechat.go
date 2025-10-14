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
		AppID:     getEnv("WECHAT_APP_ID", "your_app_id_here"),
		AppSecret: getEnv("WECHAT_APP_SECRET", "your_app_secret_here"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 