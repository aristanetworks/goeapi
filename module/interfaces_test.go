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
package module

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

func TestInterfaceFuncIsValidInterface_UnitTest(t *testing.T) {
	tests := [...]struct {
		in   string
		want bool
	}{
		{"Ethernet1", true},
		{"Token-Ring1", false},
		{"Ethernet1/1", true},
		{"Et1", false},
		{"Vlan1234", true},
		{"Management1", true},
		{"Ma1", false},
		{"Port-Channel1", true},
		{"Po1", false},
		{"Loopback2", true},
		{"Lo2", false},
		{"Vxlan1", true},
		{"Vx1", false},
		{"Vlan10", true},
		{"Vl1", false},
	}

	for _, tt := range tests {
		if got := isValidInterface(tt.in); got != tt.want {
			t.Fatalf("isValidInterface(%s) = %v; want %v", tt.in, got, tt.want)
		}
	}
}

func TestInterfaceParseShutdown_UnitTest(t *testing.T) {
	var i BaseInterfaceEntity
	shortConfig := `
interface Ethernet1
   no description
   %s
   default load-interval
   logging event link-status use-global
`
	tests := [...]struct {
		in   string
		want bool
	}{
		{"no shutdown", false},
		{"shutdown", true},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := BaseInterfaceParseShutdown(&i, testConfig); got != tt.want {
			t.Fatalf("parseShutdown() = %t; want %t", got, tt.want)
		}
	}
}

func TestInterfaceParseDescription_UnitTest(t *testing.T) {
	var i BaseInterfaceEntity
	shortConfig := `
interface Ethernet1
   %s
   no shutdown
   default load-interval
   logging event link-status use-global
`
	tests := [...]struct {
		in   string
		want string
	}{
		{"no description", ""},
		{"description the cat in the hat", "the cat in the hat"},
		{"description br549", "br549"},
		{"description dsfzsdfzdfzsfsdzfrtrtreterg", "dsfzsdfzdfzsfsdzfrtrtreterg"},
		{"description The answer to everything is 42.", "The answer to everything is 42."},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := BaseInterfaceParseDescription(&i, testConfig); got != tt.want {
			t.Fatalf("parseDescription() = %q; want %q", got, tt.want)
		}
	}
}

func TestResourceInterfaceGet_SystemTest(t *testing.T) {
}

func TestResourceInterfaceGetAll_SystemTest(t *testing.T) {
}

func TestResourceInterfaceCreate_SystemTest(t *testing.T) {
}

func TestResourceInterfaceDelete_SystemTest(t *testing.T) {
}

func TestResourceInterfaceDefault_SystemTest(t *testing.T) {
}

func TestResourceInterfaceDescription_SystemTest(t *testing.T) {
}

func TestResourceInterfaceSflowEnable_SystemTest(t *testing.T) {
}

func TestResourceInterfaceSflowDisable_SystemTest(t *testing.T) {
}

func TestPortChannelParseMinimumLinks_UnitTest(t *testing.T) {
	var p PortChannelInterfaceEntity
	shortConfig := `
interface Port-Channel5
   no description
   no shutdown
   default load-interval
   no switchport private-vlan mapping
   switchport trunk group test
   snmp trap link-status
   %s
   no port-channel lacp fallback
   port-channel lacp fallback timeout 90
   no l2 mtu
   no mlag
   no switchport port-security
   switchport port-security maximum 1
`
	tests := [...]struct {
		in   string
		want string
	}{
		{"no port-channel min-links", ""},
		{"port-channel min-links 4", "4"},
		{"port-channel min-links 6", "6"},
		{"port-channel min-links 8", "8"},
		{"port-channel min-links 15", "15"},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := PortChannelParseMinimumLinks(&p, testConfig); got != tt.want {
			t.Fatalf("parseMinimumLinks() = %q; want %q", got, tt.want)
		}
	}
}

