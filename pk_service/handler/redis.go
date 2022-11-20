package handler

import (
	"context"
	"fmt"
	"server/shared/queue"
)

type RedisWatcher struct {
	Add chan queue.UserId // 通知有用户加入了
	Ctx context.Context
	queue.UserPublisher
}

func (r *RedisWatcher) Watch() {
	defer close(r.Add)
	for {
		select {
		case uid := <-r.Add:
			//TODO: 来人了
			r.Publish(r.Ctx, uid)
		case <-r.Ctx.Done():
			// 结束了
			goto end
		}
	}
end:
	fmt.Println("end")
}
