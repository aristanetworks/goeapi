//
// Copyright (c) 2015, Arista Networks, Inc.
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
package api

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

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
	}{}

	// for each test entry
	for idx := range tests {
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
	}{}

	// for each test entry
	for idx := range tests {
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
}

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
