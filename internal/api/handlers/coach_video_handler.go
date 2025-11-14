package handlers

import (
	"classOrder-backend/internal/database"
	"classOrder-backend/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CoachVideoRequest struct {
	Title     string `json:"title" binding:"required"`
	VideoURL  string `json:"video_url" binding:"required"`
	CoverURL  string `json:"cover_url"`
	SortOrder int    `json:"sort_order"`
}

// ListCoachVideosHandler 获取指定教练的视频列表（公开）
func ListCoachVideosHandler(c *gin.Context) {
	coachID := c.Param("id")
	var videos []models.CoachVideo
	if err := database.DB.
		Where("coach_id = ?", coachID).
		Order("sort_order ASC, created_at DESC").
		Find(&videos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve coach videos"})
		return
	}
	c.JSON(http.StatusOK, videos)
}

// CreateCoachVideoHandler 新增视频
func CreateCoachVideoHandler(c *gin.Context) {
	coachIDStr := c.Param("id")
	coachID, err := strconv.Atoi(coachIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coach id"})
		return
	}

	var req CoachVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	video := models.CoachVideo{
		CoachID:   uint(coachID),
		Title:     req.Title,
		VideoURL:  req.VideoURL,
		CoverURL:  req.CoverURL,
		SortOrder: req.SortOrder,
	}

	if err := database.DB.Create(&video).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coach video"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Coach video created successfully"})
}

// UpdateCoachVideoHandler 更新视频
func UpdateCoachVideoHandler(c *gin.Context) {
	videoID := c.Param("id")
	var video models.CoachVideo
	if err := database.DB.First(&video, videoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	var req CoachVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	video.Title = req.Title
	video.VideoURL = req.VideoURL
	video.CoverURL = req.CoverURL
	video.SortOrder = req.SortOrder

	if err := database.DB.Save(&video).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update coach video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coach video updated successfully"})
}

// DeleteCoachVideoHandler 删除视频
func DeleteCoachVideoHandler(c *gin.Context) {
	videoID := c.Param("id")
	if err := database.DB.Delete(&models.CoachVideo{}, videoID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coach video"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Coach video deleted successfully"})
}
