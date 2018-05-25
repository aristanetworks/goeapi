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
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/aristanetworks/goeapi"
)

/**
 *****************************************************************************
 * Unit Tests
 *****************************************************************************
 **/
func TestBgpParseAS_UnitTest(t *testing.T) {
	var b BGPEntity

	shortConfig := `
!
ip routing
!
router bgp %s
   router-id 1.1.1.1
   neighbor test peer-group
   neighbor test remote-as 65001
   neighbor test maximum-routes 12000
   neighbor test1 peer-group
   neighbor test1 route-map RM-IN in
   neighbor test1 route-map RM-OUT out
   neighbor test1 maximum-routes 12000
   neighbor 172.16.10.1 remote-as 65000
   neighbor 172.16.10.1 maximum-routes 12000
   neighbor 172.16.10.1 peer-group test
   network 172.16.10.0/24
   network 172.17.0.0/16
!
!`
	tests := [10]struct {
		in   string
		want string
	}{
		{"", ""},
	}

	// for each test entry
	for idx := 1; idx < len(tests); idx++ {
		// get the random strings
		asNum := strconv.Itoa(RandomInt(1, 65535))
		tests[idx].in = asNum
		tests[idx].want = asNum
	}
	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		got := BgpParseAS(&b, testConfig)
		if strings.Compare(got, tt.want) != 0 {
			t.Fatalf("parseAS() = %q; want %q", got, tt.want)
		}
	}
}

func TestBgpParseRouterID_UnitTest(t *testing.T) {
	var b BGPEntity

	shortConfig := `
!
ip routing
!
router bgp 6500
   router-id %s
   neighbor test peer-group
   neighbor test remote-as 65001
   neighbor test maximum-routes 12000
   neighbor test1 peer-group
   neighbor test1 route-map RM-IN in
   neighbor test1 route-map RM-OUT out
   neighbor test1 maximum-routes 12000
   neighbor 172.16.10.1 remote-as 65000
   neighbor 172.16.10.1 maximum-routes 12000
   neighbor 172.16.10.1 peer-group test
   network 172.16.10.0/24
   network 172.17.0.0/16
!
!`
	tests := [10]struct {
		in   string
		want string
	}{
		{"", ""},
	}

	// for each test entry
	for idx := 1; idx < len(tests); idx++ {
		// get the random ip's
		ipAddr := RandomIPAddress()
		tests[idx].in = ipAddr
		tests[idx].want = ipAddr
	}
	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		got := BgpParseRouterID(&b, testConfig)
		if strings.Compare(got, tt.want) != 0 {
			t.Fatalf("parseRouterID() = %q; want %q", got, tt.want)
		}
	}
}

func TestBgpParseMaxPaths_UnitTest(t *testing.T) {
	var b BGPEntity

	shortConfig := `
!
ip routing
!
router bgp 6500
   no shutdown
   no bgp route-reflector preserve-attributes
   %s
   no bgp additional-paths install
   no bgp tie-break-on-cluster-list-length
   no bgp advertise-inactive
!
!`
	tests := [...]struct {
		in   string
		want string
	}{
		{"maximum-paths 1 ecmp 128", "maximum-paths 1 ecmp 128"},
		{"maximum-paths 4 ecmp 4", "maximum-paths 4 ecmp 4"},
		{"maximum-paths 8 ecmp 12", "maximum-paths 8 ecmp 12"},
		{"maximum-paths 12 ecmp 128", "maximum-paths 12 ecmp 128"},
		{"maximum-paths 24 ecmp 64", "maximum-paths 24 ecmp 64"},
		{"maximum-paths 64 ecmp 128", "maximum-paths 64 ecmp 128"},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		got := BgpParseMaxPaths(&b, testConfig)
		str := "maximum-paths " + got["maximum_paths"] + " ecmp " + got["maximum_ecmp_paths"]
		if strings.Compare(str, tt.want) != 0 {
			t.Fatalf("parseMaxPaths() = %q; want %q", str, tt.want)
		}
	}
}

func TestBgpParseShutdown_UnitTest(t *testing.T) {
	var b BGPEntity

	shortConfig := `
!
ip routing
!
router bgp 6500
   %s
   no router-id
   bgp convergence time 300
   bgp convergence slow-peer time 90
   no bgp confederation identifier
   no update wait-for-convergence
   no update wait-install
   bgp log-neighbor-changes
   bgp default ipv4-unicast
   no bgp default ipv6-unicast
   timers bgp 60 180
   distance bgp 200 200 200
!
!`
	tests := [...]struct {
		in   string
		want bool
	}{
		{"no shutdown", false},
		{"shutdown", true},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := BgpParseShutdown(&b, testConfig); got != tt.want {
			t.Fatalf("parseShutdown() = %t; want %t", got, tt.want)
		}
	}
}

