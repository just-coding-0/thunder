package main

import (
	"net/http"
	"time"
)

func main() {
	var c = make(chan int, 10)
	for j := 0; j < 100; j++ {
		for i := 0; i < 1000; i++ {
			c <- i
			go func() {
				res, _ := http.Get("http://127.0.0.1:8080/user/ted/get")
				if res != nil {
					res.Body.Close()
				}
				<-c
			}()
		}
		time.Sleep(time.Second)
	}
}
