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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aristanetworks/goeapi"
)

var (
	nameRegex       = regexp.MustCompile(`(?m)(?:name\s)(.*)$`)
	stateRegex      = regexp.MustCompile(`(?m)(?:state\s)(.*)$`)
	trunkGroupRegex = regexp.MustCompile(`(?m)(?:trunk\sgroup\s)(.*)$`)
	vlansRegex      = regexp.MustCompile(`(?m)^vlan\s(\d+)$`)
)

// VlanConfig represents a parsed vlan entry containing
type VlanConfig map[string]string

// Name returns the Vlan name for this VlanConfig
// Null string returned if not set
func (v VlanConfig) Name() string {
	return v["name"]
}

// State returns the Vlan state for this VlanConfig
// Null string returned if not set
func (v VlanConfig) State() string {
	return v["state"]
}

// TrunkGroups returns the Vlan trunk groups configured for this VlanConfig
// Null string returned if not set
func (v VlanConfig) TrunkGroups() string {
	return v["trunk_groups"]
}

// VlanConfigMap represents a parsed vlan entry containing
type VlanConfigMap map[string]VlanConfig

// VlanEntity provides a configuration resource for VLANs
type VlanEntity struct {
	*AbstractBaseEntity
}

// Vlan factory function to initiallize Vlans resource
// given a Node
func Vlan(node *goeapi.Node) *VlanEntity {
	return &VlanEntity{&AbstractBaseEntity{node}}
}

// findDiff helper function to find difference between two string slices
func findDiff(slice1 []string, slice2 []string) []string {
	var diff []string
	hash := map[string]int{}

	for _, val := range slice2 {
		if val == "" {
			continue
		}
		hash[val] = 1
	}

	for _, val := range slice1 {
		if val == "" {
			continue
		}
		if _, found := hash[val]; found {
			continue
		}
		diff = append(diff, val)
	}
	return diff
}

// isVlan helper function to validate vlan id
//
// Args:
//  vlan (string): vlan id
func isVlan(vlan string) bool {
	vid, _ := strconv.Atoi(vlan)
	return vid > 0 && vid < 4095
}

// Get returns the VLAN configuration as a resource object.
// Args:
//  vid (string): The vlan identifier to retrieve from the
//                running configuration.  Valid values are in the range
//                of 1 to 4095
//
// Returns:
//  VlanConfig object containing the VLAN attributes as
//  key/value pairs.
func (v *VlanEntity) Get(vlan string) VlanConfig {
	parent := "vlan " + vlan
	config, err := v.GetBlock(parent)
	if err != nil {
		return nil
	}
	var resource = make(VlanConfig)
	resource["name"] = v.parseName(config)
	resource["state"] = v.parseState(config)
	resource["trunk_groups"] = v.parseTrunkGroups(config)
	return resource
}

// GetAll Returns a VlanConfigMap object of all Vlans in the running-config
// Returns:
//  A VlanConfigMap type of all Vlan attributes
func (v *VlanEntity) GetAll() VlanConfigMap {
	config := v.Config()

	matches := vlansRegex.FindAllStringSubmatch(config, -1)
	var resources = make(VlanConfigMap)
	for _, val := range matches {
		vid := val[1]
		resources[vid] = v.Get(vid)
	}
	return resources
}

// GetSection returns the specified Vlan Entry for the name specified.
//
// Args:
//  vlan (string): The vlan id
//
// Returns:
//  Returns string representation of Vlan config entry
func (v *VlanEntity) GetSection(vlan string) string {
	parent := fmt.Sprintf(`vlan\s+(%s$)|(%s,.*)|(.*,%s,)|(.*,%s$)`,
		vlan, vlan, vlan, vlan)
	config, err := v.GetBlock(parent)
	if err != nil {
		return ""
	}
	return config
}