func TestBgpParseNetworks_UnitTest(t *testing.T) {
	var b BGPEntity

	shortConfig := `
!
ip routing
!
router bgp 6400
   router-id 1.1.1.1
   neighbor test peer-group
   neighbor test remote-as 65001
   neighbor test maximum-routes 12000
   neighbor test1 peer-group
   neighbor test1 route-map RM-IN in
   neighbor test1 route-map RM-OUT out
   neighbor test1 maximum-routes 12000
   neighbor 172.16.10.1 remote-as 65000
   neighbor 172.16.10.1 maximum-routes 12000
   neighbor 172.16.10.1 peer-group test
   network 172.16.10.0/24
   network 172.16.20.0/24 route-map Test1
   network 172.16.30.0/24 route-map Test2
   network 172.17.0.0/16
!
!`

	tests := [...]struct {
		in   []string
		want map[string]string
	}{
		{[]string{"172.16.10.0", "24", ""},
			map[string]string{"prefix": "172.16.10.0", "masklen": "24", "route_map": ""}},
		{[]string{"172.16.20.0", "24", "Test1"},
			map[string]string{"prefix": "172.16.20.0", "masklen": "24", "route_map": "Test1"}},
		{[]string{"172.16.30.0", "24", "Test2"},
			map[string]string{"prefix": "172.16.30.0", "masklen": "24", "route_map": "Test2"}},
		{[]string{"172.17.0.0", "16", ""},
			map[string]string{"prefix": "172.17.0.0", "masklen": "16", "route_map": ""}},
	}

	got := BgpParseNetworks(&b, shortConfig)
	for idx, tt := range tests {
		for k, v := range tt.want {
			if got[idx][k] != v {
				t.Fatalf("parseNetworks() = %q; want %q", got[idx], tt.want)
			}
		}
	}
	shortConfig = `
!
ip routing
!
router bgp 6400
   router-id 1.1.1.1
   neighbor test peer-group
   neighbor test remote-as 65001
   neighbor test maximum-routes 12000
   neighbor test1 peer-group
   neighbor test1 route-map RM-IN in
   neighbor test1 route-map RM-OUT out
   neighbor test1 maximum-routes 12000
   neighbor 172.16.10.1 remote-as 65000
   neighbor 172.16.10.1 maximum-routes 12000
   neighbor 172.16.10.1 peer-group test
!
`
	if got = BgpParseNetworks(&b, shortConfig); got != nil {
		t.Fatalf("parseNetworks(). No network statements should return \"\"")
	}

}

func TestBgpGet_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)
	config := bgp.Get()
	if config.BgpAs() != "65000" || config.RouterID() != "1.1.1.1" ||
		config.Shutdown() != "true" || config.MaximumPaths() != "" ||
		config.MaximumEcmpPaths() != "" || config.Networks() == nil {
		t.Fatalf("Invalid result from Get(): %#v", config)
	}
}

func TestBgpNetworkGetters_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)
	config := bgp.Get()
	networks := config.Networks()
	if networks == nil || networks[0].Prefix() != "172.16.10.0" ||
		networks[0].MaskLen() != "24" || networks[0].RouteMap() != "" {
		t.Fatalf("Invalid result from Networks(): %#v", networks)
	}
}

func TestBgpNeigborsInstance_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)
	n := bgp.Neighbors()
	if n == nil {
		t.Fatalf("Invalid result from Neighbors(): nil")
	}
}

func TestBgpGetSectionConnectionError_UnitTest(t *testing.T) {
	conn := dummyNode.GetConnection().(*DummyEapiConnection)
	conn.setReturnError(true)

	bgp := Bgp(dummyNode)
	section := bgp.GetSection()
	if section != "" {
		t.Fatalf("GetSection() should return \"\" on error")
	}
}

func TestBgpGetSection_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)
	section := bgp.GetSection()
	if section == "" {
		t.Fatalf("No section returned")
	}
}

func TestBgpCreate_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)

	cmds := []string{
		"router bgp 65000",
	}
	tests := []struct {
		val  int
		want string
		rc   bool
	}{
		{65000, "router bgp 65000", true},
		{65534, "router bgp 65534", true},
		{65535, "", false},
		{66000, "", false},
	}

	for _, tt := range tests {
		if got := bgp.Create(tt.val); got != tt.rc {
			t.Fatalf("Expected \"%t\" got \"%t\"", tt.rc, got)
		}
		if tt.rc {
			cmds[0] = tt.want
			// first two commands are 'enable', 'configure terminal'
			commands := dummyConnection.GetCommands()[2:]
			for idx, val := range commands {
				if cmds[idx] != val {
					t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
				}
			}
		}
	}
}

