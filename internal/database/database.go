package database

import (
	"classOrder-backend/config"
	"classOrder-backend/internal/models"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	type migrationScript struct {
		Name string
		Path string
	}

	scripts := []migrationScript{
		{
			Name: "user_id_type_migration",
			Path: filepath.Join("backend", "internal", "database", "migrations.sql"),
		},
		{
			Name: "venues_videos_migration",
			Path: filepath.Join("backend", "internal", "database", "migrations_venues_videos.sql"),
		},
	}

	readScript := func(preferredPath string) ([]byte, error) {
		data, err := os.ReadFile(preferredPath)
		if err == nil {
			return data, nil
		}
		alternatePath := filepath.Join("internal", "database", filepath.Base(preferredPath))
		data, err = os.ReadFile(alternatePath)
		if err != nil {
			return nil, fmt.Errorf("无法读取迁移文件，尝试的路径: %s 和 %s, 错误: %v",
				preferredPath, alternatePath, err)
		}
		return data, nil
	}

	for _, script := range scripts {
		var record MigrationRecord
		err := db.Where("name = ?", script.Name).First(&record).Error
		if err == nil {
			continue
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("查询迁移记录失败: %v", err)
		}

		migrationSQL, err := readScript(script.Path)
		if err != nil {
			return err
		}

		tx := db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("开启事务失败: %v", tx.Error)
		}

		if err := tx.Exec(string(migrationSQL)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("执行迁移 %s 失败: %v", script.Name, err)
		}

		if err := tx.Create(&MigrationRecord{
			Name:      script.Name,
			AppliedAt: time.Now(),
		}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("记录迁移历史失败: %v", err)
		}

		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("提交迁移 %s 事务失败: %v", script.Name, err)
		}

		log.Printf("迁移 %s 执行完成", script.Name)
	}

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
		&models.Venue{},
		&models.Coach{},
		&models.CoachVideo{},
		&models.Booking{},
	); err != nil {
		log.Printf("警告: 自动迁移表失败: %v", err)
		return
	}

	log.Println("数据库迁移检查成功。")
}