func TestPortChannelInterfaceGet_SystemTest(t *testing.T) {
	for _, dut := range duts {
		pc := PortChannel(dut)
		cmds := []string{
			"no interface Port-Channel1",
			"interface Port-Channel1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		got := pc.Get("Port-Channel1")
		if got["type"] != "portchannel" || got["name"] != "Port-Channel1" {
			t.Fatalf("Get(Portchannel1) = %q", got)
		}
	}
}

func TestPortChannelInterfaceGetLacpModeWithDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		pc := PortChannel(dut)
		cmds := []string{
			"no interface Port-Channel1",
			"interface Port-Channel1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		if mode := PortChannelGetLacpMode(pc, "Port-Channel1"); mode != "on" {
			t.Fatalf("parseGetLacpMode(Portchannel1) = %s; want \"on\"", mode)
		}
	}
}

func TestPortChannelInterfaceGetMembersNoMembers_SystemTest(t *testing.T) {
	for _, dut := range duts {
		pc := PortChannel(dut)
		cmds := []string{
			"no interface Port-Channel1",
			"interface Port-Channel1",
			"default interface Ethernet1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		if got := PortChannelGetMembers(pc, "Port-Channel1"); got != nil {
			t.Fatalf("parseGetMembers(Portchannel1) = %s; want 'nil'", got)
		}
	}
}

func TestPortChannelInterfaceGetMembersOneMember_SystemTest(t *testing.T) {
	for _, dut := range duts {
		pc := PortChannel(dut)
		cmds := []string{
			"no interface Port-Channel1",
			"interface Port-Channel1",
			"default interface Ethernet1",
			"interface Ethernet1",
			"channel-group 1 mode active",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		got := PortChannelGetMembers(pc, "Port-Channel1")
		if len(got) < 0 || got[0] != "Ethernet1" {
			t.Fatalf("GetMembers(Portchannel1) = %q; want '[Ethernet1]'", got)
		}
	}
}

func TestPortChannelInterfaceGetMembersTwoMembers_SystemTest(t *testing.T) {
	for _, dut := range duts {
		pc := PortChannel(dut)
		cmds := []string{
			"no interface Port-Channel1",
			"interface Port-Channel1",
			"default interface Ethernet1-2",
			"interface Ethernet1-2",
			"channel-group 1 mode active",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		got := PortChannelGetMembers(pc, "Port-Channel1")
		if len(got) < 2 || got[0] != "Ethernet1" || got[1] != "Ethernet2" {
			t.Fatalf("GetMembers(Portchannel1) = %q; want '[Ethernet1, Ethernet2]'", got)
		}
	}
}

func TestPortChannelInterfaceSetLacpMode_SystemTest(t *testing.T) {
	for _, dut := range duts {
		pc := PortChannel(dut)

		var cfgmode string
		for _, mode := range []string{"on", "active", "passive"} {
			if mode != "on" {
				cfgmode = "on"
			} else {
				cfgmode = "active"
			}
			cmds := []string{
				"no interface Port-Channel1",
				"default interface Ethernet1",
				"interface Ethernet1",
				"channel-group 1 mode " + cfgmode,
			}
			if ok := dut.Config(cmds...); !ok {
				t.Fatalf("dut.Config() failure")
			}

			if ok := pc.SetLacpMode("Port-Channel1", mode); !ok {
			}
		}
	}
}

func TestPortChannelInterfaceSetLacpModeInvalid_SystemTest(t *testing.T) {
	for _, dut := range duts {
		mode := RandomString(4, 10)

		pc := PortChannel(dut)
		if ok := pc.SetLacpMode("Port-Channel1", mode); ok {
			t.Fatalf("SetLacpMode(Port-Channel1, %s) allowed setting mode", mode)
		}
	}
}

func TestPortChannelInterfaceSetMembers_SystemTest(t *testing.T) {
}

func TestPortChannelInterfaceMinLinksValidValue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		pc := PortChannel(dut)
		minLinks := RandomInt(1, 16)

		cmds := []string{
			"no interface Port-Channel1",
			"interface Port-Channel1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := pc.SetMinimumLinks("Port-Channel1", minLinks); !ok {
			t.Fatalf("SetMinimumLinks(Port-Channel1, %d) failed", minLinks)
		}

		config, _ := pc.GetBlock("interface Port-Channel1")
		str := "port-channel min-links " + strconv.Itoa(minLinks)
		if found, _ := regexp.MatchString(str, config); !found {
			t.Fatalf("\"%s\" expected but not seen under "+
				"interface Port-Channel1 section.\n[%s]", str, config)
		}
	}
}

