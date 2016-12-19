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

// SystemEntity provides a configuration resource for System
type SystemEntity struct {
	*AbstractBaseEntity
}

// SystemConfig represents a parsed system entry containing
type SystemConfig map[string]string

// HostName returns the hostname entry in the SystemConfig
// type
func (s SystemConfig) HostName() string {
	return s["hostname"]
}

// IPRouting returns the ip route entry in the SystemConfig
// type
func (s SystemConfig) IPRouting() string {
	return s["iprouting"]
}

// System factory function to initiallize System resource
// given a Node
func System(node *goeapi.Node) *SystemEntity {
	return &SystemEntity{&AbstractBaseEntity{node}}
}

// Get Returns the system configuration abstraction
//
// The System resource returns the following:
//
//    * hostname (str): The hostname value
//
// Returns:
//  SystemConfig: Represents the node's system configuration
func (s *SystemEntity) Get() SystemConfig {
	config := s.Config()
	var resource = make(SystemConfig)
	resource["hostname"] = parseHostname(config)
	resource["iprouting"] = strconv.FormatBool(parseIPRouting(config))
	return resource
}

// parseHostname Parses the global config and returns the hostname value
//
// Args:
//  config (string): running config
// Returns:
//  string: The configured value for hostname.
func parseHostname(config string) string {
	hostName := regexp.MustCompile(`(?m)^hostname ([^\s]+)$`)
	match := hostName.FindStringSubmatch(config)
	if match != nil {
		return match[1]
	}
	return "localhost"
}

// parseHostname Parses the global config and returns the hostname value
//
// Returns:
//  string: The configured value for hostname.
func (s *SystemEntity) parseHostname() string {
	config := s.Config()
	return parseHostname(config)
}

// parseIpRouting Parses the global config and returns the ip routing value
//
// Args:
//  config (string): running config
// Returns:
//  string: The configure value for ip routing.
func parseIPRouting(config string) bool {
	if config == "" {
		return false
	}
	regex := regexp.MustCompile(`no ip routing`)
	return !(regex.MatchString(config))
}

// parseIpRouting Parses the global config and returns the ip routing value
//
// Returns:
//  string: The configure value for ip routing.
func (s *SystemEntity) parseIPRouting() bool {
	config := s.Config()
	return parseIPRouting(config)
}

// SetHostname Configures the global system hostname setting
//
// EosVersion:
//    4.13.7M
//
// Args:
//  value (str): The hostname value
//  default (bool): Controls use of the default keyword
//
// Returns:
//  bool: True if the commands are completed successfully
func (s *SystemEntity) SetHostname(hostname string) bool {
	if hostname == "" {
		return s.Configure("no hostname")
	}
	return s.Configure("hostname " + hostname)
}

// SetHostnameDefault Configures the global default system hostname setting
//
// EosVersion:
//    4.13.7M
//
// Returns:
//  bool: True if the commands are completed successfully
func (s *SystemEntity) SetHostnameDefault() bool {
	return s.Configure("default hostname")
}

// SetIPRouting Configures the state of global ip routing
//
// EosVersion:
//    4.13.7M
//
// Args:
//  value(bool): True if ip routing should be enabled or False if
//               ip routing should be disabled
//
// Returns:
//  bool: True if the commands completed successfully otherwise False
func (s *SystemEntity) SetIPRouting(value string, enable bool) bool {
	cmd := s.CommandBuilder("ip routing", value, false, enable)
	return s.Configure(cmd)
}

// SetIPRoutingDefault Configures the default tate of global ip routing
//
// EosVersion:
//    4.13.7M
//
// Returns:
//  bool: True if the commands completed successfully otherwise False
func (s *SystemEntity) SetIPRoutingDefault(value string) bool {
	cmd := s.CommandBuilder("ip routing", value, true, false)
	return s.Configure(cmd)
}
