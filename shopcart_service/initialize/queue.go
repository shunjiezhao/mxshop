package initialize

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"server/shopcart_service/global"
	"server/shopcart_service/utils/queue"
)

// 初始话 rabbitmq
func InitQueue() {
	var err error
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	global.OrderPublisher, err = queue.NewPublisher(conn)
	if err != nil {
		panic(err)
	}
	global.OrderSubscriber, err = queue.NewSubscriber(conn)
	if err != nil {
		panic(err)
	}
}
