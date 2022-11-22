package queue

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// 要改的话把 库存服务的也要改
var (
	OrderPrefix = "order."
	OrderExName = OrderPrefix + "exchange"
	// keys
	OrderReleaseQKey = OrderPrefix + "release.order" // 释放订单

	StockPrefix      = "stock."
	StockExName      = StockPrefix + "exchange"
	StockDelayQ      = StockPrefix + "delay"
	StockReleaseQ    = StockPrefix + "release"
	StockReleaseQKey = StockPrefix + "release.stock" // stock or order
	StockDelayQKey   = StockPrefix + "delay.stock"
)

type Publisher struct {
	ch       *amqp.Channel
	exchange string
}
type OrderInfo struct {
	OrderSn string `json:"order_sn"` // 订单唯一编号
	Stock   int    `json:"stock"`    //扣除多少库存
	Test    string `json:"time"`     // 测试看d
}

//延时任务发 StockDelayQKey
func (p *Publisher) Publish(c context.Context, key string, info OrderInfo) error {
	b, err := json.Marshal(info)
	if err != nil {
		log.Println("消息未发送成功， json编码失败:", err.Error())
		return err
	}
	return p.ch.PublishWithContext(c,
		p.exchange, // exchange
		key,        // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        b,
		})
}

// 创建一个延时队列
func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	// 1. 生成 一个共享的exchange
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}
	// 生成 交换
	err = deTopicExchange(ch, StockExName)

	if err != nil {
		panic(err)
	}

	// 生成订单取消的延迟队列
	q, err := declareQueue(ch, StockDelayQ, map[string]interface{}{
		//TODO:这里需要参数化
		"x-message-ttl":             6000,             // 这是 ttl的时间
		"x-dead-letter-exchange":    StockExName,      // 声明死信队列
		"x-dead-letter-routing-key": OrderReleaseQKey, // 定时检查释放
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// 延迟队列的绑定
	err = bindQueue(ch, q.Name, StockExName, StockDelayQKey)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &Publisher{
		ch:       ch,
		exchange: StockExName,
	}, nil
}

func declareQueue(ch *amqp.Channel, name string, args amqp.Table) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		args,  // arguments
	)
}
func bindQueue(ch *amqp.Channel, name, exchange, key string) error {

	return ch.QueueBind(
		name,     // queue name
		key,      // routing key
		exchange, // exchange
		false,
		nil)

}
func deTopicExchange(ch *amqp.Channel, exchange string) error {
	return ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
}
