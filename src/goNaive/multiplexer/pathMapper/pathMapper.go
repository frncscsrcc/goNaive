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

type pathMapperTree struct {
	part     string
	children map[string]*pathMapperTree
	function SimpleRunnable
}

var ALLOWED_METHODS = map[string]bool{
	"POST":   true,
	"GET":    true,
	"DELETE": true,
	"PUT":    true,
}

func New(part string) *pathMapperTree {
	var pm pathMapperTree
	pm.part = part
	pm.children = make(map[string]*pathMapperTree)
	return &pm
}

func (pm *pathMapperTree) Register(method string, path string, f SimpleRunnable) bool {
	data := pm

	// Check if method is enabled
	if enabled, ok := ALLOWED_METHODS[method]; !ok || !enabled {
		return false
	}

	// add method (eg: GET, POST, ..., ALL) in front of path
	// so we will have something like
	// 	GET/admin/user/list
	// and we will map GET as first part of pathMapperTree tree
	path = method + path

	parts := strings.Split(path, "/")
	for i, part := range parts {
		// Skip first empty element of slice
		if len(part) == 0 {
			continue
		}
		if _, ok := data.children[part]; !ok {
			var newpathMapperTree = New(part)
			data.children[part] = newpathMapperTree
		}
		if i == len(parts)-1 {
			data.children[part].function = f
		}

		data = data.children[part]
	}
	return true
}

func (pm *pathMapperTree) RegisterGet(path string, f SimpleRunnable) {
	pm.Register("GET", path, f)
}

func (pm *pathMapperTree) RegisterPost(path string, f SimpleRunnable) {
	pm.Register("POST", path, f)
}

func (pm *pathMapperTree) RegisterPut(path string, f SimpleRunnable) {
	pm.Register("PUT", path, f)
}

func (pm *pathMapperTree) RegisterDelete(path string, f SimpleRunnable) {
	pm.Register("DELETE", path, f)
}

func (pm *pathMapperTree) RegisterAll(path string, f SimpleRunnable) {
	for method, enabled := range ALLOWED_METHODS {
		if enabled {
			pm.Register(method, path, f)
		}
	}
}

func (pm *pathMapperTree) GetControllers(path string) []SimpleRunnable {
	data := pm
	var functions []SimpleRunnable

	parts := strings.Split(path, "/")
	for _, part := range parts {
		// Skip first empty element of slice
		if len(part) == 0 {
			continue
		}
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
