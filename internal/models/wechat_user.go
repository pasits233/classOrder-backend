package models

import (
	"time"
	"gorm.io/gorm"
)

// WeChatUser 微信用户模型
type WeChatUser struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	OpenID    string         `json:"open_id" gorm:"uniqueIndex;not null;comment:微信OpenID"`
	UnionID   string         `json:"union_id" gorm:"index;comment:微信UnionID"`
	NickName  string         `json:"nick_name" gorm:"comment:微信昵称"`
	AvatarURL string         `json:"avatar_url" gorm:"comment:微信头像URL"`
	Gender    int            `json:"gender" gorm:"comment:性别 0-未知 1-男 2-女"`
	Country   string         `json:"country" gorm:"comment:国家"`
	Province  string         `json:"province" gorm:"comment:省份"`
	City      string         `json:"city" gorm:"comment:城市"`
	Language  string         `json:"language" gorm:"comment:语言"`
	Role      string         `json:"role" gorm:"default:student;comment:用户角色 student-学员"`
	Status    int            `json:"status" gorm:"default:1;comment:状态 0-禁用 1-正常"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (WeChatUser) TableName() string {
	return "wechat_users"
}

// WeChatLoginRequest 微信登录请求
type WeChatLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// WeChatLoginResponse 微信登录响应
type WeChatLoginResponse struct {
	Token    string      `json:"token"`
	IsBound  bool        `json:"is_bound"`
	UserInfo *WeChatUser `json:"user_info,omitempty"`
}

// WeChatBindUserRequest 绑定微信用户请求
type WeChatBindUserRequest struct {
	NickName  string `json:"nick_name" binding:"required"`
	AvatarURL string `json:"avatar_url"`
	Gender    int    `json:"gender"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Language  string `json:"language"`
}

// WeChatBindUserResponse 绑定微信用户响应
type WeChatBindUserResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	UserInfo *WeChatUser `json:"user_info,omitempty"`
} 