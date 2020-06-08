package thunder

import (
	"net/http"
	"path"
)

type RouterGroup struct {
	Handlers HandlersChain
	basePath string
	engine   *Engine
	root     bool
}

type IRouter interface {
	IRoutes
	Group(string, ...HandlerFunc) *RouterGroup
}

type IRoutes interface {
	Use(...HandlerFunc) IRoutes
	Any(string, ...HandlerFunc) IRoutes
	GET(string, ...HandlerFunc) IRoutes
	POST(string, ...HandlerFunc) IRoutes
	DELETE(string, ...HandlerFunc) IRoutes
	PATCH(string, ...HandlerFunc) IRoutes
	PUT(string, ...HandlerFunc) IRoutes
	OPTIONS(string, ...HandlerFunc) IRoutes
	HEAD(string, ...HandlerFunc) IRoutes
}

func (Group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: Group.combineHandlers(handlers),
		basePath: joinBaesPath(Group.basePath, relativePath),
		engine:   Group.engine,
		root:     false,
	}
}

func (Group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(Group.Handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, Group.Handlers)
	copy(mergedHandlers[len(Group.Handlers):], handlers)
	return mergedHandlers
}

func joinBaesPath(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}

	return finalPath
}

func (Group *RouterGroup) Use(handlers ...HandlerFunc) IRoutes {
	Group.Handlers = append(Group.Handlers, handlers...)
	return Group.returnObj()
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (Group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	return Group.handle(http.MethodPost, relativePath, handlers)
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (Group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return Group.handle(http.MethodGet, relativePath, handlers)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (Group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	return Group.handle(http.MethodDelete, relativePath, handlers)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (Group *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
	return Group.handle(http.MethodPatch, relativePath, handlers)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (Group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	return Group.handle(http.MethodPut, relativePath, handlers)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (Group *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes {
	return Group.handle(http.MethodOptions, relativePath, handlers)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (Group *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
	return Group.handle(http.MethodHead, relativePath, handlers)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (Group *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
	Group.handle(http.MethodGet, relativePath, handlers)
	Group.handle(http.MethodPost, relativePath, handlers)
	Group.handle(http.MethodPut, relativePath, handlers)
	Group.handle(http.MethodPatch, relativePath, handlers)
	Group.handle(http.MethodHead, relativePath, handlers)
	Group.handle(http.MethodOptions, relativePath, handlers)
	Group.handle(http.MethodDelete, relativePath, handlers)
	Group.handle(http.MethodConnect, relativePath, handlers)
	Group.handle(http.MethodTrace, relativePath, handlers)
	return Group.returnObj()
}

func (Group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := joinBaesPath(Group.basePath, relativePath)
	_handlers := Group.combineHandlers(handlers)
	Group.addRoute(httpMethod, absolutePath, _handlers)
	return Group.returnObj()
}

func (Group *RouterGroup) addRoute(method, path string, handlers HandlersChain) IRoutes {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")
	root := Group.engine.trees.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		Group.engine.trees = append(Group.engine.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)

	// Update maxParams
	if paramsCount := countParams(path); paramsCount > Group.engine.maxParams {
		Group.engine.maxParams = paramsCount
	}

	return Group.engine
}

func (Group *RouterGroup) returnObj() IRoutes {
	if Group.root {
		return Group.engine
	}
	return Group
}
