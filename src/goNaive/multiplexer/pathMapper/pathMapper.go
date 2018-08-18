package pathMapper

import (
	"strings"
)

// -----------------------------------------------------------
// Note: this will be substituted with controller.Runnable
type SimpleRunnable interface {
	Run() bool
}

// -----------------------------------------------------------

type pathMapperNode struct {
	children map[string]*pathMapperNode
	function SimpleRunnable
}

var ALLOWED_METHODS = map[string]bool{
	"POST":   true,
	"GET":    true,
	"DELETE": true,
	"PUT":    true,
}

func New() *pathMapperNode {
	var pm pathMapperNode
	pm.children = make(map[string]*pathMapperNode)
	return &pm
}

func (pm *pathMapperNode) Register(method string, path string, f SimpleRunnable) {
	data := pm

	// Check if method is enabled
	if enabled, ok := ALLOWED_METHODS[method]; !ok || !enabled {
		return
	}

	// add method (eg: GET, POST, ..., ALL) in front of path
	// so we will have something like
	// 	GET/admin/user/list
	// and we will map GET as first part of pathMapperNode tree
	path = method + path

	parts := strings.Split(path, "/")
	for i, part := range parts {
		if _, ok := data.children[part]; !ok {
			var newpathMapperNode = New()
			data.children[part] = newpathMapperNode
		}
		if i == len(parts)-1 {
			data.children[part].function = f
		}

		data = data.children[part]
	}
}

func (pm *pathMapperNode) RegisterGet(path string, f SimpleRunnable) {
	pm.Register("GET", path, f)
}

func (pm *pathMapperNode) RegisterPost(path string, f SimpleRunnable) {
	pm.Register("POST", path, f)
}

func (pm *pathMapperNode) RegisterPut(path string, f SimpleRunnable) {
	pm.Register("PUT", path, f)
}

func (pm *pathMapperNode) RegisterDelete(path string, f SimpleRunnable) {
	pm.Register("DELETE", path, f)
}

func (pm *pathMapperNode) RegisterAll(path string, f SimpleRunnable) {
	for method, enabled := range ALLOWED_METHODS {
		if enabled {
			pm.Register(method, path, f)
		}
	}
}

func (pm *pathMapperNode) GetControllers(path string) []SimpleRunnable {
	data := pm
	var functions []SimpleRunnable

	parts := strings.Split(path, "/")
	for _, part := range parts {
		if _, ok := data.children[part]; !ok {
			return nil
		}

		if f := data.children[part].function; f != nil {
			functions = append(functions, f)
		}
		data = data.children[part]
	}

	return functions
}