func TestBgpDeleteConnectionError_UnitTest(t *testing.T) {
	conn := dummyNode.GetConnection().(*DummyEapiConnection)
	conn.setReturnError(true)
	bgp := Bgp(dummyNode)

	if ok := bgp.Delete(); !ok {
		t.Fatalf("Delete during failed Get() returns failure")
	}
}

func TestBgpDelete_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)

	cmds := []string{
		"no router bgp 65000",
	}
	bgp.Delete()
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpDefaultConnectionError_UnitTest(t *testing.T) {
	conn := dummyNode.GetConnection().(*DummyEapiConnection)
	conn.setReturnError(true)
	bgp := Bgp(dummyNode)

	if ok := bgp.Default(); !ok {
		t.Fatalf("Default during failed Get() returns failure")
	}
}

func TestBgpDefault_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)

	cmds := []string{
		"default router bgp 65000",
	}
	bgp.Default()
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpSetRouterIDConnectionError_UnitTest(t *testing.T) {
	conn := dummyNode.GetConnection().(*DummyEapiConnection)
	conn.setReturnError(true)
	bgp := Bgp(dummyNode)
	if ok := bgp.SetRouterID("1.1.1.1"); ok {
		t.Fatalf("SetRouterID should return false during connection error")
	}
}

func TestBgpSetRouterID_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)

	cmds := []string{
		"router bgp 65000",
		"default router-id",
	}
	tests := []struct {
		val  string
		want string
	}{
		{"", "no router-id"},
		{"1.1.1.1", "router-id 1.1.1.1"},
	}

	for _, tt := range tests {
		bgp.SetRouterID(tt.val)
		cmds[1] = tt.want
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestBgpSetRouterIDDefault_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)
	cmds := []string{
		"router bgp 65000",
		"default router-id",
	}
	bgp.SetRouterIDDefault()
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpSetMaximumPaths_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)

	cmds := []string{
		"router bgp 65000",
		"maximum-paths 20",
	}
	bgp.SetMaximumPaths(20)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpSetMaximumPathsWithEcmp_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)

	cmds := []string{
		"router bgp 65000",
		"maximum-paths 20 ecmp 20",
	}

	bgp.SetMaximumPathsWithEcmp(20, 20)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpSetMaximumPathsDefault_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)

	cmds := []string{
		"router bgp 65000",
		"default maximum-paths",
	}
	bgp.SetMaximumPathsDefault()
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpSetShutdown_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)

	cmds := []string{
		"router bgp 65000",
		"default shutdown",
	}
	tests := []struct {
		enable bool
		want   string
	}{
		{false, "no shutdown"},
		{true, "shutdown"},
	}

	for _, tt := range tests {
		bgp.SetShutdown(tt.enable)
		cmds[1] = tt.want
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestBgpSetShutdownDefault_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)

	cmds := []string{
		"router bgp 65000",
		"default shutdown",
	}
	bgp.SetShutdownDefault()
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpAddNetworkWithRouteMap_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)

	cmds := []string{
		"router bgp 65000",
		"network 172.16.10.1/24 route-map test",
	}
	bgp.AddNetworkWithRouteMap("172.16.10.1", "24", "test")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpAddNetwork_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)

	cmds := []string{
		"router bgp 65000",
		"network 172.16.10.1/24",
	}
	bgp.AddNetwork("172.16.10.1", "24")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpRemoveNetworkWithRouteMap_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)

	cmds := []string{
		"router bgp 65000",
		"no network 172.16.10.1/24 route-map test",
	}
	bgp.RemoveNetworkWithRouteMap("172.16.10.1", "24", "test")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpRemoveNetwork_UnitTest(t *testing.T) {
	bgp := Bgp(dummyNode)

	cmds := []string{
		"router bgp 65000",
		"no network 172.16.10.1/24",
	}
	bgp.RemoveNetwork("172.16.10.1", "24")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpNeigborsGet_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	neighbor := n.Get("test")

	keys := []string{
		"peer_group",
		"remote_as",
		"send_community",
		"shutdown",
		"description",
		"next_hop_self",
		"route_in_map",
		"route_out_map",
	}
	if len(keys) != len(neighbor) {
		t.Fatalf("Keys mismatch. Expect: %q got %#v", keys, neighbor)
	}
	for _, val := range keys {
		if _, found := neighbor[val]; !found {
			t.Fatalf("Key \"%s\" not found in neighbor", val)
		}
	}
}

