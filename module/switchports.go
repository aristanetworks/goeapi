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
	"strings"

	"github.com/aristanetworks/goeapi"
)

// SwitchPortConfig represents a parsed SwitchPort entry
type SwitchPortConfig map[string]string

// SwitchPortConfigMap represents a parsed SwitchPort entry
type SwitchPortConfigMap map[string]SwitchPortConfig

// SwitchPortEntity provides a configuration resource for SwitchPorts
//
// Logical layer 2 interfaces built on top of physical Ethernet and bundled
// Port-Channel interfaces can be configured and managed with an instance
// of Switchport.   The Switchport object is a resource collection and
// supports get and getall methods.  The Switchports class is derived from
// the AbstractBaseEntity object
type SwitchPortEntity struct {
	*AbstractBaseEntity
}

// Name returns the SwitchPortConfig name(string)
func (s SwitchPortConfig) Name() string {
	return s["name"]
}

// Mode returns the SwitchPortConfig mode(string)
func (s SwitchPortConfig) Mode() string {
	return s["mode"]
}

// AccessVlan returns the SwitchPortConfig access vlan(string)
func (s SwitchPortConfig) AccessVlan() string {
	return s["access_vlan"]
}

// TrunkNativeVlan returns the SwitchPortConfig native vlan(string)
func (s SwitchPortConfig) TrunkNativeVlan() string {
	return s["trunk_native_vlan"]
}

// TrunkAllowedVlans returns the SwitchPortConfig allowed trunk vlans(string)
func (s SwitchPortConfig) TrunkAllowedVlans() string {
	return s["trunk_allowed_vlans"]
}

// TrunkGroups returns the SwitchPortConfig trunk groups(string)
// comma delimited string
func (s SwitchPortConfig) TrunkGroups() string {
	return s["trunk_groups"]
}

// SwitchPort factory function to initiallize SwitchPortEntity resource
// given a Node
func SwitchPort(node *goeapi.Node) *SwitchPortEntity {
	return &SwitchPortEntity{&AbstractBaseEntity{node}}
}

// Get Returns a SwitchPortConfig object that represents a switchport
// The Switchport resource returns the following:
//    * name (string): The name of the interface
//    * mode (string): The switchport mode value
//    * access_vlan (string): The switchport access vlan value
//    * trunk_native_vlan (string): The switchport trunk native vlan vlaue
//    * trunk_allowed_vlans (string): The trunk allowed vlans value
//    * trunk_groups (string): The list of trunk groups configured
//
// Args:
//    name (string): The interface identifier to get.  Note: Switchports
//        are only supported on Ethernet and Port-Channel interfaces
//
// Returns:
//    SwitchPortConfig: An object of key/value pairs that represent
//        the switchport configuration for the interface specified  If
//        the specified argument is not a switchport then None
//        is returned
func (s *SwitchPortEntity) Get(name string) SwitchPortConfig {
	parent := "interface " + name
	config, err := s.GetBlock(parent)
	if err != nil {
		return nil
	}
	if matched, _ := regexp.MatchString(`no switchport\s*\n`, config); matched {
		return nil
	}
	return SwitchPortConfig{
		"name":                name,
		"mode":                s.parseMode(config),
		"access_vlan":         s.parseAccessVlan(config),
		"trunk_native_vlan":   s.parseTrunkNativeVlan(config),
		"trunk_allowed_vlans": s.parseTrunkAllowedVlans(config),
		"trunk_groups":        s.parseTrunkGroups(config),
	}
}

// GetAll Returns a mapped object to all Switchports
// This method will return all of the configured switchports as a
// SwitchPortConfigMap object keyed by the interface identifier.
//
// Returns:
//    A map'd SwitchPort that represents all configured
//        switchports in the current running configuration
func (s *SwitchPortEntity) GetAll() SwitchPortConfigMap {
	config := s.Config()

	re := regexp.MustCompile(`(?m)^interface\s([Et|Po].+)$`)
	matches := re.FindAllStringSubmatch(config, -1)

	response := make(SwitchPortConfigMap)

	for _, line := range matches {
		name := line[1]
		intf := s.Get(name)
		if intf != nil {
			response[name] = intf
		}
	}
	return response
}

