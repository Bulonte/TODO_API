package persistence

import (
	"TODO_API/internal/repository"
	"context"
	"gorm.io/gorm"
)

// UnitOfWork 工作单元接口
type UnitOfWork interface {
	// 获取用户仓库
	UserRepository() repository.UserRepository
	// 获取待办事项仓库
	TodoRepository() repository.TodoRepository
	// 提交事务
	Commit(ctx context.Context) error
	// 回滚事务
	Rollback(ctx context.Context) error
	// 获取数据库连接
	DB() *gorm.DB
}

// unitOfWork 工作单元实现
type unitOfWork struct {
	db           *gorm.DB
	userRepo     repository.UserRepository
	todoRepo     repository.TodoRepository
	transaction  *gorm.DB
}

// NewUnitOfWork 创建工作单元实例
func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	tx := db.Begin()
	return &unitOfWork{
		db:          db,
		transaction: tx,
		userRepo:    repository.NewUserRepository(tx),
		todoRepo:    repository.NewTodoRepository(tx),
	}
}

// UserRepository 获取用户仓库
func (u *unitOfWork) UserRepository() repository.UserRepository {
	return u.userRepo
}

// TodoRepository 获取待办事项仓库
func (u *unitOfWork) TodoRepository() repository.TodoRepository {
	return u.todoRepo
}

// Commit 提交事务
func (u *unitOfWork) Commit(ctx context.Context) error {
	return u.transaction.WithContext(ctx).Commit().Error
}

// Rollback 回滚事务
func (u *unitOfWork) Rollback(ctx context.Context) error {
	return u.transaction.WithContext(ctx).Rollback().Error
}

// DB 获取数据库连接
func (u *unitOfWork) DB() *gorm.DB {
	return u.transaction
}