func TestBgpNeigborsGetAllConnectionFailure_UnitTest(t *testing.T) {
	conn := dummyNode.GetConnection().(*DummyEapiConnection)
	conn.setReturnError(true)
	n := Bgp(dummyNode).Neighbors()
	neighbors := n.GetAll()
	if neighbors != nil {
		t.Fatalf("Expected nil on GetAll with Connection error")
	}
	if err := n.Error(); err == nil {
		t.Fatalf("Expected Connection error")
	}
	neighbors = n.GetAll()
	if neighbors == nil {
		t.Fatalf("Expected non-nil on GetAll. Connection error should be cleared")
	}
	if err := n.Error(); err != nil {
		t.Fatalf("Connection error not cleared")
	}
}

func TestBgpNeigborsGetAll_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	neighbors := n.GetAll()
	if neighbors == nil {
		t.Fatalf("GetAll on neighbors returned nil")
	}
}

func TestBgpNeigborsCreate_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"neighbor test shutdown",
	}
	n.Create("test")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpNeigborsDeleteConnectionError_UnitTest(t *testing.T) {
	conn := dummyNode.GetConnection().(*DummyEapiConnection)
	conn.setReturnError(true)

	n := Bgp(dummyNode).Neighbors()
	if ok := n.Delete("test"); !ok {
		t.Fatalf("Expected")
	}
}

func TestBgpNeigborsDelete_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"no neighbor test",
	}
	n.Delete("test")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpNeigborsSetPeerGroup_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor 172.16.10.1 peer-group",
	}

	tests := []struct {
		neighbor string
		gp       string
		want     string
		rc       bool
	}{
		{"172.16.10.1", "test", "neighbor 172.16.10.1 peer-group test", true},
		{"172.16.10.300", "test", "", false},
		{"", "test", "", false},
	}

	for idx, tt := range tests {
		if ok := n.SetPeerGroup(tt.neighbor, tt.gp); ok != tt.rc {
			t.Fatalf("Expected \"%t\" got \"%t\" for test %d", tt.rc, ok, idx)
		}
		if tt.rc {
			cmds[1] = tt.want
			// first two commands are 'enable', 'configure terminal'
			commands := dummyConnection.GetCommands()[2:]
			for idx, val := range commands {
				if cmds[idx] != val {
					t.Fatalf("test[%d] Expected \"%q\" got \"%q\"", idx, cmds, commands)
				}
			}
		}
	}
}

func TestBgpNeigborsSetPeerGroupDefault_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor 172.16.10.1 peer-group",
	}
	n.SetPeerGroupDefault("172.16.10.1")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
	if ok := n.SetPeerGroupDefault(""); ok {
		t.Fatalf("Invalid Ip should fail")
	}
}

func TestBgpNeigborsSetRemoteAS_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test remote-as",
	}

	tests := []struct {
		neighbor string
		as       string
		want     string
		rc       bool
	}{
		{"172.16.10.1", "65000", "neighbor 172.16.10.1 remote-as 65000", true},
	}

	for idx, tt := range tests {
		if ok := n.SetRemoteAS(tt.neighbor, tt.as); ok != tt.rc {
			t.Fatalf("Expected \"%t\" got \"%t\" for test %d", tt.rc, ok, idx)
		}
		if tt.rc {
			cmds[1] = tt.want
			// first two commands are 'enable', 'configure terminal'
			commands := dummyConnection.GetCommands()[2:]
			for idx, val := range commands {
				if cmds[idx] != val {
					t.Fatalf("test[%d] Expected \"%q\" got \"%q\"", idx, cmds, commands)
				}
			}
		}
	}
}

