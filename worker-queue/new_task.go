package main

import (
	"github.com/streadway/amqp"
	. "go-rabbit"
	"log"
	"os"
)

func main() {
	// 建立连接
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 创建Channel
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 定义Queue
	// durable定义为true，队列的持久化
	// durable参数在生产者和消费者中都要指定为true
	q, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)
	FailOnError(err, "Failed to declare a queue")

	// 发送消息
	// 消息的持久化
	body := BodyForm(os.Args)
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		DeliveryMode:amqp.Persistent,
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	FailOnError(err, "Failed to publish a message")
	log.Printf("[x] Sent %s", body)
}
