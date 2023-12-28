package gobroker

import (
	"errors"

	"golang.org/x/sync/errgroup"
)

type Broker interface {
	Connect(url string) error
	Disconnect()

	Init() error
	Consume(consumer Consumer) error
}

type Server struct {
	broker Broker

	consumers  map[string]*Consumer
	middleware MiddlewareChain
}

func NewServer(broker Broker) *Server {
	return &Server{
		broker:     broker,
		consumers:  make(map[string]*Consumer),
		middleware: MiddlewareChain{},
	}
}

func (s *Server) AddConsumer(queue string, handler ConsumerHandlerFunc) *Consumer {
	c := &Consumer{
		Context: ConsumerContext{
			Queue:  queue,
			Params: make(map[any]any),
		},
		Handler: handler,
	}

	s.consumers[queue] = c
	return c
}

func (s *Server) Use(middleware Middleware) {
	s.middleware.add(middleware)
}

var (
	ConnectionError = errors.New("server could not connect to broker")
	InitError       = errors.New("there was an error initializing the broker")
	ConsumingError  = errors.New("there was an error consuming a message")
)

func (s *Server) Run(url string) error {
	err := s.broker.Connect(url)
	if err != nil {
		return ConnectionError
	}
	defer s.broker.Disconnect()

	err = s.broker.Init()
	if err != nil {
		return InitError
	}

	group := errgroup.Group{}
	for _, consumer := range s.consumers {
		func(c Consumer) {
			group.Go(func() error {
				c.Handler = s.middleware.apply(c.Handler)
				return s.broker.Consume(c)
			})
		}(*consumer)
	}

	if err := group.Wait(); err != nil {
		return ConsumingError
	}

	return nil
}
