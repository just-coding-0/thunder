package main

import (
	"github.com/just-coding-0/thunder/step5-logger/thunder"
	"net/http"
)

func main() {

	t := thunder.New()
	t.Use(thunder.Logger())


	g:=t.Group("/user/:aid/", func(content *thunder.Context) {
		println(12312312312312312)
	})

	g.GET("*action",func(c *thunder.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"ok": true,
		})
	})
	http.ListenAndServe(":8080", t)
}