func TestPortChannelInterfaceMinLinksInvalidValue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		pc := PortChannel(dut)
		minLinks := RandomInt(17, 128)

		cmds := []string{
			"no interface Port-Channel1",
			"interface Port-Channel1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := pc.SetMinimumLinks("Port-Channel1", minLinks); ok {
			t.Fatalf("SetMinimumLinks(Port-Channel1, %d) should have failed", minLinks)
		}

		config, _ := pc.GetBlock("interface Port-Channel1")
		str := "port-channel min-links " + strconv.Itoa(minLinks)

		if found, _ := regexp.MatchString(str, config); found {
			t.Fatalf("\"%s\" NOT expected but seen under "+
				"interface Port-Channel1 section.\n[%s]", str, config)
		}
	}
}

func TestPortChannelInterfaceMinLinksDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		pc := PortChannel(dut)

		cmds := []string{
			"no interface Port-Channel1",
			"interface Port-Channel1",
			"port-channel min-links 4",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := pc.SetMinimumLinksDefault("Port-Channel1"); !ok {
			t.Fatalf("SetMinimumLinksDefault(Port-Channel1) failed")
		}

		config, _ := pc.GetBlock("interface Port-Channel1")
		if found, _ := regexp.MatchString("port-channel min-links 4", config); found {
			t.Fatalf("\"port-channel min-links 4\" NOT expected but seen under "+
				"interface Port-Channel1 section.\n[%s]", config)
		}
	}
}

func TestVxlanParseSourceInterface_UnitTest(t *testing.T) {
	var v VxlanInterfaceEntity
	shortConfig := `
interface Vxlan1
   no description
   no shutdown
   vxlan multicast-group 239.10.10.10
   %s
   no vxlan controller-client
   vxlan udp-port 4789
   vxlan vlan 10 vni 10
   vxlan vlan 10 flood vtep 3.3.3.3 4.4.4.4
   vxlan flood vtep 1.1.1.1 2.2.2.2
   no vxlan vlan flood vtep
   no vxlan learn-restrict vtep
   no vxlan vlan learn-restrict vtep
!
`
	tests := [...]struct {
		in   string
		want string
	}{
		{"no vxlan source-interface", ""},
		{"vxlan source-interface Loopback0", "Loopback0"},
		{"vxlan source-interface Loopback100", "Loopback100"},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := VxlanParseSourceInterface(&v, testConfig); got != tt.want {
			t.Fatalf("parseSourceInterface() = %q; want %q", got, tt.want)
		}
	}
}

func TestVxlanParseMulticastGroup_UnitTest(t *testing.T) {
	var v VxlanInterfaceEntity
	shortConfig := `
interface Vxlan1
   no description
   no shutdown
   %s
   vxlan source-interface Loopback0
   no vxlan controller-client
   vxlan udp-port 4789
   vxlan vlan 10 vni 10
   vxlan vlan 10 flood vtep 3.3.3.3 4.4.4.4
   vxlan flood vtep 1.1.1.1 2.2.2.2
   no vxlan vlan flood vtep
   no vxlan learn-restrict vtep
   no vxlan vlan learn-restrict vtep
!
`
	tests := [...]struct {
		in   string
		want string
	}{
		{"no vxlan multicast-group", ""},
		{"vxlan multicast-group 239.10.10.10", "239.10.10.10"},
		{"vxlan multicast-group 239.20.20.20", "239.20.20.20"},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := VxlanParseMulticastGroup(&v, testConfig); got != tt.want {
			t.Fatalf("parseMulticastGroup() = %q; want %q", got, tt.want)
		}
	}
}

func TestVxlanParseUdpPort_UnitTest(t *testing.T) {
	var v VxlanInterfaceEntity
	shortConfig := `
interface Vxlan1
   no description
   no shutdown
   vxlan multicast-group 239.10.10.10
   vxlan source-interface Loopback0
   no vxlan controller-client
   %s
   vxlan vlan 10 vni 10
   vxlan vlan 10 flood vtep 3.3.3.3 4.4.4.4
   vxlan flood vtep 1.1.1.1 2.2.2.2
   no vxlan vlan flood vtep
   no vxlan learn-restrict vtep
   no vxlan vlan learn-restrict vtep
!
`
	tests := [...]struct {
		in   string
		want string
	}{
		{"vxlan udp-port 4789", "4789"},
		{"vxlan udp-port 1024", "1024"},
		{"vxlan udp-port 65534", "65534"},
		{"vxlan udp-port 1024", "1024"},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := VxlanParseUDPPort(&v, testConfig); got != tt.want {
			t.Fatalf("parseUDPPort() = %q; want %q", got, tt.want)
		}
	}
}

