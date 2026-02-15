package response

import "time"

// TodoResponse 待办事项响应
type TodoResponse struct {
	ID           uint       `json:"id"`
	UserID       uint       `json:"user_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description,omitempty"`
	Status       uint8      `json:"status"`
	StatusText   string     `json:"status_text"`
	Priority     uint8      `json:"priority"`
	PriorityText string     `json:"priority_text"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	CompletedAt  *time.Time `json:"completed,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	IsOverdue    bool       `json:"is_overdue"`
}

// Pagination 分页信息
type Pagination struct {
	Page       uint `json:"page"`
	PageSize   uint `json:"page_size"`
	Total      uint `json:"total"`
	TotalPages uint `json:"total_pages"`
}

// Statistics 统计信息
type Statistics struct {
	TotalCount      uint `json:"total_count"`
	PendingCount    uint `json:"pending_count"`
	InProgressCount uint `json:"in_progress_count"`
	CompletedCount  uint `json:"completed_count"`
}

// TodoListResponse 待办事项列表响应
type TodoListResponse struct {
	Todos      []TodoResponse `json:"todos"`
	Pagination Pagination     `json:"pagination"`
	Statistics Statistics     `json:"statistics,omitempty"`
}

// TodoStatsResponse 待办事项统计响应
type TodoStatsResponse struct {
	Statistics Statistics `json:"statistics"`
}