func TestBgpNeigborsSetRemoteASDefault_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test remote-as",
	}
	n.SetRemoteASDefault("test")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpNeigborsSetShutdown_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test shutdown",
	}

	tests := []struct {
		neighbor string
		enable   bool
		want     string
	}{
		{"test", true, "neighbor test shutdown"},
		{"test", false, "no neighbor test shutdown"},
	}

	for _, tt := range tests {
		n.SetShutdown(tt.neighbor, tt.enable)
		cmds[1] = tt.want
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestBgpNeigborsSetShutdownDefault_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test shutdown",
	}
	n.SetShutdownDefault("test")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpNeigborsSetSendCommunity_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test send-community",
	}

	tests := []struct {
		neighbor string
		val      bool
		want     string
	}{
		{"test", true, "neighbor test send-community"},
		{"test", false, "no neighbor test send-community"},
	}

	for _, tt := range tests {
		n.SetSendCommunity(tt.neighbor, tt.val)
		cmds[1] = tt.want
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestBgpNeigborsSetSendCommunityDefault_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test send-community",
	}
	n.SetSendCommunityDefault("test")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpNeigborsSetNextHopSelf_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test send-community",
	}

	tests := []struct {
		neighbor string
		val      bool
		want     string
	}{
		{"test", true, "neighbor test next-hop-self"},
		{"test", false, "no neighbor test next-hop-self"},
	}

	for _, tt := range tests {
		n.SetNextHopSelf(tt.neighbor, tt.val)
		cmds[1] = tt.want
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestBgpNeigborsSetNextHopSelfDefault_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test next-hop-self",
	}
	n.SetNextHopSelfDefault("test")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpNeigborsSetRouteMapIn_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test route-map in",
	}

	tests := []struct {
		neighbor string
		val      string
		want     string
	}{
		{"test", "TEST_RM", "neighbor test route-map TEST_RM in"},
	}

	for _, tt := range tests {
		n.SetRouteMapIn(tt.neighbor, tt.val)
		cmds[1] = tt.want
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestBgpNeigborsSetRouteMapInDefault_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test route-map in",
	}
	n.SetRouteMapInDefault("test")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpNeigborsSetRouteMapOut_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test route-map out",
	}
	tests := []struct {
		neighbor string
		val      string
		want     string
	}{
		{"test", "TEST_RM", "neighbor test route-map TEST_RM out"},
	}

	for _, tt := range tests {
		n.SetRouteMapOut(tt.neighbor, tt.val)
		cmds[1] = tt.want
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestBgpNeigborsSetRouteMapOutDefault_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test route-map out",
	}
	n.SetRouteMapOutDefault("test")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestBgpNeigborsSetDescription_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test description",
	}

	tests := []struct {
		in   string
		want string
	}{
		{"this is a test", "neighbor test description this is a test"},
	}

	for _, tt := range tests {
		n.SetDescription("test", tt.in)
		cmds[1] = tt.want
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestBgpNeigborsSetDescriptionDefault_UnitTest(t *testing.T) {
	n := Bgp(dummyNode).Neighbors()
	cmds := []string{
		"router bgp 65000",
		"default neighbor test description",
	}
	n.SetDescriptionDefault("test")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

/**
 *****************************************************************************
 * System Tests
 *****************************************************************************
 **/
func TestBgpGet_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		bgp.Get()
	}
}

func TestBgpCreate_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		section := bgp.GetSection()
		if found, _ := regexp.MatchString("router bgp 100", section); found {
			t.Fatalf("Found BGP config for AS 100.\n[%s]", section)
		}

		if ok := bgp.Create(100); !ok {
			t.Fatalf("Failure to Create AS 100")
		}

		section = bgp.GetSection()
		if found, _ := regexp.MatchString("router bgp 100", section); !found {
			t.Fatalf("\"router bgp 100\" expected but not seen under "+
				"router section.\n[%s]", section)
		}
	}
}

func TestBgpDelete_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		section := bgp.GetSection()
		if found, _ := regexp.MatchString("router bgp 100", section); !found {
			t.Fatalf("BGP not found in config for AS 100.\n[%s]", section)
		}

		if ok := bgp.Delete(); !ok {
			t.Fatalf("Failure to Delete AS 100")
		}

		section = bgp.GetSection()
		if found, _ := regexp.MatchString("router bgp 100", section); found {
			t.Fatalf("\"router bgp 100\" NOT expected but not seen\n[%s]",
				section)
		}
	}
}

func TestBgpDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		section := bgp.GetSection()
		if found, _ := regexp.MatchString("router bgp 100", section); !found {
			t.Fatalf("BGP not found in config for AS 100.\n[%s]", section)
		}

		if ok := bgp.Default(); !ok {
			t.Fatalf("Failure to Default AS 100")
		}

		section = bgp.GetSection()
		if found, _ := regexp.MatchString("router bgp 100", section); found {
			t.Fatalf("\"router bgp 100\" NOT expected but not seen\n[%s]",
				section)
		}
	}
}

func TestBgpSetRouterID_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		section := bgp.GetSection()
		if found, _ := regexp.MatchString("router bgp 100", section); !found {
			t.Fatalf("BGP not found in config for AS 100.\n[%s]", section)
		}

		if ok := bgp.SetRouterID("1.2.3.4"); !ok {
			t.Fatalf("Failure to Default AS 100")
		}

		section = bgp.GetSection()
		if found, _ := regexp.MatchString("router-id 1.2.3.4", section); !found {
			t.Fatalf("\"router-id 1.2.3.4\" expected but not seen under "+
				"router bgp 100 section.\n[%s]", section)
		}
	}
}

