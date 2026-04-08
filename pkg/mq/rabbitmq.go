package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient RabbitMQ 客户端接口
type RabbitMQClient interface {
	// Publish 发布消息
	Publish(ctx context.Context, exchange, routingKey string, message []byte) error
	// Consume 消费消息
	Consume(ctx context.Context, queue, consumer string, autoAck bool) (<-chan amqp091.Delivery, error)
	// DeclareQueue 声明队列
	DeclareQueue(ctx context.Context, queue string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) (amqp091.Queue, error)
	// DeclareExchange 声明交换机
	DeclareExchange(ctx context.Context, exchange, kind string, durable, autoDelete, internal, noWait bool, args amqp091.Table) error
	// BindQueue 绑定队列
	BindQueue(ctx context.Context, queue, routingKey, exchange string, noWait bool, args amqp091.Table) error
	// Close 关闭连接
	Close() error
}

// rabbitMQClient RabbitMQ 客户端实现
type rabbitMQClient struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

// NewRabbitMQClient 创建 RabbitMQ 客户端实例
func NewRabbitMQClient(addr, username, password, vhost string) (RabbitMQClient, error) {
	// 连接 RabbitMQ
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, addr, vhost))
	if err != nil {
		return nil, err
	}

	// 创建通道
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &rabbitMQClient{
		conn:    conn,
		channel: channel,
	}, nil
}

// Publish 发布消息
func (r *rabbitMQClient) Publish(ctx context.Context, exchange, routingKey string, message []byte) error {
	return r.channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp091.Persistent,
			Timestamp:    time.Now(),
		},
	)
}

// Consume 消费消息
func (r *rabbitMQClient) Consume(ctx context.Context, queue, consumer string, autoAck bool) (<-chan amqp091.Delivery, error) {
	return r.channel.Consume(
		queue,
		consumer,
		autoAck,
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
}

// DeclareQueue 声明队列
func (r *rabbitMQClient) DeclareQueue(ctx context.Context, queue string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) (amqp091.Queue, error) {
	return r.channel.QueueDeclare(
		queue,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)
}

// DeclareExchange 声明交换机
func (r *rabbitMQClient) DeclareExchange(ctx context.Context, exchange, kind string, durable, autoDelete, internal, noWait bool, args amqp091.Table) error {
	return r.channel.ExchangeDeclare(
		exchange,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		args,
	)
}

// BindQueue 绑定队列
func (r *rabbitMQClient) BindQueue(ctx context.Context, queue, routingKey, exchange string, noWait bool, args amqp091.Table) error {
	return r.channel.QueueBind(
		queue,
		routingKey,
		exchange,
		noWait,
		args,
	)
}

// Close 关闭连接
func (r *rabbitMQClient) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
