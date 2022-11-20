package initialize

import (
	"github.com/streadway/amqp"
	"server/pk_service/global"
	"server/pk_service/utils/queue"
	queue2 "server/shared/queue"
)

// 初始话 rabbitmq
func InitQueue() {
	var err error
	//TODO: 将地址配置化
	amqpConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	global.UserWaitQueue, err = queue.NewPublisher(amqpConn, queue2.UserQueueExchangeName)
	if err != nil {
		panic(err)
	}
}
