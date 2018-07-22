package controller1

import (
	"io"
	"net/http"
)

type exampleController struct{}

func New() *exampleController {
	var c exampleController
	return &c
}

func (c *exampleController) Run(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is controller1!\n")
}
