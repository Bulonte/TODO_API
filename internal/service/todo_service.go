package service

import (
	"TODO_API/internal/app/dto/request"
	"TODO_API/internal/app/dto/response"
	"TODO_API/internal/domain/model"
	"TODO_API/internal/repository"
	"context"
	"errors"
	"time"
)

// TodoService 待办事项服务接口
type TodoService interface {
	Create(ctx context.Context, userID uint, req *request.CreateTodoRequest) (*response.TodoResponse, error)
	GetTodoByID(ctx context.Context, id, userID uint) (*response.TodoResponse, error)
	GetTodos(ctx context.Context, userID uint, query *request.TodoQueryRequest) (*response.TodoListResponse, error)
	UpdateTodo(ctx context.Context, id, userID uint, req *request.UpdateTodoRequest) (*response.TodoResponse, error)
	DeleteTodo(ctx context.Context, id, userID uint) error
	UpdateTodoStatus(ctx context.Context, id, userID uint, status uint8) (*response.TodoResponse, error)
	BatchUpdateStatus(ctx context.Context, userID uint, req *request.BatchUpdateTodoRequest) error
}

type todoService struct {
	todoRepo repository.TodoRepository
}

// NewTodoService 创建待办事项服务实例
func NewTodoService(todoRepo repository.TodoRepository) TodoService {
	return &todoService{todoRepo: todoRepo}
}

// todoToResponse 将Todo模型转换为响应格式
func (s *todoService) todoToResponse(todo *model.Todo) *response.TodoResponse {
	var description string
	if todo.Description != nil {
		description = *todo.Description
	}
	isOverdue := false
	if todo.DueDate != nil && todo.Status != 2 {
		isOverdue = todo.DueDate.Before(time.Now())
	}

	return &response.TodoResponse{
		ID:           todo.ID,
		UserID:       todo.UserID,
		Title:        todo.Title,
		Description:  description,
		Status:       uint8(todo.Status),
		StatusText:   s.getStatusText(todo.Status),
		Priority:     uint8(todo.Priority),
		PriorityText: s.getPriorityText(todo.Priority),
		DueDate:      todo.DueDate,
		CompletedAt:  todo.CompletedAt,
		CreatedAt:    todo.CreatedAt,
		UpdatedAt:    todo.UpdatedAt,
		IsOverdue:    isOverdue,
	}
}

// statsToResponse 转换统计信息
func (s *todoService) statsToResponse(stats *model.TodoStatistics) response.Statistics {
	if stats == nil {
		return response.Statistics{}
	}

	return response.Statistics{
		TotalCount:      stats.TotalCount,
		PendingCount:    stats.PendingCount,
		InProgressCount: stats.ProgressCount,
		CompletedCount:  stats.CompletedCount,
	}
}

// getStatusText 获取状态文本
func (s *todoService) getStatusText(status model.TodoStatus) string {
	switch status {
	case 0:
		return "待办"
	case 1:
		return "进行中"
	case 2:
		return "已完成"
	default:
		return "未知"
	}
}

// getPriorityText 获取优先级文本
func (s *todoService) getPriorityText(priority model.TodosPriority) string {
	switch priority {
	case 1:
		return "低"
	case 2:
		return "中"
	case 3:
		return "高"
	case 4:
		return "紧急"
	default:
		return "未知"
	}
}

// CreateTodo 创建待办事项
func (s *todoService) Create(ctx context.Context, userID uint, req *request.CreateTodoRequest) (*response.TodoResponse, error) {
	todo := &model.Todo{
		UserID:   userID,
		Title:    req.Title,
		Status:   model.TodoStatus(req.Status),
		Priority: model.TodosPriority(req.Priority),
		DueDate:  req.DueDate,
	}
	if req.Description != "" {
		todo.Description = &req.Description
	}

	if err := s.todoRepo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return s.todoToResponse(todo), nil
}

