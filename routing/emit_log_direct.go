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
	err = ch.ExchangeDeclare("logs_direct", "direct", true, false, false, false, nil)
	FailOnError(err, "Failed to declare a exchange")

	// 发送消息
	body := BodyForm2(os.Args)
	// binding_key
	bindingKey := SeverityForm(os.Args)
	err = ch.Publish("logs_direct", bindingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	FailOnError(err, "Failed to publish a message")
	log.Printf("[x] Sent %s", body)
}
