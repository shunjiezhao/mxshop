package divide

import (
	"context"
	"fmt"
	"sync"
	"web-api/shared/queue"
	userpb "web-api/user-web/proto"
)

type UserDivide struct {
	sub queue.UserSubscriber
	sync.RWMutex
	store    map[queue.UserId]chan queue.UserId
	PKClient userpb.PKClient
}

func NewDivide(PKClient userpb.PKClient, sub queue.UserSubscriber) *UserDivide {
	u := &UserDivide{
		sub:      sub,
		RWMutex:  sync.RWMutex{},
		store:    make(map[queue.UserId]chan queue.UserId),
		PKClient: PKClient,
	}
	go u.Divide()
	return u
}
func (u *UserDivide) Register(id queue.UserId, ch chan queue.UserId) error {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.store[id]; !ok {
		u.store[id] = ch
		return nil
	}
	return fmt.Errorf("已经注册过了")
}
func (u *UserDivide) UnRegister(id queue.UserId) {
	u.Lock()
	defer u.Unlock()
	delete(u.store, id)
}

func (u *UserDivide) Divide() {
	cnt := 0                         // 统计当前拿了有几个用户 超过两个就会将匹配成功， 给对方的管道发送自己的id
	userIdStore := [2]queue.UserId{} //存放两个用户id
	ch, cleanUp, err := u.sub.Subscribe(context.Background())
	defer cleanUp()
	if err != nil {
		panic(err)
		return
	}
	for userId := range ch {
		userIdStore[cnt] = userId
		cnt++
		if cnt == 2 {
			// 互相传递匹配的到来
			u.store[userIdStore[0]] <- userIdStore[1]
			u.store[userIdStore[1]] <- userIdStore[0]
			//TODO:处理错误
			u.PKClient.Create(context.Background(), &userpb.CreateRequest{
				Id1: int32(userIdStore[0]),
				Id2: int32(userIdStore[1]),
			})
			delete(u.store, userIdStore[0])
			delete(u.store, userIdStore[1])
			cnt = 0
		}
	}
}
