package mola

import (
	"log"
	"time"
)

type MiddlewareFunc func(HandlerFunc) HandlerFunc

func LoggerMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		start := time.Now()
		next(ctx)
		log.Printf("[%s] %s %v\n", ctx.Request.Method, ctx.Request.URL.Path, time.Since(start))
	}
}

func ApplyMiddleware(handler HandlerFunc, middlewares ...MiddlewareFunc) HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
