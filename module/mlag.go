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

// GlobalMlagConfig represents a parsed Global Mlag entry
// {
//		"domain_id"      : "2",
//		"local_interface": "1.1.1.1",
//		"peer_address"   : "2.2.2.2",
//		"peer_link"      : "Ethernet1",
//		"shutdown"       : "false",
// }
type GlobalMlagConfig map[string]string

// InterfaceMlagConfig represents the parsed Mlag config for all
// interfaces
// {
//		"Port-Channel1" : "2",
//		"Port-Channel10" : "5",
//		...
// }
type InterfaceMlagConfig map[string]string

// MlagConfig represents an Mlag Config entry
type MlagConfig struct {
	config     GlobalMlagConfig
	interfaces InterfaceMlagConfig
}

// DomainID returns the mlag global domain id
// Empty string is returned if not found
func (m MlagConfig) DomainID() string {
	return m.config["domain_id"]
}

// LocalInterface returns the mlag local inteface
// Empty string is returned if not found
func (m MlagConfig) LocalInterface() string {
	return m.config["local_interface"]
}

// PeerAddress returns the mlag peer address
// Empty string is returned if not found
func (m MlagConfig) PeerAddress() string {
	return m.config["peer_address"]
}

// PeerLink returns configured peer-link
// Empty string is returned if not found
func (m MlagConfig) PeerLink() string {
	return m.config["peer_link"]
}

// Shutdown returns string 'true' if mlag shutdown
func (m MlagConfig) Shutdown() string {
	return m.config["shutdown"]
}

// InterfaceConfig returns the mlag ID for the given interface.
// Empty string is returned if not found
func (m MlagConfig) InterfaceConfig(intf string) string {
	return m.interfaces[intf]
}

// isEqual compares two MlagConfig objects
func (m MlagConfig) isEqual(dest MlagConfig) bool {
	if (len(m.config) != len(dest.config)) ||
		(len(m.interfaces) != len(dest.interfaces)) {
		return false
	}
	for k, v := range m.config {
		val, found := dest.config[k]
		if !found || val != v {
			return false
		}
	}
	for k, v := range m.interfaces {
		val, found := dest.interfaces[k]
		if !found || val != v {
			return false
		}
	}
	return true
}

// MlagEntity provides a configuration resource for Mlags
type MlagEntity struct {
	*AbstractBaseEntity
}

// Mlag factory function to initiallize MlagEntity resource
// given a Node
func Mlag(node *goeapi.Node) *MlagEntity {
	return &MlagEntity{&AbstractBaseEntity{node}}
}

// Get ...
func (m *MlagEntity) Get() *MlagConfig {
	config := m.parseConfig()
	interfaces := m.parseInterfaces()
	return &MlagConfig{config: config, interfaces: interfaces}
}

// Parses the mlag global configuration
// Returns: A GlobalMlagConfig object key/value
func (m *MlagEntity) parseConfig() GlobalMlagConfig {
	config, _ := m.GetBlock("mlag configuration")
	return GlobalMlagConfig{
		"domain_id":       m.parseDomainID(config),
		"local_interface": m.parseLocalInterface(config),
		"peer_address":    m.parsePeerAddress(config),
		"peer_link":       m.parsePeerLink(config),
		"shutdown":        strconv.FormatBool(m.parseShutdown(config)),
	}
}

