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
	"strings"

	"github.com/aristanetworks/goeapi"
)

var validInterfaces = map[string]bool{
	"Ethernet":     true,
	"Management":   true,
	"Loopback":     true,
	"Vlan":         true,
	"Port-Channel": true,
	"Vxlan":        true,
}

// InterfaceConfig represents the parsed Interface config
// {
//		"name"			: "ethernet1",
//		"type"			: "generic",
//		"shutdown"		: "false",
//		"description"	: "Backhaul-to-East",
// }
type InterfaceConfig map[string]string

// BaseInterfaceEntity provides a configuration resource for Interface
type BaseInterfaceEntity struct {
	*AbstractBaseEntity
}

// Interface factory function to initiallize BaseInterfaceEntity resource
// given a Node
func Interface(node *goeapi.Node) *BaseInterfaceEntity {
	return &BaseInterfaceEntity{&AbstractBaseEntity{node}}
}

// isValidInterface provides some first level checking of interface
// name validity
func isValidInterface(value string) bool {
	var validIntf = regexp.MustCompile(`([EPVLM][a-z-C]+)`)
	match := validIntf.FindString(value)
	if match == "" {
		return false
	}
	_, found := validInterfaces[match]
	return found
}

// Get returns the interface config for the given interface name(string)
func (i *BaseInterfaceEntity) Get(name string) InterfaceConfig {
	parent := `interface\s+` + name
	config, _ := i.GetBlock(parent)

	return InterfaceConfig{
		"name":        name,
		"type":        "generic",
		"shutdown":    strconv.FormatBool(i.parseShutdown(config)),
		"description": i.parseDescription(config),
	}
}

// parseShutdown returns true if shutdown is seen in interface config
// or false is not shutdown.
func (i *BaseInterfaceEntity) parseShutdown(config string) bool {
	matched, _ := regexp.MatchString(`(?m)no shutdown`, config)
	return !matched
}

