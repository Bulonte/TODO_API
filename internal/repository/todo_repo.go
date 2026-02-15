package repository

import (
	"TODO_API/internal/domain/model"
	"context"

	"gorm.io/gorm"
)

// TodoRepository 待办事项仓储接口
type TodoRepository interface {
	Create(ctx context.Context, todo *model.Todo) error
	GetByID(ctx context.Context, id uint) (*model.Todo, error)
	GetByUserID(ctx context.Context, userID uint, page, pageSize uint,
		status, priority *uint8, keyword string) ([]model.Todo, int64, error)
	Update(ctx context.Context, todo *model.Todo) error
	Delete(ctx context.Context, id uint) error
	BatchUpdateStatus(ctx context.Context, userID uint, todoIDs []uint, status model.TodoStatus) error
	GetStatistics(ctx context.Context, userID uint) (*model.TodoStatistics, error)
}

type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository 创建待办事项仓储实例
func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

// Create 创建待办事项
func (r *todoRepository) Create(ctx context.Context, todo *model.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

// GetByID 根据ID获取待办事项
func (r *todoRepository) GetByID(ctx context.Context, id uint) (*model.Todo, error) {
	var todo model.Todo
	err := r.db.WithContext(ctx).Preload("User").First(&todo, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &todo, nil
}

// GetByUserID 根据用户ID获取待办事项列表
func (r *todoRepository) GetByUserID(ctx context.Context, userID uint,
	page, pageSize uint, status, priority *uint8,
	keyword string) ([]model.Todo, int64, error) {
	var todos []model.Todo
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&model.Todo{}).Where("user_id = ?", userID)
	// 条件筛选
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if priority != nil {
		query = query.Where("priority = ?", *priority)
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	//获取总数
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	// 分页查询
	offset := int((page - 1) * pageSize)
	err := query.Order("CASE WHEN due_date IS NOT NULL THEN due_date ELSE '9999-12-31' END ASC").
		Order("priority DESC").
		Order("created_at DESC").
		Offset(offset).Limit(int(pageSize)).Find(&todos).Error

	return todos, totalCount, err
}

// Update 更新待办事项
func (r *todoRepository) Update(ctx context.Context, todo *model.Todo) error {
	return r.db.WithContext(ctx).Save(todo).Error
}

// Delete 删除待办事项
func (r *todoRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Todo{}, id).Error
}

// BatchUpdateStatus 批量更新状态
func (r *todoRepository) BatchUpdateStatus(ctx context.Context, userID uint,
	todoIDs []uint, status model.TodoStatus) error {
	return r.db.WithContext(ctx).Model(&model.Todo{}).
		Where("user_id = ? AND id IN (?)", userID, todoIDs).
		Update("status", status).Error
}

// GetStatistics 获取统计信息
func (r *todoRepository) GetStatistics(ctx context.Context,
	userID uint) (*model.TodoStatistics, error) {
	var sta model.TodoStatistics

	err := r.db.WithContext(ctx).Model(&model.Todo{}).
		Select("COUNT(*) as total_count, "+
			"SUM(CASE WHEN status = 0 THEN 1 ELSE 0 END) as pending_count, "+
			"SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) as progress_count, "+
			"SUM(CASE WHEN status = 2 THEN 1 ELSE 0 END) as completed_count").
		Where("user_id = ?", userID).
		Scan(&sta).Error

	if err != nil {
		return nil, err
	}
	return &sta, nil
}
