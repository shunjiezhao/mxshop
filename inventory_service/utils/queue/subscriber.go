package queue

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Subscriber struct {
	Release chan OrderInfo
}

// 创建 订单完成支付消息队列 和 订单释放队列
func NewSubscriber(conn *amqp.Connection) (*Subscriber, error) {
	sub := &Subscriber{
		Release: make(chan OrderInfo),
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = deTopicExchange(ch, StockExName)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	q, err := deAndBind(ch, StockReleaseQ, nil, StockExName, StockReleaseQKey)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	q, err = deAndBind(ch, StockReleaseQ, nil, OrderExName, OrderReleaseQKey)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	finish, err := ch.Consume(
		q.Name,          // queue
		"stock.release", //
		false,           // auto ack
		false,           // exclusive
		false,           // no local
		false,           // no wait
		nil,             // args
	)

	go func() {
		var info OrderInfo
		for msg := range finish {
			err := json.Unmarshal(msg.Body, &info)
			if err != nil {
				log.Println("消息未接受成功， json编码失败:", err.Error())
				continue
			}

			log.Println("消息接受成功", string(msg.Body))
			sub.Release <- info
		}
	}()

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
