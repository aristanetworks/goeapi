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
	"net"
	"regexp"
	"strconv"

	"github.com/aristanetworks/goeapi"
)

// AclEntry represents a parsed Acl entry of the form
//  {
//      "action"  : ""
//      "srcaddr" : ""
//      "srclen"  : ""
//      "log"     : ""
//  }
type AclEntry map[string]string

// AclEntryMap is a key/value
type AclEntryMap map[string]AclEntry

// AclConfig represents an Acl Config entry with a name, type,
// and individual entries in AclEntryMap
type AclConfig struct {
	aclName    string
	aclType    string
	aclEntries AclEntryMap
}

// Name returns a string of the ACL Config
// name
func (a AclConfig) Name() string {
	return a.aclName
}

// Type returns a string of the ACL type
func (a AclConfig) Type() string {
	return a.aclType
}

// Entries returns a map[string] of AclConfigs
// each keyed entry being the name/label of the
// AclConfig
func (a AclConfig) Entries() AclEntryMap {
	return a.aclEntries
}

// Action returns the action for the AclEntry
func (a AclEntry) Action() string {
	return a["action"]
}

// SrcAddr returns the source address for the AclEntry
func (a AclEntry) SrcAddr() string {
	return a["srcaddr"]
}

// SrcLen returns the source mask length for the AclEntry
func (a AclEntry) SrcLen() string {
	return a["srclen"]
}

// Log returns log is configured for the AclEntry
func (a AclEntry) Log() string {
	return a["log"]
}

// AclEntity provides a configuration resource for Acl
type AclEntity struct {
	*AbstractBaseEntity
}

// Acl factory function to initiallize Acl resource
// given a Node
func Acl(node *goeapi.Node) *AclEntity {
	return &AclEntity{&AbstractBaseEntity{node}}
}

// maskToPrefixlen Converts a subnet mask from dotted decimal to bit length
// If mask is not in canonical form (i.e. ones followed by zeros) then
// returns 0
func maskToPrefixlen(mask string) string {
	if mask == "" {
		mask = "255.255.255.255"
	}
	prefixSize, _ := net.IPMask(net.ParseIP(mask).To4()).Size()
	return strconv.Itoa(prefixSize)
}

// prefixlenToMask Converts a prefix length to a dotted decimal subnet mask
func prefixlenToMask(prefixlen string) string {
	if prefixlen == "" {
		prefixlen = "32"
	}
	addr := "0.0.0.0/" + prefixlen
	_, ipNet, _ := net.ParseCIDR(addr)
	return net.IP(ipNet.Mask).String()
}

// Get returns the specified AclEntity from the nodes current configuration.
//
// Args:
//  name (string): The ACL name
//
// Returns:
//  Returns AclConfig object
func (a *AclEntity) Get(name string) (*AclConfig, error) {
	parent := "ip access-list standard " + name
	config, err := a.GetBlock(parent)
	if err != nil {
		return nil, err
	}
	return &AclConfig{
		aclName:    name,
		aclType:    "standard",
		aclEntries: a.parseEntries(config),
	}, nil
}

// GetAll returns the collection of ACLs from the nodes running
// configuration as a hash. The ACL resource collection hash is
// keyed by the ACL name.
//
// Returns:
// Returns a hash that represents the entire ACL collection from
//  the nodes running configuration. If there are no ACLs configured,
//  this method will return an empty hash.
func (a *AclEntity) GetAll() map[string]*AclConfig {
	config := a.Config()
	re := regexp.MustCompile(`(?m)ip access-list standard ([^\s]+)`)
	matches := re.FindAllStringSubmatch(config, -1)
	aclConfigs := make(map[string]*AclConfig)

	for _, acl := range matches {
		name := acl[1]
		aclConfigs[name], _ = a.Get(name)
	}
	return aclConfigs
}

// defined constants for accessing matched entries
const (
	seqnum = iota + 1
	action
	anyip
	host
	ip
	mlen
	mask
	log
)

// parseEntries scans the nodes configurations and parses
// the entries within an AclEntity.
//
// Args:
//  config (string): The switch config.
//
// Return:
func (a *AclEntity) parseEntries(config string) AclEntryMap {
	//        1    2     3     4    5    6     7     8
	//      (seq, act, anyip, host, ip, mlen, mask, log)
	entryRegex := regexp.MustCompile(`(\d+)` +
		`(?: ([p|d]\w+))` +
		`(?: (any))?` +
		`(?: (host))?` +
		`(?: ([0-9]+(?:\.[0-9]+){3}))?` +
		`(?:/([0-9]{1,2}))?` +
		`(?: ([0-9]+(?:\.[0-9]+){3}))?` +
		`(?: (log))?`)

	itemRegex := regexp.MustCompile(`(?m)\d+ [p|d].*$`)
	matches := itemRegex.FindAllString(config, -1)

	entries := make(AclEntryMap)

	for _, item := range matches {
		match := entryRegex.FindStringSubmatch(item)
		if match == nil {
			continue
		}
		result := make(AclEntry)
		result["action"] = match[action]

		if result["srcaddr"] = match[ip]; match[ip] == "" {
			result["srcaddr"] = "0.0.0.0"
		}

		if result["srclen"] = match[mlen]; match[mlen] == "" {
			result["srclen"] = maskToPrefixlen(match[mask])
		}
		result["log"] = match[log]
		entries[match[seqnum]] = result
	}
	return entries
}

