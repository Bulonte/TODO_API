package event

import (
	"time"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	// EventType 获取事件类型
	EventType() string
	// EventID 获取事件ID
	EventID() string
	// OccurredOn 获取事件发生时间
	OccurredOn() time.Time
}

// BaseEvent 基础事件结构
type BaseEvent struct {
	ID        string    `json:"id"`
	OccurredAt time.Time `json:"occurred_at"`
}

// EventID 获取事件ID
func (e *BaseEvent) EventID() string {
	return e.ID
}

// OccurredOn 获取事件发生时间
func (e *BaseEvent) OccurredOn() time.Time {
	return e.OccurredAt
}

// UserRegistered 用户注册事件
type UserRegistered struct {
	BaseEvent
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// EventType 获取事件类型
func (e *UserRegistered) EventType() string {
	return "user.registered"
}

// TodoCreated 待办事项创建事件
type TodoCreated struct {
	BaseEvent
	TodoID    uint   `json:"todo_id"`
	UserID    uint   `json:"user_id"`
	Title     string `json:"title"`
	Priority  uint8  `json:"priority"`
	DueDate   *time.Time `json:"due_date,omitempty"`
}

// EventType 获取事件类型
func (e *TodoCreated) EventType() string {
	return "todo.created"
}

// TodoCompleted 待办事项完成事件
type TodoCompleted struct {
	BaseEvent
	TodoID      uint      `json:"todo_id"`
	UserID      uint      `json:"user_id"`
	CompletedAt time.Time `json:"completed_at"`
}

// EventType 获取事件类型
func (e *TodoCompleted) EventType() string {
	return "todo.completed"
}

// TodoDueSoon 待办事项即将到期事件
type TodoDueSoon struct {
	BaseEvent
	TodoID   uint      `json:"todo_id"`
	UserID   uint      `json:"user_id"`
	Title    string    `json:"title"`
	DueDate  time.Time `json:"due_date"`
	DaysLeft int       `json:"days_left"`
}

// EventType 获取事件类型
func (e *TodoDueSoon) EventType() string {
	return "todo.due_soon"
}
