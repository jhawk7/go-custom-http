package main

import (
	"chs/internal/customhttp"
	"net"
)

const (
	NETWORK = "tcp"
	ADDR    = "0.0.0.0:8000"
)

func main() {
	l, lErr := net.Listen(NETWORK, ADDR)
	if lErr != nil {
		panic(lErr)
	}

	defer l.Close()
	routes := map[string]string{
		"GET":    "/test",
		"POST":   "/test",
		"DELETE": "/test",
	}

	for {
		conn, cErr := l.Accept()
		if cErr != nil {
			panic(cErr)
		}

		chttp := customhttp.InitHttp(conn, routes)
		go chttp.HandleRequest()
	}
}
