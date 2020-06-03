// Copyright 2020 just-codeding-0 . All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package thunder

import (
	"fmt"
	"net/http"
	"regexp"
)

type Engine struct {
	router map[string]HandlersChain
}

func New() *Engine{
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

// Http Handler interface
func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	key := fmt.Sprintf("%s_%s", request.Method, request.URL.Path)
	if handler, ok := engine.router[key]; ok {
		for _,val:=range handler{
			val(writer,request)
		}
	} else {
		fmt.Fprintf(writer, "404 nout found")
	}
}

type HandlerFunc func(writer http.ResponseWriter, request *http.Request)

type HandlersChain []HandlerFunc

func (engine *Engine) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	if matches, err := regexp.MatchString("^[A-Z]+$", httpMethod); !matches || err != nil {
		panic("http method " + httpMethod + " is not valid")
	}
	return engine.handle(httpMethod, relativePath, handlers)
}

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
	engine.addRoute(httpMethod,relativePath,handlers)
	return engine
}

func   (engine *Engine) addRoute(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	assert1(relativePath[0] == '/', "path must begin with '/'")
	assert1(httpMethod != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")
	k := fmt.Sprintf("%s_%s", httpMethod, relativePath)
	engine.router[k] = handlers
	return engine
}

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}
