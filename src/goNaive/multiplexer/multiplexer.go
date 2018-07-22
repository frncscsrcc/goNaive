package multiplexer

import (
	"fmt"
	"goNaive/controller"
	"net/http"
)

type multiplexer struct {
	pathsToController map[string]controller.Runnable
}

func NewMultiplexer() *multiplexer {
	var m multiplexer
	m.pathsToController = make(map[string]controller.Runnable)
	return &m
}

func (m *multiplexer) RegisterController(path string, r controller.Runnable) {
	m.pathsToController[path] = r
}

func (m *multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Multipexing %v\n", r.RequestURI)
	if controller, ok := m.pathsToController[r.RequestURI]; ok {
		controller.Run(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 - Not found")
		return
	}
}
