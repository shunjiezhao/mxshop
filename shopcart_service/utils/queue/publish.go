package queue

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	proto "server/inventory_service/proto/gen/v1/inventory"
)

// 要改的话把 库存服务的也要改
var (
	OrderPrefix   = "order."
	OrderExName   = OrderPrefix + "exchange"
	OrderDelayQ   = OrderPrefix + "delay"
	OrderFinishQ  = OrderPrefix + "finish"
	OrderReleaseQ = OrderPrefix + "release"

	// keys
	OrderDelayQKey   = OrderPrefix + "delay.order"   // 当创建成功向消息队列发送一个延时任务
	OrderReleaseQKey = OrderPrefix + "release.order" // 释放订单
	OrderFinishQKey  = OrderPrefix + "finish.order"
	DelayOrderTimeMs = 3000 // 3000ms
)

type Publisher struct {
	ch       *amqp.Channel
	exchange string
}

//TODO:这个proto需要从外面拉取嘛？
type OrderInfo struct {
	OrderSn   string                `protobuf:"bytes,2,opt,name=orderSn,proto3" json:"orderSn,omitempty"`
	GoodsInfo []*proto.GoodsInvInfo `protobuf:"bytes,1,rep,name=goodsInfo,proto3" json:"goodsInfo,omitempty"`
	Test      string                `json:"time"` // 测试看d
}

//延时任务发 OrderDelayQKey
//完成任务发 OrderFinishQKey
//释放订单发送 OrderReleaseQKey
func (p *Publisher) Publish(c context.Context, key string, info *OrderInfo) error {
	b, err := json.Marshal(info)
	if err != nil {
		log.Println("消息未发送成功， json编码失败:", err.Error())
		return err
	}
	fmt.Println(key, ":发送消息 ", string(b))
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
	err = deTopicExchange(ch, OrderExName)
	if err != nil {
		panic(err)
	}

	// 生成订单取消的延迟队列

	q, err := declareQueue(ch, OrderDelayQ, map[string]interface{}{
		//TODO:这里需要参数化
		"x-message-ttl":             DelayOrderTimeMs, // 这是 ttl的时间
		"x-dead-letter-exchange":    OrderExName,      // 声明死信队列
		"x-dead-letter-routing-key": OrderReleaseQKey, // 定时检查释放
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// 延迟队列的绑定
	err = bindQueue(ch, q.Name, OrderExName, OrderDelayQKey)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &Publisher{
		ch:       ch,
		exchange: OrderExName,
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
