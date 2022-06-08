package engine

import (
	"net/http"
)

type Routes interface {
	Run(address string) error
	Route(method, pattern string, handler ...HandlerFunc)
}

type Engine struct {
	trees    methodTrees
	handlers HandlersChain
}

func New() *Engine {
	return &Engine{
		trees: make(methodTrees, 0, 9),
	}
}

type HandlerFunc func(ctx *Context)

type HandlersChain []HandlerFunc

func (e *Engine) Run(address string) error {
	return http.ListenAndServe(address, e.Handler())
}

func (e *Engine) Route(method, path string, handlers ...HandlerFunc) {
	e.route(method, path, handlers)
}

// ServeHTTP conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := NewContext(w, req)
	engine.handleHTTPRequest(ctx)
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	method := c.Request.Method
	path := c.Request.URL.Path

	trees := engine.trees
	for i := 0; i < len(trees); i++ {
		if trees[i].method != method {
			continue
		}

		root := trees[i].root

		value := root.getNodeValue(path)
		if value.handlers != nil {
			c.handlers = value.handlers
			c.fullPath = value.fullPath
			c.Next()
			return
		}
	}
	serveError(c, http.StatusNotFound, default404Body)
}

func (engine *Engine) Handler() http.Handler {
	return engine
}

func (engine *Engine) route(httpMethod, path string, handlers HandlersChain) Routes {
	handlers = engine.combineHandlers(handlers)
	engine.addRoute(httpMethod, path, handlers)
	return engine
}

func (engine *Engine) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(engine.handlers) + len(handlers)
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, engine.handlers)
	copy(mergedHandlers[len(engine.handlers):], handlers)
	return mergedHandlers
}

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {

	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		engine.trees = append(engine.trees, &methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)
}