// GetTodoByID 根据ID获取待办事项
func (s *todoService) GetTodoByID(ctx context.Context, id, userID uint) (*response.TodoResponse, error) {
	todo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if todo == nil {
		return nil, errors.New("待办事项不存在")
	}

	//检查权限
	if todo.UserID != userID {
		return nil, errors.New("无权限访问此待办事项")
	}

	return s.todoToResponse(todo), nil
}

// GetTodos 获取待办事项列表
func (s *todoService) GetTodos(ctx context.Context, userID uint, query *request.TodoQueryRequest) (*response.TodoListResponse, error) {
	var statusPtr, priorityPtr *uint8
	if query.Status != nil {
		statusPtr = query.Status
	}
	if query.Priority != nil {
		priorityPtr = query.Priority
	}

	todos, totalCount, err := s.todoRepo.GetByUserID(ctx, userID, query.Page, query.PageSize, statusPtr, priorityPtr, query.KeyWord)
	if err != nil {
		return nil, err
	}

	//转换为响应格式
	todosResponses := make([]response.TodoResponse, len(todos))
	for i, t := range todos {
		todosResponses[i] = *s.todoToResponse(&t)
	}
	// 计算分页信息
	totalPage := (uint(totalCount) + query.PageSize - 1) / query.PageSize
	//获取统计信息
	stats, _ := s.todoRepo.GetStatistics(ctx, userID)

	return &response.TodoListResponse{
		Todos: todosResponses,
		Pagination: response.Pagination{
			Page:       query.Page,
			PageSize:   query.PageSize,
			Total:      uint(totalCount),
			TotalPages: totalPage,
		},
		Statistics: s.statsToResponse(stats),
	}, nil
}

// UpdateTodo 更新待办事项
func (s *todoService) UpdateTodo(ctx context.Context, id, userID uint, req *request.UpdateTodoRequest) (*response.TodoResponse, error) {
	todo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if todo == nil {
		return nil, errors.New("待办事项不存在")
	}

	//检查权限
	if todo.UserID != userID {
		return nil, errors.New("无权限修改此待办事项")
	}

	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.Description != "" {
		todo.Description = &req.Description
	}
	if req.Status != nil {
		todo.Status = model.TodoStatus(*req.Status)
		if *req.Status == 2 {
			now := time.Now()
			todo.CompletedAt = &now
		} else {
			todo.CompletedAt = nil
		}
	}

	if req.Priority != nil {
		todo.Priority = model.TodosPriority(*req.Priority)
	}
	if req.DueDate != nil {
		todo.DueDate = req.DueDate
	}

	if err := s.todoRepo.Update(ctx, todo); err != nil {
		return nil, err
	}
	return s.todoToResponse(todo), nil
}

// DeleteTodo 删除待办事项
func (s *todoService) DeleteTodo(ctx context.Context, id, userID uint) error {
	todo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if todo == nil {
		return errors.New("待办事项不存在")
	}

	if todo.UserID != userID {
		return errors.New("无权限删除此待办事项")
	}

	return s.todoRepo.Delete(ctx, id)
}

// UpdateTodoStatus 更新待办事项状态
func (s *todoService) UpdateTodoStatus(ctx context.Context, id, userID uint, status uint8) (*response.TodoResponse, error) {
	todo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if todo == nil {
		return nil, errors.New("待办事项不存在")
	}

	if todo.UserID != userID {
		return nil, errors.New("无权限修改此待办事项")
	}

	todo.Status = model.TodoStatus(status)
	if status == 2 {
		now := time.Now()
		todo.CompletedAt = &now
	} else {
		todo.CompletedAt = nil
	}

	if err := s.todoRepo.Update(ctx, todo); err != nil {
		return nil, err
	}

	return s.todoToResponse(todo), nil
}

// BatchUpdateStatus 批量更新状态
func (s *todoService) BatchUpdateStatus(ctx context.Context, userID uint, req *request.BatchUpdateTodoRequest) error {
	return s.todoRepo.BatchUpdateStatus(ctx, userID, req.TodoIDs, model.TodoStatus(*req.Status))
}
