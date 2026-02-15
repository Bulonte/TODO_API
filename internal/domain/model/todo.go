package model

import (
	"time"

	"gorm.io/gorm"
)

// 待办事项状态
type TodoStatus uint8

const (
	todosPending    TodoStatus = 0 //待办
	todosInProgress TodoStatus = 1 //进行中
	todosCompleted  TodoStatus = 2 //已完成
)

// 待办事项优先级
type TodosPriority uint8

const (
	todosLowPriority    TodosPriority = 1 //低
	todosMedPriority    TodosPriority = 2 //中
	todosHighPriority   TodosPriority = 3 //高
	todosUrgentPriority TodosPriority = 4 //紧急
)

type Todo struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Title       string         `gorm:"type:varchar(200);not null" json:"title"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	Status      TodoStatus     `gorm:"type:tinyint;default:0" json:"status"`   // 0-待办,1-进行中,2-已完成
	Priority    TodosPriority  `gorm:"type:tinyint;default:1" json:"priority"` // 1-低,2-中,3-高,4-紧急
	DueDate     *time.Time     `gorm:"index" json:"due_date,omitempty"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	CreatedAt   *time.Time     `json:"created_at"`
	UpdatedAt   *time.Time     `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联用户
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Todo) TableName() string {
	return "todos"
}

// IsCompleted 检查是否已完成
func IsCompleted(t *Todo) bool {
	return t.Status == todosCompleted
}

// MarkCompleted 标记为已完成
func MarkCompleted(t *Todo) {
	t.Status = todosCompleted
	now := time.Now()
	t.CompletedAt = &now
}

// MarkInProgress 标记为进行中
func MarkProgress(t *Todo) {
	t.Status = todosInProgress
	t.CompletedAt = nil
}

// MarkPending 标记为待办
func MarkPending(t *Todo) {
	t.Status = todosPending
	t.CompletedAt = nil
}
