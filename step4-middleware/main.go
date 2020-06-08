package main

import (
	"fmt"
	"github.com/just-coding-0/thunder/step4-middleware/thunder"
	"net/http"
)

func main() {

	t := thunder.New()


	g:=t.Group("/user/:aid/", func(content *thunder.Context) {
		println(12312312312312312)
	})

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
