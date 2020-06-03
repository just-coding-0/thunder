// Copyright 2020 just-codeding-0 . All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"github.com/just-coding-0/thunder/step1-http/base3/thunder"
	"net/http"
)

func main() {
	e := thunder.New()
	e.GET("/helloworld", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "hello word")
	})
	http.ListenAndServe(":8080", e)
}
