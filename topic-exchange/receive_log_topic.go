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

	// 创建Exchange
	err = ch.ExchangeDeclare("logs_topic", "topic", true, false, false, false, nil)
	FailOnError(err, "Failed to declare a exchange")

	// 定义临时Queue，生成随机名字
	// exclusive设置为true，当连接断开，队列将会被删除
	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	FailOnError(err, "Failed to declare a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warnint] [error]", os.Args[0])
		os.Exit(0)
	}

	// Bind
	for _, routingKey := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "logs_topic", routingKey)

		err = ch.QueueBind(q.Name, routingKey, "logs_topic", false, nil)
		FailOnError(err, "Failed to bind q queue")
	}

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
