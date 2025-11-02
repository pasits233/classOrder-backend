package router

import (
	"classOrder-backend/internal/api/handlers"
	"classOrder-backend/middleware"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有API路由
func SetupRouter() *gin.Engine {
	// 使用默认配置创建一个Gin引擎
	r := gin.Default()

	// 提供静态文件服务，用于访问上传的头像
	// 例如 /uploads/avatar.png
	// 尝试多个可能的路径来找到 uploads 目录
	var uploadsDir string
	possiblePaths := []string{
		"/root/williamcai/classOrder-backend/uploads", // 服务器实际路径
		filepath.Join(os.Getenv("PWD"), "uploads"),    // 从环境变量获取
		filepath.Join(os.Getenv("HOME"), "classOrder-backend", "uploads"),
	}
	
	// 如果 Getwd() 成功，也添加该路径
	if wd, err := os.Getwd(); err == nil {
		possiblePaths = append(possiblePaths, filepath.Join(wd, "uploads"))
	}
	
	// 尝试找到存在的 uploads 目录
	found := false
	for _, path := range possiblePaths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			uploadsDir = path
			found = true
			gin.DefaultWriter.Write([]byte(fmt.Sprintf("[INFO] Using uploads directory: %s\n", uploadsDir)))
			break
		}
	}
	
	// 如果都没找到，使用相对路径并尝试创建
	if !found {
		uploadsDir = "./uploads"
		if wd, err := os.Getwd(); err == nil {
			uploadsPath := filepath.Join(wd, "uploads")
			if err := os.MkdirAll(uploadsPath, 0755); err == nil {
				uploadsDir = uploadsPath
				gin.DefaultWriter.Write([]byte(fmt.Sprintf("[INFO] Created uploads directory: %s\n", uploadsDir)))
			}
		}
	}
	
	r.Static("/uploads", uploadsDir)

	// 公开的登录路由
	r.POST("/api/login", handlers.LoginHandler)
	
	// 微信登录相关路由
	wechatHandler := handlers.NewWeChatHandler()
	r.POST("/api/wechat/login", wechatHandler.LoginHandler)
	r.POST("/api/wechat/bind-user", middleware.JWTAuthMiddleware(), wechatHandler.BindUserHandler)
	r.GET("/api/wechat/user-info", middleware.JWTAuthMiddleware(), wechatHandler.GetUserInfoHandler)

	// 微信预约相关路由
	wechatBookingHandler := handlers.WeChatCreateBookingHandler
	r.POST("/api/wechat/bookings", middleware.JWTAuthMiddleware(), wechatBookingHandler)
	r.GET("/api/wechat/bookings", middleware.JWTAuthMiddleware(), handlers.WeChatListBookingsHandler)
	r.DELETE("/api/wechat/bookings/:id", middleware.JWTAuthMiddleware(), handlers.WeChatDeleteBookingHandler)

	// API路由组
	api := r.Group("/api")
	{
		// 上传文件路由 (需要登录)
		// 任何登录用户都可以上传，但在教练创建/更新时由管理员使用
		api.POST("/upload", middleware.JWTAuthMiddleware(), handlers.UploadHandler)

		// 教练管理路由
		coaches := api.Group("/coaches")
		{
			coaches.GET("", handlers.ListCoachesHandler)      // 获取教练列表 (公开)
			coaches.GET("/:id", handlers.GetCoachHandler)     // 获取单个教练信息 (公开)
			
			// 以下操作需要管理员权限
			adminCoaches := coaches.Group("", middleware.JWTAuthMiddleware(), middleware.AdminAuthMiddleware())
			{
				adminCoaches.POST("", handlers.CreateCoachHandler)
				adminCoaches.PUT("/:id", handlers.UpdateCoachHandler)
				adminCoaches.DELETE("/:id", handlers.DeleteCoachHandler)
				adminCoaches.POST("/:id/reset-password", handlers.ResetCoachPasswordHandler) // 重置密码
			}
		}

		// 预约管理路由
		bookings := api.Group("/bookings", middleware.JWTAuthMiddleware())
		{
			bookings.GET("", handlers.ListBookingsHandler)
			bookings.POST("", handlers.CreateBookingHandler)
			bookings.PUT(":id", handlers.UpdateBookingHandler)
			bookings.DELETE(":id", handlers.DeleteBookingHandler)
		}

		// 教练自助管理个人信息（仅需登录）
		api.GET("/coach/profile", middleware.JWTAuthMiddleware(), handlers.GetOwnCoachProfileHandler)
		api.PUT("/coach/profile", middleware.JWTAuthMiddleware(), handlers.UpdateOwnCoachProfileHandler)
		api.POST("/coach/reset-password", middleware.JWTAuthMiddleware(), handlers.ResetOwnCoachPasswordHandler) // 教练重置自己的密码
	}

	return r
} 