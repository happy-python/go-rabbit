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

	// 消费者
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message : %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}
