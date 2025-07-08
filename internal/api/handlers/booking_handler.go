package handlers

import (
	"classOrder-backend/internal/database"
	"classOrder-backend/internal/models"
	"net/http"
	"time"
	"strings"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	// 并发锁+冲突检测
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		var existing []models.Booking
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("coach_id = ? AND booking_date = ?", req.CoachID, bookingDate).Find(&existing).Error
		if err != nil {
			return err
		}
		newRanges := parseTimeRanges(req.TimeSlots)
		for _, e := range existing {
			existRanges := parseTimeRanges(e.TimeSlot)
			for _, nr := range newRanges {
				for _, er := range existRanges {
					if timeRangeOverlap(nr, er) {
						c.JSON(http.StatusConflict, gin.H{"error": "所选时间段已被预约，请选择其他时间段"})
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
			db = db.Where("booking_date = ?", date)
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
		})
	}
	c.JSON(http.StatusOK, resp)
} 