package router

import (
	"classOrder-backend/internal/api/handlers"
	"classOrder-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有API路由
func SetupRouter() *gin.Engine {
	// 使用默认配置创建一个Gin引擎
	r := gin.Default()

	// 提供静态文件服务，用于访问上传的头像
	// 例如 /uploads/avatar.png
	r.Static("/uploads", "./uploads")

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