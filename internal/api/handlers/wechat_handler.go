package handlers

import (
	"net/http"

	"classOrder-backend/internal/models"
	"classOrder-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// WeChatHandler 微信相关处理器
type WeChatHandler struct {
	wechatService *services.WeChatService
}

// NewWeChatHandler 创建微信处理器实例
func NewWeChatHandler() *WeChatHandler {
	return &WeChatHandler{
		wechatService: services.NewWeChatService(),
	}
}

// LoginHandler 微信登录处理
func (h *WeChatHandler) LoginHandler(c *gin.Context) {
	var req models.WeChatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 调用微信服务登录
	response, err := h.wechatService.Login(req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// BindUserHandler 绑定微信用户信息
func (h *WeChatHandler) BindUserHandler(c *gin.Context) {
	var req models.WeChatBindUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 从 JWT 中获取 openID
	openID, exists := c.Get("open_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	openIDStr, ok := openID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息格式错误"})
		return
	}

	// 调用微信服务绑定用户信息
	response, err := h.wechatService.BindUser(openIDStr, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUserInfoHandler 获取微信用户信息
func (h *WeChatHandler) GetUserInfoHandler(c *gin.Context) {
	// 从 JWT 中获取 openID
	openID, exists := c.Get("open_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	openIDStr, ok := openID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息格式错误"})
		return
	}

	// 调用微信服务获取用户信息
	userInfo, err := h.wechatService.GetUserInfo(openIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userInfo,
	})
} 