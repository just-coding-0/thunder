package main

import (
	"fmt"
	"github.com/just-coding-0/thunder/step3-group/thunder"
	"net/http"
)

func main() {

	t := thunder.New()

	t.GET("/test/:aid/*action", func(c *thunder.Context) {
		k, ok :=c.Params.Get("aid")
		fmt.Println(k, ok)
		k, ok =c.Params.Get("action")
		fmt.Println(k, ok)
		c.JSON(http.StatusOK, map[string]interface{}{
			"ok": true,
		})
	})

	g:=t.Group("/user/:aid/")

	g.GET("*action",func(c *thunder.Context) {
		k, ok :=c.Params.Get("aid")
		fmt.Println(k, ok)
		k, ok =c.Params.Get("action")
		fmt.Println(k, ok)
		c.JSON(http.StatusOK, map[string]interface{}{
			"ok": true,
		})
	})
	http.ListenAndServe(":8080", t)
}
