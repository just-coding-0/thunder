package thunder

import (
	"net/http"
)

var (
	default404Body = []byte("404 page not found")
	default405Body = []byte("405 method not allowed")
)

type Engine struct {
	trees     methodTrees
	maxParams uint16
}

func New() *Engine {
	return &Engine{
	}
}

type IRoutes interface {
	Any(string, ...HandlerFunc) IRoutes
	GET(string, ...HandlerFunc) IRoutes
	POST(string, ...HandlerFunc) IRoutes
	DELETE(string, ...HandlerFunc) IRoutes
	PATCH(string, ...HandlerFunc) IRoutes
	PUT(string, ...HandlerFunc) IRoutes
	OPTIONS(string, ...HandlerFunc) IRoutes
	HEAD(string, ...HandlerFunc) IRoutes
}

// ServeHTTP conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c:=engine.allocateContext()
	c.writermem.reset(w)
	c.Request  = req
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

// POST is a shortcut for router.Handle("POST", path, handle).
func (engine *Engine) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine.handle(http.MethodPost, relativePath, handlers)
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (engine *Engine) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine.handle(http.MethodGet, relativePath, handlers)
}


// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (engine *Engine) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine.handle(http.MethodDelete, relativePath, handlers)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (engine *Engine) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine.handle(http.MethodPatch, relativePath, handlers)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (engine *Engine) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine.handle(http.MethodPut, relativePath, handlers)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (engine *Engine) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine.handle(http.MethodOptions, relativePath, handlers)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (engine *Engine) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
	return engine.handle(http.MethodHead, relativePath, handlers)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (engine *Engine) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
	engine.handle(http.MethodGet, relativePath, handlers)
	engine.handle(http.MethodPost, relativePath, handlers)
	engine.handle(http.MethodPut, relativePath, handlers)
	engine.handle(http.MethodPatch, relativePath, handlers)
	engine.handle(http.MethodHead, relativePath, handlers)
	engine.handle(http.MethodOptions, relativePath, handlers)
	engine.handle(http.MethodDelete, relativePath, handlers)
	engine.handle(http.MethodConnect, relativePath, handlers)
	engine.handle(http.MethodTrace, relativePath, handlers)
	return engine
}

func (engine *Engine) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	engine.addRoute(httpMethod, relativePath, handlers)
	return engine
}

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) IRoutes {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")
	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)

	// Update maxParams
	if paramsCount := countParams(path); paramsCount > engine.maxParams {
		engine.maxParams = paramsCount
	}

	return engine
}

func (engine *Engine) allocateContext() *Context {
	v := make(Params, 0, engine.maxParams)
	return &Context{engine: engine, params: &v}
}