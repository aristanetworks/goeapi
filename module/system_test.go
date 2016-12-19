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
	"testing"
)

var initFile bool

/**
 *****************************************************************************
 * Unit Tests
 *****************************************************************************
 **/
func TestSystemHostNameParser_UnitTest(t *testing.T) {

	shortConfig := `
default mrp leave-all-timer
!
no msrp streams load-file
!
hostname %s
no ip domain lookup source-interface
ip name-server vrf default 192.168.1.32
ip domain-name stormcontrol.net
`

	tests := [10]struct {
		in   string
		want string
	}{}

	var ts string

	for idx := range tests {
		ts = RandomString(2, 14)
		testConfig := fmt.Sprintf(shortConfig, ts)
		tests[idx].in = testConfig
		tests[idx].want = ts
	}

	for _, tt := range tests {
		if got := parseHostname(tt.in); got != tt.want {
			t.Fatalf("parseHostname() = %q; want %q", got, tt.want)
		}
	}
	if got := parseHostname(""); got != "localhost" {
		t.Fatalf("parseHostname() = %q; want \"localhost\"", got)
	}
}

func TestSystemIpRoutingParser_UnitTest(t *testing.T) {

	shortConfig := `
ip route 0.0.0.0/0 192.68.1.254 1 tag 0
!
ip icmp redirect
%s
!
no ip multicast-routing
no ip multicast-routing static
no ip multicast multipath none
ip mfib activity polling-interval 60
ip mfib max-fastdrops 1024
ip mfib cache-entries unresolved max 4000
ip mfib packet-buffers unresolved max 3
`

	tests := [...]struct {
		in   string
		want bool
	}{
		{"no ip routing", false},
		{"ip routing", true},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := parseIPRouting(testConfig); got != tt.want {
			t.Fatalf("parseIPRouting() = %v; want %v. Config: %s", got, tt.want, testConfig)
		}
	}
	if got := parseIPRouting(""); got != false {
		t.Fatalf("parseIPRouting(\"\") = %v; want false.", got)
	}
}

func TestSystemGetKeysReturned_UnitTest(t *testing.T) {
	sys := System(dummyNode)
	sysConfig := sys.Get()
	for _, val := range []string{"hostname", "iprouting"} {
		if _, found := sysConfig[val]; !found {
			t.Fatalf("Get() missing key %s", val)
		}
	}
}

func TestSystemCheckSystemConfigGetters_UnitTest(t *testing.T) {
	sys := System(dummyNode)
	sysConfig := sys.Get()

	if sysConfig.HostName() != "veos01" ||
		sysConfig.IPRouting() != "true" {
		t.Fatalf("Mismatch in values. Got: %#v", sysConfig)
	}
}

func TestSystemSetHostname_UnitTest(t *testing.T) {
	sys := System(dummyNode)
	newName1 := RandomString(3, 16)
	newName2 := RandomString(3, 16)
	newName3 := RandomString(3, 16)
	tests := []struct {
		in   string
		want string
	}{
		{newName1, "hostname " + newName1},
		{newName2, "hostname " + newName2},
		{newName3, "hostname " + newName3},
		{"", "no hostname"},
	}
	for _, tt := range tests {
		if ok := sys.SetHostname(tt.in); !ok {
			t.Fatalf("SetHostname")
		}
		cmds := dummyConnection.GetCommands()
		if cmds[len(cmds)-1] != tt.want {
			t.Errorf("Expected \"%s\" got \"%s\"", tt.want, cmds[len(cmds)-1])
		}
	}
}

func TestSystemSetHostnameDefault_UnitTest(t *testing.T) {
	sys := System(dummyNode)
	if ok := sys.SetHostnameDefault(); !ok {
		t.Fatalf("SetHostnameDefault")
	}
	cmds := dummyConnection.GetCommands()
	if cmds[len(cmds)-1] != "default hostname" {
		t.Errorf("Expected \"default hostname\" got \"%s\"", cmds[len(cmds)-1])
	}
}

func TestSystemSetIPRouting_UnitTest(t *testing.T) {
	sys := System(dummyNode)
	tests := []struct {
		in1  string
		in2  bool
		want string
	}{
		{"", true, "ip routing"},
		{"", false, "no ip routing"},
		{"vrf TEST", true, "ip routing vrf TEST"},
		{"vrf TEST", false, "no ip routing vrf TEST"},
	}
	for _, tt := range tests {
		if ok := sys.SetIPRouting(tt.in1, tt.in2); !ok {
			t.Fatalf("SetIPRouting")
		}
		cmds := dummyConnection.GetCommands()
		if cmds[len(cmds)-1] != tt.want {
			t.Errorf("Expected \"%s\" got \"%s\"", tt.want, cmds[len(cmds)-1])
		}
	}
}

func TestSystemSetIPRoutingDefault_UnitTest(t *testing.T) {
	sys := System(dummyNode)
	tests := []struct {
		in   string
		want string
	}{
		{"", "default ip routing"},
		{"vrf TEST", "default ip routing vrf TEST"},
	}
	for _, tt := range tests {
		if ok := sys.SetIPRoutingDefault(tt.in); !ok {
			t.Fatalf("SetIPRoutingDefault")
		}
		cmds := dummyConnection.GetCommands()
		if cmds[len(cmds)-1] != tt.want {
			t.Errorf("Expected \"%s\" got \"%s\"", tt.want, cmds[len(cmds)-1])
		}
	}
}

