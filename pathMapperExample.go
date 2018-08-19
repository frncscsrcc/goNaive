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
	fmt.Printf("\t%v\n", ec.message)
	return true
}

// --------------------------------------------------------------------

func main() {

	pm := pathMapper.New()

	// Register some controllers
	pm.RegisterGet("/admin/user/add", NewExampleController("I am the controller for /admin/user/add (GET method)"))
	pm.RegisterGet("/admin/user/get/:userID/all", NewExampleController("I am the controller for/admin/user/get/:userID/all (GET method)"))
	pm.RegisterGet("/admin/role/list", NewExampleController("I am the controller for /admin/role/list (GET method)"))
	pm.RegisterDelete("/admin/user/delete", NewExampleController("I am the controller for /admin/user/delete  (DELETE method)"))
	pm.RegisterGet("/admin/role/list", NewExampleController("I am the controller for /admin/role/list (GET method)"))
	pm.RegisterAll("/admin", NewExampleController("I am the controller for /admin (ALL methods)"))
	pm.RegisterPost("/powerUser/user/list", NewExampleController("I am the controller for /powerUser/user/list (POST method)"))

	var fa []pathMapper.SimpleRunnable

	// Simulate an http request
	fmt.Printf("Requesting GET/admin/user/get/test/all\n")
	fa = pm.GetControllers("GET/admin/user/get/test/all")
	for _, f := range fa {
		ok := f.Run()
		// If controller does not return true, stop the execution and handle the error
		// message
		if !ok {
			// ...
			break
		}
	}

	fmt.Printf("%v", pm)

}
