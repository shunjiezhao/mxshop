package initialize

import (
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	queue2 "web-api/shared/queue"
	"web-api/user-web/global"
	"web-api/user-web/utils/queue"
)

// 初始话 rabbitmq
func InitQueue(logger *zap.Logger) {
	var err error
	//TODO: 将地址配置化
	amqpConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	global.UserEnterSubscriber, err = queue.NewSubscriber(amqpConn, queue2.UserEnterQExName, logger)
	if err != nil {
		panic(err)
	}

	global.UserCompleteSubscriber, err = queue.NewSubscriber(amqpConn, queue2.UserCmlQExName, logger)
	if err != nil {
		panic(err)
	}
}
