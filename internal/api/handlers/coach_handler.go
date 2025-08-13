package handlers

import (
	"classOrder-backend/internal/database"
	"classOrder-backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateCoachRequest 定义了创建教练的请求结构
type CreateCoachRequest struct {
	Username    string `json:"username" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
}

// CreateCoachHandler 创建一个新的教练及其关联的用户账户
func CreateCoachHandler(c *gin.Context) {
	var req CreateCoachRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// 使用默认密码 coach123
	defaultPassword := "coach123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 使用事务确保原子性
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 创建User
		newUser := models.User{
			Username:     req.Username,
			PasswordHash: string(hashedPassword),
			Role:         "coach",
		}
		if err := tx.Create(&newUser).Error; err != nil {
			return err // 返回错误以回滚事务
		}

		// 创建Coach
		newCoach := models.Coach{
			UserID:      newUser.ID,
			Name:        req.Name,
			Description: req.Description,
			AvatarURL:   req.AvatarURL,
		}
		if err := tx.Create(&newCoach).Error; err != nil {
			return err // 返回错误以回滚事务
		}

		// 事务成功，自动提交
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coach: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Coach created successfully with default password: " + defaultPassword})
}

// ListCoachesHandler 获取所有教练的列表
func ListCoachesHandler(c *gin.Context) {
	var coaches []models.Coach
	// Preload("User") 会同时加载关联的User信息
	if err := database.DB.Preload("User").Find(&coaches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve coaches"})
		return
	}

	// 为了安全，我们不应该返回密码哈希值。可以定义一个新的struct来返回安全的数据。
	type SafeCoachResponse struct {
		ID          uint   `json:"id"`
		UserID      uint   `json:"user_id"`
		Username    string `json:"username"`
		Name        string `json:"name"`
		Description string `json:"description"`
		AvatarURL   string `json:"avatar_url"`
	}

	var response []SafeCoachResponse
	for _, coach := range coaches {
		response = append(response, SafeCoachResponse{
			ID:          coach.ID,
			UserID:      coach.UserID,
			Username:    coach.User.Username,
			Name:        coach.Name,
			Description: coach.Description,
			AvatarURL:   coach.AvatarURL,
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetCoachHandler 获取单个教练的详细信息
func GetCoachHandler(c *gin.Context) {
	id := c.Param("id")
	var coach models.Coach
	if err := database.DB.Preload("User").First(&coach, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coach not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve coach"})
		}
		return
	}

	// 同样，返回安全的数据
	response := gin.H{
		"id":          coach.ID,
		"user_id":     coach.UserID,
		"username":    coach.User.Username,
		"name":        coach.Name,
		"description": coach.Description,
		"avatar_url":  coach.AvatarURL,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateCoachRequest 定义了更新教练的请求结构
// 移除密码相关字段，密码通过单独的重置接口处理
type UpdateCoachRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
}

// UpdateCoachHandler 更新教练信息
func UpdateCoachHandler(c *gin.Context) {
	id := c.Param("id")
	var coach models.Coach
	// 预加载 User
	if err := database.DB.Preload("User").First(&coach, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coach not found"})
		return
	}

	var req UpdateCoachRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 更新教练表
	coach.Name = req.Name
	coach.Description = req.Description
	coach.AvatarURL = req.AvatarURL
	if err := database.DB.Save(&coach).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update coach"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coach updated successfully"})
}

// ResetCoachPasswordHandler 重置教练密码为默认密码
func ResetCoachPasswordHandler(c *gin.Context) {
	id := c.Param("id")
	var coach models.Coach
	// 预加载 User
	if err := database.DB.Preload("User").First(&coach, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coach not found"})
		return
	}

	// 重置密码为 coach123
	defaultPassword := "coach123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 更新用户密码
	if err := database.DB.Model(&models.User{}).Where("id = ?", coach.UserID).Update("password_hash", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully to: " + defaultPassword})
}

// DeleteCoachHandler 删除一个教练及其关联的用户
func DeleteCoachHandler(c *gin.Context) {
	id := c.Param("id")
	var coach models.Coach
	if err := database.DB.First(&coach, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coach not found"})
		return
	}

	// 使用事务确保原子性
	dbErr := database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除Coach记录会导致User记录被级联删除（如果DB支持或GORM处理得当）
		// 但为了保险，我们显式删除User
		if err := tx.Delete(&models.User{}, coach.UserID).Error; err != nil {
			return err
		}
		// 删除coach会因为外键约束失败，所以先删除user
		// GORM的级联删除在Delete时可能不会按预期工作，显式删除更安全
		// 实际上，由于我们在模型中设置了级联删除，所以删除User就足够了。
		// 删除Coach记录是不必要的。

		return nil
	})

	// 由于外键约束 ON DELETE CASCADE，删除users表中的记录会自动删除coaches表中对应的记录
	// 所以我们只需要删除user即可
	if dbErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coach and associated user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coach deleted successfully"})
}

// 新增：教练自助获取个人信息
func GetOwnCoachProfileHandler(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}
	userID, ok := userIDVal.(float64) // JWT 默认 float64
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
		return
	}

	var coach models.Coach
	if err := database.DB.Preload("User").Where("user_id = ?", uint(userID)).First(&coach).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coach not found"})
		return
	}

	response := gin.H{
		"id":          coach.ID,
		"user_id":     coach.UserID,
		"username":    coach.User.Username,
		"name":        coach.Name,
		"description": coach.Description,
		"avatar_url":  coach.AvatarURL,
	}
	c.JSON(http.StatusOK, response)
}

// 新增：教练自助修改个人信息
func UpdateOwnCoachProfileHandler(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}
	userID, ok := userIDVal.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
		return
	}

	var coach models.Coach
	if err := database.DB.Preload("User").Where("user_id = ?", uint(userID)).First(&coach).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coach not found"})
		return
	}

	var req UpdateCoachRequest
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	_ = c.ShouldBindJSON(&body)

	if req.Name != "" {
		coach.Name = req.Name
	}
	if req.Description != "" {
		coach.Description = req.Description
	}
	if req.AvatarURL != "" {
		coach.AvatarURL = req.AvatarURL
	}
	if err := database.DB.Save(&coach).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update coach"})
		return
	}

	// 如果有新密码，校验原密码
	if req.Password != "" {
		if req.OldPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "原密码不能为空"})
			return
		}
		var user models.User
		if err := database.DB.First(&user, coach.UserID).Error; err == nil {
			if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)) != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "原密码错误"})
				return
			}
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			user.PasswordHash = string(hashedPassword)
			database.DB.Save(&user)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