// parseName scans the provided configuration block and extracts
// the vlan name.  The config block is expected to always return the
// vlan name.
//
// Args:
//  config (string): The vlan configuration block from the nodes running
//                configuration
// Returns:
//  string value of name
func (v *VlanEntity) parseName(config string) string {
	if config == "" {
		return ""
	}
	match := nameRegex.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseState scans the provided configuration block and extracts
// the vlan state value.  The config block is expected to always return
// the vlan state config.
//
// Args:
//  config (string): The vlan configuration block from the nodes
//                running configuration
// Returns:
//  string: state of the vlan
func (v *VlanEntity) parseState(config string) string {
	if config == "" {
		return ""
	}
	match := stateRegex.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseTrunkGroups scans the provided configuration block and
// extracts all the vlan trunk groups.  If no trunk groups are configured
// an empty string is returned as the vlaue.
//
// Args:
//  config (string): The vlan configuration block form the node's
//                running configuration
// Returns:
//  string: comma separated list of trunkgroups
func (v *VlanEntity) parseTrunkGroups(config string) string {
	trunkGroups := []string{}

	matches := trunkGroupRegex.FindAllStringSubmatch(config, -1)
	if matches == nil {
		return ""
	}

	for _, match := range matches {
		for idx, tid := range match {
			if idx == 0 || tid == "" {
				continue
			}
			trunkGroups = append(trunkGroups, tid)
		}
	}
	return strings.Join(trunkGroups, ",")
}

// Create Creates a new VLAN resource
//
// Args:
//  vid (string): The VLAN ID to create
//
// Returns:
//  True if create was successful otherwise False
func (v *VlanEntity) Create(vid string) bool {
	var commands = []string{"vlan " + vid}
	if isVlan(vid) {
		return v.Configure(commands...)
	}
	return false
}

// Delete Deletes a VLAN from the running configuration
//
// Args:
//  vid (string): The VLAN ID to delete
//
// Returns:
//  True if the operation was successful otherwise False
func (v *VlanEntity) Delete(vid string) bool {
	var commands = []string{"no vlan " + vid}
	if isVlan(vid) {
		return v.Configure(commands...)
	}
	return false
}

// Default Defaults the VLAN configuration
//
// .. code-block:: none
//
//    default vlan <vlanid>
//
// Args:
//  vid (string): The VLAN ID to default
//
// Returns:
//  True if the operation was successful otherwise False
func (v *VlanEntity) Default(vid string) bool {
	var commands = []string{"default vlan " + vid}
	if isVlan(vid) {
		return v.Configure(commands...)
	}
	return false
}

// ConfigureVlan Configures the specified Vlan using commands
//
// Args:
//  vid (string): The VLAN ID to configure
//  commands: The list of commands to configure
//
// Returns:
//  True if the commands completed successfully
func (v *VlanEntity) ConfigureVlan(vid string, cmds ...string) bool {
	var commands = []string{"vlan " + vid}
	commands = append(commands, cmds...)
	return v.Configure(commands...)
}

// SetName Configures the VLAN name
//
// EosVersion:
//    4.13.7M
//
// Args:
//  vid (string): The VLAN ID to Configures
//  name (string): The value to configure the vlan name
//
// Returns:
//  True if the operation was successful otherwise False
func (v *VlanEntity) SetName(vid string, name string) bool {
	return v.ConfigureVlan(vid, "name "+name)
}

// SetNameDefault Configures the VLAN name
//
// EosVersion:
//    4.13.7M
//
// Args:
//  vid (string): The VLAN ID to Configures
//
// Returns:
//  True if the operation was successful otherwise False
func (v *VlanEntity) SetNameDefault(vid string) bool {
	return v.ConfigureVlan(vid, "default name")
}

// SetState Configures the VLAN state
//
// EosVersion:
//    4.13.7M
//
// Args:
//  vid (string): The VLAN ID to configure
//  value (string): The value to set the vlan state to
//
// Returns:
//  True if the operation was successful otherwise False
func (v *VlanEntity) SetState(vid string, value string) bool {
	if value == "" {
		return v.ConfigureVlan(vid, "no state")
	}
	return v.ConfigureVlan(vid, "state "+value)
}

// SetStateDefault Configures the VLAN state
//
// EosVersion:
//    4.13.7M
//
// Args:
//  vid (string): The VLAN ID to configure
//
// Returns:
//  True if the operation was successful otherwise False
func (v *VlanEntity) SetStateDefault(vid string) bool {
	return v.ConfigureVlan(vid, "default state")
}

// SetTrunkGroup Configures the list of trunk groups support on a vlan
//
// This method handles configuring the vlan trunk group value to default
// if the default flag is set to True.  If the default flag is set
// to False, then this method will calculate the set of trunk
// group names to be added and to be removed.
//
// EosVersion:
//    4.13.7M
//
// Args:
//  vid (string): The VLAN ID to configure
//  value (string): The list of trunk groups that should be configured
//               for this vlan id.
//
// Returns:
//  True if the operation was successful otherwise False
func (v *VlanEntity) SetTrunkGroup(vid string, value []string) bool {
	var failure = false

	currentValue := strings.Split(v.Get(vid)["trunk_groups"], ",")

	diff := findDiff(value, currentValue)
	for _, name := range diff {
		if ok := v.AddTrunkGroup(vid, name); !ok {
			failure = true
		}
	}
	diff = findDiff(currentValue, value)
	for _, name := range diff {
		if ok := v.RemoveTrunkGroup(vid, name); !ok {
			failure = true
		}
	}
	return !failure
}

// SetTrunkGroupDefault Configures the default list of trunk groups support on a vlan
//
// This method handles configuring the vlan trunk group value to default
// if the default flag is set to True.  If the default flag is set
// to False, then this method will calculate the set of trunk
// group names to be added and to be removed.
//
// EosVersion:
//    4.13.7M
//
// Args:
//  vid (string): The VLAN ID to configure
//
// Returns:
//  True if the operation was successful otherwise False
func (v *VlanEntity) SetTrunkGroupDefault(vid string) bool {
	return v.ConfigureVlan(vid, "default trunk group")
}

// AddTrunkGroup Adds a new trunk group to the Vlan in the running-config
//
// EosVersion:
//    4.13.7M
//
// Args:
//  vid (string): The VLAN ID to configure
//  name (string): The trunk group to add to the list
//
// Returns:
//  True if the operation was successful otherwise False
func (v *VlanEntity) AddTrunkGroup(vid string, name string) bool {
	var commands = []string{"trunk group " + name}
	return v.ConfigureVlan(vid, commands...)
}

// RemoveTrunkGroup Removes a trunk group from the list of configured trunk
// groups for the specified VLAN ID
//
// EosVersion:
//    4.13.7M
//
// Args:
//  vid (string): The VLAN ID to configure
//  name (string): The trunk group to add to the list
//
// Returns:
//  True if the operation was successful otherwise False
func (v *VlanEntity) RemoveTrunkGroup(vid string, name string) bool {
	var commands = []string{"no trunk group " + name}
	return v.ConfigureVlan(vid, commands...)
}
