package queue

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"strconv"
	"web-api/shared/queue"
)

type Publisher struct {
	ch       *amqp.Channel
	exchange string
}

//Publish publishes a message.
func (p *Publisher) Publish(c context.Context, uid queue.UserId) error {
	msg := strconv.Itoa(int(uid))
	if msg != "" {
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

type Subscriber struct {
	conn     *amqp.Connection
	exchange string
	logger   *zap.Logger
}

func NewSubscriber(conn *amqp.Connection, exchange string, logger *zap.Logger) (*Subscriber, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("can not allocate channel:%v", err)
	}
	defer channel.Close()
	err = declareExchange(channel, exchange)
	if err != nil {
		return nil, fmt.Errorf("can not declare exchange:%v", err)
	}
	return &Subscriber{
		conn:     conn,
		exchange: exchange,
		logger:   logger,
	}, nil

}
func (s *Subscriber) Subscribe(ctx context.Context) (chan queue.UserId, func(), error) {
	raw, cleanUp, err := s.SubscribeRaw(ctx)
	if err != nil {
		return nil, cleanUp, err
	}
	carCh := make(chan queue.UserId)
	go func() {
		for msg := range raw {
			var uid uint64
			if uid, err = strconv.ParseUint(string(msg.Body), 10, 32); err != nil {
				s.logger.Info("解析uid错误", zap.Error(err))
			}
			carCh <- queue.UserId(uid)
		}
		close(carCh)
	}()
	return carCh, cleanUp, nil
}

func (s *Subscriber) SubscribeRaw(ctx context.Context) (<-chan amqp.Delivery, func(), error) {
	ch, err := s.conn.Channel()
	if err != nil {
		return nil, func() {}, fmt.Errorf("can not allocate channel:%v", err)
	}

	closeCh := func() {
		err := ch.Close()
		if err != nil {
			s.logger.Error("can not close channel", zap.Error(err))
		}
	}
	q, err := ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil)
	if err != nil {
		return nil, closeCh, fmt.Errorf("can not allocate queue:%v", err)
	}

	cleanUp := func() {
		// 最后关闭channel
		defer closeCh()
		// 先关闭 删除队列
		if _, err := ch.QueueDelete(q.Name, false, false, false); err != nil {
			s.logger.Error("can not delete queue", zap.String("queueName", q.Name))
		}
	}
	err = ch.QueueBind(
		q.Name,
		"",
		s.exchange,
		false,
		nil)
	if err != nil {
		return nil, cleanUp, fmt.Errorf("can not  bind queue:%v", err)
	}
	consume, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return nil, cleanUp, fmt.Errorf("can not allocate consume:%v", err)
	}
	return consume, cleanUp, nil
}
