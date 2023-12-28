package brokers

import (
	"github.com/pmoura-dev/gobroker"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQMessage struct {
	delivery amqp091.Delivery
}

func (m RabbitMQMessage) GetCorrelationID() string {
	return m.delivery.CorrelationId
}

func (m RabbitMQMessage) GetBody() []byte {
	return m.delivery.Body
}

type RabbitMQQueue struct {
	name       string
	routingKey string
	exchange   string
}

func (q *RabbitMQQueue) Bind(exchange string, routingKey string) {
	q.exchange = exchange
	q.routingKey = routingKey
}

type RabbitMQBroker struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel

	queues    []*RabbitMQQueue
	exchanges []string
}

func NewRabbitMQBroker() *RabbitMQBroker {
	return &RabbitMQBroker{}
}

func (b *RabbitMQBroker) AddQueue(queue string) *RabbitMQQueue {
	q := &RabbitMQQueue{
		name: queue,
	}

	b.queues = append(b.queues, q)
	return q
}

func (b *RabbitMQBroker) AddExchange(exchange string) {
	b.exchanges = append(b.exchanges, exchange)
}

func (b *RabbitMQBroker) Connect(url string) error {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	b.conn = conn
	b.ch = ch
	return nil
}

func (b *RabbitMQBroker) Disconnect() {
	_ = b.ch.Close()
	_ = b.conn.Close()
}

func (b *RabbitMQBroker) Init() error {
	for _, exchange := range b.exchanges {
		err := b.ch.ExchangeDeclare(
			exchange,
			amqp091.ExchangeTopic,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	for _, queue := range b.queues {
		_, err := b.ch.QueueDeclare(
			queue.name,
			false,
			true,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		if queue.exchange != "" {
			err = b.ch.ExchangeDeclare(
				queue.exchange,
				amqp091.ExchangeTopic,
				true,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				return err
			}
		}

		if queue.routingKey != "" {
			err = b.ch.QueueBind(
				queue.name,
				queue.routingKey,
				queue.exchange,
				false,
				nil,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *RabbitMQBroker) Consume(consumer gobroker.Consumer) error {
	deliveryChan, err := b.ch.Consume(
		consumer.Context.Queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for delivery := range deliveryChan {
		message := RabbitMQMessage{delivery: delivery}
		err := consumer.Handler(consumer.Context, message)
		if err != nil {
			return err
		}
	}

	return nil
}
