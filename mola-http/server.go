package mola

import (
	"net/http"
)

type Server struct {
	router *Router
}

func NewServer() *Server {
	return &Server{router: NewRouter()}
}

func (s *Server) GET(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	s.router.Handle(http.MethodGet, path, ApplyMiddleware(handler, middlewares...))
}

func (s *Server) POST(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	s.router.Handle(http.MethodPost, path, ApplyMiddleware(handler, middlewares...))
}

func (s *Server) Start(port string) error {
	return http.ListenAndServe(":"+port, s.router)
}
