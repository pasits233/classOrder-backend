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
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null;unique" json:"name"`
	Address     string    `gorm:"type:varchar(255)" json:"address"`
	ManagerName string    `gorm:"type:varchar(255)" json:"manager_name"`
	Contact     string    `gorm:"type:varchar(100)" json:"contact"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Bookings    []Booking `gorm:"foreignKey:VenueID;constraint:OnDelete:RESTRICT;" json:"bookings,omitempty"`
}

// CoachVideo 教练视频
type CoachVideo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CoachID   uint      `gorm:"not null;index" json:"coach_id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	VideoURL  string    `gorm:"type:varchar(255);not null" json:"video_url"`
	CoverURL  string    `gorm:"type:varchar(255)" json:"cover_url"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
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
