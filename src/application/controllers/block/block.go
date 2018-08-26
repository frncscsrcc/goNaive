package block

import (
	"net/http"
)

type controller struct{}

func New() *controller {
	var c controller
	return &c
}

func (c *controller) Run(w http.ResponseWriter, r *http.Request) bool {
	http.Error(w, "Not authenticated", 403)
	return false
}
