package handlers

import (
	"classOrder-backend/internal/database"
	"classOrder-backend/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VenueRequest struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address"`
	ManagerName string `json:"manager_name"`
	Contact     string `json:"contact"`
}

// ListVenuesHandler 返回全部场地，供前台/小程序使用
func ListVenuesHandler(c *gin.Context) {
	var venues []models.Venue
	if err := database.DB.Order("name ASC").Find(&venues).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve venues"})
		return
	}
	c.JSON(http.StatusOK, venues)
}

// CreateVenueHandler 新增场地
func CreateVenueHandler(c *gin.Context) {
	var req VenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	venue := models.Venue{
		Name:        req.Name,
		Address:     req.Address,
		ManagerName: req.ManagerName,
		Contact:     req.Contact,
	}

	if err := database.DB.Create(&venue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create venue: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Venue created successfully"})
}

// UpdateVenueHandler 更新场地
func UpdateVenueHandler(c *gin.Context) {
	id := c.Param("id")
	var venue models.Venue
	if err := database.DB.First(&venue, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Venue not found"})
		return
	}

	var req VenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	venue.Name = req.Name
	venue.Address = req.Address
	venue.ManagerName = req.ManagerName
	venue.Contact = req.Contact

	if err := database.DB.Save(&venue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update venue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Venue updated successfully"})
}

// DeleteVenueHandler 删除场地（无预约记录时）
func DeleteVenueHandler(c *gin.Context) {
	id := c.Param("id")
	var venue models.Venue
	if err := database.DB.First(&venue, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Venue not found"})
		return
	}

	var count int64
	if err := database.DB.Model(&models.Booking{}).Where("venue_id = ?", venue.ID).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check bookings for venue"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "当前场地仍有关联预约，无法删除"})
		return
	}

	if err := database.DB.Delete(&venue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete venue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Venue deleted successfully"})
}

// GetVenueHandler 获取单个场地详情（可用于前端回填）
func GetVenueHandler(c *gin.Context) {
	id := c.Param("id")
	var venue models.Venue
	if err := database.DB.First(&venue, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Venue not found"})
		return
	}
	c.JSON(http.StatusOK, venue)
}

// ListVenueBookingsHandler 可选：查询某场地的所有预约
func ListVenueBookingsHandler(c *gin.Context) {
	id := c.Param("id")
	venueID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid venue id"})
		return
	}

	var bookings []models.Booking
	if err := database.DB.
		Preload("Coach").
		Where("venue_id = ?", venueID).
		Order("booking_date ASC, time_slot ASC").
		Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}