// parseDomainID Scans the config block and parses the domain-id value
//
// Args:
//  config (str): The config block to scan
//
// Returns:
//  string value of the domain-id. "" if not found
func (m *MlagEntity) parseDomainID(config string) string {
	if config == "" {
		return ""
	}
	re := regexp.MustCompile(`(?m)domain-id (\w+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseLocalInterface Scans the config block and parses the local-interface
// value
//
// Args:
//  config (str): The config block to scan
//
// Returns:
//  string value of the local-interface. "" if not found
func (m *MlagEntity) parseLocalInterface(config string) string {
	if config == "" {
		return ""
	}
	re := regexp.MustCompile(`(?m)local-interface (\w+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parsePeerAddress Scans the config block and parses the peer-address value
//
// Args:
//  config (str): The config block to scan
//
// Returns:
//  string value of peer address. "" if not found
func (m *MlagEntity) parsePeerAddress(config string) string {
	if config == "" {
		return ""
	}
	re := regexp.MustCompile(`(?m)peer-address ([^\s]+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parsePeerLink Scans the config block and parses the peer-link value
//
// Args:
//  config (str): The config block to scan
//
// Returns:
//  String value of peer-link config. "" if not found
func (m *MlagEntity) parsePeerLink(config string) string {
	if config == "" {
		return ""
	}
	re := regexp.MustCompile(`(?m)peer-link (\S+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseShutdown Scans the config block and parses the shutdown value
//
// Args:
//  config (str): The config block to scan
//
// Returns:
//  bool: True if interface is in shutdown state
func (m *MlagEntity) parseShutdown(config string) bool {
	if config == "" {
		return false
	}
	matched, _ := regexp.MatchString(`(?m)no shutdown`, config)
	if matched {
		return false
	}
	return !matched
}

// parseInterfaces Scans the global config and returns the configured interfaces
//
// Returns:
//  InterfaceMlagConfig object
func (m *MlagEntity) parseInterfaces() InterfaceMlagConfig {
	config := m.Config()

	var resource = make(InterfaceMlagConfig)

	reIntf := regexp.MustCompile(`(?m)^interface (Po.+)$`)
	reMlag := regexp.MustCompile(`(?m)mlag (\d+)`)

	interfaces := reIntf.FindAllStringSubmatch(config, -1)

	for _, intf := range interfaces {
		parent := "interface " + intf[1]
		config, _ := m.GetBlock(parent)

		match := reMlag.FindStringSubmatch(config)
		if match != nil {
			resource[intf[1]] = match[1]
		}
	}
	return resource
}

// GetSection returns the specified Mlag Entry for the name specified.
//
// Args:
//  name (string): The Mlag name
//
// Returns:
//  Returns string representation of Mlag config entry
func (m *MlagEntity) GetSection() string {
	config, err := m.GetBlock("mlag configuration")
	if err != nil {
		return ""
	}
	return config
}

// ConfigureMlag is a config wrapper for initial mlag config command
//
// Args:
//  cmd (string): command to issue
//  value (string): The value to configure the mlag
//  default (bool): Configures using the default keyword
//
// Returns:
//  bool: True if the commands complete successfully
func (m *MlagEntity) ConfigureMlag(cmd string, value string, def bool, enable bool) bool {
	cfg := m.CommandBuilder(cmd, value, def, enable)
	var commands = []string{"mlag configuration", cfg}
	return m.Configure(commands...)
}

// SetDomainID Configures the mlag domain-id value
//
// Args:
//  value (str): The value to configure the domain-id
//
// Returns:
//  bool: Returns True if the commands complete successfully
func (m *MlagEntity) SetDomainID(value string) bool {
	if value == "" {
		return m.ConfigureMlag("domain-id", value, false, false)
	}
	return m.ConfigureMlag("domain-id", value, false, true)
}

// SetDomainIDDefault Configures the default mlag domain-id value
//
// Returns:
//  bool: Returns True if the commands complete successfully
func (m *MlagEntity) SetDomainIDDefault() bool {
	return m.ConfigureMlag("domain-id", "", true, false)
}

// SetLocalInterface Configures the mlag local-interface value
//
// Args:
//  value (str): The value to configure the local-interface
//
// Returns:
//  bool: Returns True if the commands complete successfully
func (m *MlagEntity) SetLocalInterface(value string) bool {
	if value == "" {
		return m.ConfigureMlag("local-interface", value, false, false)
	}
	return m.ConfigureMlag("local-interface", value, false, true)
}

// SetLocalInterfaceDefault Configures the default mlag local-interface value
//
// Returns:
//  bool: Returns True if the commands complete successfully
func (m *MlagEntity) SetLocalInterfaceDefault() bool {
	return m.ConfigureMlag("local-interface", "", true, false)
}

// SetPeerAddress Configures the mlag peer-address value
//
// Args:
//  value (str): The value to configure the peer-address
//  default (bool): Configures the peer-address using the
//
// Returns:
//  bool: Returns True if the commands complete successfully
func (m *MlagEntity) SetPeerAddress(value string) bool {
	if value == "" {
		return m.ConfigureMlag("peer-address", value, false, false)
	}
	return m.ConfigureMlag("peer-address", value, false, true)
}

// SetPeerAddressDefault Configures the default mlag peer-address value
//
// Returns:
//  bool: Returns True if the commands complete successfully
func (m *MlagEntity) SetPeerAddressDefault() bool {
	return m.ConfigureMlag("peer-address", "", true, false)
}

// SetPeerLink Configures the mlag peer-link value
//
// Args:
//  value (str): The value to configure the peer-link
//  default (bool): Configures the peer-link using the
//                  default keyword
//
// Returns:
//  bool: Returns True if the commands complete successfully
func (m *MlagEntity) SetPeerLink(value string) bool {
	if value == "" {
		return m.ConfigureMlag("peer-link", value, false, false)
	}
	return m.ConfigureMlag("peer-link", value, false, true)
}

// SetPeerLinkDefault Configures the default mlag peer-link value
//
// Returns:
//  bool: Returns True if the commands complete successfully
func (m *MlagEntity) SetPeerLinkDefault() bool {
	return m.ConfigureMlag("peer-link", "", true, false)
}

// SetShutdown Configures the mlag shutdown value
// Args:
//  enable (bool): true for enabled, false for shutdown
//
// Returns:
//  bool: Returns True if the commands complete successfully
func (m *MlagEntity) SetShutdown(enable bool) bool {
	return m.ConfigureMlag("shutdown", "", false, enable)
}

// SetShutdownDefault Configures the mlag default shutdown value
//
// Returns:
//  bool: Returns True if the commands complete successfully
func (m *MlagEntity) SetShutdownDefault() bool {
	return m.ConfigureMlag("shutdown", "", true, false)
}

// SetMlagID Configures the interface mlag value for the specified interface
//
// Args:
//  name (str): The interface to configure.  Valid values for the
//              name arg include Port-Channel*
//  value (str): The mlag identifier to cofigure on the interface
//
// Returns:
//  bool: Returns True if the commands complete successfully
func (m *MlagEntity) SetMlagID(name string, value string) bool {
	var cmd string
	if value == "" {
		cmd = m.CommandBuilder("mlag", value, false, false)
	} else {
		cmd = m.CommandBuilder("mlag", value, false, true)
	}
	var commands = []string{cmd}
	return m.ConfigureInterface(name, commands...)
}

// SetMlagIDDefault Configures the default interface mlag value for the
// specified interface
//
// Args:
//  name (str): The interface to configure.  Valid values for the
//              name arg include Port-Channel*
//  value (str): The mlag identifier to cofigure on the interface
// Returns:
//  bool: Returns True if the commands complete successfully
func (m *MlagEntity) SetMlagIDDefault(name string) bool {
	return m.ConfigureInterface(name, []string{"default mlag"}...)
}
