package thunder

import (
	"github.com/just-coding-0/thunder/step2-prefix-tree-router/base3-context/thunder/render"
	"net/http"
	"net/url"
)

type Context struct {
	writermem responseWriter
	Request   *http.Request
	Writer    ResponseWriter

	Params   Params
	handlers HandlersChain
	fullPath string

	engine *Engine

	queryCache url.Values

	Keys map[string]interface{}

	index    int8
}

func newContext(request *http.Request, res responseWriter) *Context {
	return &Context{
		writermem: res,
		Request:   request,
		Params:    Params{},
		Writer:    &res,
		Keys:      make(map[string]interface{},10),
	}
}

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (ps Params) Get(name string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return entry.Value, true
		}
	}
	return "", false
}

func (ps Params) ByName(name string) (va string) {
	va, _ = ps.Get(name)
	return
}

func (c *Context) Next() {
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Keys[key]
	return
}

func (c *Context) Query(key string) (value string, exists bool) {
	if c.queryCache == nil {
		c.queryCache = c.Request.URL.Query()
	}
	value = c.queryCache.Get(key)
	exists = len(value) > 0
	return
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) render(code int, r render.Render) {
	c.Status(code)
	if err := r.Render(c.Writer); err != nil {
		panic(err)
	}
}

func (c *Context) JSON(code int, obj interface{}) {
	c.render(code, render.JSON{Data: obj})
}