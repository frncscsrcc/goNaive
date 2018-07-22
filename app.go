package main

import (
	"application/controllers/controller1"
	"application/controllers/controller2"
	"fmt"
	"goNaive/multiplexer"
	"net/http"
)

func main() {
	// Create a naive multiplexer
	naiveMux := multiplexer.NewMultiplexer()

	c1 := controller1.New()
	naiveMux.RegisterController("/test1", c1)

	c2 := controller2.New()
	naiveMux.RegisterController("/test2", c2)

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: naiveMux,
	}

	fmt.Printf("Listening on port 8080\n")

	server.ListenAndServe()
}
