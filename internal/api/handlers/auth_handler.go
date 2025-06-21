package handlers

import (
	"classorder-backend/config"
	"classorder-backend/internal/database"
	"classorder-backend/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest 定义了登录请求的JSON结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 定义了成功登录后返回的JSON结构
type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

// LoginHandler 处理用户登录请求
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 1. 从数据库中查找用户
	var user models.User
	result := database.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 2. 比较密码哈希值
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 3. 生成JWT
	token, err := generateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		Role:  user.Role,
	})
}

// generateJWT 为指定用户生成JWT令牌
func generateJWT(user models.User) (string, error) {
	cfg := config.Cfg.JWT
	
	// 创建JWT的claims
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * time.Duration(cfg.Expiration)).Unix(),
		"iat":      time.Now().Unix(),
	}

	// 使用HS256签名算法创建一个新的token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 使用密钥签名并获取完整的编码后的字符串token
	return token.SignedString([]byte(cfg.Secret))
} 