package mola

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrNotFound         = errors.New("404 Not Found")
	ErrMethodNotAllowed = errors.New("405 Method Not Allowed")
)

type tree struct {
	node *node
}

type actions struct {
	handler http.Handler
}

type node struct {
	label     string
	actions   map[string]*actions
	children  map[string]*node
	isDynamic bool
	paramName string // Stores dynamic param name like ":id"
}

type result struct {
	actions *actions
}

const (
	pathRoot      = "/"
	pathDelimiter = "/"
)

func newResult() *result {
	return &result{}
}

func NewTree() *tree {
	return &tree{
		node: &node{
			label:    pathRoot,
			actions:  make(map[string]*actions),
			children: make(map[string]*node),
		},
	}
}

func (t *tree) Insert(methods []string, path string, handler http.Handler) error {
	curNode := t.node
	if path == pathRoot {
		curNode.label = path
		for _, method := range methods {
			curNode.actions[method] = &actions{handler: handler}
		}
		return nil
	}

	ep := explodePath(path)
	for i, p := range ep {
		isDynamic := strings.HasPrefix(p, ":")
		var key = p
		if isDynamic {
			key = ":" // Single identifier for all dynamic routes
		}

		nextNode, exists := curNode.children[key]
		if !exists {
			nextNode = &node{
				label:     p,
				isDynamic: isDynamic,
				paramName: p[1:], // Remove ":" prefix
				actions:   make(map[string]*actions),
				children:  make(map[string]*node),
			}
			curNode.children[key] = nextNode
		}

		curNode = nextNode
		if i == len(ep)-1 {
			for _, method := range methods {
				curNode.actions[method] = &actions{handler: handler}
			}
		}
	}

	return nil
}

func (t *tree) Search(method string, path string) (*result, map[string]string, error) {
	curNode := t.node
	params := make(map[string]string)

	if path != pathRoot {
		ep := explodePath(path)
		for _, p := range ep {
			if nextNode, ok := curNode.children[p]; ok {
				curNode = nextNode
			} else if dynamicNode, ok := curNode.children[":"]; ok { // Check dynamic route
				curNode = dynamicNode
				params[curNode.paramName] = p // Store param value
			} else {
				return nil, nil, ErrNotFound
			}
		}
	}

	if curNode.actions[method] == nil {
		return nil, nil, ErrMethodNotAllowed
	}

	return &result{actions: curNode.actions[method]}, params, nil
}

// explodePath removes an empty value in slice.
func explodePath(path string) []string {
	s := strings.Split(path, pathDelimiter)
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