// GetSection returns the specified Acl Entry for the name specified.
//
// Args:
//  name (string): The ACL name
//
// Returns:
//  Returns string representation of Acl config entry
func (a *AclEntity) GetSection(name string) string {
	parent := "ip access-list standard " + name
	config, err := a.GetBlock(parent)
	if err != nil {
		return ""
	}
	return config
}

// Create will create a new ACL resource in the nodes current
// configuration with the specified name.  If the create method
// is called and the ACL already exists, this method will still
// return true. The ACL will not have any entries. Use add_entry
// to add entries to the ACL.
//
//  EosVersion
//      4.13.7M
//
// Args:
//      name (string): The ACL name to create on the node. Must begin
//                 with an alphabetic character. Cannot contain spaces or
//                 quotation marks.
//
// Returns:
//  returns true if the command completed successfully
func (a *AclEntity) Create(name string) bool {
	var commands = []string{"ip access-list standard " + name}
	return a.Configure(commands...)
}

// Delete will delete an existing ACL resource from the nodes current
// running configuration.  If the delete method is called and the ACL
// does not exist, this method will succeed.
//
//  EosVersion
//      4.13.7M
//
//  Args:
//      name (string): The ACL name to delete on the node.
//
// Returns:
//  returns true if the command completed successfully
func (a *AclEntity) Delete(name string) bool {
	var commands = []string{"no ip access-list standard " + name}
	return a.Configure(commands...)
}

// Default will configure the ACL using the default keyword.  This
// command has the same effect as deleting the ACL from the nodes
// running configuration.
//
//  EosVersion
//      4.13.7M
//
// Args:
//  name (string): The ACL name to set to the default value
//                 on the node.
//
// Returns:
//  returns true if the command complete successfully
func (a *AclEntity) Default(name string) bool {
	var commands = []string{"default ip access-list standard " + name}
	return a.Configure(commands...)
}

// UpdateEntry will update an entry, identified by the seqno
// in the ACL specified by name, with the passed in parameters.
//
//  EosVersion
//      4.13.7M
//
//  name (string): The ACL name to update on the node.
//  seqno (string): The sequence number of the entry in the ACL to update.
//  action (string): The action triggered by the ACL. Valid
//                   values are 'permit', 'deny', or 'remark'
//  addr (string): The IP address to permit or deny.
//  prefixlen (string):  The prefixlen for the IP address.
//  log (bool): Triggers an informational log message to the console
//              about the matching packet.
//
// Returns:
//  returns true if the command complete successfully
func (a *AclEntity) UpdateEntry(name string, seqno string, action string, addr string,
	prefixlen string, log bool) bool {

	commands := []string{"ip access-list standard " + name}
	commands = append(commands, "no "+seqno)

	entry := seqno + " " + action + " " + addr + "/" + prefixlen
	if log {
		entry = entry + " log"
	}
	commands = append(commands, entry)
	commands = append(commands, "exit")
	return a.Configure(commands...)
}

// AddEntry will add an entry to the specified ACL with the
// passed in parameters.
//
//  EosVersion
//      4.13.7M
//
//  name (string): The ACL name to update on the node.
//  action (string): The action triggered by the ACL. Valid
//                   values are 'permit', 'deny', or 'remark'
//  addr (string): The IP address to permit or deny.
//  prefixlen (string):  The prefixlen for the IP address.
//  log (bool): Triggers an informational log message to the console
//              about the matching packet.
//
// Returns:
//  returns true if the command complete successfully
func (a *AclEntity) AddEntry(name string, action string, addr string,
	prefixlen string, log bool) bool {

	commands := []string{"ip access-list standard " + name}
	entry := action + " " + addr + "/" + prefixlen
	if log {
		entry = entry + " log"
	}
	commands = append(commands, entry)
	commands = append(commands, "exit")
	return a.Configure(commands...)
}

// RemoveEntry will remove the entry specified by the seqno for
// the ACL specified by name.
//
//  EosVersion:
//      4.13.7M
//
// Args:
//  name (string): The ACL name to update on the node.
//  seqno (int): The sequence number of the entry in the ACL to remove.
//
// Returns:
//  returns true if the command complete successfully
func (a *AclEntity) RemoveEntry(name string, seqno int) bool {
	var commands = []string{
		"ip access-list standard " + name,
		"no " + strconv.Itoa(seqno),
		"exit",
	}
	return a.Configure(commands...)
}
