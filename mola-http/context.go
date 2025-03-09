package mola

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writter http.ResponseWriter
	Request *http.Request
	Params  map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writter: w,
		Request: r,
		Params:  make(map[string]string),
	}
}

// JSON response
func (c *Context) JSON(status int, data interface{}) {
	c.Writter.Header().Set("Content-Type", "application/json")
	c.Writter.WriteHeader(status)
	json.NewEncoder(c.Writter).Encode(data)
}

// Plain text Response
func (c *Context) Text(status int, message string) {
	c.Writter.Header().Set("Content-Type", "text/plain")
	c.Writter.WriteHeader(status)
	c.Writter.Write([]byte(message))
}
