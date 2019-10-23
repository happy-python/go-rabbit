package main

import (
	"github.com/streadway/amqp"
	. "go-rabbit"
	"log"
	"strconv"
)

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

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
	q, err := ch.QueueDeclare("rpc_queue", false, false, false, false, nil)
	FailOnError(err, "Failed to declare a queue")

	// 给队列设置预取数为1，它告诉RabbitMQ不要一次性分发超过1个的消息给某一个消费者，换句话说，就是当分发给该消费者的前一个消息还没有收到ack确认时，RabbitMQ将不会再给它派发消息，而是寻找下一个空闲的消费者目标进行分发。
	// 默认是轮询调度
	err = ch.Qos(1, 0, false)
	FailOnError(err, "Failed to set Qos")

	// 等待接受请求
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	FailOnError(err, "Failed to set Qos")

	forever := make(chan bool)

	// 开启协程处理请求
	go func() {
		for d := range msgs {
			n, err := strconv.Atoi(string(d.Body))
			FailOnError(err, "Failed to convert body to an integer")

			log.Printf("fib(%d)", n)

			response := fib(n)

			// 发送响应
			err = ch.Publish("", d.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          []byte(strconv.Itoa(response)),
			})
			FailOnError(err, "Failed to publish a message")
			log.Printf("Sent %d", response)

			d.Ack(false)
		}
	}()

	log.Printf("Awaiting RPC reqeusts")
	<-forever
}
