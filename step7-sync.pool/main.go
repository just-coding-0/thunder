package main

import (
	"github.com/just-coding-0/thunder/step6-recover/thunder"
	"net/http"
	"runtime"
	"time"
)

var r runtime.MemStats

func main() {
	runtime.GOMAXPROCS(1)

	t := thunder.New()
	t.Use(thunder.Logger())
	t.Use(thunder.Recovery())

	g := t.Group("/user/:aid/")

	g.GET("*action", func(c *thunder.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"ok": true,
		})
	})
	go func() {
		t := time.NewTicker(time.Second * 2)
		for {
			select {
			case <-t.C:
				runtime.ReadMemStats(&r)
				println(r.HeapObjects, r.TotalAlloc, r.NumGC)
			}
		}
	}()

	http.ListenAndServe(":8080", t)
}
