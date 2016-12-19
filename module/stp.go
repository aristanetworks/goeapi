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
	"regexp"
	"strconv"

	"github.com/aristanetworks/goeapi"
)

// STPEntity provides a configuration resource for STP
type STPEntity struct {
	instances  *STPInstanceEntity
	interfaces *STPInterfaceEntity
	*AbstractBaseEntity
}

// Stp factory function to initiallize STPEntity resource
// given a Node
func Stp(node *goeapi.Node) *STPEntity {
	return &STPEntity{AbstractBaseEntity: &AbstractBaseEntity{node}}
}

// Get ...
func (s *STPEntity) Get() {
}

// Interfaces returns the STPInterfaces instance
func (s *STPEntity) Interfaces() *STPInterfaceEntity {
	if s.interfaces == nil {
		s.interfaces = STPInterfaces(s.AbstractBaseEntity.node)
	}
	return s.interfaces
}

// Instances returns the STPInstances resource for STP instances
func (s *STPEntity) Instances() *STPInstanceEntity {
	if s.instances == nil {
		s.instances = STPInstance(s.AbstractBaseEntity.node)
	}
	return s.instances
}

// SetMode sets the spanning tree mode. Valid parameters are:
// 	mstp or none.
// Returns:
//  bool: Returns True if the commands complete successfully
func (s *STPEntity) SetMode(value string) bool {
	if value == "" {
		return s.Configure("no spanning-tree mode")
	}
	if value != "mstp" && value != "none" {
		return false
	}
	return s.Configure("spanning-tree mode " + value)
}

// STPInstanceEntity provides a configuration resource for STPInstance
type STPInstanceEntity struct {
	*AbstractBaseEntity
}

// STPInstance factory function to initiallize STPInstanceEntity resource
// given a Node
func STPInstance(node *goeapi.Node) *STPInstanceEntity {
	return &STPInstanceEntity{&AbstractBaseEntity{node}}
}

//func (s *STPInstanceEntity) GetAll() {
//}

// STPInterfaceConfig represents the parsed STP interface config
// {
//		"bpduguard"     : "false",
//		"portfast"      : "true",
//		"portfast_type" : "edge",
// }
type STPInterfaceConfig map[string]string

// STPInterfaceCollection is a collection of STPInterfaceConfigs
// mapped by interface name:
// {
//		"Port-Channel1" : STPInterfaceConfig {
//							"bpduguard"     : "false",
//							"portfast"      : "true",
//							"portfast_type" : "edge",
//						  },
//		   "Ethernet1" : ...
// }
type STPInterfaceCollection map[string]STPInterfaceConfig

// STPInterfaceEntity provides a configuration resource for STP
type STPInterfaceEntity struct {
	*AbstractBaseEntity
}

// STPInterfaces factory function to initiallize STPInterfaceEntity resource
// given a Node
func STPInterfaces(node *goeapi.Node) *STPInterfaceEntity {
	return &STPInterfaceEntity{&AbstractBaseEntity{node}}
}

// Get returns an STPInterfaceConfig type for a given interface name(string).
func (s *STPInterfaceEntity) Get(name string) STPInterfaceConfig {
	parent := `interface\s+` + name
	config, _ := s.GetBlock(parent)

	return STPInterfaceConfig{
		"bpduguard":     strconv.FormatBool(s.parseBPDUGuard(config)),
		"portfast":      strconv.FormatBool(s.parsePortfast(config)),
		"portfast_type": s.parsePortfastType(config),
	}
}

// GetAll returns a collection of STPInterfaceConfigs key'd by
// interface name.
func (s *STPInterfaceEntity) GetAll() STPInterfaceCollection {
	config := s.Config()

	re := regexp.MustCompile(`(?m)^interface\s(Eth.+|Po.+)$`)
	interfaces := re.FindAllStringSubmatch(config, -1)

	collection := make(STPInterfaceCollection)

	for _, matchedLine := range interfaces {
		intf := matchedLine[1]
		if tmp := s.Get(intf); tmp != nil {
			collection[intf] = tmp
		}
	}
	return collection
}

// parseBPDUGuard parses the provided interface config returning
// true(bool) if bpdu guard is enabled
func (s *STPInterfaceEntity) parseBPDUGuard(config string) bool {
	matched, _ := regexp.MatchString("spanning-tree bpduguard enable", config)
	return matched
}

// parsePortfast parses the provided interface config returning
// true(bool) if spanning-tree portfast is enabled
func (s *STPInterfaceEntity) parsePortfast(config string) bool {
	matched, _ := regexp.MatchString("no spanning-tree portfast", config)
	return !matched
}

// parsePortfastType parses the provided interface config returning the
// portfast type (network, edge, or normal)
func (s *STPInterfaceEntity) parsePortfastType(config string) string {
	found, _ := regexp.MatchString("spanning-tree portfast network", config)
	if found {
		return "network"
	}
	found, _ = regexp.MatchString("no spanning-tree portfast", config)
	if found {
		return "normal"
	}
	return "edge"
}

// ConfigureInterface (redefined from Base)
// Returns true(bool) if configuration successful
func (s *STPInterfaceEntity) ConfigureInterface(name string, cmds ...string) bool {
	if !isValidStpInterface(name) {
		return false
	}
	return s.AbstractBaseEntity.ConfigureInterface(name, cmds...)
}

// SetPortfastType sets the spanning-tree portfast type for the interface name(string) to
// one of the following valid args:
//	network
// 	edge
//	normal
// Returns true(bool) if configuration successful
func (s *STPInterfaceEntity) SetPortfastType(name string, value string) bool {
	validTypes := map[string]bool{
		"network": true,
		"edge":    true,
		"normal":  true,
	}
	if _, found := validTypes[value]; !found {
		return false
	}

	cmds := []string{"spanning-tree portfast " + value}
	if value == "edge" {
		cmds = append(cmds, "spanning-tree portfast auto")
	}
	return s.ConfigureInterface(name, cmds...)
}

// SetPortfast sets the spanning-tree portfast for the interface name(string) to
// be enabled(true) or disabled(false). Returns true(bool) if configuration successful
func (s *STPInterfaceEntity) SetPortfast(name string, enable bool) bool {
	cmds := s.CommandBuilder("spanning-tree portfast", "", false, enable)
	return s.ConfigureInterface(name, cmds)
}

// SetPortfastDefault sets the spanning-tree portfast for the interface name(string)
// back to the default config. Returns true(bool) if configuration successful
func (s *STPInterfaceEntity) SetPortfastDefault(name string) bool {
	cmd := s.CommandBuilder("spanning-tree portfast", "", true, false)
	return s.ConfigureInterface(name, cmd)
}

// SetBPDUGuard eables(true) or disables(false) spanning-tree bpduguard for the
// interface name(string). Returns true(bool) if configuration successful
func (s *STPInterfaceEntity) SetBPDUGuard(name string, enable bool) bool {
	param := "disable"
	if enable {
		param = "enable"
	}
	cmd := s.CommandBuilder("spanning-tree bpduguard", param, false, true)
	return s.ConfigureInterface(name, cmd)
}

// SetBPDUGuardDefault sets the spanning-tree bpduguard for the interface name(string)
// back to the default config. Returns true(bool) if configuration successful
func (s *STPInterfaceEntity) SetBPDUGuardDefault(name string) bool {
	cmd := s.CommandBuilder("spanning-tree bpduguard", "", true, false)
	return s.ConfigureInterface(name, cmd)
}

// isValidStpInterface
func isValidStpInterface(value string) bool {
	valid, _ := regexp.MatchString(`^Eth.+|^Po.+`, value)
	return valid
}
