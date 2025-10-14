package utils

import (
	"time"

	"classOrder-backend/config"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT 声明结构
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	OpenID string `json:"open_id,omitempty"` // 微信用户的 OpenID
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成 JWT token
// 支持传统用户和微信用户
func GenerateJWT(userID uint, openID string, role string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 如果是微信用户，添加 OpenID
	if openID != "" {
		claims.OpenID = openID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWT.Secret))
}

// GenerateTraditionalJWT 为传统用户生成 JWT token
func GenerateTraditionalJWT(userID uint, role string) (string, error) {
	return GenerateJWT(userID, "", role)
}

// GenerateWeChatJWT 为微信用户生成 JWT token
func GenerateWeChatJWT(userID uint, openID string, role string) (string, error) {
	return GenerateJWT(userID, openID, role)
}

// ParseJWT 解析 JWT token
func ParseJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
} 