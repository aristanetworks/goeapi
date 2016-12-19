//
// Copyright (c) 2015-2016, Arista Networks, Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//   * Redistributions of source code must retain the above copyright notice,
//   this list of conditions and the following disclaimer.
//
//   * Redistributions in binary form must reproduce the above copyright
//   notice, this list of conditions and the following disclaimer in the
//   documentation and/or other materials provided with the distribution.
//
//   * Neither the name of Arista Networks nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL ARISTA NETWORKS
// BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR
// BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
// OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN
// IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package module

import (
	"github.com/aristanetworks/goeapi"
)

// ShowVersion defined data structure for mapping JSON response
// of 'show version' to manageable object
type ShowVersion struct {
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

// GetCmd returns the command type this EapiCommand relates to
func (s ShowVersion) GetCmd() string {
	return "show version"
}

// ShowInterface defined data structure for mapping JSON response
// of 'show interface' to manageable object
type ShowInterface struct {
	Interfaces map[string]SwitchInterface
}

// SwitchInterface defined data structure for mapping JSON response
// of 'show interface' to manageable object
type SwitchInterface struct {
	Bandwidth                 int
	BurnedInAddress           string
	Description               string //`json:"description"`
	ForwardingModel           string
	Hardware                  string
	InterfaceAddress          []InterfaceAddress
	InterfaceCounters         EthInterfaceCounters
	InterfaceMembership       string
	InterfaceStatistics       InterfaceStatistics
	InterfaceStatus           string
	L2Mtu                     int
	LastStatusChangeTimestamp float64
	LineProtocolStatus        string
	Mtu                       int
	Name                      string
	PhysicalAddress           string
}

// InterfaceAddress defined data structure for mapping JSON response
// of 'show interface' to manageable object
type InterfaceAddress struct {
	BroadcastAddress       string
	PrimaryIP              IPAddress
	SecondaryIPs           interface{}
	SecondaryIPOrderedList []IPAddress
	VirtualIP              IPAddress
}

// IPAddress defined data structure for mapping JSON response
// of 'show interface' to manageable object
type IPAddress struct {
	Address string
	MaskLen int
}

// InterfaceStatistics defined data structure for mapping JSON response
// of 'show interface' to manageable object
type InterfaceStatistics struct {
	InBitsRate     float64
	OutBitsRate    float64
	InPktsRate     float64
	UpdateInterval float64
	OutPktsRate    float64
}

// EthInterfaceCounters defined data structure for mapping JSON response
// of 'show interface' to manageable object
type EthInterfaceCounters struct {
	CounterRefreshTime float64
	InBroadcastPkts    int
	InDiscards         int
	InMulticastPkts    int
	InOctets           int
	InUcastPkts        int
	InputErrorsDetail  PhysicalInputErrors
	LastClear          float64
	LinkStatusChanges  int
	OutBroadcastPkts   int
	OutDiscards        int
	OutMulticastPkts   int
	OutOctets          int
	OutUcastPkts       int
	OutErrorsDetail    PhysicalOutputErrors
	TotalInErrors      int
	TotalOutErrors     int
}

// PhysicalInputErrors defined data structure for mapping JSON response
// of 'show interface' to manageable object
type PhysicalInputErrors struct {
	AlignmentErrots int
	FcsErrors       int
	GiantFrames     int
	RuntFrames      int
	RxPause         int
	SymbolErrors    int
}

// PhysicalOutputErrors defined data structure for mapping JSON response
// of 'show interface' to manageable object
type PhysicalOutputErrors struct {
	Collisions            int
	DeferredTransmissions int
	LateCollisions        int
	TxPause               int
}

// GetCmd returns the command type this EapiCommand relates to
func (s ShowInterface) GetCmd() string {
	return "show interfaces"
}

// ShowTrunkGroup defined data structure for mapping JSON response
// of 'show vlan trunk group' to manageable object
type ShowTrunkGroup struct {
	TrunkGroups map[string]struct {
		Names []string
	}
}

// GetCmd returns the command type this EapiCommand relates to
func (s ShowTrunkGroup) GetCmd() string {
	return "show vlan trunk group"
}

// ShowEntity provides a configuration resource for VLANs
type ShowEntity struct {
	*AbstractBaseEntity
}

// Show factory function to initiallize Show resource
// given a Node
func Show(node *goeapi.Node) *ShowEntity {
	return &ShowEntity{&AbstractBaseEntity{node}}
}

// ShowVersion returns the pre-defined structure
// (with "json" key in the struct field's tag value) for the
// decoded response from 'show version' command
func (s *ShowEntity) ShowVersion() ShowVersion {
	handle, _ := s.node.GetHandle("json")
	var showversion ShowVersion
	handle.AddCommand(&showversion)
	handle.Call()
	handle.Close()
	return showversion
}

// ShowInterfaces returns the pre-defined structure
// (with "json" key in the struct field's tag value) for the
// decoded response from 'show interfaces' command
func (s *ShowEntity) ShowInterfaces() ShowInterface {
	handle, _ := s.node.GetHandle("json")
	var showinterface ShowInterface
	handle.AddCommand(&showinterface)
	handle.Call()
	handle.Close()
	return showinterface
}

// ShowTrunkGroups returns the pre-defined structure
// (with "json" key in the struct field's tag value) for the
// decoded response from 'show vlan trunk group' command
func (s *ShowEntity) ShowTrunkGroups() ShowTrunkGroup {
	handle, _ := s.node.GetHandle("json")
	var showTrunkGroups ShowTrunkGroup
	handle.AddCommand(&showTrunkGroups)
	handle.Call()
	handle.Close()
	return showTrunkGroups
}
