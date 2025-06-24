package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 是整个项目的配置结构体
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `yaml:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	DBName    string `yaml:"dbname"`
	Charset   string `yaml:"charset"`
	ParseTime string `yaml:"parseTime"`
	Loc       string `yaml:"loc"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	Expiration int    `yaml:"expiration"`
}

// Cfg 是一个全局可访问的配置实例
var Cfg *Config

// InitConfig 初始化配置
func InitConfig() {
	// 尝试多个可能的配置文件路径
	configPaths := []string{
		"config/config.yaml",
		"backend/config/config.yaml",
		"../config/config.yaml",
	}

	var configPath string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}
	}

	if configPath == "" {
		log.Fatal("Could not find config file in any of the expected locations")
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// 解析配置
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	Cfg = &config
	log.Printf("Configuration loaded successfully from %s", configPath)
} 