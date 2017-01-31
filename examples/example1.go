package main

import (
	"fmt"

	"github.com/aristanetworks/goeapi"
	"github.com/aristanetworks/goeapi/module"
)

func main() {
	node, err := goeapi.ConnectTo("dut")
	if err != nil {
		panic(err)
	}
	conf := node.RunningConfig()
	fmt.Println(conf)

	var showversion module.ShowVersion
	handle, _ := node.GetHandle("json")
	if err := handle.Enable(&showversion); err != nil {
		panic(err)
	}

	fmt.Println("\nVersion:", showversion.Version)

	s := module.Show(node)
	showData := s.ShowVersion()
	fmt.Printf("\nModelname         : %s\n", showData.ModelName)
	fmt.Printf("Internal Version  : %s\n", showData.InternalVersion)
	fmt.Printf("System MAC        : %s\n", showData.SystemMacAddress)
	fmt.Printf("Serial Number     : %s\n", showData.SerialNumber)
	fmt.Printf("Mem Total         : %d\n", showData.MemTotal)
	fmt.Printf("Bootup Timestamp  : %.2f\n", showData.BootupTimestamp)
	fmt.Printf("Mem Free          : %d\n", showData.MemFree)
	fmt.Printf("Version           : %s\n", showData.Version)
	fmt.Printf("Architecture      : %s\n", showData.Architecture)
	fmt.Printf("Internal Build ID : %s\n", showData.InternalBuildID)
	fmt.Printf("Hardware Revision : %s\n", showData.HardwareRevision)

	sys := module.System(node)
	if ok := sys.SetHostname("Ladie"); !ok {
		fmt.Printf("SetHostname Failed\n")
	}
	sysInfo := sys.Get()
	fmt.Printf("\nSysinfo: %#v\n", sysInfo.HostName())
}
