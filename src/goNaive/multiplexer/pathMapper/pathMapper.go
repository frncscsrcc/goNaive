package pathMapper

import (
	"fmt"
	"goNaive/controller"
	"regexp"
	"strings"
)

type pathVariablePart struct {
	name  string
	child *PathMapperNode
}

type PathMapperNode struct {
	children map[string]*PathMapperNode
	variable pathVariablePart
	function controller.Runnable
}

var ALLOWED_METHODS = map[string]bool{
	"POST":   true,
	"GET":    true,
	"DELETE": true,
	"PUT":    true,
}

func New() *PathMapperNode {
	var pm PathMapperNode
	pm.children = make(map[string]*PathMapperNode)
	return &pm
}

func isVariablePart(part string) (variable string, ok bool) {
	re := regexp.MustCompile("^:(.+)$")
	match := re.FindStringSubmatch(part)
	if len(match) != 2 {
		return "", false
	} else {
		return match[1], true
	}
}

func (pm *PathMapperNode) Register(method string, path string, f controller.Runnable) {
	data := pm

	if method == "ALL"{
		for method, enabled := range ALLOWED_METHODS {
			if enabled {
				pm.Register(method, path, f)
			}
		}
	}

	// Check if method is enabled
	if enabled, ok := ALLOWED_METHODS[method]; !ok || !enabled {
		return
	}

	// add method (eg: GET, POST, ..., ALL) in front of path
	// so we will have something like
	// 	GET/admin/user/list
	// and we will map GET as first part of PathMapperNode tree
	path = method + path
	parts := strings.Split(path, "/")
	for i, part := range parts {
		// Is the current part a variable (eg ":abcd")
		if variableName, ok := isVariablePart(part); ok {
			if data.variable.name != "" {
				return
			}
			var newPathMapperNode = New()
			data.variable.name = variableName
			data.variable.child = newPathMapperNode
			data = data.variable.child
		} else {
			if _, ok := data.children[part]; !ok {
				var newPathMapperNode = New()
				data.children[part] = newPathMapperNode
			}
			data = data.children[part]
		}

		if i == len(parts)-1 {
			data.function = f
		}
	}
}

func (pm *PathMapperNode) RegisterGet(path string, f controller.Runnable) {
	pm.Register("GET", path, f)
}

func (pm *PathMapperNode) RegisterPost(path string, f controller.Runnable) {
	pm.Register("POST", path, f)
}

func (pm *PathMapperNode) RegisterPut(path string, f controller.Runnable) {
	pm.Register("PUT", path, f)
}

func (pm *PathMapperNode) RegisterDelete(path string, f controller.Runnable) {
	pm.Register("DELETE", path, f)
}

func (pm *PathMapperNode) RegisterAll(path string, f controller.Runnable) {
	for method, enabled := range ALLOWED_METHODS {
		if enabled {
			pm.Register(method, path, f)
		}
	}
}

func (pm *PathMapperNode) GetControllers(path string) []controller.Runnable {
	data := pm
	var functions []controller.Runnable

	parts := strings.Split(path, "/")
	for _, part := range parts {
		if _, ok := data.children[part]; ok {
			data = data.children[part]
			if f := data.function; f != nil {
				functions = append(functions, f)
			}
		} else if data.variable.name != "" {
			data = data.variable.child
			if f := data.function; f != nil {
				functions = append(functions, f)
			}
		} else {
			return nil
		}
	}

	return functions
}

func ident(val int) string {
	var i int
	s := ""
	for i = 0; i < val; i++ {
		s = s + "  "
	}
	return s
}

func showNode(pm *PathMapperNode, level int) string {
	s := ""
	for k, v := range pm.children {
		s += ident(level) + k + "\n"
		s += showNode(v, level+1)
	}
	if pm.variable.name != "" {
		s += ident(level) + ":" + pm.variable.name + " (VAR)\n"
		s += showNode(pm.variable.child, level+1)
	}
	return s
}

func (pm *PathMapperNode) String() string {
	return fmt.Sprintf("%v", showNode(pm, 0))
}
