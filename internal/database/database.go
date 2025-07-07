package database

import (
	"classOrder-backend/config"
	"classOrder-backend/internal/models"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DB *gorm.DB

// MigrationRecord 用于记录迁移历史
type MigrationRecord struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null;unique"`
	AppliedAt time.Time `gorm:"not null"`
}

// ExecuteMigrations 执行自定义SQL迁移
func ExecuteMigrations(db *gorm.DB) error {
	// 创建迁移记录表
	if err := db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("创建迁移记录表失败: %v", err)
	}

	// 检查是否已经执行过此次迁移
	var record MigrationRecord
	if err := db.Where("name = ?", "user_id_type_migration").First(&record).Error; err == nil {
		log.Println("迁移已经执行过，跳过...")
		return nil
	}

	// 使用相对于可执行文件的路径
	migrationPath := filepath.Join("backend", "internal", "database", "migrations.sql")
	
	// 尝试读取迁移文件
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		alternatePath := filepath.Join("internal", "database", "migrations.sql")
		migrationSQL, err = os.ReadFile(alternatePath)
		if err != nil {
			return fmt.Errorf("无法读取迁移文件，尝试的路径: %s 和 %s, 错误: %v", 
				migrationPath, alternatePath, err)
		}
	}

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开启事务失败: %v", tx.Error)
	}

	// 执行迁移SQL
	if err := tx.Exec(string(migrationSQL)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("执行迁移失败: %v", err)
	}

	// 记录迁移完成
	if err := tx.Create(&MigrationRecord{
		Name:      "user_id_type_migration",
		AppliedAt: time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("记录迁移历史失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	log.Println("成功执行并记录迁移")
	return nil
}

// InitDB 初始化数据库连接并执行自动迁移
func InitDB() {
	var err error
	cfg := config.Cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s&multiStatements=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用GORM的外键约束处理
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	log.Println("数据库连接成功。")

	// 首先执行自定义迁移
	if err := ExecuteMigrations(DB); err != nil {
		log.Printf("警告: 自定义迁移失败: %v", err)
		return
	}

	log.Println("自定义迁移成功完成。")

	// 自动迁移模型（不处理外键约束）
	if err := DB.Migrator().AutoMigrate(
		&models.User{},
		&models.Coach{},
		&models.Booking{},
	); err != nil {
		log.Printf("警告: 自动迁移表失败: %v", err)
		return
	}

	log.Println("数据库迁移检查成功。")
} 