func TestVxlanParseVlans_UnitTest(t *testing.T) {
	var v VxlanInterfaceEntity
	shortConfig := `
interface Vxlan1
   no description
   no shutdown
   vxlan multicast-group 239.10.10.10
   vxlan source-interface Loopback0
   no vxlan controller-client
   vxlan udp-port 4789
   vxlan vlan 10 vni 10
   vxlan vlan 10 flood vtep 3.3.3.3 4.4.4.4
   vxlan flood vtep 1.1.1.1 2.2.2.2
   no vxlan vlan flood vtep
   no vxlan learn-restrict vtep
   no vxlan vlan learn-restrict vtep
!
`
	VxlanParseVlans(&v, shortConfig)
}

func TestVxlanParseFloodList_UnitTest(t *testing.T) {
	var v VxlanInterfaceEntity
	shortConfig := `
interface Vxlan1
   no description
   no shutdown
   vxlan multicast-group 239.10.10.10
   vxlan source-interface Loopback0
   no vxlan controller-client
   vxlan udp-port 4789
   vxlan vlan 10 vni 10
   vxlan vlan 10 flood vtep 3.3.3.3 4.4.4.4
   %s
   no vxlan vlan flood vtep
   no vxlan learn-restrict vtep
   no vxlan vlan learn-restrict vtep
!
`
	tests := [...]struct {
		in   string
		want string
	}{
		{"no vxlan flood vtep", ""},
		{"vxlan flood vtep 1.1.1.1", "1.1.1.1"},
		{"vxlan flood vtep 1.1.1.1 2.2.2.2", "1.1.1.1,2.2.2.2"},
		{"vxlan flood vtep 1.1.1.1 2.2.2.2 3.4.5.6", "1.1.1.1,2.2.2.2,3.4.5.6"},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := VxlanParseFloodList(&v, testConfig); got != tt.want {
			t.Fatalf("parseSourceInterface() = %q; want %q", got, tt.want)
		}
	}
}

func TestVxlanInterfaceGet_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vxlan.Get("Vxlan1")
	}
}

