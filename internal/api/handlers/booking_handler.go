package handlers

import (
	"classOrder-backend/internal/database"
	"classOrder-backend/internal/models"
	"net/http"
	"time"
	"strings"
	"sort"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type CreateBookingRequest struct {
	StudentName string `json:"student_name" binding:"required"`
	CoachID     uint   `json:"coach_id" binding:"required"`
	Date        string `json:"date" binding:"required"` // YYYY-MM-DD
	TimeSlots   string `json:"time_slots" binding:"required"`
}

type UpdateBookingRequest struct {
	StudentName string `json:"student_name"`
	CoachID     uint   `json:"coach_id"`
	Date        string `json:"date"`
	TimeSlots   string `json:"time_slots"`
}

func parseTimeRanges(slots string) [][2]string {
	ranges := strings.Split(slots, ",")
	var result [][2]string
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		parts := strings.Split(r, "-")
		if len(parts) == 2 {
			result = append(result, [2]string{parts[0], parts[1]})
		}
	}
	return result
}

func timeRangeOverlap(a, b [2]string) bool {
	return !(a[1] <= b[0] || a[0] >= b[1])
}

// CreateBookingHandler 创建预约
func CreateBookingHandler(c *gin.Context) {
	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}
	bookingDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	log.Printf("[CreateBooking] coach_id=%d, date=%s, time_slots=%s", req.CoachID, bookingDate.Format("2006-01-02"), req.TimeSlots)

	// 并发锁+冲突检测
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		var existing []models.Booking
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("coach_id = ? AND booking_date = ?", req.CoachID, bookingDate).Find(&existing).Error
		if err != nil {
			return err
		}
		log.Printf("[CreateBooking] existing bookings: %+v", existing)
		newRanges := parseTimeRanges(req.TimeSlots)
		for _, e := range existing {
			existRanges := parseTimeRanges(e.TimeSlot)
			for _, nr := range newRanges {
				for _, er := range existRanges {
					if timeRangeOverlap(nr, er) {
						c.JSON(http.StatusConflict, gin.H{"error": "所选时间段已被预约，请选择其他时间段"})
						log.Printf("[CreateBooking] conflict: new=%v, exist=%v", nr, er)
						return errors.New("time slot conflict")
					}
				}
			}
		}
		booking := models.Booking{
			CoachID:     req.CoachID,
			BookingDate: bookingDate,
			TimeSlot:    req.TimeSlots,
			ClientInfo:  req.StudentName,
		}
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err.Error() == "time slot conflict" {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Booking created successfully"})
}

// UpdateBookingHandler 更新预约
func UpdateBookingHandler(c *gin.Context) {
	id := c.Param("id")
	var booking models.Booking
	if err := database.DB.First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	var req UpdateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	// 记录原始信息
	if req.StudentName != "" {
		booking.ClientInfo = req.StudentName
	}
	if req.CoachID != 0 {
		booking.CoachID = req.CoachID
	}
	if req.Date != "" {
		if date, err := time.Parse("2006-01-02", req.Date); err == nil {
			booking.BookingDate = date
		}
	}
	if req.TimeSlots != "" {
		booking.TimeSlot = req.TimeSlots
	}
	// 冲突校验：查找同教练同天除自己外的所有预约，判断时间段是否重叠
	var existing []models.Booking
	database.DB.Where("coach_id = ? AND booking_date = ? AND id != ?", booking.CoachID, booking.BookingDate, booking.ID).Find(&existing)
	newRanges := parseTimeRanges(booking.TimeSlot)
	for _, e := range existing {
		existRanges := parseTimeRanges(e.TimeSlot)
		for _, nr := range newRanges {
			for _, er := range existRanges {
				if timeRangeOverlap(nr, er) {
					c.JSON(http.StatusConflict, gin.H{"error": "所选时间段已被预约，请选择其他时间段"})
					return
				}
			}
		}
	}
	if err := database.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Booking updated successfully"})
}

