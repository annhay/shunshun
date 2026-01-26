package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// RabbitMQ 连接和通道
var (
	rmqConn    *amqp091.Connection
	rmqChannel *amqp091.Channel
	rmqConfig  RabbitMQConfig
)

// RabbitMQConfig RabbitMQ 配置

type RabbitMQConfig struct {
	Host       string
	Port       int
	User       string
	Password   string
	VHost      string
	Exchange   string
	Queue      string
	RoutingKey string
}

// OrderNotification 订单通知消息结构
type OrderNotification struct {
	OrderID     int64  `json:"order_id"`
	UserID      int64  `json:"user_id"`
	DriverID    int64  `json:"driver_id,omitempty"`
	OrderStatus string `json:"order_status"`
	Message     string `json:"message"`
	Timestamp   int64  `json:"timestamp"`
}

// InitRabbitMQ 初始化 RabbitMQ 连接
func InitRabbitMQ(config RabbitMQConfig) error {
	rmqConfig = config

	// 构建连接字符串
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.User, config.Password, config.Host, config.Port, config.VHost)

	// 连接 RabbitMQ
	var err error
	rmqConn, err = amqp091.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// 创建通道
	rmqChannel, err = rmqConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	// 声明交换机
	err = rmqChannel.ExchangeDeclare(
		config.Exchange, // 交换机名称
		"direct",       // 交换机类型
		true,            // 持久化
		false,           // 自动删除
		false,           // 内部
		false,           // 无等待
		nil,             // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	// 声明队列
	_, err = rmqChannel.QueueDeclare(
		config.Queue, // 队列名称
		true,         // 持久化
		false,        // 自动删除
		false,        // 独占
		false,        // 无等待
		nil,          // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// 绑定队列到交换机
	err = rmqChannel.QueueBind(
		config.Queue,      // 队列名称
		config.RoutingKey, // 路由键
		config.Exchange,   // 交换机名称
		false,             // 无等待
		nil,               // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to bind a queue: %w", err)
	}

	log.Println("RabbitMQ initialized successfully")
	return nil
}

// CloseRabbitMQ 关闭 RabbitMQ 连接
func CloseRabbitMQ() {
	if rmqChannel != nil {
		rmqChannel.Close()
	}
	if rmqConn != nil {
		rmqConn.Close()
	}
	log.Println("RabbitMQ connection closed")
}

// SendOrderNotification 发送订单通知
func SendOrderNotification(notification OrderNotification) error {
	if rmqChannel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	// 设置时间戳
	notification.Timestamp = time.Now().Unix()

	// 序列化消息
	messageBody, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// 发送消息
	err = rmqChannel.PublishWithContext(
		context.Background(),
		rmqConfig.Exchange,   // 交换机
		rmqConfig.RoutingKey, // 路由键
		false,                // 强制
		false,                // 立即
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         messageBody,
			DeliveryMode: amqp091.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Order notification sent: OrderID=%d, Status=%s", notification.OrderID, notification.OrderStatus)
	return nil
}

// ConsumeOrderNotifications 消费订单通知
func ConsumeOrderNotifications(callback func(notification OrderNotification) error) error {
	if rmqChannel == nil {
		return fmt.Errorf("RabbitMQ channel not initialized")
	}

	// 消费消息
	messages, err := rmqChannel.Consume(
		rmqConfig.Queue, // 队列名称
		"",             // 消费者标签
		false,          // 自动确认
		false,          // 独占
		false,          // 无本地
		false,          // 无等待
		nil,            // 参数
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// 处理消息
	go func() {
		for msg := range messages {
			var notification OrderNotification
			err := json.Unmarshal(msg.Body, &notification)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				msg.Nack(false, false)
				continue
			}

			// 调用回调函数处理通知
			err = callback(notification)
			if err != nil {
				log.Printf("Failed to process notification: %v", err)
				msg.Nack(false, true) // 重新入队
				continue
			}

			// 确认消息
			msg.Ack(false)
		}
	}()

	log.Println("Order notification consumer started")
	return nil
}
