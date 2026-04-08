package model

import (
	"TODO_API/internal/domain/event"
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

	// 事件列表（不存储到数据库）
	Events []event.DomainEvent `gorm:"-" json:"-"`
}

func (User) TableName() string {
	return "users"
}

// AddEvent 添加事件
func (u *User) AddEvent(e event.DomainEvent) {
	u.Events = append(u.Events, e)
}

// GetEvents 获取事件列表
func (u *User) GetEvents() []event.DomainEvent {
	return u.Events
}

// ClearEvents 清空事件列表
func (u *User) ClearEvents() {
	u.Events = []event.DomainEvent{}
}

