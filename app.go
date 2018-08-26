package main

import (
	"application/controllers"
	"fmt"
	"goNaive/multiplexer"
	"net/http"
)

func main() {
	// Create a naive multiplexer
	naiveMux := multiplexer.NewMultiplexer()

	controllers.RegisterControllers(naiveMux)

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: naiveMux,
	}

	fmt.Printf("Listening on port 8080\n")

	server.ListenAndServe()
}
