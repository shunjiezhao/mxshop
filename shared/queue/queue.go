package queue

import "context"

// 负责定义接口
type UserSubscriber interface {
	//TODO：可以路由进行优化
	Subscribe(ctx context.Context) (ch chan UserId, cleanUp func(), err error)
}

type UserPublisher interface {
	Publish(context.Context, UserId) error
}

var UserQueueExchangeName = "UserWaitQueueName"

type UserId int32
