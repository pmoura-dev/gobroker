package gobroker

type Message interface {
	GetCorrelationID() string
	GetTopic() string
	GetBody() []byte
}

type Consumer struct {
	Context ConsumerContext
	Handler ConsumerHandlerFunc
}

func (c *Consumer) AddParam(key any, value any) *Consumer {
	c.Context.Params[key] = value
	return c
}

type ConsumerContext struct {
	Queue  string
	Params map[any]any

	Publisher Publisher
}

type ConsumerHandlerFunc func(ctx ConsumerContext, message Message) error
