package mq

import (
	"TODO_API/internal/domain/event"
	"context"
	"encoding/json"
	"fmt"
)

// EventPublisher 事件发布器接口
type EventPublisher interface {
	// PublishEvent 发布领域事件
	PublishEvent(ctx context.Context, e event.DomainEvent)
}

// RabbitMQEventPublisher RabbitMQ 事件发布器实现
type RabbitMQEventPublisher struct {
	rmqClient RabbitMQClient
	exchange  string
}

// NewRabbitMQEventPublisher 创建 RabbitMQ 事件发布器实例
func NewRabbitMQEventPublisher(rmqClient RabbitMQClient, exchange string) EventPublisher {
	// 声明交换机
	ctx := context.Background()
	rmqClient.DeclareExchange(ctx, exchange, "topic", true, false, false, false, nil)

	return &RabbitMQEventPublisher{
		rmqClient: rmqClient,
		exchange:  exchange,
	}
}

// PublishEvent 发布领域事件
func (p *RabbitMQEventPublisher) PublishEvent(ctx context.Context, e event.DomainEvent) {
	// 序列化事件
	eventJSON, err := json.Marshal(e)
	if err != nil {
		fmt.Printf("序列化事件失败: %v\n", err)
		return
	}

	// 发布事件
	routingKey := e.EventType()
	err = p.rmqClient.Publish(ctx, p.exchange, routingKey, eventJSON)
	if err != nil {
		fmt.Printf("发布事件失败: %v\n", err)
		return
	}

	fmt.Printf("事件发布成功: %s\n", e.EventType())
}