func TestSystemParseHostname_UnitTest(t *testing.T) {
	sys := System(dummyNode)
	hostname := sys.parseHostname()
	if hostname != "veos01" {
		t.Fatalf("ParseHostname failed. expected \"localhost\" got %s", hostname)
	}
}

func TestSystemIPRouting_UnitTest(t *testing.T) {
	sys := System(dummyNode)
	if ok := sys.parseIPRouting(); !ok {
		t.Fatalf("ParseIPRouting expected \"true\"")
	}
}

/**
 *****************************************************************************
 * System Tests
 *****************************************************************************
 **/
func TestSystemParseHostName_SystemTest(t *testing.T) {
	for _, dut := range duts {
		dut.SetAutoRefresh(true)
		sys := System(dut)

		currName := sys.parseHostname()
		newName := RandomString(3, 16)

		if ok := sys.SetHostname(newName); !ok {
			t.Error("Test1: Sethostname failed")
		}
		dut.RunningConfig()
		hname := sys.parseHostname()
		if hname != newName {
			t.Fatalf("Test2: Sethostname failed %s != %s", newName, hname)
		}

		if ok := sys.SetHostname(currName); !ok {
			t.Error("Test3: Sethostname failed reapplying original hostname")
		}
		dut.RunningConfig()
		hname = sys.parseHostname()
		if hname == newName {
			t.Fatalf("Test4: parseHostname failed %s != %s", currName, hname)
		}
	}
}

func TestSystemGet_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"default hostname",
			"ip routing",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		system := System(dut)
		sysConfig := system.Get()

		if sysConfig.HostName() != "localhost" ||
			sysConfig.IPRouting() != "true" {
			t.Fatalf("Mismatch in values. Got: %#v", sysConfig)
		}
	}
}

func TestSystemWithPeriod_SystemTest(t *testing.T) {
	for _, dut := range duts {
		name := "host.domain.net"
		cmds := []string{
			"hostname " + name,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		system := System(dut)
		sysConfig := system.Get()

		if sysConfig.HostName() != name {
			t.Fatalf("Expecting %s for hostname. Got \"%s\"", name, sysConfig["hostname"])
		}
	}
}

func TestSystemCheckName_SystemTest(t *testing.T) {
	for _, dut := range duts {
		name := "teststring"
		cmds := []string{
			"hostname " + name,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		system := System(dut)
		sysConfig := system.Get()

		if sysConfig["hostname"] != name {
			t.Fatalf("Expecting %s for hostname. Got \"%s\"", name, sysConfig["hostname"])
		}
	}
}

func TestSystemSetHostnameWithVal_SystemTest(t *testing.T) {
	for _, dut := range duts {
		name := RandomString(2, 14)
		cmds := []string{
			"default hostname",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		system := System(dut)
		if ok := system.SetHostname(name); !ok {
			t.Fatalf("SetHostname failed.")
		}

		nameCurrent := system.parseHostname()
		if name != nameCurrent {
			t.Fatalf("Expecting \"%s\" for hostname. Got \"%s\"", name, nameCurrent)
		}
	}
}

func TestSystemSetHostnameNoVal_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"default hostname",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		system := System(dut)
		if ok := system.SetHostname(""); !ok {
			t.Fatalf("SetHostname failed.")
		}

		config := dut.RunningConfig()
		if found, _ := regexp.MatchString("no hostname", config); !found {
			t.Fatalf("Expecting \"no hostname\" in running config.")
		}
	}
}

func TestSystemSetHostnameDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"hostname test",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		system := System(dut)
		if ok := system.SetHostnameDefault(); !ok {
			t.Fatalf("SetHostnameDefault failed.")
		}

		config := dut.RunningConfig()
		if found, _ := regexp.MatchString("no hostname", config); !found {
			t.Fatalf("Expecting \"no hostname\" in running config.")
		}
	}
}

func TestSystemSetIpRoutingTrue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		system := System(dut)
		if ok := system.SetIPRouting("", true); !ok {
			t.Fatalf("SetIPRouting failed.")
		}

		if system.parseIPRouting() != true {
			t.Fatalf("expecting ip routing to be configured.\n")
		}
	}
}

func TestSystemSetIpRoutingFalse_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		system := System(dut)
		if ok := system.SetIPRouting("", false); !ok {
			t.Fatalf("SetIPRouting failed.")
		}
		if system.parseIPRouting() != false {
			t.Fatalf("expecting no ip routing to be configured.")
		}
	}
}

func TestSystemSetIpRoutingDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"ip routing",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		system := System(dut)
		if ok := system.SetIPRoutingDefault(""); !ok {
			t.Fatalf("SetIPRoutingDefault failed.")
		}
		if system.parseIPRouting() != false {
			t.Fatalf("expecting no ip routing to be configured.")
		}
	}
}

func TestSystemSetHostnameWithPeriod_SystemTest(t *testing.T) {
	for _, dut := range duts {
		name := "localhost"
		cmds := []string{
			"hostname " + name,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		system := System(dut)
		if ok := system.SetHostname("host.domain.net"); !ok {
			t.Fatalf("SetHostname failed.")
		}
		nameCurrent := system.parseHostname()
		if "host.domain.net" != nameCurrent {
			t.Fatalf("Expecting \"host.domain.net\" for hostname. Got \"%s\"", nameCurrent)
		}
	}
}
