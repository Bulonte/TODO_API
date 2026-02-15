package request

import "time"

// CreateTodoRequest 创建待办事项请求
type CreateTodoRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=200"`
	Description string     `json:"description,omitempty"`
	Status      uint8      `json:"status,omitempty" binding:"oneof=0 1 2"`
	Priority    uint8      `json:"priority,omitempty" binding:"oneof=1 2 3 4"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// UpdateTodoRequest 更新待办事项请求
type UpdateTodoRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=200"`
	Description string     `json:"description,omitempty"`
	Status      *uint8     `json:"status,omitempty" binding:"oneof=0 1 2"`
	Priority    *uint8     `json:"priority,omitempty" binding:"oneof=1 2 3 4"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// TodoQueryRequest 待办事项查询请求
type TodoQueryRequest struct {
	Page     uint   `form:"page,default=1" binding:"required,min=1"`
	PageSize uint   `form:"page,default=10" binding:"required,min=1,max=100"`
	Status   *uint8 `form:"status,omitempty" binding:"oneof=0 1 2"`
	Priority *uint8 `form:"priority,omitempty" binding:"oneof=1 2 3 4"`
	KeyWord  string `form:"keyword"`
}

// UpdateTodoStatusRequest 更新状态请求
type UpdateTodoStatusRequest struct {
	Status *uint8 `json:"status" binding:"oneof=0 1 2"`
}

// BatchUpdateRequest 批量操作请求
type BatchUpdateTodoRequest struct {
	TodoIDs []uint `json:"todo_ids" binding:"required,min=1"`
	Status  *uint8 `json:"status,omitempty" binding:"oneof=0 1 2"`
}