func TestBgpSetRouterIDNoValue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
			"router-id 1.2.3.4",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		section := bgp.GetSection()
		if found, _ := regexp.MatchString("router-id 1.2.3.4", section); !found {
			t.Fatalf("router-id 1.2.3.4 not found in config for AS 100.\n[%s]",
				section)
		}

		if ok := bgp.SetRouterID(""); !ok {
			t.Fatalf("Failure to Default AS 100")
		}

		section = bgp.GetSection()
		if found, _ := regexp.MatchString("router-id 1.2.3.4", section); found {
			t.Fatalf("\"router-id 1.2.3.4\" NOT expected but seen under "+
				"router bgp 100 section.\n[%s]", section)
		}
	}
}

func TestBgpSetRouterIdDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
			"router-id 1.2.3.4",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		section := bgp.GetSection()
		if found, _ := regexp.MatchString("router-id 1.2.3.4", section); !found {
			t.Fatalf("router-id 1.2.3.4 not found in config for AS 100.\n[%s]",
				section)
		}

		if ok := bgp.SetRouterIDDefault(); !ok {
			t.Fatalf("Failure to Default AS 100")
		}

		section = bgp.GetSection()
		if found, _ := regexp.MatchString("router-id 1.2.3.4", section); found {
			t.Fatalf("\"router-id 1.2.3.4\" NOT expected but seen under "+
				"router bgp 100 section.\n[%s]", section)
		}
	}
}

func TestBgpSetMaximumPaths_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		if ok := bgp.SetMaximumPaths(12); !ok {
			t.Fatalf("Failure to set MaximumPaths in AS 100")
		}

		section := bgp.GetSection()
		if found, _ := regexp.MatchString("maximum-paths 12", section); !found {
			t.Fatalf("\"maximum-paths 12\" expected but seen under "+
				"router bgp 100 section.\n[%s]", section)
		}
	}
}

func TestBgpSetMaximumPathsWithEcmp_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		if ok := bgp.SetMaximumPathsWithEcmp(18, 32); !ok {
			t.Fatalf("Failure to set MaximumPathsWithEcmp in AS 100")
		}

		section := bgp.GetSection()
		if found, _ := regexp.MatchString("maximum-paths 18 ecmp 32", section); !found {
			t.Fatalf("\"maximum-paths 12\" expected but seen under "+
				"router bgp 100 section.\n[%s]", section)
		}
	}
}

func TestBgpSetMaximumPathsDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
			"maximum-paths 12 ecmp 18",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		if ok := bgp.SetMaximumPathsDefault(); !ok {
			t.Fatalf("Failure to set MaximumPathsDefault in AS 100")
		}

		section := bgp.GetSection()
		if found, _ := regexp.MatchString("maximum-paths 12", section); found {
			t.Fatalf("\"maximum-paths 12\" NOT expected but seen under "+
				"router bgp 100 section.\n[%s]", section)
		}
	}
}

func TestBgpShutdownTrue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		section := bgp.GetSection()
		if found, _ := regexp.MatchString(`(?m)no shutdown`, section); !found {
			t.Fatalf("no shutdown not found in config for AS 100.\n[%s]",
				section)
		}

		if ok := bgp.SetShutdown(true); !ok {
			t.Fatalf("Failure to Default AS 100")
		}

		section = bgp.GetSection()
		if found, _ := regexp.MatchString(`(?m)no shutdown`, section); found {
			t.Fatalf("\"no shutdown\" NOT expected but seen under "+
				"router bgp 100 section.\n[%s]", section)
		}
	}
}

func TestBgpShutdownFalse_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
			"shutdown",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		section := bgp.GetSection()
		if found, _ := regexp.MatchString(`(?m)no shutdown`, section); found {
			t.Fatalf("no shutdown found in config for AS 100.\n[%s]",
				section)
		}

		if ok := bgp.SetShutdown(false); !ok {
			t.Fatalf("Failure to Default AS 100")
		}

		section = bgp.GetSection()
		if found, _ := regexp.MatchString(`(?m)no shutdown`, section); !found {
			t.Fatalf("\"no shutdown\" expected but seen under "+
				"router bgp 100 section.\n[%s]", section)
		}
	}
}

func TestBgpShutdownDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
			"shutdown",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		section := bgp.GetSection()
		if found, _ := regexp.MatchString(`(?m)no shutdown`, section); found {
			t.Fatalf("no shutdown found in config for AS 100.\n[%s]",
				section)
		}

		if ok := bgp.SetShutdownDefault(); !ok {
			t.Fatalf("Failure to Default AS 100")
		}

		section = bgp.GetSection()
		if found, _ := regexp.MatchString(`(?m)no shutdown`, section); !found {
			t.Fatalf("\"no shutdown\" expected but seen under "+
				"router bgp 100 section.\n[%s]", section)
		}
	}
}

func TestBgpAddNetworkWithRouteMap_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		section := bgp.GetSection()
		if found, _ := regexp.MatchString("network 172.16.10.0/24 route-map Test1",
			section); found {
			t.Fatalf("network found in config for AS 100.\n[%s]",
				section)
		}

		if ok := bgp.AddNetworkWithRouteMap("172.16.10.0", "24", "Test1"); !ok {
			t.Fatalf("Failure to Add Network With RouteMap under bgp 100")
		}

		section = bgp.GetSection()
		if found, _ := regexp.MatchString("network 172.16.10.0/24 route-map Test1", section); !found {
			t.Fatalf("\"network 172.16.10.0/24 route-map Test1\" expected but seen under "+
				"router bgp 100 section.\n[%s]", section)
		}
	}
}

func TestBgpAddNetwork_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		section := bgp.GetSection()
		if found, _ := regexp.MatchString("network 172.16.20.0/24",
			section); found {
			t.Fatalf("network found in config for AS 100.\n[%s]",
				section)
		}

		if ok := bgp.AddNetwork("172.16.20.0", "24"); !ok {
			t.Fatalf("Failure to Add Network under bgp 100")
		}

		section = bgp.GetSection()
		if found, _ := regexp.MatchString("network 172.16.20.0/24", section); !found {
			t.Fatalf("\"network 172.16.20.0/24\" expected but seen under "+
				"router bgp 100 section.\n[%s]", section)
		}
	}
}

func TestBgpRemoveNetworkWithRouteMap_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
			"network 172.16.10.0/24 route-map Test1",
			"network 172.16.30.0/24 route-map Test2",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		section := bgp.GetSection()
		if found, _ := regexp.MatchString("network 172.16.30.0/24 route-map Test2",
			section); !found {
			t.Fatalf("network not found in config for bgp 100.\n[%s]",
				section)
		}

		if ok := bgp.RemoveNetworkWithRouteMap("172.16.30.0", "24", "Test2"); !ok {
			t.Fatalf("Failure to Remove Network With RouteMap in AS 100")
		}

		section = bgp.GetSection()
		if found, _ := regexp.MatchString("network 172.16.30.0/24 route-map Test2", section); found {
			t.Fatalf("\"network 172.16.30.0/24 route-map Test2\" NOT expected but seen under "+
				"router bgp 100 section.\n[%s]", section)
		}
	}
}

func TestBgpRemoveNetwork_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
			"default router bgp",
			"router bgp 100",
			"network 172.16.10.0/24",
			"network 172.16.40.0/24",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		bgp := Bgp(dut)

		section := bgp.GetSection()
		if found, _ := regexp.MatchString("network 172.16.40.0/24",
			section); !found {
			t.Fatalf("network not found in config for bgp 100.\n[%s]",
				section)
		}

		if ok := bgp.RemoveNetwork("172.16.40.0", "24"); !ok {
			t.Fatalf("Failure to Remove Network With RouteMap in AS 100")
		}

		section = bgp.GetSection()
		if found, _ := regexp.MatchString("network 172.16.40.0/24", section); found {
			t.Fatalf("\"network 172.16.30.0/24\" NOT expected but seen under "+
				"router bgp 100 section.\n[%s]", section)
		}
	}
}

