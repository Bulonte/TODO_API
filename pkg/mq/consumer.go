package mq

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

// EventHandler 事件处理器接口
type EventHandler interface {
	// Handle 处理事件
	Handle(ctx context.Context, eventType string, message []byte)
}

// Consumer 消费者接口
type Consumer interface {
	// Start 启动消费者
	Start(ctx context.Context)
	// Stop 停止消费者
	Stop()
}

// RabbitMQConsumer RabbitMQ 消费者实现
type RabbitMQConsumer struct {
	rmqClient    RabbitMQClient
	exchange     string
	queue        string
	routingKey   string
	eventHandler EventHandler
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewRabbitMQConsumer 创建 RabbitMQ 消费者实例
func NewRabbitMQConsumer(rmqClient RabbitMQClient, exchange, queue, routingKey string, eventHandler EventHandler) Consumer {
	ctx, cancel := context.WithCancel(context.Background())

	// 声明队列
	rmqClient.DeclareQueue(ctx, queue, true, false, false, false, nil)

	// 绑定队列
	rmqClient.BindQueue(ctx, queue, routingKey, exchange, false, nil)

	return &RabbitMQConsumer{
		rmqClient:    rmqClient,
		exchange:     exchange,
		queue:        queue,
		routingKey:   routingKey,
		eventHandler: eventHandler,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start 启动消费者
func (c *RabbitMQConsumer) Start(ctx context.Context) {
	// 消费消息
	deliveries, err := c.rmqClient.Consume(ctx, c.queue, "", false)
	if err != nil {
		fmt.Printf("消费消息失败: %v\n", err)
		return
	}

	// 处理消息
	go func() {
		for d := range deliveries {
			var delivery amqp091.Delivery = d
			// 处理消息
			c.eventHandler.Handle(ctx, delivery.RoutingKey, delivery.Body)

			// 确认消息
			delivery.Ack(false)
		}
	}()

	fmt.Printf("消费者启动成功: %s\n", c.queue)
}

// Stop 停止消费者
func (c *RabbitMQConsumer) Stop() {
	c.cancel()
	fmt.Printf("消费者停止成功: %s\n", c.queue)
}
