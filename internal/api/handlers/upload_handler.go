package handlers

import (
	"fmt"
	"net/http"
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
	
	// 定义保存路径
	// 注意: 这里的 'uploads' 目录需要存在于项目的根目录下
	dst := filepath.Join("uploads", newFileName)

	// 保存文件到服务器
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 返回文件的可访问URL
	// 这里的URL是相对于服务器根目录的
	fileURL := fmt.Sprintf("/%s", dst)
	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"file_url": fileURL,
	})
} 