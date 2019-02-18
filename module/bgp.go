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

// BgpNetworkEntry represents a BGP network config entry
type BgpNetworkEntry map[string]string

// BgpConfig represents the parsed Bgp config for an interface
type BgpConfig struct {
	asNumber     string
	routerID     string
	shutdown     string
	maxPaths     string
	maxEcmpPaths string
	networks     []BgpNetworkEntry
}

// BgpAs returns the BGP config AS number
// Empty string is returned if no value set.
func (b *BgpConfig) BgpAs() string {
	return b.asNumber
}

// RouterID returns the BGP config Router-id
// Empty string is returned if no value set.
func (b *BgpConfig) RouterID() string {
	return b.routerID
}

// Shutdown returns 'true'(string) if shutdown.
// Otherwise 'false'(string)
func (b *BgpConfig) Shutdown() string {
	return b.shutdown
}

// MaximumPaths returns the BGP config value for
// maximum paths. Empty string is returned if no value set.
func (b *BgpConfig) MaximumPaths() string {
	return b.maxPaths
}

// MaximumEcmpPaths returns the configured value for
// BGP max ECMP paths. Empty string is returned if no value set.
func (b *BgpConfig) MaximumEcmpPaths() string {
	return b.maxEcmpPaths
}

// Networks returns a list of configured network statements
// Each entry represents a bgp network entry. Entry formed as
// follows:
// [
// 	  	{
//			"prefix":"",
//			"masklen":"",
//			"route_map":"",
//		},
//		{..
// ]
func (b *BgpConfig) Networks() []BgpNetworkEntry {
	return b.networks
}

// Prefix returns the prefix data for this BgpNetworkEntry
func (b BgpNetworkEntry) Prefix() string {
	return b["prefix"]
}

// MaskLen returns the masklen data for this BgpNetworkEntry
func (b BgpNetworkEntry) MaskLen() string {
	return b["masklen"]
}

// RouteMap returns the routemap data for this BgpNetworkEntry
func (b BgpNetworkEntry) RouteMap() string {
	return b["route_map"]
}

// BGPEntity provides a configuration resource for Bgp
type BGPEntity struct {
	neighbors *BgpNeighborsEntity
	*AbstractBaseEntity
}

// Bgp factory function to initiallize BGPEntity resource
// given a Node
func Bgp(node *goeapi.Node) *BGPEntity {
	return &BGPEntity{AbstractBaseEntity: &AbstractBaseEntity{node}}
}

// Get the BGP Config for the current entity.
// Returns a BgpConfig object
func (b *BGPEntity) Get() *BgpConfig {
	config, err := b.GetBlock(`^router bgp .*`)
	if err != nil {
		return nil
	}
	var bgp = new(BgpConfig)
	bgp.asNumber = b.parseAS(config)
	bgp.routerID = b.parseRouterID(config)
	bgp.shutdown = strconv.FormatBool(b.parseShutdown(config))
	bgp.networks = b.parseNetworks(config)

	m := b.parseMaxPaths(config)
	bgp.maxPaths = m["maximum_paths"]
	bgp.maxEcmpPaths = m["maximum_ecmp_paths"]

	return bgp
}

// Neighbors returns the instance of the BgpNeighborsEntity
// for this BGPEntity
func (b *BGPEntity) Neighbors() *BgpNeighborsEntity {
	if b.neighbors == nil {
		b.neighbors = BgpNeighbors(b.AbstractBaseEntity.node)
	}
	return b.neighbors
}

