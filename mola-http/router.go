package mola

import (
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

type Router struct {
	routes *tree
}

func NewRouter() *Router {
	return &Router{
		routes: NewTree(), // Use NewTree() instead of NewTrieNode()
	}
}

func (r *Router) Handle(method string, path string, handler HandlerFunc) {
	// Convert HandlerFunc to http.Handler
	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := NewContext(w, req)
		handler(ctx)
	})

	// Insert route into the tree
	r.routes.Insert([]string{method}, path, httpHandler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := NewContext(w, req)

	log.Printf("Handling request: %s %s", req.Method, req.URL.Path)

	// Search for route
	res, params, err := r.routes.Search(req.Method, req.URL.Path)
	if err != nil {
		if err == ErrNotFound {
			NotFound(w, req)
			return
		} else if err == ErrMethodNotAllowed {
			MethodNotAllowed(w, req)
			return
		}
		log.Printf("Routing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set params and execute the handler
	ctx.Params = params
	res.actions.handler.ServeHTTP(w, req)
}

func NotFound(w http.ResponseWriter, req *http.Request) {
	log.Printf("Route Not Found: %v %v\n", req.Method, req.URL.Path)
	http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
}

func MethodNotAllowed(w http.ResponseWriter, req *http.Request) {
	log.Printf("Method Not Allowed: %v %v\n", req.Method, req.URL.Path)
	http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
}
