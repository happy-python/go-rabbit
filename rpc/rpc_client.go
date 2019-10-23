package main

import (
	"github.com/streadway/amqp"
	. "go-rabbit"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func fibRPC(n int) (res int, err error) {
	// 建立连接
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 创建Channel
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 定义Queue
	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	FailOnError(err, "Failed to declare a queue")

	// 等待接受响应
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	FailOnError(err, "Faield to register a consumer")

	corrId := RandomString(32)

	// 发送请求
	// CorrelationId用来表示request和response的关联关系
	err = ch.Publish("", "rpc_queue", false, false, amqp.Publishing{
		ContentType:   "text/plain",
		CorrelationId: corrId,
		ReplyTo:       q.Name,
		Body:          []byte(strconv.Itoa(n)),
	})

	FailOnError(err, "Failed to publish a message")

	for d := range msgs {
		// 判断是否与请求匹配
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			FailOnError(err, "Failed to convert body to integer")
			break
		}
	}

	return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	n := BodyForm3(os.Args)

	log.Printf("Requesting fib(%d)", n)
	res, err := fibRPC(n)
	FailOnError(err, "Failed to handle RPC request")

	log.Printf("Got %d", res)
}
