package database

import (
	"classorder-backend/config"
	"classorder-backend/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接并执行自动迁移
func InitDB() {
	var err error
	cfg := config.Cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("Database connection successful.")

	// 自动迁移模式, 这将创建表、缺失的外键、约束、列和索引。
	// 注意：它不会删除未使用的列，以保护您的数据。
	err = DB.AutoMigrate(&models.User{}, &models.Coach{}, &models.Booking{})
	if err != nil {
		log.Fatalf("failed to auto migrate tables: %v", err)
	}
	log.Println("Database migration check successful.")
} 