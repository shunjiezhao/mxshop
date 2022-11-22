package queue

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"server/shared/queue"
	"strconv"
)

type Publisher struct {
	ch       *amqp.Channel
	exchange string
}

//Publish publishes a message.
func (p *Publisher) Publish(c context.Context, uid queue.UserId) error {
	msg := strconv.Itoa(int(uid))
	if msg == "" {
		return fmt.Errorf("user_id:%d 转换为 '' ", uid)
	}
	return p.ch.Publish(
		p.exchange,
		"",    //Key
		false, //mandatory
		false, //immedaiiote
		amqp.Publishing{
			Body: []byte(msg),
		},
	)
}

func NewPublisher(conn *amqp.Connection, exchange string) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}

	err = declareExchange(ch, exchange)

	if err != nil {
		return nil, fmt.Errorf("cannot declare exchange: %v", err)
	}
	return &Publisher{
		ch:       ch,
		exchange: exchange,
	}, nil

}
func declareExchange(ch *amqp.Channel, exchange string) error {
	return ch.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
}

type pk_serviceriber struct {
	conn     *amqp.Connection
	exchange string
	logger   *zap.Logger
}
