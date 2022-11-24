package initialize

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"server/inventory_service/global"
	"server/inventory_service/utils/queue"
)

func InitQueue() {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		panic(err)
	}

	global.StockRebackPublisher, err = queue.NewPublisher(conn)
	if err != nil {
		panic(err)
	}

	global.StockRebackSubscriber, err = queue.NewSubscriber(conn)
	if err != nil {
		panic(err)
	}
}