// GetSection returns the specified SwitchPort Entry for the name specified.
//
// Args:
//  name (string): The port name
//
// Returns:
//  Returns string representation of SwitchPort config entry
func (s *SwitchPortEntity) GetSection(name string) string {
	parent := "interface " + name
	config, err := s.GetBlock(parent)
	if err != nil {
		return ""
	}
	return config
}

// parseMode Scans the specified config and parses the switchport mode value
//
// Args:
//    config (string): The interface configuration block to scan
//
// Returns:
//    A string value of switchport mode.
func (s *SwitchPortEntity) parseMode(config string) string {
	re := regexp.MustCompile(`(?m)switchport mode (\w+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// paresTrunkGroup Scans the specified config and parses the trunk group values
//
// Args:
//    config (string): The interface configuraiton blcok
//
// Returns:
//    A comma delimitedd string with all trunk group values
func (s *SwitchPortEntity) parseTrunkGroups(config string) string {
	trunkGroups := []string{}
	re := regexp.MustCompile(`(?m)switchport trunk group ([^\s]+)`)
	matches := re.FindAllStringSubmatch(config, -1)
	if matches == nil {
		return ""
	}

	for _, match := range matches {
		for idx, name := range match {
			if idx == 0 || name == "" {
				continue
			}
			trunkGroups = append(trunkGroups, name)
		}
	}
	return strings.Join(trunkGroups, ",")
}

// parseAccessVlan Scans the specified config and parse the access-vlan value
//
// Args:
//    config (string): The interface configuration block to scan
//
// Returns:
//    A string value of switchport access value.
func (s *SwitchPortEntity) parseAccessVlan(config string) string {
	re := regexp.MustCompile(`(?m)switchport access vlan (\d+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseTrunkNativeVlan Scans the specified config and parse the trunk native
// vlan value
//
// Args:
//    config (string): The interface configuration block to scan
//
// Returns:
//    A string value of switchport trunk native vlan value.
func (s *SwitchPortEntity) parseTrunkNativeVlan(config string) string {
	re := regexp.MustCompile(`(?m)switchport trunk native vlan (\d+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseTrunkAllowedVlans Scans the specified config and parse the trunk
// allowed vlans value
//
// Args:
//    config (string): The interface configuration block to scan
//
// Returns:
//    A string value of switchport trunk allowed vlans value.
func (s *SwitchPortEntity) parseTrunkAllowedVlans(config string) string {
	re := regexp.MustCompile(`(?m)switchport trunk allowed vlan (.+)$`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// Create Creates a new logical layer 2 interface
// This method will create a new switchport for the interface specified
// in the arguments (name).  If the logical switchport already exists
// then this command will have no effect
//
// Args:
//    name (string): The interface identifier to create the logical
//        layer 2 switchport for.  The name must be the full interface
//        name and not an abbreviated interface name (eg Ethernet1, not
//        Et1)
//
// Returns:
//    True if the create operation succeeds otherwise False.  If the
//        interface specified in args is already a switchport then this
//        method will have no effect but will still return True
func (s *SwitchPortEntity) Create(name string) bool {
	var commands = []string{"interface " + name,
		"no ip address",
		"switchport",
	}
	return s.Configure(commands...)
}

// Delete Deletes the logical layer 2 interface
// This method will delete the logical switchport for the interface
// specified in the arguments.  If the interface doe not have a logical
// layer 2 interface defined, then this method will have no effect.
//
// Args:
//    name (string): The interface identifier to create the logical
//        layer 2 switchport for.  The name must be the full interface
//        name and not an abbreviated interface name (eg Ethernet1, not
//        Et1)
//
// Returns:
//    True if the create operation succeeds otherwise False.  If the
//        interface specified in args is already a switchport then this
//        method will have no effect but will still return True
func (s *SwitchPortEntity) Delete(name string) bool {
	var commands = []string{"interface " + name,
		"no switchport",
	}
	return s.Configure(commands...)
}

// Default Defaults the configuration of the switchport interface
// This method will default the configuration state of the logical
// layer 2 interface.
//
// Args:
//    name (string): The interface identifier to create the logical
//        layer 2 switchport for.  The name must be the full interface
//        name and not an abbreviated interface name (eg Ethernet1, not
//        Et1)
//
// Returns:
//    True if the create operation succeeds otherwise False.  If the
//        interface specified in args is already a switchport then this
//        method will have no effect but will still return True
func (s *SwitchPortEntity) Default(name string) bool {
	var commands = []string{"interface " + name,
		"no ip address",
		"default switchport",
	}
	return s.Configure(commands...)
}

// SetMode Configures the switchport mode
//
// Args:
//    name (string): The interface identifier to create the logical
//        layer 2 switchport for.  The name must be the full interface
//        name and not an abbreviated interface name (eg Ethernet1, not
//        Et1)
//     value (string): The value to set the mode to.  Accepted values
//        for this argument are access or trunk
//     default (bool): Configures the mode parameter to its default
//        value using the EOS CLI
//
// Returns:
//    True if the create operation succeeds otherwise False.
func (s *SwitchPortEntity) SetMode(name string, value string) bool {
	command := s.CommandBuilder("switchport mode", value, false, true)
	return s.ConfigureInterface(name, command)
}

// SetModeDefault Configures the switchport mode
//
// Args:
//    name (string): The interface identifier to create the logical
//        layer 2 switchport for.  The name must be the full interface
//        name and not an abbreviated interface name (eg Ethernet1, not
//        Et1)
//
// Returns:
//    True if the create operation succeeds otherwise False.
func (s *SwitchPortEntity) SetModeDefault(name string) bool {
	return s.ConfigureInterface(name, "default switchport mode")
}

// SetAccessVlan Configures the switchport access vlan
//
// Args:
//    name (string): The interface identifier to create the logical
//        layer 2 switchport for.  The name must be the full interface
//        name and not an abbreviated interface name (eg Ethernet1, not
//        Et1)
//     value (string): The value to set the access vlan to.  The value
//        must be a valid VLAN ID in the range of 1 to 4094.
//
// Returns:
//    True if the create operation succeeds otherwise False.
func (s *SwitchPortEntity) SetAccessVlan(name string, value string) bool {
	command := s.CommandBuilder("switchport access vlan", value, false, true)
	return s.ConfigureInterface(name, command)
}

// SetAccessVlanDefault Configures the default switchport access vlan
//
// Args:
//    name (string): The interface identifier to create the logical
//        layer 2 switchport for.  The name must be the full interface
//        name and not an abbreviated interface name (eg Ethernet1, not
//        Et1)
//
// Returns:
//    True if the create operation succeeds otherwise False.
func (s *SwitchPortEntity) SetAccessVlanDefault(name string) bool {
	return s.ConfigureInterface(name, "default switchport access vlan")
}

// SetTrunkNativeVlan Configures the switchport trunk native vlan value
//
// Args:
//    name (string): The interface identifier to create the logical
//        layer 2 switchport for.  The name must be the full interface
//        name and not an abbreviated interface name (eg Ethernet1, not
//        Et1)
//     value (string): The value to set the trunk nativevlan to.  The
//        value must be a valid VLAN ID in the range of 1 to 4094.
//
// Returns:
//    True if the create operation succeeds otherwise False.
func (s *SwitchPortEntity) SetTrunkNativeVlan(name string, value string) bool {
	command := s.CommandBuilder("switchport trunk native vlan", value, false, true)
	return s.ConfigureInterface(name, command)
}

// SetTrunkNativeVlanDefault Configures the default switchport trunk native
// vlan value
//
// Args:
//    name (string): The interface identifier to create the logical
//        layer 2 switchport for.  The name must be the full interface
//        name and not an abbreviated interface name (eg Ethernet1, not
//        Et1)
//
// Returns:
//    True if the operation succeeds otherwise False.
func (s *SwitchPortEntity) SetTrunkNativeVlanDefault(name string) bool {
	return s.ConfigureInterface(name, "default switchport trunk native vlan")
}

// SetTrunkAllowedVlans Configures the switchport trunk allowed vlans value
//
// Args:
//    name (string): The interface identifier to create the logical
//        layer 2 switchport for.  The name must be the full interface
//        name and not an abbreviated interface name (eg Ethernet1, not
//        Et1)
//     value (string): The value to set the trunk allowed vlans to.  The
//        value must be a valid VLAN ID in the range of 1 to 4094.
//
// Returns:
//    True if the create operation succeeds otherwise False.
func (s *SwitchPortEntity) SetTrunkAllowedVlans(name string, value string) bool {
	command := s.CommandBuilder("switchport trunk allowed vlan", value, false, true)
	return s.ConfigureInterface(name, command)
}

// SetTrunkAllowedVlansDefault Configures the default switchport trunk allowed
//
// Args:
//    name (string): The interface identifier to create the logical
//        layer 2 switchport for.  The name must be the full interface
//        name and not an abbreviated interface name (eg Ethernet1, not
//        Et1)
//
// Returns:
//    True if the create operation succeeds otherwise False.
func (s *SwitchPortEntity) SetTrunkAllowedVlansDefault(name string) bool {
	return s.ConfigureInterface(name, "default switchport trunk allowed vlan")
}

// SetTrunkGroups Configures the switchport trunk group value
//
// Args:
//    intf (string): The interface identifier to configure.
//    value (string): The set of values to configure the trunk group
//
// Returns:
//    True if the config operation succeeds otherwise False
func (s *SwitchPortEntity) SetTrunkGroups(intf string, value []string) bool {
	var failure = false

	currentValue := strings.Split(s.Get(intf)["trunk_groups"], ",")

	diff := findDiff(value, currentValue)
	for _, name := range diff {
		if ok := s.AddTrunkGroup(intf, name); !ok {
			failure = true
		}
	}

	diff = findDiff(currentValue, value)
	for _, name := range diff {
		if ok := s.RemoveTrunkGroup(intf, name); !ok {
			failure = true
		}
	}
	return !failure
}

// SetTrunkGroupsDefault Configures default switchport trunk group value
//
// Args:
//    intf (string): The interface identifier to configure.
//
// Returns:
//    True if the config operation succeeds otherwise False
func (s *SwitchPortEntity) SetTrunkGroupsDefault(intf string) bool {
	return s.ConfigureInterface(intf, "default switchport trunk group")
}

// AddTrunkGroup Adds the specified trunk group to the interface
//
// Args:
//    intf (string): The interface name to apply the trunk group to
//    value (string): The trunk group value to apply to the interface
//
// Returns:
//    True if the operation as successfully applied otherwise false
func (s *SwitchPortEntity) AddTrunkGroup(intf string, value string) bool {
	str := "switchport trunk group " + value
	return s.ConfigureInterface(intf, str)
}

// RemoveTrunkGroup Removes a specified trunk group to the interface
//
// Args:
//    intf (string): The interface name to remove the trunk group from
//    value (string): The trunk group value
//
// Returns:
//    True if the operation as successfully applied otherwise false
func (s *SwitchPortEntity) RemoveTrunkGroup(intf string, value string) bool {
	str := "no switchport trunk group " + value
	return s.ConfigureInterface(intf, str)
}
