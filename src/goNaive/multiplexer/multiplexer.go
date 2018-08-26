package multiplexer

import (
	"fmt"
	"goNaive/controller"
	"goNaive/multiplexer/pathMapper"
	"net/http"
)

type Multiplexer struct {
	pathMapper *pathMapper.PathMapperNode
}

func NewMultiplexer() *Multiplexer {
	var m Multiplexer
	m.pathMapper = pathMapper.New()
	return &m
}

func (m *Multiplexer) RegisterController(method string, path string, r controller.Runnable) {
	m.pathMapper.Register(method, path, r)
}

func (m *Multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Multipexing %v\n", r.RequestURI)

	controllers := m.pathMapper.GetControllers(r.Method + r.RequestURI)
	if len(controllers) > 0 {
		for _, controller := range controllers {
			ok := controller.Run(w, r)
			// If controller does not return true, stop the execution
			if !ok {
				break
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 - Not found")
		return
	}
}
