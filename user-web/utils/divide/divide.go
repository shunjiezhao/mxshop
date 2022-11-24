package divide

import (
	"context"
	"fmt"
	"sync"
	"web-api/shared/queue"
	userpb "web-api/user-web/proto"
)

type Message struct {
	Uid int32
	Msg []byte
}
type Result struct {
	OtherID queue.UserId  // 另一个 user id
	Read    chan *Message // 其他用户发来的消息
	Write   chan *Message // 发往其他一个人的消息管道
}
type UserDivide struct {
	enterSub    queue.UserSubscriber
	completeSub queue.UserSubscriber
	sync.RWMutex
	store    map[queue.UserId]chan *Result
	msg      map[queue.UserId]chan []byte
	PKClient userpb.PKClient
}

func NewDivide(PKClient userpb.PKClient, enterSub queue.UserSubscriber, completeSub queue.UserSubscriber) *UserDivide {
	u := &UserDivide{
		enterSub:    enterSub,
		completeSub: completeSub,
		RWMutex:     sync.RWMutex{},
		store:       make(map[queue.UserId]chan *Result),
		msg:         make(map[queue.UserId]chan []byte),
		PKClient:    PKClient,
	}
	go u.Divide()
	return u
}
func (u *UserDivide) Register(id queue.UserId, ch chan *Result, receive chan []byte) error {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.store[id]; !ok {
		u.store[id] = ch
		u.msg[id] = receive
		return nil
	}
	return fmt.Errorf("已经注册过了")
}
func (u *UserDivide) UnRegister(id queue.UserId) {
	u.Lock()
	defer u.Unlock()
	fmt.Println("delete", id)
	delete(u.store, id)
	delete(u.msg, id)
}

func (u *UserDivide) Divide() {
	cnt := 0                         // 统计当前拿了有几个用户 超过两个就会将匹配成功， 给对方的管道发送自己的id
	userIdStore := [2]queue.UserId{} //存放两个用户id
	ch, cleanUp, err := u.enterSub.Subscribe(context.Background())
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
			id0 := userIdStore[0]
			id1 := userIdStore[1]
			if id0 > id1 {
				id0, id1 = id1, id0
			}

			ch0Read := make(chan *Message)
			ch1Read := make(chan *Message)
			fmt.Println("%v %v", ch0Read, ch1Read)

			// 消息发送的流程
			// 1. a 收到消息 写入 read 管道
			// 2. b 从 read 管道收到消息 写入websocket
			u.store[id0] <- &Result{OtherID: id1, Read: ch1Read, Write: ch0Read}
			u.store[id1] <- &Result{OtherID: id0, Read: ch0Read, Write: ch1Read}

			//TODO:处理错误
			u.PKClient.Create(context.Background(), &userpb.CreateRequest{
				Id1: int32(id0),
				Id2: int32(id1),
			})
			delete(u.store, id0)
			delete(u.store, id1)
			cnt = 0
		}
	}
}

func (u *UserDivide) SendMsg(uid queue.UserId, msg []byte) {
	u.msg[uid] <- msg
}
