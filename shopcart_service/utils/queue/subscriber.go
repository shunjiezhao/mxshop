package queue

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Subscriber struct {
	Finish  chan OrderInfo
	Release chan OrderInfo
}

// 创建 订单完成支付消息队列 和 订单释放队列
func NewSubscriber(conn *amqp.Connection) (*Subscriber, error) {
	sub := &Subscriber{
		Finish:  make(chan OrderInfo),
		Release: make(chan OrderInfo),
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	err = deTopicExchange(ch, OrderExName)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	q, err := deAndBind(ch, OrderFinishQ, nil, OrderExName, OrderFinishQKey)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	finish, err := ch.Consume(
		q.Name,         // queue
		"order.finish", // consumer
		true,           // auto ack
		false,          // exclusive
		false,          // no local
		false,          // no wait
		nil,            // args
	)

	go func() {
		for msg := range finish {
			var info OrderInfo
			err := json.Unmarshal(msg.Body, &info)
			if err != nil {
				log.Println("消息未接受成功， json编码失败:", err.Error())
				continue
			}
			fmt.Println("接受订单完成消息", string(msg.Body))

			sub.Finish <- info
		}
	}()

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	q, err = deAndBind(ch, OrderReleaseQ, nil, OrderExName, OrderReleaseQKey)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	release, err := ch.Consume(
		q.Name,          // queue
		"order.release", // consumer
		true,            // auto ack
		false,           // exclusive
		false,           // no local
		false,           // no wait
		nil,             // args
	)

	go func() {
		for msg := range release {
			var info OrderInfo
			err := json.Unmarshal(msg.Body, &info)

			if err != nil {
				log.Println("消息未接受成功， json编码失败:", err.Error())
				continue
			}

			fmt.Println("发送释放消息 ", string(msg.Body))
			sub.Release <- info
		}
	}()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return sub, nil
}

func deAndBind(ch *amqp.Channel, name string, args amqp.Table, exchange, key string) (amqp.Queue, error) {
	// 声明 完成订单支付消息队列
	q, err := declareQueue(ch, name, args)
	if err != nil {
		log.Println(err.Error())
		return q, err
	}
	// 绑定
	err = bindQueue(ch, q.Name, exchange, key)
	if err != nil {
		return q, err
	}
	return q, nil
}
