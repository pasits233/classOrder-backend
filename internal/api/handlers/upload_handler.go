package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadHandler 处理文件上传请求
func UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// 生成一个唯一的文件名以避免冲突
	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension
	
	// 尝试多个可能的路径来找到 uploads 目录
	var uploadsDir string
	possiblePaths := []string{
		"/root/williamcai/classOrder-backend/uploads", // 服务器实际路径
		filepath.Join(os.Getenv("PWD"), "uploads"),    // 从环境变量获取
		filepath.Join(os.Getenv("HOME"), "classOrder-backend", "uploads"),
	}
	
	// 如果 Getwd() 成功，也添加该路径
	if wd, err := os.Getwd(); err == nil {
		possiblePaths = append(possiblePaths, filepath.Join(wd, "uploads"))
	}
	
	// 尝试找到存在的 uploads 目录
	found := false
	for _, path := range possiblePaths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			uploadsDir = path
			found = true
			break
		}
	}
	
	// 如果都没找到，尝试创建（使用第一个可能的路径）
	if !found {
		if wd, err := os.Getwd(); err == nil {
			uploadsDir = filepath.Join(wd, "uploads")
		} else {
			uploadsDir = "/root/williamcai/classOrder-backend/uploads"
		}
	}
	
	// 确保 uploads 目录存在
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create uploads directory: " + err.Error()})
		return
	}
	
	// 定义保存路径（使用绝对路径）
	dst := filepath.Join(uploadsDir, newFileName)

	// 保存文件到服务器
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + err.Error()})
		return
	}

	// 返回文件的可访问URL
	fileURL := fmt.Sprintf("/uploads/%s", newFileName)
	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"file_url": fileURL,
	})
} 