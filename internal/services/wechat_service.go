package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"classOrder-backend/internal/config"
	"classOrder-backend/internal/database"
	"classOrder-backend/internal/models"
	"classOrder-backend/internal/utils"
)

// WeChatService 微信服务
type WeChatService struct{}

// NewWeChatService 创建微信服务实例
func NewWeChatService() *WeChatService {
	return &WeChatService{}
}

// WeChatCode2SessionResponse 微信登录凭证校验响应
type WeChatCode2SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// Login 微信登录
func (s *WeChatService) Login(code string) (*models.WeChatLoginResponse, error) {
	// 1. 通过 code 获取 openid 和 session_key
	sessionInfo, err := s.code2Session(code)
	if err != nil {
		return nil, fmt.Errorf("获取微信登录凭证失败: %v", err)
	}

	if sessionInfo.ErrCode != 0 {
		return nil, fmt.Errorf("微信登录失败: %s", sessionInfo.ErrMsg)
	}

	// 2. 查找或创建微信用户
	wechatUser, err := s.findOrCreateWeChatUser(sessionInfo)
	if err != nil {
		return nil, fmt.Errorf("处理微信用户信息失败: %v", err)
	}

	// 3. 生成 JWT token
	token, err := utils.GenerateJWT(wechatUser.ID, wechatUser.OpenID, "student")
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	// 4. 判断用户是否已绑定完整信息
	isBound := s.isUserBound(wechatUser)

	response := &models.WeChatLoginResponse{
		Token:   token,
		IsBound: isBound,
	}

	// 如果用户已绑定，返回用户信息
	if isBound {
		response.UserInfo = wechatUser
	}

	return response, nil
}

// code2Session 通过 code 获取 openid 和 session_key
func (s *WeChatService) code2Session(code string) (*WeChatCode2SessionResponse, error) {
	wechatConfig := config.GetWeChatConfig()
	
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		wechatConfig.AppID,
		wechatConfig.AppSecret,
		code,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sessionInfo WeChatCode2SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionInfo); err != nil {
		return nil, err
	}

	return &sessionInfo, nil
}

// findOrCreateWeChatUser 查找或创建微信用户
func (s *WeChatService) findOrCreateWeChatUser(sessionInfo *WeChatCode2SessionResponse) (*models.WeChatUser, error) {
	var wechatUser models.WeChatUser

	// 先尝试查找现有用户
	err := database.DB.Where("open_id = ?", sessionInfo.OpenID).First(&wechatUser).Error
	if err == nil {
		// 用户存在，更新登录时间
		wechatUser.UpdatedAt = time.Now()
		database.DB.Save(&wechatUser)
		return &wechatUser, nil
	}

	// 用户不存在，创建新用户
	wechatUser = models.WeChatUser{
		OpenID:  sessionInfo.OpenID,
		UnionID: sessionInfo.UnionID,
		Role:    "student",
		Status:  1,
	}

	if err := database.DB.Create(&wechatUser).Error; err != nil {
		return nil, err
	}

	return &wechatUser, nil
}

// isUserBound 判断用户是否已绑定完整信息
func (s *WeChatService) isUserBound(user *models.WeChatUser) bool {
	return user.NickName != "" && user.AvatarURL != ""
}

// BindUser 绑定微信用户信息
func (s *WeChatService) BindUser(openID string, bindData *models.WeChatBindUserRequest) (*models.WeChatBindUserResponse, error) {
	var wechatUser models.WeChatUser

	// 查找用户
	if err := database.DB.Where("open_id = ?", openID).First(&wechatUser).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 更新用户信息
	wechatUser.NickName = bindData.NickName
	wechatUser.AvatarURL = bindData.AvatarURL
	wechatUser.Gender = bindData.Gender
	wechatUser.Country = bindData.Country
	wechatUser.Province = bindData.Province
	wechatUser.City = bindData.City
	wechatUser.Language = bindData.Language
	wechatUser.UpdatedAt = time.Now()

	if err := database.DB.Save(&wechatUser).Error; err != nil {
		return nil, fmt.Errorf("保存用户信息失败: %v", err)
	}

	return &models.WeChatBindUserResponse{
		Success:  true,
		Message:  "用户信息绑定成功",
		UserInfo: &wechatUser,
	}, nil
}

// GetUserInfo 获取微信用户信息
func (s *WeChatService) GetUserInfo(openID string) (*models.WeChatUser, error) {
	var wechatUser models.WeChatUser

	if err := database.DB.Where("open_id = ?", openID).First(&wechatUser).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	return &wechatUser, nil
}

// GetUserByID 根据ID获取微信用户信息
func (s *WeChatService) GetUserByID(userID uint) (*models.WeChatUser, error) {
	var wechatUser models.WeChatUser

	if err := database.DB.First(&wechatUser, userID).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	return &wechatUser, nil
} 