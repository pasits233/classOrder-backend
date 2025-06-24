package main

import (
	"classorder-backend/config"
	"classorder-backend/internal/database"
	"classorder-backend/internal/router"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库连接
	database.InitDB()

	// 设置路由
	r := router.SetupRouter()

	// 添加静态文件服务
	r.StaticFS("/", http.Dir(filepath.Join("frontend", "dist")))

	// 启动服务器
	port := ":9528"
	log.Printf("Server is running on port%s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
} 