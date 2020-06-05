package thunder

import (
	"net/http"
	"strings"
)

type Engine struct {
	trees methodTrees
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
	f := responseWriter{
		ResponseWriter: w,
		size:           0,
		status:         0,
	}
	context := newContext(req, f)
	context.fullPath = req.URL.Path
	engine.handleHTTPRequest(context)
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	method := strings.ToUpper(c.Request.Method)
	path := c.Request.URL.Path
	tree := engine.trees.get(method)
	if tree == nil {
		c.JSON(http.StatusNotFound, "404 not found")
		return
	}

	/*
		1.当前节点是wildChild 那么必定只有一个子节点  类型是参数或者是catchAll
		/v1/:id/hello
		/v1/:id/hello/ 这是两个接口

	*/

walk:
	for {
		i := longestCommonPrefix(tree.path, path)

		// 已经找到节点
		if tree.path == path {
			c.handlers = tree.handlers
			c.Next()
			return
		}

		path = path[i:]

		if tree.wildChild {
			tree = tree.children[0]
			continue walk
		}

		//  /v1/:id
		//  /v1/:id/

		if tree.nType == param {
			idx := strings.Index(path, "/")
			if idx > 0 {
				c.Keys[tree.path[1:]] = path[:idx]
				path = path[idx:]
				if len(tree.children) == 0 {
					c.JSON(http.StatusNotFound, "404 not found")
					return
				}
				tree = tree.children[0]
				continue walk
			}

			c.Keys[tree.path[1:]] = path
			c.handlers = tree.handlers
			c.Next()
		}

		// 如果是catchAll
		if tree.nType == catchAll {
			c.Keys[tree.path[2:]] = path
			if tree.handlers == nil {
				c.JSON(http.StatusNotFound, "404 not found")
				return
			}

			c.handlers = tree.handlers
			c.Next()
			return
		}

		if tree.nType <= root {

			for idx := range tree.indices {
				if path[0] == tree.indices[idx] {
					tree = tree.children[idx]
					continue walk
				}
			}
		}
		c.JSON(http.StatusNotFound, "404 not found")
		return
	}

	c.JSON(http.StatusNotFound, "404 not found")

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

func (engine *Engine) Print() {

	for _, v := range engine.trees {
		println(v.root.path, v.root.fullPath)
		for _, v1 := range v.root.children {
			println(v1.path, v1.fullPath)

			for _, v2 := range v1.children {
				println(v2.path, v2.fullPath)
			}
		}

	}
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
	return engine
}
