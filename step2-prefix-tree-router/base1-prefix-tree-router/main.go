package main

import (
	"fmt"
	"github.com/just-coding-0/thunder/step2-prefix-tree-router/base1-prefix-tree-router/thunder"
	"net/http"
)

func main() {

	t := thunder.New()

	t.GET("/test/12", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, fmt.Sprintf("2"))
	})

	t.GET("/test/123", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, fmt.Sprintf("1"))
	})

	http.ListenAndServe(":8080",t)
}
