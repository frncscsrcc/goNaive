package pathMapper

import (
	"strings"
)

// -----------------------------------------------------------
// Note: this will be substituted with controller.Runnable
type simpleRunnable interface {
	Run() bool
}

// -----------------------------------------------------------

type pathMapper struct {
	part     string
	children map[string]*pathMapper
	function simpleRunnable
}

func New(part string) *pathMapper {
	var pm pathMapper
	pm.part = part
	pm.children = make(map[string]*pathMapper)
	return &pm
}

func (pm *pathMapper) Add(path string, f simpleRunnable) {
	data := pm

	parts := strings.Split(path, "/")
	for i, part := range parts {
		// Skip first empty element of slice
		if len(part) == 0 {
			continue
		}
		if _, ok := data.children[part]; !ok {
			var newPathMapper = New(part)
			data.children[part] = newPathMapper
		}
		if i == len(parts)-1 {
			data.children[part].function = f
		}

		data = data.children[part]
	}
}

func (pm *pathMapper) GetControllers(path string) []simpleRunnable {
	data := pm
	var functions []simpleRunnable

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
