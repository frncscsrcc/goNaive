package main

import (
	"fmt"
	"goNaive/multiplexer/pathMapper"
)

// --------------------------------------------------------------------
// This block creates a fake Runnable controller (note that for simplicity
// we are not passing http.ResponseWriter, *http.Request).
// The functionality of controllers is very simple: just print on STDOUT
// a message passed in the "constructor"
// --------------------------------------------------------------------
type exampleController struct {
	message string
}

func NewExampleController(message string) exampleController {
	var ec exampleController
	ec.message = message
	return ec
}

// Run returns a boolean, if is true, the framework will continue with
// the next controller in the chain. In this case we will return always
// true.
func (ec exampleController) Run() bool {
	fmt.Printf("%v\n", ec)
	return true
}

// --------------------------------------------------------------------

func main() {

	pathMapper := pathMapper.New("/")

	// Register some controllers
	pathMapper.Add("/admin/user/add", NewExampleController("I am the controller for /admin/user/add"))
	pathMapper.Add("/admin/user/delete", NewExampleController("I am the controller for /admin/user/delete"))
	pathMapper.Add("/admin/role/list", NewExampleController("I am the controller for /admin/role/list"))
	pathMapper.Add("/admin", NewExampleController("I am the controller for /admin"))
	pathMapper.Add("/powerUser/user/list", NewExampleController("I am the controller for /powerUser/user/list"))

	// Example of http request
	fmt.Printf("Requesting /admin/user/delete\n")
	fa := pathMapper.GetControllers("/admin/user/delete")
	for _, f := range fa {
		ok := f.Run()
		// If controller does not retur true, stop the execution and handle the error
		// message
		if !ok {
			// ...
			break
		}
	}

	// Example of http request
	fmt.Printf("Requesting /powerUser/user/list\n")
	fa = pathMapper.GetControllers("/powerUser/user/list")
	for _, f := range fa {
		ok := f.Run()
		// If controller does not retur true, stop the execution and handle the error
		// message
		if !ok {
			// ...
			break
		}
	}

}
