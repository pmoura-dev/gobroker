package middleware

import (
	"log"

	"github.com/pmoura-dev/gobroker"
)

func Logging(next gobroker.ConsumerHandlerFunc) gobroker.ConsumerHandlerFunc {
	return func(ctx gobroker.ConsumerContext, message gobroker.Message) error {
		log.Printf("New message from: %s\n", ctx.Queue)
		return next(ctx, message)
	}
}
