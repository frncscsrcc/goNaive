package controllers

import (
	multiplexer "goNaive/multiplexer"

	"application/controllers/block"
	"application/controllers/controller1"
	"application/controllers/controller2"
)

func RegisterControllers(m *multiplexer.Multiplexer){
	m.RegisterController("ALL", "/block", block.New());
	m.RegisterController("GET", "/block/test1", controller1.New());
	m.RegisterController("GET", "/test2", controller2.New());
}