package main

import (
	"classorder-backend/config"
	"classorder-backend/internal/api/routes"
	"classorder-backend/internal/database"
	"fmt"
	"log"
)

func main() {
	// 1. 加载配置
	config.LoadConfig("config/config.yaml")
	
	// 2. 初始化数据库连接
	database.InitDB()

	// 3. 设置路由
	r := routes.SetupRouter()

	// 4. 从配置启动服务器
	port := config.Cfg.Server.Port
	fmt.Printf("Server is running on http://127.0.0.1%s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
} 