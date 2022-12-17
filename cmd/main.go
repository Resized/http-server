package main

import (
	http_server "github.com/Resized/http-server/pkg/http-server"
	"net/http"
)

func main() {
	http_server.Run()
	get, err := http.Get("http://localhost:8090/add?data=hello,world")
	if err != nil {
		return
	}
	println(get.StatusCode)
}
