package main

import (
	"fmt"
	"github.com/just-coding-0/thunder/step2-prefix-tree-router/base4-optimize-tree-serach/thunder"
	"net/http"
)

func main() {

	t := thunder.New()

	t.GET("/test/:aid/", func(c *thunder.Context) {
		k, ok :=c.Params.Get("aid")
		fmt.Println(k, ok)
		c.JSON(http.StatusOK, map[string]interface{}{
			"ok": true,
		})
	})

	t.GET("/user/:aid/*action", func(c *thunder.Context) {
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
