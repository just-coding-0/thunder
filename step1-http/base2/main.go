// Copyright 2020 just-codeding-0 . All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"net/http"
)

type Engine struct{}

// Http Handler interface
func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request){
	switch request.URL.Path {
	case "/":
		fmt.Fprintf(writer, "[PATH] %s", request.URL.Path)
	case "/hello":
		fmt.Fprint(writer, "hello")
	}

}


// basic HTTP server
func main() {

	var engine = new(Engine)

	http.ListenAndServe(":8080", engine)
}
