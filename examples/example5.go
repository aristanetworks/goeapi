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

type showVlan struct {
	SourceDetail string
	Vlans        map[string]vlan
}

func (s *showVlan) GetCmd() string {
	return "show vlan"
}

type vlan struct {
	Status     string
	Name       string
	Interfaces map[string]devInterface
	Dynamic    bool
}

func (v *vlan) String() string {
	interfaces := make([]string, 0, len(v.Interfaces))
	for intf := range v.Interfaces {
		interfaces = append(interfaces, intf)
	}
	return fmt.Sprintf("Name:%s Status:%s Interfaces:%s", v.Name, v.Status, interfaces)
}

type devInterface struct {
	Annotation      string
	PrivatePromoted bool
}

type showHostName struct {
	Fqdn     string
	HostName string
}

func (s *showHostName) GetCmd() string {
	return "show hostname"
}

type showInterfacesStatus struct {
	InterfaceStatuses map[string]interfaceStatus
}

func (s *showInterfacesStatus) GetCmd() string {
	return "show interfaces status"
}

type interfaceStatus struct {
	Bandwidth           uint32
	InterfaceType       string
	Description         string
	AutoNegotiateActive bool
	Duplex              string
	LinkStatus          string
	LineProtocolStatus  string
	VlanInformation     vlanInfo
}

type vlanInfo struct {
	InterfaceMode            string
	VlanID                   int
	InterfaceForwardingModel string
}

func main() {
	node, err := goeapi.ConnectTo("dut")
	if err != nil {
		panic(err)
	}

	shVerRsp := &showVersionResp{}
	shHostnameRsp := &showHostName{}
	shVlanRsp := &showVlan{}
	shIntRsp := &showInterfacesStatus{}

	handle, _ := node.GetHandle(goeapi.Parameters{Format: "json"})
	handle.AddCommand(shVerRsp)
	handle.AddCommand(shHostnameRsp)
	handle.AddCommand(shVlanRsp)
	handle.AddCommand(shIntRsp)

	if err := handle.Call(); err != nil {
		panic(err)
	}
	fmt.Printf("Hostname          : %s\n", shHostnameRsp.HostName)
	fmt.Printf("Version           : %s\n", shVerRsp.Version)
	fmt.Printf("System MAC        : %s\n", shVerRsp.SystemMacAddress)
	fmt.Printf("Serial Number     : %s\n", shVerRsp.SerialNumber)

	fmt.Println("\nInterface        Description        Line     Link")
	fmt.Println("-------------------------------------------------")
	for intf, v := range shIntRsp.InterfaceStatuses {
		fmt.Printf("%-14s   %-16s   %-6s   %-12s\n",
			intf, v.Description, v.LineProtocolStatus, v.LinkStatus)
	}
	fmt.Println("")
	for vlanID, v := range shVlanRsp.Vlans {
		fmt.Printf("VlanID:%s %s\n", vlanID, v.String())
	}

}