func TestShowIPBGPSummary_UnitTest(t *testing.T) {
	var dummyNode *goeapi.Node
	var dummyConnection *DummyConnection

	dummyConnection = &DummyConnection{}

	dummyNode = &goeapi.Node{}
	dummyNode.SetConnection(dummyConnection)

	show := Show(dummyNode)
	showBgpSummary, err := show.ShowIPBGPSummary()
	if err != nil {
		t.Errorf("Error during show ip bgp summary, %s", err)
	}

	// Test VRFs inside the BGP Summary
	var vrfScenarios = []struct {
		Name     string
		NumPeers int
		RouterID string
		ASN      int64
	}{
		{
			Name:     "default",
			NumPeers: 4,
			RouterID: "10.10.10.38",
			ASN:      65004,
		},
	}

	for _, vrf := range vrfScenarios {
		if _, ok := showBgpSummary.VRFs[vrf.Name]; !ok {
			t.Errorf("VRF Name %s does not exist", vrf.Name)
		} else {
			vrfSummary := showBgpSummary.VRFs[vrf.Name]
			if vrf.RouterID != vrfSummary.RouterID {
				t.Errorf("RouterID does not match expected %s, got %s", vrfSummary.RouterID, vrf.RouterID)
			}

			if vrf.ASN != vrfSummary.ASN {
				t.Errorf("ASN does not match expected %d, got %d", vrfSummary.ASN, vrf.ASN)
			}

			if len(vrfSummary.Peers) != vrf.NumPeers {
				t.Errorf("Number of Peers does not match expected %d, got %d", vrf.NumPeers, len(vrfSummary.Peers))
			}
		}
	}

	// Test Peers inside a VRF
	var peerScenarios = []struct {
		ASN                 int64
		PeerIP              string
		PeerState           string
		PeerStateIdleReason string
		PrefixAccepted      int
		PrefixReceived      int
		UnderMaintenance    bool
		UpDownTime          float64
		Version             int
	}{
		{
			PeerIP:              "10.10.10.33",
			ASN:                 65000,
			PeerState:           "Idle",
			PeerStateIdleReason: "NoInterface",
			PrefixAccepted:      0,
			PrefixReceived:      0,
			UnderMaintenance:    false,
			UpDownTime:          1524094401.78999,
			Version:             4,
		},
		{
			PeerIP:           "10.10.10.14",
			ASN:              65001,
			PeerState:        "Active",
			PrefixAccepted:   0,
			PrefixReceived:   0,
			UnderMaintenance: false,
			UpDownTime:       1524094401.78899,
			Version:          4,
		},
		{
			PeerIP:           "10.10.10.216",
			ASN:              65002,
			PeerState:        "Established",
			PrefixAccepted:   429,
			PrefixReceived:   444,
			UnderMaintenance: false,
			UpDownTime:       1525451375.739572,
			Version:          4,
		},
		{
			PeerIP:           "10.10.10.4",
			ASN:              65003,
			PeerState:        "Established",
			PrefixAccepted:   429,
			PrefixReceived:   444,
			UnderMaintenance: false,
			UpDownTime:       1525451380.74079,
			Version:          4,
		},
	}

	for _, peer := range peerScenarios {
		if _, ok := showBgpSummary.VRFs["default"].Peers[peer.PeerIP]; !ok {
			t.Errorf("Peer IP %s not found in Peers", peer.PeerIP)
		} else {
			peerSummary := showBgpSummary.VRFs["default"].Peers[peer.PeerIP]

			if peer.Version != peerSummary.Version {
				t.Errorf("Peer Version does not match expected %d, got %d", peerSummary.Version, peer.Version)
			}

			if peer.ASN != peerSummary.ASN {
				t.Errorf("Peer ASN does not match expected %d, got %d", peerSummary.ASN, peer.ASN)
			}

			if peer.PeerStateIdleReason != peerSummary.PeerStateIdleReason {
				t.Errorf("Peer PeerStateIdleReason does not match expected %s, got %s", peerSummary.PeerStateIdleReason, peer.PeerStateIdleReason)
			}

			if peer.PrefixAccepted != peerSummary.PrefixAccepted {
				t.Errorf("Peer PrefixAccepted does not match expected %d, got %d", peerSummary.PrefixAccepted, peer.PrefixAccepted)
			}

			if peer.PrefixReceived != peerSummary.PrefixReceived {
				t.Errorf("Peer PrefixReceived does not match expected %d, got %d", peerSummary.PrefixReceived, peer.PrefixReceived)
			}

			if peer.PeerState != peerSummary.PeerState {
				t.Errorf("Peer PeerState does not match expected %s, got %s", peerSummary.PeerState, peer.PeerState)
			}

			if peer.UpDownTime != peerSummary.UpDownTime {
				t.Errorf("Peer UpDownTime does not match expected %f, got %f", peerSummary.UpDownTime, peer.UpDownTime)
			}

			if peer.UnderMaintenance != peerSummary.UnderMaintenance {
				t.Errorf("Peer UnderMaintenance does not match expected %t, got %t", peerSummary.UnderMaintenance, peer.UnderMaintenance)
			}
		}
	}
}

func TestShowIPBGPSummaryErrorDuringCall_UnitTest(t *testing.T) {
	dummyConnection := &DummyConnection{err: errors.New("error during connection")}
	dummyNode := &goeapi.Node{}
	dummyNode.SetConnection(dummyConnection)

	show := Show(dummyNode)
	_, err := show.ShowIPBGPSummary()
	if err == nil {
		t.Errorf("Error expected during show ip bgp summary")
	}
}