// parseDescription returns the description specified in the interface config.
func (i *BaseInterfaceEntity) parseDescription(config string) string {
	re := regexp.MustCompile(`(?m)description (.+)$`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// Create creates a new interface on the node
func (i *BaseInterfaceEntity) Create(name string) bool {
	return i.Configure("interface " + name)
}

// Delete removes an interface from the node
func (i *BaseInterfaceEntity) Delete(name string) bool {
	return i.Configure("no interface " + name)
}

// Default reverts back to default config for interface
func (i *BaseInterfaceEntity) Default(name string) bool {
	return i.Configure("default interface " + name)
}

// SetDescription sets the description on the interface name(sting) to value(string)
func (i *BaseInterfaceEntity) SetDescription(name string, value string) bool {
	cmd := i.CommandBuilder("description", value, false, true)
	return i.ConfigureInterface(name, cmd)
}

// SetDescriptionDefault reverts back to the default description value
func (i *BaseInterfaceEntity) SetDescriptionDefault(name string) bool {
	cmd := i.CommandBuilder("description", "", true, false)
	return i.ConfigureInterface(name, cmd)
}

// SetShutdown sets the interface name(string) to shutdown(true)
// or no-shutdown(false)
func (i *BaseInterfaceEntity) SetShutdown(name string, shut bool) bool {
	cmd := i.CommandBuilder("shutdown", "", false, shut)
	return i.ConfigureInterface(name, cmd)
}

// SetShutdownDefault reverts back to the default shutdown config for interface
func (i *BaseInterfaceEntity) SetShutdownDefault(name string) bool {
	cmd := i.CommandBuilder("shutdown", "", true, false)
	return i.ConfigureInterface(name, cmd)
}

///////////////////////////////

// EthernetInterfaceEntity provides a configuration resource for
// Ethernet Interface
type EthernetInterfaceEntity struct {
	*BaseInterfaceEntity
}

// EthernetInterface factory function to initiallize EthernetInterfaceEntity resource
// given a Node
func EthernetInterface(node *goeapi.Node) *EthernetInterfaceEntity {
	return &EthernetInterfaceEntity{&BaseInterfaceEntity{&AbstractBaseEntity{node}}}
}

// Get returns interface as a set of key/value pairs in InterfaceConfig
func (e *EthernetInterfaceEntity) Get(name string) InterfaceConfig {
	parent := `interface\s+` + name
	config, _ := e.GetBlock(parent)

	resource := e.BaseInterfaceEntity.Get(name)
	resource["type"] = "ethernet"
	resource["sflow"] = strconv.FormatBool(parseSflow(config))
	resource["flowcontrol_send"] = parseFlowControlSend(config)
	resource["flowcontrol_receive"] = parseFlowControlReceive(config)
	return resource
}

// parseSflow parses the given config(string) and returns true(bool)
// if sflow configured
func parseSflow(config string) bool {
	matched, _ := regexp.MatchString("no sflow", config)
	return !matched
}

// parseFlowControlSend parses the given config and returns the flowcontrol
// send operation parameter
func parseFlowControlSend(config string) string {
	re := regexp.MustCompile(`(?m)flowcontrol send (\w+)$`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return "off"
	}
	return match[1]
}

// parseFlowControlReceive parses the given config and returns the flowcontrol
// receive operation parameter
func parseFlowControlReceive(config string) string {
	re := regexp.MustCompile(`(?m)flowcontrol receive (\w+)$`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return "off"
	}
	return match[1]
}

// Create not supported
func (e *EthernetInterfaceEntity) Create(name string) bool {
	return false // not implemented
}

// Delete not supported
func (e *EthernetInterfaceEntity) Delete(name string) bool {
	return false // not implemented
}

// SetFlowcontrolSend configures the interface flowcontrol send value(true: on)
func (e *EthernetInterfaceEntity) SetFlowcontrolSend(name string, value bool) bool {
	return e.setFlowcontrol(name, "send", value)
}

// SetFlowcontrolReceive configures the interface flowcontrol receive value(true: on)
func (e *EthernetInterfaceEntity) SetFlowcontrolReceive(name string, value bool) bool {
	return e.setFlowcontrol(name, "receive", value)
}

// setFlowcontrol configures the interface flowcontrol value
func (e *EthernetInterfaceEntity) setFlowcontrol(name string, direction string, value bool) bool {
	var str string
	if value {
		str = "flowcontrol " + direction + " on"
	} else {
		str = "flowcontrol " + direction + " off"
	}
	cmds := []string{
		"interface " + name,
		str,
	}
	return e.Configure(cmds...)
}

// DisableFlowcontrolSend disables the interface flowcontrol send value
func (e *EthernetInterfaceEntity) DisableFlowcontrolSend(name string) bool {
	return e.disableFlowcontrol(name, "send")
}

// DisableFlowcontrolReceive disables the interface flowcontrol receive value
func (e *EthernetInterfaceEntity) DisableFlowcontrolReceive(name string) bool {
	return e.disableFlowcontrol(name, "receive")
}

// DisableFlowcontrol disables the interface flowcontrol
func (e *EthernetInterfaceEntity) disableFlowcontrol(name string, direction string) bool {
	cmds := []string{
		"interface " + name,
		"no flowcontrol " + direction,
	}
	return e.Configure(cmds...)
}

// SetSflow configures the sFlow state (true:enable, false:disable) on the
// interface name(string)
func (e *EthernetInterfaceEntity) SetSflow(name string, value bool) bool {
	str := "no sflow enable"
	if value {
		str = "sflow enable"
	}
	cmds := []string{
		"interface " + name,
		str,
	}
	return e.Configure(cmds...)
}

// SetSflowDefault configures the defalt sFlow state on the
// interface name(string)
func (e *EthernetInterfaceEntity) SetSflowDefault(name string) bool {
	cmds := []string{
		"interface " + name,
		"default sflow",
	}
	return e.Configure(cmds...)
}

///////////////////////////////

const defaultLacpMode = "on"

// PortChannelInterfaceEntity provides a configuration resource for
// PortChannel
type PortChannelInterfaceEntity struct {
	*BaseInterfaceEntity
}

// PortChannel factory function to initiallize VxlanInterfaceEntity resource
// given a Node
func PortChannel(node *goeapi.Node) *PortChannelInterfaceEntity {
	return &PortChannelInterfaceEntity{&BaseInterfaceEntity{&AbstractBaseEntity{node}}}
}

// Get returns the PortChannel interface config for the interface name(string) given.
// Returned is a InterfaceConfig type
func (p *PortChannelInterfaceEntity) Get(name string) InterfaceConfig {
	parent := `interface\s+` + name
	config, _ := p.GetBlock(parent)

	resource := p.BaseInterfaceEntity.Get(name)
	resource["type"] = "portchannel"
	resource["lacp_mode"] = p.getLacpMode(name)
	resource["minimum_links"] = p.parseMinimumLinks(config)
	resource["members"] = strings.Join(p.getMembers(name), ",")
	return resource
}

// parseMinimumLinks returns the configured min-links for the specified Port-Channel
// config
func (p *PortChannelInterfaceEntity) parseMinimumLinks(config string) string {
	re := regexp.MustCompile(`port-channel min-links (\d+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// getLacpMode returns the LACP mode for the specified Port-Channel interface
func (p *PortChannelInterfaceEntity) getLacpMode(name string) string {

	members := p.getMembers(name)
	re := regexp.MustCompile(`channel-group\s\d+\smode\s(.+)`)
	for _, member := range members {
		parent := `interface\s+` + member
		config, _ := p.GetBlock(parent)

		match := re.FindStringSubmatch(config)
		if match == nil {
			return "NaN"
		}
		return match[1]
	}
	return defaultLacpMode
}

// getMembers returns a list of member interfaces for the specified
// Port-Channel
func (p *PortChannelInterfaceEntity) getMembers(name string) []string {
	re := regexp.MustCompile(`(\d+)`)
	match := re.FindStringSubmatch(name)
	if match == nil {
		return nil
	}
	grPid := match[1]
	command := []string{"show port-channel " + grPid + " all-ports"}
	config, _ := p.node.Enable(command)
	tmp := config[0]["result"]

	re = regexp.MustCompile(`\bEthernet[\d]*\b`)
	matches := re.FindAllString(tmp, -1)
	if matches == nil {
		return nil
	}
	return matches
}

// SetMembers configures the array of member interfaces for the Port-Channel
func (p *PortChannelInterfaceEntity) SetMembers(name string, members ...string) bool {
	re := regexp.MustCompile(`(\d+)$`)
	match := re.FindStringSubmatch(name)
	if match == nil {
		return false
	}
	grpID := match[1]
	currentMembers := p.getMembers(name)
	lacpMode := p.getLacpMode(name)

	var commands []string

	// delete
	diff := findDiff(currentMembers, members)
	for _, member := range diff {
		commands = append(commands, "interface "+member)
		commands = append(commands, "no channel-group "+grpID)
	}

	// add
	diff = findDiff(members, currentMembers)
	for _, member := range diff {
		commands = append(commands, "interface "+member)
		commands = append(commands, "channel-group "+grpID+" mode "+lacpMode)
	}
	return p.Configure(commands...)
}

// SetLacpMode configures the LACP mode of the member interfaces
func (p *PortChannelInterfaceEntity) SetLacpMode(name string, mode string) bool {
	validModes := map[string]bool{
		"on":      true,
		"passive": true,
		"active":  true,
	}
	if _, found := validModes[mode]; !found {
		return false
	}
	re := regexp.MustCompile(`(\d+)$`)
	match := re.FindStringSubmatch(name)
	if match == nil {
		return false
	}
	grpID := match[1]

	var removeCommands []string
	var addCommands []string
	for _, member := range p.getMembers(name) {
		removeCommands = append(removeCommands, "interface "+member)
		removeCommands = append(removeCommands, "no channel-group "+grpID)
		addCommands = append(addCommands, "interface "+member)
		addCommands = append(addCommands, "channel-group "+grpID+" mode "+mode)
	}
	return p.Configure(append(removeCommands, addCommands...)...)
}

// SetMinimumLinks configures the Port-Channel min-links value
func (p *PortChannelInterfaceEntity) SetMinimumLinks(name string, value int) bool {
	if value < 1 || value > 16 {
		return false
	}
	cmd := "port-channel min-links " + strconv.Itoa(value)
	commands := []string{
		"interface " + name,
		cmd,
	}
	return p.Configure(commands...)
}

// SetMinimumLinksDefault returns the specified interface min-links config to it's
// default configuration.
func (p *PortChannelInterfaceEntity) SetMinimumLinksDefault(name string) bool {
	commands := []string{
		"interface " + name,
		"default port-channel min-links",
	}
	return p.Configure(commands...)
}

// VxlanInterfaceConfig represents the parsed Vxlan interface config
// {
//
//		"name"				: "Vxlan1"
//		"type"				: "vxlan"
//		"shutdown"			: "false"
//		"description"		: ""
//		"source_interface"	: "Ethernet1",
//		"multicast_group"	: "",
//		"udp_port"			: "1024",
//		"flood_list"		: "",
// }
type VxlanInterfaceConfig map[string]string

// VxlanConfigCollection is a collection of Vxlan interfaces
// {
//		"1" : VxlanInterfaceConfig{},
// }
type VxlanConfigCollection map[string]VxlanInterfaceConfig

// VxlanInterfaceEntity provides a configuration resource for Vxlan
type VxlanInterfaceEntity struct {
	*BaseInterfaceEntity
}

// Vxlan factory function to initiallize VxlanInterfaceEntity resource
// given a Node
func Vxlan(node *goeapi.Node) *VxlanInterfaceEntity {
	return &VxlanInterfaceEntity{&BaseInterfaceEntity{&AbstractBaseEntity{node}}}
}

// Get returns the Vxlan interface config for the interface name(string) given.
// Returned is a InterfaceConfig type
func (v *VxlanInterfaceEntity) Get(name string) InterfaceConfig {
	parent := `interface\s+` + name
	config, _ := v.GetBlock(parent)

	resource := v.BaseInterfaceEntity.Get(name)
	resource["type"] = "vxlan"
	resource["source_interface"] = v.parseSourceInterface(config)
	resource["multicast_group"] = v.parseMulticastGroup(config)
	resource["udp_port"] = v.parseUDPPort(config)
	//resource["vlans"]            = v.parseVlans(config)
	resource["flood_list"] = v.parseFloodList(config)
	return resource
}

// parseSourceInterface parses given vxlan interface config and returns the source
// inteface or empty string if not found
func (v *VxlanInterfaceEntity) parseSourceInterface(config string) string {
	re := regexp.MustCompile(`vxlan source-interface ([^\s]+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseMulticastGroup parses given vxlan interface config and returns the
// multicast group associated with it or empty string if not found.
func (v *VxlanInterfaceEntity) parseMulticastGroup(config string) string {
	re := regexp.MustCompile(`vxlan multicast-group ([^\s]+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseUDPPort parses given vxlan interface config and returns the udp port
// configured or empty string if not found.
func (v *VxlanInterfaceEntity) parseUDPPort(config string) string {
	re := regexp.MustCompile(`vxlan udp-port (\d+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseVlans parses given vxlan interface config and returns the collection
// of vlans configured. Returns VxlanConfigCollection.
func (v *VxlanInterfaceEntity) parseVlans(config string) VxlanConfigCollection {
	re := regexp.MustCompile(`(?m)vxlan vlan (\d+)`)
	vlans := re.FindAllStringSubmatch(config, -1)

	hash := map[string]struct{}{}

	collection := make(VxlanConfigCollection)

	for _, matchedLine := range vlans {
		vid := matchedLine[1]
		if _, found := hash[vid]; !found {
			hash[vid] = struct{}{}

			reStr := `vxlan vlan ` + vid + ` vni (\d+)`
			reGlobalFld := regexp.MustCompile(reStr)
			match := reGlobalFld.FindStringSubmatch(config)

			collection[vid] = make(VxlanInterfaceConfig)
			if match == nil {
				collection[vid]["vni"] = ""
			} else {
				collection[vid]["vni"] = match[1]
			}

			reStr = `(?m)vxlan vlan ` + vid + ` flood vtep (.*)$`
			reLocalFld := regexp.MustCompile(reStr)
			match = reLocalFld.FindStringSubmatch(config)
			if match == nil {
				collection[vid]["flood_list"] = ""
			} else {
				collection[vid]["flood_list"] = match[1]
			}

		}
	}
	return collection
}

// parseFloodList parses given vxlan interface config and returns the flood list
// in the form of comma delimited string.
func (v *VxlanInterfaceEntity) parseFloodList(config string) string {
	re := regexp.MustCompile(`(?m)vxlan flood vtep (.+)$`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return strings.Join(strings.Split(match[1], " "), ",")
}

// SetSourceInterface sets the vxlan interface to the given value(string).
// If empty string is specified, then default setting is used.
func (v *VxlanInterfaceEntity) SetSourceInterface(name string, value string) bool {
	var cmd string
	if value == "" {
		cmd = v.CommandBuilder("vxlan source-interface", value, false, false)
	} else {
		cmd = v.CommandBuilder("vxlan source-interface", value, false, true)
	}
	return v.ConfigureInterface(name, cmd)
}

// SetSourceInterfaceDefault sets the vxlan interface source-interface back to
// default settings
func (v *VxlanInterfaceEntity) SetSourceInterfaceDefault(name string) bool {
	cmd := v.CommandBuilder("vxlan source-interface", "", true, false)
	return v.ConfigureInterface(name, cmd)
}

// SetMulticastGroup sets the vxlan interface multicast-group configuration to the
// value specified.
// If empty string is specified, then default setting is used.
func (v *VxlanInterfaceEntity) SetMulticastGroup(name string, value string) bool {
	var cmd string
	if value == "" {
		cmd = v.CommandBuilder("vxlan multicast-group", value, false, false)
	} else {
		cmd = v.CommandBuilder("vxlan multicast-group", value, false, true)
	}
	return v.ConfigureInterface(name, cmd)
}

// SetMulticastGroupDefault sets the vxlan interface multicast-group configuration
// back to default settings.
func (v *VxlanInterfaceEntity) SetMulticastGroupDefault(name string) bool {
	cmd := v.CommandBuilder("vxlan multicast-group", "", true, false)
	return v.ConfigureInterface(name, cmd)
}

// SetUDPPort sets the vxlan interface udp port to the provided port value(int)
func (v *VxlanInterfaceEntity) SetUDPPort(name string, value int) bool {
	if value < 1024 || value > 65535 {
		return v.SetUDPPortDefault(name)
	}
	cmd := v.CommandBuilder("vxlan udp-port", strconv.Itoa(value), false, true)
	return v.ConfigureInterface(name, cmd)
}

// SetUDPPortDefault sets the vxlan interface udp port configuration to the default
// settings
func (v *VxlanInterfaceEntity) SetUDPPortDefault(name string) bool {
	cmd := v.CommandBuilder("vxlan udp-port", "", true, false)
	return v.ConfigureInterface(name, cmd)
}

// AddVtepGlobalFlood adds to interface name(string) a vtep(string) endpoint with the
// global flood list
func (v *VxlanInterfaceEntity) AddVtepGlobalFlood(name string, vtep string) bool {
	cmd := "vxlan flood vtep add " + vtep
	return v.ConfigureInterface(name, cmd)
}

// AddVtepLocalFlood adds to interface name(string) a vtep(string) endpoint with the
// local vlan(int) flood list
func (v *VxlanInterfaceEntity) AddVtepLocalFlood(name string, vtep string, vlan int) bool {
	cmd := "vxlan vlan " + strconv.Itoa(vlan) + " flood vtep add " + vtep
	return v.ConfigureInterface(name, cmd)
}

// RemoveVtepGlobalFlood removes from interface name(string) a global vtep(string) flood list
func (v *VxlanInterfaceEntity) RemoveVtepGlobalFlood(name string, vtep string) bool {
	cmd := "vxlan flood vtep remove " + vtep
	return v.ConfigureInterface(name, cmd)
}

// RemoveVtepLocalFlood removes from interface name(string) a vtep(string) endpoint from the
// local vlan(int) flood list
func (v *VxlanInterfaceEntity) RemoveVtepLocalFlood(name string, vtep string, vlan int) bool {
	cmd := "vxlan vlan " + strconv.Itoa(vlan) + " flood vtep remove " + vtep
	return v.ConfigureInterface(name, cmd)
}

// UpdateVlan adds a new vlan vid(int) to vni(int) for the interface name(string)
func (v *VxlanInterfaceEntity) UpdateVlan(name string, vid int, vni int) bool {
	cmd := "vxlan vlan " + strconv.Itoa(vid) + " vni " + strconv.Itoa(vni)
	return v.ConfigureInterface(name, cmd)
}

// RemoveVlan removes a vlan vid(int) to vni mapping from a given
// interface name(string).
func (v *VxlanInterfaceEntity) RemoveVlan(name string, vid int) bool {
	cmd := "no vxlan vlan " + strconv.Itoa(vid) + " vni"
	return v.ConfigureInterface(name, cmd)
}
