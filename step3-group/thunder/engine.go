package thunder

import (
	"net/http"
)

var (
	default404Body = []byte("404 page not found")
	default405Body = []byte("405 method not allowed")
)

type Engine struct {
	RouterGroup
	trees     methodTrees
	maxParams uint16
}

func New() *Engine {
	e := &Engine{}
	e.RouterGroup = RouterGroup{
		Handlers: nil,
		basePath: "/",
		engine:   e,
		root:     true,
	}

	return e
}

// ServeHTTP conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.allocateContext()
	c.writermem.reset(w)
	c.Request = req
	c.reset()

	engine.handleHTTPRequest(c)
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path

	// 寻找指定的根
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		// Find route in tree
		value := root.getValue(rPath, c.params, false)
		if value.params != nil {
			c.Params = *value.params
		}
		if value.handlers != nil {
			c.handlers = value.handlers
			c.fullPath = value.fullPath
			c.Next()
			c.writermem.WriteHeaderNow()
			return
		}

		break
	}
	serveError(c, http.StatusNotFound, default404Body)
}

func serveError(c *Context, code int, defaultMessage []byte) {
	c.writermem.status = code
	c.Next()
	if c.writermem.Written() {
		return
	}
	if c.writermem.Status() == code {
		c.writermem.Header()["Content-Type"] = []string{"text/plain"}
		c.Writer.Write(defaultMessage)
		return
	}
	c.writermem.WriteHeaderNow()
}

type HandlerFunc func(content *Context)

type HandlersChain []HandlerFunc

func (engine *Engine) allocateContext() *Context {
	v := make(Params, 0, engine.maxParams)
	return &Context{engine: engine, params: &v}
}
