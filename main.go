package main

import (
	"classOrder-backend/config"
	"classOrder-backend/internal/database"
	"classOrder-backend/internal/router"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// 如果不是API请求，则尝试提供静态文件
		if !strings.HasPrefix(path, "/api") {
			// 构建前端文件的绝对路径
			distPath := filepath.Join("frontend", "dist")
			if !filepath.IsAbs(distPath) {
				// 如果是相对路径，转换为绝对路径
				wd, err := os.Getwd()
				if err == nil {
					distPath = filepath.Join(wd, "frontend", "dist")
				}
			}

			// 如果是根路径或者文件不存在，默认返回index.html
			requestedFile := filepath.Join(distPath, path)
			if path == "/" || !fileExists(requestedFile) {
				indexPath := filepath.Join(distPath, "index.html")
				if fileExists(indexPath) {
					c.File(indexPath)
					return
				}
				// 如果index.html也不存在，记录错误
				log.Printf("警告: index.html不存在于路径: %s", indexPath)
				c.String(http.StatusNotFound, "File not found")
				return
			}
			c.File(requestedFile)
		}
	})

	// 启动服务器
	port := ":9528"
	log.Printf("Server is running on port%s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
} 