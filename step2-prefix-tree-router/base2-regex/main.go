package main

import (
	"fmt"
	"github.com/just-coding-0/thunder/step2-prefix-tree-router/base2-regex/thunder"
	"net/http"
)

func main() {

	t := thunder.New()

	t.GET("/test/:aid/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, fmt.Sprintf("2"))
	})

	t.GET("/user/:aid/*action", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, fmt.Sprintf("1"))
	})

	http.ListenAndServe(":8080",t)
}
