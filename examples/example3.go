package main

import (
	"fmt"

	"github.com/aristanetworks/goeapi"
)

func main() {
	node, err := goeapi.ConnectTo("dut")
	if err != nil {
		panic(err)
	}
	commands := []string{"show version"}
	conf, _ := node.Enable(commands)
	for k, v := range conf[0] {
		fmt.Println("k:", k, "v:", v)
	}
}
