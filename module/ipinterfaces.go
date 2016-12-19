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

// IPInterfaceConfig represents a parsed IPInterface entry
type IPInterfaceConfig map[string]string

// IPInterfaceConfigMap represents a hash of various IPInterfaceConfig's
type IPInterfaceConfigMap map[string]IPInterfaceConfig

// IPInterfaceEntity provides a configuration resource for IpInterfaces
type IPInterfaceEntity struct {
	*AbstractBaseEntity
}

// Name returns the name of the Ip Interface
func (i IPInterfaceConfig) Name() string {
	return i["name"]
}

// Address returns the address of the Ip Interface Config
func (i IPInterfaceConfig) Address() string {
	return i["address"]
}

// Mtu returns the mtu within Ip Interface Config
func (i IPInterfaceConfig) Mtu() string {
	return i["mtu"]
}

// IPInterface factory function to initiallize IPInterface resource
// given a Node
func IPInterface(node *goeapi.Node) *IPInterfaceEntity {
	return &IPInterfaceEntity{&AbstractBaseEntity{node}}
}

// isValidMtu validates the MTU size
func isValidMtu(value int) bool {
	if value >= 68 && value <= 65535 {
		return true
	}
	return false
}

// Get Returns the specific IP interface properties
// The IPinterface resource returns the following:
//
//  * name (str): The name of the interface
//  * address (str): The IP address of the interface in the form
//    of A.B.C.D/E
//  * mtu (int): The configured value for IP MTU.
//
// Args:
//  name (string): The interface identifier to retrieve the
//                 configuration for
//
// Return:
//     An IPInterfaceConfig object of key/value pairs that represents
//     the current configuration of the node.  If the specified
//     interface does not exist then nil is returned.
func (i *IPInterfaceEntity) Get(name string) (IPInterfaceConfig, error) {
	parent := "interface " + name
	config, _ := i.GetBlock(parent)

	matchedInt, _ := regexp.MatchString("Et|Po", name)
	noSwport, _ := regexp.MatchString("(?m)no switchport$", config)

	if matchedInt && !noSwport {
		return nil, nil
	}

	return IPInterfaceConfig{
		"name":    name,
		"address": i.parseAddress(config),
		"mtu":     i.parseMtu(config),
	}, nil
}

// GetAll Returns all of the IP interfaces found in the running-config
// Example:
//    {
//        'Ethernet1': {...},
//        'Ethernet2': {...}
//    }
//
// Returns:
//  A map'd object of key/value pairs keyed by interface
//  name that represents all of the IP interfaces on
//  the current node.
func (i *IPInterfaceEntity) GetAll() IPInterfaceConfigMap {
	var re = regexp.MustCompile(`(?m)^interface\s(.+)`)

	config := i.Config()

	interfaces := re.FindAllStringSubmatch(config, -1)
	response := make(IPInterfaceConfigMap)

	for _, name := range interfaces {
		intf, _ := i.Get(name[1])
		if intf != nil {
			response[name[1]] = intf
		}
	}
	return response
}

// parseAddress Parses the config block and returns the ip address value
// The provided configuration block is scaned and the configured value
// for the IP address is returned as a string.  If the IP address
// value is not configured, then None is returned for the value
//
// Args:
//  config (str): The interface configuration block to parse
//
// Return:
//  string: address of ip interface
func (i *IPInterfaceEntity) parseAddress(config string) string {
	var re = regexp.MustCompile(`ip address ([^\s]+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseMtu Parses the config block and returns the configured IP MTU value
// The provided configuration block is scanned and the configured value
// for the IP MTU is returned as as string object.  The IP MTU value is
// expected to always be present in the provided config block
//
// Args:
//  config (str): The interface configuration block to parse
//
// Return:
//  string representation of MTU size
func (i *IPInterfaceEntity) parseMtu(config string) string {
	var re = regexp.MustCompile(`mtu (\d+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// GetEthInterfaces Returns all of the Interfaces found in the running-config
// Returns:
//  []string of interfaces
func (i *IPInterfaceEntity) GetEthInterfaces() []string {
	var re = regexp.MustCompile(`(?m)^interface\s(Eth.+)`)
	config := i.Config()

	interfaces := re.FindAllStringSubmatch(config, -1)

	response := make([]string, len(interfaces))

	for idx, name := range interfaces {
		response[idx] = name[1]
	}
	return response
}

// Create Creates a new IP interface instance
// This method will create a new logical IP interface for the specified
// physical interface.   If a logical IP interface already exists then
// this operation will have no effect.
//
// Note:
//  Configuring a logical IP interface on a physical interface will
//  remove any existing logical switchports have have been created
//
// Args:
//  name (string): The interface identifier to create the logical
//  layer 3 IP interface for.  The name must be the full interface
//  name and not an abbreviated interface name (eg Ethernet1, not
//  Et1).
//
// Returns:
//  True if the create operation succeeds otherwise False.  If the
//  specified interface is already created the this method will
//  have no effect but will still return True
func (i *IPInterfaceEntity) Create(name string) bool {
	commands := []string{
		"interface " + name,
		"no switchport",
	}
	return i.Configure(commands...)
}

// Delete Deletes an IP interface instance from the running configuration
// This method will delete the logical IP interface for the specified
// physical interface.  If the interface does not have a logical
// IP interface defined, then this method will have no effect.
//
// Args:
//  name (string): The interface identifier to create the logical
//  layer 3 IP interface for.  The name must be the full interface
//  name and not an abbreviated interface name (eg Ethernet1, not
//  Et1).
//
// Returns:
//  True if the delete operation succeeds otherwise False.
func (i *IPInterfaceEntity) Delete(name string) bool {
	commands := []string{
		"interface " + name,
		"no ip address",
		"switchport",
	}
	return i.Configure(commands...)
}

// SetAddress Configures the interface IP address
// Args:
//  name (string): The interface identifier to apply the interface
//                 config to
//  value (string): The IP address and mask to set the interface to.
//                  The value should be in the format of A.B.C.D/E
//                  Value of "" deconfigured ip address
//
// Returns:
//  True if the operation succeeds
func (i *IPInterfaceEntity) SetAddress(name string, value string) bool {
	commands := []string{"interface " + name}
	if value != "" {
		commands = append(commands, "ip address "+value)
	} else {
		commands = append(commands, "no ip address")
	}
	return i.Configure(commands...)
}

// SetAddressDefault Configures the default interface IP address
// Args:
//  name (string): The interface identifier to apply the interface
//                 config to
// Returns:
//  True if the operation succeeds
func (i *IPInterfaceEntity) SetAddressDefault(name string) bool {
	commands := []string{
		"interface " + name,
		"default ip address",
	}
	return i.Configure(commands...)
}

// SetMtu Configures the interface IP MTU
// Args:
//  name (string): The interface identifier to apply the interface
//                 config to
//  value (integer): The MTU value to set the interface to.  Accepted
//                   values include 68 to 65535
//
// Returns:
//  True if the operation succeeds otherwise False.
func (i *IPInterfaceEntity) SetMtu(name string, value int) bool {
	if !isValidMtu(value) {
		return false
	}
	commands := []string{
		"interface " + name,
		"mtu " + strconv.Itoa(value),
	}
	return i.Configure(commands...)
}

// SetMtuDefault Configures the default interface IP MTU
// Args:
//  name (string): The interface identifier to apply the interface
//                 config to
// Returns:
//  True if the operation succeeds otherwise False.
func (i *IPInterfaceEntity) SetMtuDefault(name string) bool {
	commands := []string{
		"interface " + name,
		"default mtu",
	}
	return i.Configure(commands...)
}