func TestVxlanInterfaceSetSourceInterface_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.SetSourceInterface("Vxlan1", "Loopback0"); !ok {
			t.Fatalf("SetSourceInterface(Vxlan1, Loopback0) failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("vxlan source-interface Loopback0", config); !found {
			t.Fatalf("\"vxlan source-interface Loopback0\" expected but not seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceSetSourceInterfaceNegate_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
			"vxlan source-interface Loopback0",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.SetSourceInterface("Vxlan1", ""); !ok {
			t.Fatalf("SetSourceInterface(Vxlan1, \"\") failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("vxlan source-interface Loopback0", config); found {
			t.Fatalf("\"vxlan source-interface Loopback0\" NOT expected but seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceSetSourceInterfaceDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
			"vxlan source-interface Loopback0",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.SetSourceInterfaceDefault("Vxlan1"); !ok {
			t.Fatalf("SetSourceInterfaceDefault(Vxlan1) failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("no vxlan source-interface", config); !found {
			t.Fatalf("\"no vxlan source-interface\" expected but not seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceSetMulticastGroup_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.SetMulticastGroup("Vxlan1", "239.10.10.10"); !ok {
			t.Fatalf("SetMulticastGroup(Vxlan1, 239.10.10.10) failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("vxlan multicast-group 239.10.10.10", config); !found {
			t.Fatalf("\"vxlan multicast-group 239.10.10.10\" expected but not seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceSetMulticastGroupNegate_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
			"vxlan multicast-group 239.10.10.10",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.SetMulticastGroup("Vxlan1", ""); !ok {
			t.Fatalf("SetMulticastGroup(Vxlan1, \"\") failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("no vxlan multicast-group", config); !found {
			t.Fatalf("\"no vxlan multicast-group\" expected but not seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceSetMulticastGroupDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
			"vxlan multicast-group 239.10.10.10",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.SetMulticastGroupDefault("Vxlan1"); !ok {
			t.Fatalf("SetMulticastGroupDefault(Vxlan1) failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("no vxlan multicast-group", config); !found {
			t.Fatalf("\"no vxlan multicast-group\" expected but not seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceSetUdpPort_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.SetUDPPort("Vxlan1", 1024); !ok {
			t.Fatalf("SetUDPPortDefault(Vxlan1) failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("vxlan udp-port 1024", config); !found {
			t.Fatalf("\"vxlan udp-port 1024\" expected but not seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceSetUdpPortInvalid_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
			"vxlan udp-port 1024",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.SetUDPPort("Vxlan1", 0); !ok {
			t.Fatalf("SetUDPPort(Vxlan1, ) failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("vxlan udp-port 4789", config); !found {
			t.Fatalf("\"vxlan udp-port 4789\" expected but not seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceSetUdpPortDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
			"vxlan udp-port 1024",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.SetUDPPortDefault("Vxlan1"); !ok {
			t.Fatalf("SetUDPPortDefault(Vxlan1) failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("vxlan udp-port 4789", config); !found {
			t.Fatalf("\"vxlan udp-port 4789\" expected but not seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceAddVtep_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.AddVtepGlobalFlood("Vxlan1", "1.1.1.1"); !ok {
			t.Fatalf("AddVtepGlobalFlood(Vxlan1, 1.1.1.1) failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("vxlan flood vtep 1.1.1.1", config); !found {
			t.Fatalf("\"vxlan flood vtep 1.1.1.1\" expected but not seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceAddVtepToVlan_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.AddVtepLocalFlood("Vxlan1", "1.1.1.1", 10); !ok {
			t.Fatalf("AddVtepLocalFlood(Vxlan1, 1.1.1.1, 10) failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("vxlan vlan 10 flood vtep 1.1.1.1", config); !found {
			t.Fatalf("\"vxlan vlan 10 flood vtep 1.1.1.1\" expected but not seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceRemoveVtep_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
			"vxlan flood vtep add 1.1.1.1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.RemoveVtepGlobalFlood("Vxlan1", "1.1.1.1"); !ok {
			t.Fatalf("RemoveVtepGlobalFlood(Vxlan1, 1.1.1.1, 10) failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("vxlan flood vtep 1.1.1.1", config); found {
			t.Fatalf("\"vxlan flood vtep 1.1.1.1\" NOT expected but seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceRemoveVtepFromVlan_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
			"vxlan vlan 10 flood vtep add 1.1.1.1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.RemoveVtepLocalFlood("Vxlan1", "1.1.1.1", 10); !ok {
			t.Fatalf("RemoveVtepLocalFlood(Vxlan1, 1.1.1.1, 10) failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("vxlan vlan 10 flood vtep 1.1.1.1", config); found {
			t.Fatalf("\"vxlan vlan 10 flood vtep 1.1.1.1\" NOT expected but seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceUpdateVlan_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
			"vxlan vlan 10 vni 10",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.UpdateVlan("Vxlan1", 10, 10); !ok {
			t.Fatalf("UpdateVlan(Vxlan1, 10, 10) failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("vxlan vlan 10 vni 10", config); !found {
			t.Fatalf("\"vxlan vlan 10 vni 10\" expected but not seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}

func TestVxlanInterfaceRemoveVlan_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vxlan := Vxlan(dut)

		cmds := []string{
			"no interface vxlan1",
			"interface vxlan1",
			"vxlan vlan 10 vni 10",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := vxlan.RemoveVlan("Vxlan1", 10); !ok {
			t.Fatalf("RemoveVlan(Vxlan1, 10) failed")
		}

		config, _ := vxlan.GetBlock("interface Vxlan1")
		if found, _ := regexp.MatchString("vxlan vlan 10 vni 10", config); found {
			t.Fatalf("\"vxlan vlan 10 vni 10\" Not expected but seen under "+
				"interface Vxlan1 section.\n[%s]", config)
		}
	}
}
