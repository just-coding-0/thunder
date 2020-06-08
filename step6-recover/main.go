package main

import (
	"github.com/just-coding-0/thunder/step6-recover/thunder"
	"net/http"
)

func main() {

	t := thunder.New()
	t.Use(thunder.Logger())
	t.Use(thunder.Recovery())


	g:=t.Group("/user/:aid/", func(content *thunder.Context) {
		panic(123123)
	})

	g.GET("*action",func(c *thunder.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"ok": true,
		})
	})
	http.ListenAndServe(":8080", t)
}
