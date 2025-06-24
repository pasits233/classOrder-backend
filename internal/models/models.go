package models

import (
	"time"
)

// User 对应于 'users' 表
type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"type:varchar(255);not null;unique"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	Role         string    `gorm:"type:varchar(50);not null"` // 'admin' 或 'coach'
	CreatedAt    time.Time
	Coach        Coach     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // 一对一关系
}

// Coach 对应于 'coaches' 表
type Coach struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null;unique"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	AvatarURL   string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	Bookings    []Booking `gorm:"foreignKey:CoachID;constraint:OnDelete:CASCADE;"` // 一对多关系
	User        User      `gorm:"foreignKey:UserID"` // 新增字段
}

// Booking 对应于 'bookings' 表
type Booking struct {
	ID          uint      `gorm:"primaryKey"`
	CoachID     uint      `gorm:"not null"`
	BookingDate time.Time `gorm:"type:date;not null"`
	TimeSlot    string    `gorm:"type:varchar(50);not null"`
	ClientInfo  string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
} 