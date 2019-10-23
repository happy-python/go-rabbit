package main

import (
	"github.com/streadway/amqp"
	. "go-rabbit"
	"log"
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
	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	FailOnError(err, "Failed to declare a queue")

	// 发送消息
	body := "hello"
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	FailOnError(err, "Failed to publish a message")
	log.Printf("[x] Sent %s", body)
}
