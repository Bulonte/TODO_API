package mq

import (
	"context"
	"encoding/json"
	"fmt"
)

// DefaultEventHandler 默认事件处理器
type DefaultEventHandler struct {}

// NewDefaultEventHandler 创建默认事件处理器
func NewDefaultEventHandler() EventHandler {
	return &DefaultEventHandler{}
}

// Handle 处理事件
func (h *DefaultEventHandler) Handle(ctx context.Context, eventType string, message []byte) {
	// 解析事件
	var eventData map[string]interface{}
	if err := json.Unmarshal(message, &eventData); err != nil {
		fmt.Printf("解析事件失败: %v\n", err)
		return
	}

	// 根据事件类型处理
	switch eventType {
	case "UserRegistered":
		h.handleUserRegistered(eventData)
	case "TodoCreated":
		h.handleTodoCreated(eventData)
	case "TodoCompleted":
		h.handleTodoCompleted(eventData)
	default:
		fmt.Printf("未知事件类型: %s\n", eventType)
	}
}

// handleUserRegistered 处理用户注册事件
func (h *DefaultEventHandler) handleUserRegistered(eventData map[string]interface{}) {
	fmt.Printf("处理用户注册事件: %v\n", eventData)
	// 这里可以添加具体的处理逻辑，比如发送欢迎邮件等
}

// handleTodoCreated 处理待办事项创建事件
func (h *DefaultEventHandler) handleTodoCreated(eventData map[string]interface{}) {
	fmt.Printf("处理待办事项创建事件: %v\n", eventData)
	// 这里可以添加具体的处理逻辑，比如发送通知等
}

// handleTodoCompleted 处理待办事项完成事件
func (h *DefaultEventHandler) handleTodoCompleted(eventData map[string]interface{}) {
	fmt.Printf("处理待办事项完成事件: %v\n", eventData)
	// 这里可以添加具体的处理逻辑，比如更新统计信息等
}
