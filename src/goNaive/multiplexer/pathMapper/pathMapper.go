package pathMapper

import (
	"regexp"
	"strings"
	"fmt"
)

// -----------------------------------------------------------
// Note: this will be substituted with controller.Runnable
type SimpleRunnable interface {
	Run() bool
}

// -----------------------------------------------------------

type pathVariablePart struct {
	name  string
	child *pathMapperNode
}

type pathMapperNode struct {
	children map[string]*pathMapperNode
	variable pathVariablePart
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

func isVariablePart(part string) (variable string, ok bool) {
	re := regexp.MustCompile("^:(.+)$")
	match := re.FindStringSubmatch(part)
	if len(match) != 2 {
		return "", false
	} else {
		return match[1], true
	}
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
		// Is the current part a variable (eg ":abcd")
		if variableName, ok := isVariablePart(part); ok {
			if(data.variable.name != ""){
				return;
			}
			var newpathMapperNode = New()
			data.variable.name = variableName
			data.variable.child = newpathMapperNode
			data = data.variable.child
		} else {
			if _, ok := data.children[part]; !ok {
				var newpathMapperNode = New()
				data.children[part] = newpathMapperNode
			}
			data = data.children[part]
		}

		if i == len(parts)-1 {
			data.function = f
		}
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

func ident(val int) string{
	var i int
	s := ""
	for i = 0; i < val; i++{
		s = s + "  ";
	}
	return s;
}

func showNode(pm *pathMapperNode, level int) string{
	s := ""
	for k, v := range pm.children {
		s += ident(level) + k + "\n"
		s += showNode(v, level + 1)
	}
	if(pm.variable.name != ""){
		s += ident(level) + ":" + pm.variable.name + " (VAR)\n"
		s += showNode(pm.variable.child, level + 1)
	}
	return s;
}


func (pm *pathMapperNode) String() string {
        return fmt.Sprintf("%v", showNode(pm, 0))
}
