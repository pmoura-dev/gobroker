package gobroker

type Middleware func(handler ConsumerHandlerFunc) ConsumerHandlerFunc

type MiddlewareChain struct {
	middlewares []Middleware
}

func (c *MiddlewareChain) add(middleware Middleware) {
	c.middlewares = append(c.middlewares, middleware)
}

func (c *MiddlewareChain) apply(handler ConsumerHandlerFunc) ConsumerHandlerFunc {
	for _, middleware := range c.middlewares {
		handler = middleware(handler)
	}
	return handler
}
