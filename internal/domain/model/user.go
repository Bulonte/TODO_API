package model

import (
	"time"

	"gorm.io/gorm"
)

// User用户模型
type User struct {
	ID           uint           `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Username     string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	AvatarURL    *string        `gorm:"type:varchar(255)" json:"avatar_url,omitempty"`
	Status       uint8          `gorm:"type:tinyint;default:1" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}
