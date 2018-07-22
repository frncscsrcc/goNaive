package controller

import (
	"net/http"
)

type Runnable interface {
	Run(http.ResponseWriter, *http.Request)
}
