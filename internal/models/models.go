package models

import (
	"time"
)

// User 对应于 'users' 表
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"type:varchar(255);not null;unique"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Role         string `gorm:"type:varchar(50);not null"` // 'admin' 或 'coach'
	CreatedAt    time.Time
	// Coach        Coach     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"` // 移除递归引用
}

// Coach 对应于 'coaches' 表
type Coach struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null;unique"`
	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	AvatarURL   string `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	Bookings    []Booking    `gorm:"foreignKey:CoachID;constraint:OnDelete:CASCADE;"` // 一对多关系
	User        User         `gorm:"foreignKey:UserID"`                               // 新增字段
	Videos      []CoachVideo `gorm:"foreignKey:CoachID;constraint:OnDelete:CASCADE;"`
}

// Venue 对应于 'venues' 表
type Venue struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(255);not null;unique"`
	Address     string `gorm:"type:varchar(255)"`
	ManagerName string `gorm:"type:varchar(255)"`
	Contact     string `gorm:"type:varchar(100)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Bookings    []Booking `gorm:"foreignKey:VenueID;constraint:OnDelete:RESTRICT;"`
}

// CoachVideo 教练视频
type CoachVideo struct {
	ID        uint   `gorm:"primaryKey"`
	CoachID   uint   `gorm:"not null;index"`
	Title     string `gorm:"type:varchar(255);not null"`
	VideoURL  string `gorm:"type:varchar(255);not null"`
	CoverURL  string `gorm:"type:varchar(255)"`
	SortOrder int    `gorm:"default:0"`
	CreatedAt time.Time
}

// Booking 对应于 'bookings' 表
type Booking struct {
	ID          uint      `gorm:"primaryKey"`
	CoachID     uint      `gorm:"not null;index"`
	VenueID     uint      `gorm:"not null;index:uniq_venue_date_slot,unique"`
	StudentID   uint      `gorm:"index;comment:学员ID（微信用户ID）"`
	BookingDate time.Time `gorm:"type:date;not null;index:uniq_venue_date_slot,unique"`
	TimeSlot    string    `gorm:"type:varchar(50);not null;index:uniq_venue_date_slot,unique"`
	ClientInfo  string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	Venue       Venue `gorm:"foreignKey:VenueID"`
	Coach       Coach `gorm:"foreignKey:CoachID"`
}
