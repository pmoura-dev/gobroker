package main

import (
	"fmt"
	"log"

	"github.com/pmoura-dev/gobroker"
	"github.com/pmoura-dev/gobroker/brokers"
	"github.com/pmoura-dev/gobroker/middleware"
)

func fooHandler(ctx gobroker.ConsumerContext, message gobroker.Message) error {
	param := ctx.Params["foo"].(string)

	fmt.Println("Message received, here is your parameter:", param)
	fmt.Println("Body:", string(message.GetBody()))

	return ctx.Publisher.Publish([]byte("response"), "foo.response_key", map[string]any{
		"exchange": "foo.exchange",
	})
}

func responseHandler(ctx gobroker.ConsumerContext, message gobroker.Message) error {
	fmt.Println("Body:", string(message.GetBody()))
	fmt.Println("Response received")
	return nil
}

func main() {

	b := brokers.NewRabbitMQBroker()
	b.AddExchange("foo.exchange")
	b.AddQueue("foo.queue").Bind("foo.exchange", "foo.key")
	b.AddQueue("foo.responses").Bind("foo.exchange", "foo.response_key")

	s := gobroker.NewServer(b)

	s.Use(middleware.Logging)

	s.AddConsumer("foo.queue", fooHandler).AddParam("foo", "bar")
	s.AddConsumer("foo.responses", responseHandler)

	if err := s.Run("amqp://guest:guest@localhost:5672"); err != nil {
		log.Fatal("error")
	}
}
