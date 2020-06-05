package main

import (
	"fmt"
	"net/http"
)

// basic HTTP server
func main() {

	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "hello")
	})

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "[PATH] %s", request.URL.Path)
	})

	http.ListenAndServe(":8080", nil)
}
