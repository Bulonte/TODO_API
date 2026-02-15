package model

type TodoStatistics struct {
	TotalCount     uint `json:"total_count" gorm:"column:total_count"`
	PendingCount   uint `json:"pending_count" gorm:"column:pending_count"`
	ProgressCount  uint `json:"progress_count" gorm:"column:progress_count"`
	CompletedCount uint `json:"completed_count" gorm:"column:completed_count"`
}