// DeleteBookingHandler 删除预约
func DeleteBookingHandler(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Booking{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
}

// ListBookingsHandler 查询预约
func ListBookingsHandler(c *gin.Context) {
	var bookings []models.Booking
	coachID := c.Query("coach_id")
	dateStr := c.Query("date")
	db := database.DB
	if coachID != "" {
		db = db.Where("coach_id = ?", coachID)
	}
	if dateStr != "" {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			db = db.Where("DATE(booking_date) = ?", date.Format("2006-01-02"))
		}
	}
	if err := db.Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookings"})
		return
	}
	// 返回前端需要的字段
	var resp []gin.H
	for _, b := range bookings {
		resp = append(resp, gin.H{
			"id":           b.ID,
			"coach_id":     b.CoachID,
			"date":         b.BookingDate.Format("2006-01-02"),
			"time_slots":   b.TimeSlot,
			"student_name": b.ClientInfo,
			"student_id":   b.StudentID,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// BookingAvailabilityHandler 公开接口：查询某个教练在某个日期的已预约时间段（用于显示可用性）
func BookingAvailabilityHandler(c *gin.Context) {
	coachID := c.Query("coach_id")
	dateStr := c.Query("date")
	
	if coachID == "" || dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "coach_id and date are required"})
		return
	}
	
	var bookings []models.Booking
	db := database.DB.Where("coach_id = ?", coachID)
	
	if date, err := time.Parse("2006-01-02", dateStr); err == nil {
		db = db.Where("DATE(booking_date) = ?", date.Format("2006-01-02"))
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}
	
	if err := db.Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookings"})
		return
	}
	
	// 添加日志
	log.Printf("[BookingAvailability] coach_id=%s, date=%s, found %d bookings", coachID, dateStr, len(bookings))
	
	// 返回所有已预约的时间段，按 30 分钟粒度展开
	// 例如：09:00-11:00 会被展开为 [09:00-09:30, 09:30-10:00, 10:00-10:30, 10:30-11:00]
	expandToHalfHour := func(rangeStr string) []string {
		parts := strings.Split(rangeStr, "-")
		if len(parts) != 2 {
			return nil
		}
		startStr := strings.TrimSpace(parts[0])
		endStr := strings.TrimSpace(parts[1])
		// 解析为当天固定日期，便于做时间加减
		baseDate := "2006-01-02"
		// 使用固定日期进行组合，日期值不重要
		startTime, err1 := time.Parse(baseDate+" 15:04", "2006-01-02 "+startStr)
		endTime, err2 := time.Parse(baseDate+" 15:04", "2006-01-02 "+endStr)
		if err1 != nil || err2 != nil {
			return nil
		}
		// 容错：如果结束早于开始，直接返回
		if !endTime.After(startTime) {
			return nil
		}
		var expanded []string
		cursor := startTime
		for cursor.Before(endTime) {
			next := cursor.Add(30 * time.Minute)
			if next.After(endTime) {
				next = endTime
			}
			expanded = append(expanded, cursor.Format("15:04")+"-"+next.Format("15:04"))
			cursor = next
		}
		return expanded
	}

	// 去重集合，避免重复时间段
	occupiedSet := make(map[string]struct{})
	for _, b := range bookings {
		log.Printf("[BookingAvailability] booking ID=%d, time_slot=%s", b.ID, b.TimeSlot)
		ranges := strings.Split(b.TimeSlot, ",")
		for _, r := range ranges {
			r = strings.TrimSpace(r)
			if r == "" {
				continue
			}
			// 尝试展开为半小时粒度；如果解析失败，则按原样加入
			expanded := expandToHalfHour(r)
			if len(expanded) == 0 {
				occupiedSet[r] = struct{}{}
				continue
			}
			for _, s := range expanded {
				occupiedSet[s] = struct{}{}
			}
		}
	}

	// 转为有序数组（按时间排序）
	var occupiedSlots []string
	for s := range occupiedSet {
		occupiedSlots = append(occupiedSlots, s)
	}
	// 简单按字符串排序，因统一为 HH:MM-HH:MM，字典序等同时间顺序
	if len(occupiedSlots) > 1 {
		sort.Strings(occupiedSlots)
	}

	// 确保返回空数组而不是 nil
	if occupiedSlots == nil {
		occupiedSlots = []string{}
	}

	log.Printf("[BookingAvailability] occupied_slots=%v (len=%d)", occupiedSlots, len(occupiedSlots))
	
	c.JSON(http.StatusOK, gin.H{
		"coach_id": coachID,
		"date":     dateStr,
		"occupied_slots": occupiedSlots,
	})
} 