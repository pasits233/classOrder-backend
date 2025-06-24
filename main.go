package main

import (
	"classOrder-backend/config"
	"classOrder-backend/internal/database"
	"classOrder-backend/internal/router"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库连接
	database.InitDB()

	// 设置路由
	r := router.SetupRouter()

	// 添加静态文件服务
	// 将静态文件服务移到最后，并使用一个特定的前缀
	r.NoRoute(func(c *gin.Context) {
		// 如果不是API请求，则尝试提供静态文件
		if c.Request.URL.Path[:4] != "/api" {
			c.FileFromFS(c.Request.URL.Path, http.Dir(filepath.Join("frontend", "dist")))
		}
	})

	// 启动服务器
	port := ":9528"
	log.Printf("Server is running on port%s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
} 