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

// LoadConfig 从指定的路径加载配置到全局实例
func LoadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read config file at path %s: %v", path, err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
	Cfg = &config
	log.Println("Configuration loaded successfully.")
} 