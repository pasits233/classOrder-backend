package handlers

import (
	"classOrder-backend/internal/database"
	"classOrder-backend/internal/models"
	"classOrder-backend/internal/services"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// WeChatCreateBookingRequest 微信用户创建预约请求
type WeChatCreateBookingRequest struct {
	CoachID  uint    `json:"coach_id" binding:"required"`
	VenueID  uint    `json:"venue_id" binding:"required"`
	Date     string  `json:"date" binding:"required"`
	TimeSlot string  `json:"time_slot" binding:"required"`
	Duration float64 `json:"duration" binding:"required"`
}

// WeChatCreateBookingHandler 微信用户创建预约
func WeChatCreateBookingHandler(c *gin.Context) {
	var req WeChatCreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 从 JWT 中获取微信用户信息
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

	// 获取微信用户信息
	wechatService := services.NewWeChatService()
	wechatUser, err := wechatService.GetUserInfo(openIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 验证日期格式
	bookingDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误，请使用 YYYY-MM-DD 格式"})
		return
	}

	// 验证日期不能是过去的时间
	if bookingDate.Before(time.Now().Truncate(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能预约过去的日期"})
		return
	}

	log.Printf("[WeChatCreateBooking] user_id=%d, coach_id=%d, venue_id=%d, date=%s, time_slot=%s",
		wechatUser.ID, req.CoachID, req.VenueID, req.Date, req.TimeSlot)

	// 并发锁+冲突检测
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 检查教练和场地是否存在
		var coach models.Coach
		if err := tx.First(&coach, req.CoachID).Error; err != nil {
			return errors.New("教练不存在")
		}
		var venue models.Venue
		if err := tx.First(&venue, req.VenueID).Error; err != nil {
			return errors.New("场地不存在")
		}

		// 创建预约记录（昵称缺失时做回退）；并依赖唯一索引防止并发下重复
		displayName := wechatUser.NickName
		if displayName == "" {
			displayName = fmt.Sprintf("微信用户#%d", wechatUser.ID)
		}
		booking := models.Booking{
			CoachID:     req.CoachID,
			VenueID:     req.VenueID,
			BookingDate: bookingDate,
			TimeSlot:    req.TimeSlot,
			ClientInfo:  displayName,
			StudentID:   wechatUser.ID,
		}

		if err := tx.Create(&booking).Error; err != nil {
			// 处理并发唯一约束冲突（Duplicate entry）
			if strings.Contains(err.Error(), "Duplicate entry") {
				return errors.New("该时间段已被预约，请选择其他时间段")
			}
			return err
		}

		return nil
	})

	if err != nil {
		if err.Error() == "该时间段已被预约，请选择其他时间段" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "教练不存在" || err.Error() == "场地不存在" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建预约失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "预约创建成功",
		"data": gin.H{
			"coach_id":   req.CoachID,
			"venue_id":   req.VenueID,
			"date":       req.Date,
			"time_slot":  req.TimeSlot,
			"duration":   req.Duration,
			"student_id": wechatUser.ID,
		},
	})
}

// WeChatListBookingsHandler 微信用户查询预约列表
func WeChatListBookingsHandler(c *gin.Context) {
	// 从 JWT 中获取微信用户信息
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

	// 获取微信用户信息
	wechatService := services.NewWeChatService()
	wechatUser, err := wechatService.GetUserInfo(openIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 查询该用户的所有预约
	var bookings []models.Booking
	db := database.DB.Preload("Venue").Where("student_id = ?", wechatUser.ID)

	// 支持按教练ID和日期筛选
	coachID := c.Query("coach_id")
	dateStr := c.Query("date")

	if coachID != "" {
		db = db.Where("coach_id = ?", coachID)
	}
	if dateStr != "" {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			db = db.Where("DATE(booking_date) = ?", date.Format("2006-01-02"))
		}
	}

	if err := db.Order("booking_date DESC, created_at DESC").Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取预约列表失败"})
		return
	}

	// 返回前端需要的字段
	var resp []gin.H
	for _, b := range bookings {
		venueName := ""
		if b.Venue.ID != 0 {
			venueName = b.Venue.Name
		}
		resp = append(resp, gin.H{
			"id":         b.ID,
			"coach_id":   b.CoachID,
			"venue_id":   b.VenueID,
			"venue_name": venueName,
			"date":       b.BookingDate.Format("2006-01-02"),
			"time_slot":  b.TimeSlot,
			"student_id": b.StudentID,
			"created_at": b.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// WeChatDeleteBookingHandler 微信用户删除预约
func WeChatDeleteBookingHandler(c *gin.Context) {
	bookingID := c.Param("id")

	// 从 JWT 中获取微信用户信息
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

	// 获取微信用户信息
	wechatService := services.NewWeChatService()
	wechatUser, err := wechatService.GetUserInfo(openIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 查找预约记录
	var booking models.Booking
	if err := database.DB.First(&booking, bookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "预约记录不存在"})
		return
	}

	// 验证是否为该用户的预约
	if booking.StudentID != wechatUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此预约"})
		return
	}

	// 删除预约
	if err := database.DB.Delete(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除预约失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "预约删除成功",
	})
}