// parseBgpAS parses the given BGP config for the AS value
func (b *BGPEntity) parseAS(config string) string {
	re := regexp.MustCompile(`(?m)^router bgp (\d+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseRouterID parses the given BGP config for the router-id
func (b *BGPEntity) parseRouterID(config string) string {
	re := regexp.MustCompile(`(?m)router-id ([^\s]+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseMaxPaths parses the given BGP config for maximum paths entries.
// Returned is a map[string]string with keys maximum_paths or maximum_ecmp_paths
//
//  {
//	      "maximum_paths"      : "3",
//	      "maximum_ecmp_paths" : "4",
//  }
func (b *BGPEntity) parseMaxPaths(config string) map[string]string {
	var maxPaths string
	var maxEcmp string
	re := regexp.MustCompile(`maximum-paths\s+(\d+)(?:\s+ecmp\s+(\d+))?`)
	if match := re.FindStringSubmatch(config); match != nil {
		maxPaths = match[1]
		maxEcmp = match[2]
	}
	return map[string]string{
		"maximum_paths":      maxPaths,
		"maximum_ecmp_paths": maxEcmp,
	}
}

// parseShutdown scans the config block and parses the shutdown value
// Returns true (bool) if interface is in shutdown state.
func (b *BGPEntity) parseShutdown(config string) bool {
	matched, _ := regexp.MatchString(`(?m)no shutdown`, config)
	if matched {
		return false
	}
	return !matched
}

// parseNetworks parses the BGP config for any configured network statement.
// Returned is a []map[sring]string...which represents a list of network
// configuration statements each with a mapped entry containing the 'prefix',
// 'masklen', and 'route_map'.
func (b *BGPEntity) parseNetworks(config string) []BgpNetworkEntry {
	reNet := regexp.MustCompile(`(?m)network (.+)/(\d+)(?: route-map (\w+))*`)

	matches := reNet.FindAllStringSubmatch(config, -1)

	if matches == nil {
		return nil
	}

	resource := make([]BgpNetworkEntry, len(matches))

	for idx, line := range matches {
		resource[idx] = BgpNetworkEntry{
			"prefix":    line[1],
			"masklen":   line[2],
			"route_map": line[3],
		}
	}
	return resource
}

// GetSection returns the BGP config section as a string.
func (b *BGPEntity) GetSection() string {
	config, err := b.GetBlock(`^router bgp .*`)
	if err != nil {
		return ""
	}
	return config
}

// ConfigureBgp configures the BGP Entity with the given
// command. Returns true (bool) if the commands complete
// successfully
func (b *BGPEntity) ConfigureBgp(cmd string) bool {
	config := b.Get()
	if config == nil {
		return false
	}
	commands := []string{
		"router bgp " + config.BgpAs(),
		cmd,
	}
	return b.Configure(commands...)
}

// Create creates a BGP instance on the node using the given
// AS value. Returns true(bool) if the commands complete successfully
func (b *BGPEntity) Create(bgpAS int) bool {
	if !(0 < bgpAS && bgpAS < 65535) {
		return false
	}
	cmd := "router bgp " + strconv.Itoa(bgpAS)
	return b.Configure(cmd)
}

// Delete deletes the BGP instance on the node.
// Returns true(bool) if the commands complete successfully
func (b *BGPEntity) Delete() bool {
	config := b.Get()
	if config == nil {
		return true
	}
	cmd := "no router bgp " + config.BgpAs()
	return b.Configure(cmd)
}

// Default sets the default config for BGP instance on the node.
// Returns true(bool) if the commands complete successfully
func (b *BGPEntity) Default() bool {
	config := b.Get()
	if config == nil {
		return true
	}
	cmd := "default router bgp " + config.BgpAs()
	return b.Configure(cmd)
}

// SetRouterID configures the router-id using the provided value.
// Returns true(bool) if the commands complete successfully
func (b *BGPEntity) SetRouterID(value string) bool {
	if value == "" {
		return b.ConfigureBgp("no router-id")
	}
	return b.ConfigureBgp("router-id " + value)
}

// SetRouterIDDefault sets the default router-id value
// Returns true(bool) if the commands complete successfully
func (b *BGPEntity) SetRouterIDDefault() bool {
	return b.ConfigureBgp("default router-id")
}

// SetMaximumPaths sets the BGP maximum path using the provided maxPath
// value. Returns true(bool) if the commands complete successfully
func (b *BGPEntity) SetMaximumPaths(maxPath int) bool {
	command := "maximum-paths " + strconv.Itoa(maxPath)
	return b.ConfigureBgp(command)
}

// SetMaximumPathsWithEcmp set the BGP maximum path / max Ecmp configuration
// for this entity.
// Returns true(bool) if the commands complete successfully
func (b *BGPEntity) SetMaximumPathsWithEcmp(maxPath int, maxEcmp int) bool {
	cmd := "maximum-paths " + strconv.Itoa(maxPath) + " ecmp " + strconv.Itoa(maxEcmp)
	return b.ConfigureBgp(cmd)
}

// SetMaximumPathsDefault resets the maximum paths configuration to its
// default values
// Returns true(bool) if the commands complete successfully
func (b *BGPEntity) SetMaximumPathsDefault() bool {
	return b.ConfigureBgp("default maximum-paths")
}

// SetShutdown configures this BGP entity to be 'shutdown' (true), or
// 'no shutdown' (false)
// Returns true(bool) if the commands complete successfully
func (b *BGPEntity) SetShutdown(enable bool) bool {
	cmd := b.CommandBuilder("shutdown", "", false, enable)
	return b.ConfigureBgp(cmd)
}

// SetShutdownDefault configures the default shutdown configuration
// for this BGPEntity
// Returns true(bool) if the commands complete successfully
func (b *BGPEntity) SetShutdownDefault() bool {
	return b.ConfigureBgp("default shutdown")
}

// AddNetworkWithRouteMap configures BGP network using supplied network
// prefix, mask length, and route-map
// Returns true(bool) if the commands complete successfully
func (b *BGPEntity) AddNetworkWithRouteMap(prefix string, maskLen string, routeMap string) bool {
	command := "network " + prefix + "/" + maskLen
	if routeMap != "" {
		command = command + " route-map " + routeMap
	}
	return b.ConfigureBgp(command)
}

// AddNetwork configures BGP network using supplied network prefix and
// mask length.
// Returns true(bool) if the commands complete successfully
func (b *BGPEntity) AddNetwork(prefix string, maskLen string) bool {
	return b.AddNetworkWithRouteMap(prefix, maskLen, "")
}

// RemoveNetworkWithRouteMap removes the configured BGP network config
// using the supplied network prefix, mask length, and route-map
// Returns true(bool) if the commands complete successfully
func (b *BGPEntity) RemoveNetworkWithRouteMap(prefix string, maskLen string, routeMap string) bool {
	command := "no network " + prefix + "/" + maskLen
	if routeMap != "" {
		command = command + " route-map " + routeMap
	}
	return b.ConfigureBgp(command)
}

// RemoveNetwork removes the configured BGP network config
// using the supplied network prefix and mask length
// Returns true(bool) if the commands complete successfully
func (b *BGPEntity) RemoveNetwork(prefix string, maskLen string) bool {
	return b.RemoveNetworkWithRouteMap(prefix, maskLen, "")
}

// BgpNeighborConfig represents the parsed Bgp neighbor config
// {
//		"peer_group"		: "peer1",
//		"remote_as"			: "99",
//		"send_community"	: "true",
//		"shutdown"			: "false",
//		"description"		: "This is a bgp entity",
//		"next_hop_self"		: "1.1.1.1",
//		"route_in_map"		: "in-map",
//		"route_out_map"		: "out-map",
// }
type BgpNeighborConfig map[string]string

// BgpNeighborCollection is a collection of BgpNeighborConfigs.
// Each key entry of the collection is a unique neighbor(key:string)
// mapping to its respective BGPNeighborConfig:
// Example:
//	{
//		"172.16.10.1" : BgpNeighborConfig
//					{
//						"peer_group"    : "",
//						"remote_as"     : "",
//						...
//					},
//	}
type BgpNeighborCollection map[string]BgpNeighborConfig

// BgpNeighborsEntity provides a configuration resource for Bgp
// neighbors
type BgpNeighborsEntity struct {
	*AbstractBaseEntity
}

// BgpNeighbors factory function to initiallize BgpNeighborsEntity resource
// given a Node
func BgpNeighbors(node *goeapi.Node) *BgpNeighborsEntity {
	return &BgpNeighborsEntity{&AbstractBaseEntity{node}}
}

// Get the BGP Neighbot Config for the current entity.
// Returns a BgpNeighborConfig object
func (b *BgpNeighborsEntity) Get(name string) BgpNeighborConfig {
	config, _ := b.GetBlock(`^router bgp .*`)

	return BgpNeighborConfig{
		"peer_group":     b.parsePeerGroup(config, name),
		"remote_as":      b.parseRemoteAS(config, name),
		"send_community": strconv.FormatBool(b.parseSendCommunity(config, name)),
		"shutdown":       strconv.FormatBool(b.parseShutdown(config, name)),
		"description":    b.parseDescription(config, name),
		"next_hop_self":  strconv.FormatBool(b.parseNextHopSelf(config, name)),
		"route_in_map":   b.parseRouteMapIn(config, name),
		"route_out_map":  b.parseRouteMapOut(config, name),
	}
}

// GetAll returns the BGP Neighbor Collection
func (b *BgpNeighborsEntity) GetAll() BgpNeighborCollection {
	config, _ := b.GetBlock(`^router bgp .*`)
	if config == "" {
		return nil
	}

	re := regexp.MustCompile(`(?m)^\s+neighbor ([^\s]+)`)
	neighbors := re.FindAllStringSubmatch(config, -1)

	// hash map to mark previously seen neighbor entries
	hash := map[string]struct{}{}
	collection := make(BgpNeighborCollection)

	for _, matchedLine := range neighbors {
		neighbor := matchedLine[1]
		if _, found := hash[neighbor]; !found {
			hash[neighbor] = struct{}{}
			collection[neighbor] = b.Get(neighbor)
		}
	}
	return collection
}

// parsePeerGroup parses the provided BGP config(string) looking for neighbor entries with
// given name(string) and returns the peer-group label(string) associated with the neighbor.
// Empty string is returned on no peer-goup seen
func (b *BgpNeighborsEntity) parsePeerGroup(config string, name string) string {
	regex := `neighbor ` + name + ` peer-group ([^\s]+)`
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseRemoteAS parses the provided BGP config looking for neighbor entries with
// given name(string) and returns the remote AS(string) associated with the neighbor.
// Empty string is returned on no remote-as seen
func (b *BgpNeighborsEntity) parseRemoteAS(config string, name string) string {
	regex := `neighbor ` + name + ` remote-as (\d+)`
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseSendCommunity parses the provided BGP config looking for neighbor entries with
// given name(string) and returns true(bool) on neighbor send-community configured.
func (b *BgpNeighborsEntity) parseSendCommunity(config string, name string) bool {
	regex := "neighbor " + name + " send-community"
	matched, _ := regexp.MatchString(regex, config)
	return matched
}

// parseShutdown parses the provided BGP config looking for neighbor entries with
// given name(string) and returns true(bool) on neighbor shutdown configured.
func (b *BgpNeighborsEntity) parseShutdown(config string, name string) bool {
	regex := "no neighbor " + name + " shutdown"
	matched, _ := regexp.MatchString(regex, config)
	return !matched
}

// parseDescription parses the provided BGP config looking for neighbor entries with
// given name(string) and returns the description(string) associated with the neighbor.
// Empty string is returned on no description seen
func (b *BgpNeighborsEntity) parseDescription(config string, name string) string {
	regex := `(?m)neighbor ` + name + ` description (.*)$`
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseNextHopSelf parses the provided BGP config looking for neighbor entries with
// given name(string) and returns true(bool) if next-hop-self configured.
func (b *BgpNeighborsEntity) parseNextHopSelf(config string, name string) bool {
	regex := "neighbor " + name + " next-hop-self"
	matched, _ := regexp.MatchString(regex, config)
	return matched
}

// parseRouteMapIn parses the provided BGP config looking for neighbor entries with
// given name(string) and returns the inbound route-map reference(string) associated
// with the neighbor. Empty string is returned on no peer-goup seen
func (b *BgpNeighborsEntity) parseRouteMapIn(config string, name string) string {
	regex := `(?m)neighbor ` + name + ` route-map ([^\s]+) in`
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// parseRouteMapOut parses the provided BGP config looking for neighbor entries with
// given name(string) and returns the outbound route-map reference(string) associated
// with the neighbor. Empty string is returned on no peer-goup seen
func (b *BgpNeighborsEntity) parseRouteMapOut(config string, name string) string {
	regex := `(?m)neighbor ` + name + ` route-map ([^\s]+) out`
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// Create creates a neighbor entry in the shutdown state.
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) Create(name string) bool {
	return b.SetShutdown(name, true)
}

// Delete removes the neighbor name(string) entry.
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) Delete(name string) bool {
	resp := b.Configure("no neighbor " + name)
	if !resp {
		resp = b.Configure("no neighbor " + name + " peer-group")
	}
	return resp
}

// Configure (redefined from base) Configures router bgp instance.
// Returns true(bool) if the commands complete successfully, false if
// configure fails or device BGP instance doesn't exsist.
func (b *BgpNeighborsEntity) Configure(cmd string) bool {
	config, _ := b.GetBlock(`^router bgp .*`)
	re := regexp.MustCompile(`(?m)^router bgp (\d+)`)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return false
	}
	commands := []string{
		"router bgp " + match[1],
		cmd,
	}
	return b.AbstractBaseEntity.Configure(commands...)
}

// CommandBuilder (redefined from base) Builds proper bgp neighbot
// configuration command based on provided arguments:
// 	name(string) - neighbor name
// 	cmd(string) - command to use
// 	value(string) - additional parameters for command
// 	def(bool) - If true, Default configuration needed
// 	shut(bool) - If true, negate configure command
// Returns compiled command. If def is 'true', command is returned
// with 'default' prepended. If shut is 'true', command is returned
// with 'no' prepended.
func (b *BgpNeighborsEntity) CommandBuilder(name string, cmd string,
	value string, def bool, shut bool) string {
	str := "neighbor " + name + " " + cmd
	return b.AbstractBaseEntity.CommandBuilder(str, value, def, shut)
}

// SetPeerGroup sets the neighbor(string) peer-group value(string)
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetPeerGroup(name string, value string) bool {
	if net.ParseIP(name) == nil {
		return false
	}
	cmd := b.CommandBuilder(name, "peer-group", value, false, true)
	return b.Configure(cmd)
}

// SetPeerGroupDefault sets the default configuration value for neighbor
// peer-group configuration
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetPeerGroupDefault(name string) bool {
	if net.ParseIP(name) == nil {
		return false
	}
	cmd := b.CommandBuilder(name, "peer-group", "", true, false)
	return b.Configure(cmd)
}

// SetRemoteAS sets the neighbor name(string) remote-as configuration to
// value(string)
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetRemoteAS(name string, value string) bool {
	cmd := b.CommandBuilder(name, "remote-as", value, false, true)
	return b.Configure(cmd)
}

// SetRemoteASDefault sets the default configuration value for the neighbor
// remote-as configuration
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetRemoteASDefault(name string) bool {
	cmd := b.CommandBuilder(name, "remote-as", "", true, false)
	return b.Configure(cmd)
}

// SetShutdown set the neighbor name(string) shutdown state to
//	shut(boo) - true:shutdown,    false:no shutdown
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetShutdown(name string, shut bool) bool {
	cmd := b.CommandBuilder(name, "shutdown", "", false, shut)
	return b.Configure(cmd)
}

// SetShutdownDefault sets the default configuration value for the neighbor
// shutdown configuration
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetShutdownDefault(name string) bool {
	cmd := b.CommandBuilder(name, "shutdown", "", true, false)
	return b.Configure(cmd)
}

// SetSendCommunity sets the neighbor name(string) send-community configuration to
// value(string).
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetSendCommunity(name string, enable bool) bool {
	cmd := b.CommandBuilder(name, "send-community", "", false, enable)
	return b.Configure(cmd)
}

// SetSendCommunityDefault sets the default configuration value for the neighbor
// send-community configuration
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetSendCommunityDefault(name string) bool {
	cmd := b.CommandBuilder(name, "send-community", "", true, false)
	return b.Configure(cmd)
}

// SetNextHopSelf sets the neighbor name(string) next-hop-self to enabled(true)
// or disabled(false)
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetNextHopSelf(name string, enabled bool) bool {
	cmd := b.CommandBuilder(name, "next-hop-self", "", false, enabled)
	return b.Configure(cmd)
}

// SetNextHopSelfDefault sets the default configuration value for the neighbor
// next-hop-self configuration
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetNextHopSelfDefault(name string) bool {
	cmd := b.CommandBuilder(name, "next-hop-self", "", true, false)
	return b.Configure(cmd)
}

// SetRouteMapIn sets the neighbor name(string) inbound route-map entry using
// value(string)
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetRouteMapIn(name string, value string) bool {
	cmd := b.CommandBuilder(name, "route-map", value, false, true)
	cmd = cmd + " in"
	return b.Configure(cmd)
}

// SetRouteMapInDefault sets the default configuration value for the neighbor
// inbound route-map configuration
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetRouteMapInDefault(name string) bool {
	cmd := b.CommandBuilder(name, "route-map", "", true, false)
	cmd = cmd + " in"
	return b.Configure(cmd)
}

// SetRouteMapOut sets the neighbor name(string) outbound route-map entry using
// value(string)
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetRouteMapOut(name string, value string) bool {
	cmd := b.CommandBuilder(name, "route-map", value, false, true)
	cmd = cmd + " out"
	return b.Configure(cmd)
}

// SetRouteMapOutDefault sets the default configuration value for the neighbor
// outbound route-map configuration
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetRouteMapOutDefault(name string) bool {
	cmd := b.CommandBuilder(name, "route-map", "", true, false)
	cmd = cmd + " out"
	return b.Configure(cmd)
}

// SetDescription sets the neighbor name(string) using the provided value(string)
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetDescription(name string, value string) bool {
	cmd := b.CommandBuilder(name, "description", value, false, true)
	return b.Configure(cmd)
}

// SetDescriptionDefault sets the default configuration value for the neighbor
// description configuration
// Returns true(bool) if the commands complete successfully
func (b *BgpNeighborsEntity) SetDescriptionDefault(name string) bool {
	cmd := b.CommandBuilder(name, "description", "", true, false)
	return b.Configure(cmd)
}

type ShowIPBGPSummary struct {
	VRFs map[string]VRF
}

type VRF struct {
	RouterID string                        `json:"routerId"`
	Peers    map[string]BGPNeighborSummary `json:"peers"`
	VRF      string                        `json:"vrf"`
	ASN      int64                         `json:"asn"`
}

type BGPNeighborSummary struct {
	MsgSent             int     `json:"msgSent"`
	InMsgQueue          int     `json:"inMsgQueue"`
	PrefixReceived      int     `json:"prefixReceived"`
	UpDownTime          float64 `json:"upDownTime"`
	Version             int     `json:"version"`
	MsgReceived         int     `json:"msgReceived"`
	PrefixAccepted      int     `json:"prefixAccepted"`
	PeerState           string  `json:"peerState"`
	PeerStateIdleReason string  `json:"peerStateIdleReason,omitempty"`
	OutMsgQueue         int     `json:"outMsgQueue"`
	UnderMaintenance    bool    `json:"underMaintenance"`
	ASN                 int64   `json:"asn"`
}

func (b *ShowIPBGPSummary) GetCmd() string {
	return "show ip bgp summary"
}

func (s *ShowEntity) ShowIPBGPSummary() (ShowIPBGPSummary, error) {
	handle, _ := s.node.GetHandle("json")
	var showipbgpsummary ShowIPBGPSummary
	handle.AddCommand(&showipbgpsummary)

	if err := handle.Call(); err != nil {
		return showipbgpsummary, err
	}

	handle.Close()
	return showipbgpsummary, nil
}
