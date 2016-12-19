package main

import (
	"fmt"

	"github.com/aristanetworks/goeapi"
)

// MyShowVlan ...
type MyShowVlan struct {
	SourceDetail string
	Vlans        map[string]Vlan
}

// Vlan ...
type Vlan struct {
	Status     string
	Name       string
	Interfaces map[string]Interface
	Dynamic    bool
}

// Interface ...
type Interface struct {
	Annotation      string
	PrivatePromoted bool
}

// GetCmd ...
func (s *MyShowVlan) GetCmd() string {
	return "show vlan configured-ports"
}

func main() {
	node, err := goeapi.ConnectTo("dut")
	if err != nil {
		panic(err)
	}

	sv := &MyShowVlan{}

	handle, _ := node.GetHandle("json")
	handle.AddCommand(sv)
	if err := handle.Call(); err != nil {
		panic(err)
	}

	for k, v := range sv.Vlans {
		fmt.Printf("Vlan:%s\n", k)
		fmt.Printf("  Name  : %s\n", v.Name)
		fmt.Printf("  Status: %s\n", v.Status)
	}
}
