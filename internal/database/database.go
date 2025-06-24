package database

import (
	"classorder-backend/config"
	"classorder-backend/internal/models"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ExecuteMigrations 执行自定义SQL迁移
func ExecuteMigrations(db *gorm.DB) error {
	// 使用相对于可执行文件的路径
	migrationPath := filepath.Join("backend", "internal", "database", "migrations.sql")
	
	// 尝试读取迁移文件
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		// 如果第一次尝试失败，尝试不同的路径
		alternatePath := filepath.Join("internal", "database", "migrations.sql")
		migrationSQL, err = os.ReadFile(alternatePath)
		if err != nil {
			return fmt.Errorf("无法读取迁移文件，尝试的路径: %s 和 %s, 错误: %v", 
				migrationPath, alternatePath, err)
		}
	}

	// 分别执行每条SQL语句
	err = db.Exec(string(migrationSQL)).Error
	if err != nil {
		return fmt.Errorf("执行迁移失败: %v", err)
	}

	return nil
}

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
		log.Fatalf("连接数据库失败: %v", err)
	}

	log.Println("数据库连接成功。")

	// 首先执行自定义迁移
	if err := ExecuteMigrations(DB); err != nil {
		log.Printf("警告: 自定义迁移失败: %v", err)
		// 即使自定义迁移失败，继续执行自动迁移
	}

	// 自动迁移模型
	err = DB.AutoMigrate(&models.User{}, &models.Coach{}, &models.Booking{})
	if err != nil {
		log.Fatalf("自动迁移表失败: %v", err)
	}
	log.Println("数据库迁移检查成功。")
} 