// Copyright 2020 just-codeding-0 . All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package thunder

import (
	"fmt"
	"net/http"
)

type Engine struct {
	router map[string]HandlersChain
	trees  methodTrees
}

func New() *Engine {
	return &Engine{
		router: make(map[string]HandlersChain),
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
	engine.handleHTTPRequest(w, req)
}

func (engine *Engine) handleHTTPRequest(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path
	tree := engine.trees.get(method)
	if tree == nil {
		fmt.Fprintf(w, " handler not found ")
		return
	}
walk:
	for {

		i := longestCommonPrefix(path, tree.path)
		path = path[i:]

		// 如果路径相等,执行
		if len(path) == 0 && tree.fullPath == req.URL.Path {
			for _,v:=range tree.handlers{
				v(w,req)
			}
			return
		}

		// 寻找子节点
		for idx := range tree.indices {
			if tree.indices[idx] == path[0] {
				tree = tree.children[idx]
				continue walk
			}
		}

		break
	}

	fmt.Fprintf(w, " handler not found ")
}

type HandlerFunc func(writer http.ResponseWriter, request *http.Request)

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
	return engine
}
