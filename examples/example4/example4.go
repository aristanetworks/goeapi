package main

import (
	"fmt"

	"github.com/aristanetworks/goeapi"
)

type showVersionResp struct {
	ModelName        string
	InternalVersion  string
	SystemMacAddress string
	SerialNumber     string
	MemTotal         int
	BootupTimestamp  float64
	MemFree          int
	Version          string
	Architecture     string
	InternalBuildID  string
	HardwareRevision string
}

func (s *showVersionResp) GetCmd() string {
	return "show version"
}

func main() {
	node, err := goeapi.ConnectTo("dut")
	if err != nil {
		panic(err)
	}

	svRsp := &showVersionResp{}

	handle, _ := node.GetHandle("json")
	handle.AddCommand(svRsp)
	if err := handle.Call(); err != nil {
		panic(err)
	}
	fmt.Printf("Version           : %s\n", svRsp.Version)
	fmt.Printf("System MAC        : %s\n", svRsp.SystemMacAddress)
	fmt.Printf("Serial Number     : %s\n", svRsp.SerialNumber)
}
