package main

import (
	"bytes"
	"github.com/streadway/amqp"
	. "go-rabbit"
	"log"
	"time"
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
	q, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)
	FailOnError(err, "Failed to declare a queue")

	// 给队列设置预取数为1，它告诉RabbitMQ不要一次性分发超过1个的消息给某一个消费者，换句话说，就是当分发给该消费者的前一个消息还没有收到ack确认时，RabbitMQ将不会再给它派发消息，而是寻找下一个空闲的消费者目标进行分发。
	// 默认是轮询调度
	err = ch.Qos(1, 0, false)
	FailOnError(err, "Failed to set Qos")

	// 消费者
	// 自动确认
	//msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	//FailOnError(err, "Failed to register a consumer")

	// 手动确认
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message : %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t*time.Second)
			log.Println("Done")
			// 手动确认
			d.Ack(false)
		}
	}()

	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}
