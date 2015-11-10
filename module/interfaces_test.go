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

func TestResourceInterfaceGet_UnitTest(t *testing.T) {
	initFixture()
	i := Interface(dummyNode)

	keys := []string{
		"name",
		"type",
		"shutdown",
		"description",
	}

	config := i.Get("Loopback0")

	for _, key := range keys {
		if _, found := config[key]; !found {
			t.Fatalf("Get(Loopback0) key mismatch expect: %q got %#v", keys, config)
		}
	}
}

func TestResourceInterfaceCreate_UnitTest(t *testing.T) {
	i := Interface(dummyNode)

	for _, intf := range interfaceList {
		cmds := []string{
			"interface " + intf,
		}
		i.Create(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestResourceInterfaceDelete_UnitTest(t *testing.T) {
	i := Interface(dummyNode)

	for _, intf := range interfaceList {
		cmds := []string{
			"no interface " + intf,
		}
		i.Delete(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestResourceInterfaceDefault_UnitTest(t *testing.T) {
	i := Interface(dummyNode)

	for _, intf := range interfaceList {
		cmds := []string{
			"default interface " + intf,
		}
		i.Default(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestResourceInterfaceDescription_UnitTest(t *testing.T) {
	i := Interface(dummyNode)

	for _, intf := range interfaceList {
		cmds := []string{
			"interface " + intf,
			"default description",
		}
		tests := []struct {
			in   string
			want string
		}{
			{"Test description", "description Test description"},
		}

		for _, tt := range tests {
			i.SetDescription(intf, tt.in)
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
}

func TestResourceInterfaceDescriptionDefault_UnitTest(t *testing.T) {
	i := Interface(dummyNode)

	for _, intf := range interfaceList {
		cmds := []string{
			"interface " + intf,
			"default description",
		}
		i.SetDescriptionDefault(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestResourceInterfaceSetShutDown_UnitTest(t *testing.T) {
	i := Interface(dummyNode)

	for _, intf := range interfaceList {
		cmds := []string{
			"interface " + intf,
			"default shutdown",
		}
		tests := []struct {
			shut bool
			want string
		}{
			{true, "shutdown"},
			{false, "no shutdown"},
		}

		for _, tt := range tests {
			i.SetShutdown(intf, tt.shut)
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
}

func TestResourceInterfaceSetShutDownDefault_UnitTest(t *testing.T) {
	i := Interface(dummyNode)

	for _, intf := range interfaceList {
		cmds := []string{
			"interface " + intf,
			"default shutdown",
		}
		i.SetShutdownDefault(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestEthernetInterfaceParseSflow_UnitTest(t *testing.T) {
	shortConfig := `
interface Ethernet1
   no shutdown
   default load-interval
   logging event link-status use-global
   uc-tx-queue 7
      priority strict
      no bandwidth percent
      no shape rate
      no bandwidth guaranteed
   %s
   no storm-control broadcast
   no storm-control multicast
   no storm-control all
`
	tests := [...]struct {
		in   string
		want bool
	}{
		{"sflow enable", true},
		{"no sflow", false},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := parseSflow(testConfig); got != tt.want {
			t.Fatalf("parseSflow() = %t; want %t", got, tt.want)
		}
	}
}

func TestEthernetInterfaceParseFlowControlSend_UnitTest(t *testing.T) {
	shortConfig := `
interface Ethernet1
   no shutdown
   default load-interval
   logging event link-status use-global
   no dcbx mode
   no mac-address
   no link-debounce
   %s
   no flowcontrol receive
   no mac timestamp
   no speed
`
	tests := [...]struct {
		in   string
		want string
	}{
		{"flowcontrol send on", "on"},
		{"flowcontrol send off", "off"},
		{"no flowcontrol send", "off"},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := parseFlowControlSend(testConfig); got != tt.want {
			t.Fatalf("parseFlowControlSend() = %q; want %q", got, tt.want)
		}
	}
}

func TestEthernetInterfaceParseFlowControlReceive_UnitTest(t *testing.T) {
	shortConfig := `
interface Ethernet1
   no shutdown
   default load-interval
   logging event link-status use-global
   no dcbx mode
   no mac-address
   no link-debounce
   no flowcontrol send
   %s
   no mac timestamp
   no speed
`
	tests := [...]struct {
		in   string
		want string
	}{
		{"flowcontrol receive on", "on"},
		{"flowcontrol receive off", "off"},
		{"no flowcontrol receive", "off"},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := parseFlowControlReceive(testConfig); got != tt.want {
			t.Fatalf("parseFlowControlReceive() = %q; want %q", got, tt.want)
		}
	}
}

func TestEthernetInterfaceGetCheckKeys_UnitTest(t *testing.T) {
	initFixture()
	i := EthernetInterface(dummyNode)

	keys := []string{
		"name",
		"shutdown",
		"description",
		"sflow",
		"flowcontrol_send",
		"flowcontrol_receive",
		"type",
	}

	intf := i.Get("Ethernet1")

	if len(keys) != len(intf) {
		t.Fatalf("Keys mismatch. Expect: %q got %#v", keys, intf)
	}
	for _, val := range keys {
		if _, found := intf[val]; !found {
			t.Fatalf("Key \"%s\" not found in neighbor", val)
		}
	}
}

func TestEthernetInterfaceCreate_UnitTest(t *testing.T) {
	i := EthernetInterface(dummyNode)
	if ok := i.Create("Ethernet1"); ok {
		t.Fatalf("Should not allow")
	}
}

func TestEthernetInterfaceDelete_UnitTest(t *testing.T) {
	i := EthernetInterface(dummyNode)
	if ok := i.Delete("Ethernet1"); ok {
		t.Fatalf("Should not allow")
	}
}

func TestEthernetInterfaceSetFlowcontrolSend_UnitTest(t *testing.T) {
	i := EthernetInterface(dummyNode)

	for _, intf := range interfaceList {
		cmds := []string{
			"interface " + intf,
			"default flowcontrol send",
		}
		tests := []struct {
			enable bool
			want   string
		}{
			{true, "flowcontrol send on"},
			{false, "flowcontrol send off"},
		}

		for _, tt := range tests {
			i.SetFlowcontrolSend(intf, tt.enable)
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
}

func TestEthernetInterfaceSetFlowcontrolReceive_UnitTest(t *testing.T) {
	i := EthernetInterface(dummyNode)

	for _, intf := range interfaceList {
		cmds := []string{
			"interface " + intf,
			"default flowcontrol send",
		}
		tests := []struct {
			enable bool
			want   string
		}{
			{true, "flowcontrol receive on"},
			{false, "flowcontrol receive off"},
		}

		for _, tt := range tests {
			i.SetFlowcontrolReceive(intf, tt.enable)
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
}

func TestEthernetInterfaceDisableFlowcontrolSend_UnitTest(t *testing.T) {
	i := EthernetInterface(dummyNode)

	for _, intf := range interfaceList {
		cmds := []string{
			"interface " + intf,
			"no flowcontrol send",
		}
		i.DisableFlowcontrolSend(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestEthernetInterfaceDisableFlowcontrolReceive_UnitTest(t *testing.T) {
	i := EthernetInterface(dummyNode)

	for _, intf := range interfaceList {
		cmds := []string{
			"interface " + intf,
			"no flowcontrol receive",
		}
		i.DisableFlowcontrolReceive(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestEthernetInterfaceSetSflow_UnitTest(t *testing.T) {
	i := EthernetInterface(dummyNode)

	for _, intf := range interfaceList {
		cmds := []string{
			"interface " + intf,
			"default sflow",
		}
		tests := []struct {
			enable bool
			want   string
		}{
			{true, "sflow enable"},
			{false, "no sflow enable"},
		}

		for _, tt := range tests {
			i.SetSflow(intf, tt.enable)
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
}

func TestEthernetInterfaceSetSflowDefault_UnitTest(t *testing.T) {
	i := EthernetInterface(dummyNode)

	for _, intf := range interfaceList {
		cmds := []string{
			"interface " + intf,
			"default sflow",
		}
		i.SetSflowDefault(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
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

func TestPortChannelGetCheckKeys_UnitTest(t *testing.T) {
	initFixture()
	p := PortChannel(dummyNode)

	keys := []string{
		"name",
		"shutdown",
		"description",
		"type",
		"lacp_mode",
		"minimum_links",
		"members",
	}

	pc := p.Get("Port-Channel1")

	if len(keys) != len(pc) {
		t.Fatalf("Keys mismatch. Expect: %q got %#v", keys, pc)
	}
	for _, val := range keys {
		if _, found := pc[val]; !found {
			t.Fatalf("Key \"%s\" not found in neighbor", val)
		}
	}
}

func TestPortChannelSetMinimumLinks_UnitTest(t *testing.T) {
	p := PortChannel(dummyNode)
	cmds := []string{
		"interface Port-Channel1",
		"default port-channel min-links",
	}
	tests := [...]struct {
		val  int
		want string
		rc   bool
	}{
		{0, "", false},
		{4, "port-channel min-links 4", true},
		{8, "port-channel min-links 8", true},
		{17, "", false},
	}

	for test, tt := range tests {
		if got := p.SetMinimumLinks("Port-Channel1", tt.val); got != tt.rc {
			t.Fatalf("Test[%d] Expected \"%t\" got \"%t\"", test, tt.rc, got)
		}
		if tt.rc {
			cmds[1] = tt.want
			// first two commands are 'enable', 'configure terminal'
			commands := dummyConnection.GetCommands()[2:]
			for idx, val := range commands {
				if cmds[idx] != val {
					t.Fatalf("Test[%d] Expected \"%q\" got \"%q\"", test, cmds, commands)
				}
			}
		}
	}
}

func TestPortChannelSetMembers_UnitTest(t *testing.T) {
	p := PortChannel(dummyNode)
	cmds := []string{
		"interface Ethernet6",
		"no channel-group 1",
		"interface Ethernet7",
		"channel-group 1 mode on",
	}
	members := []string{"Ethernet5", "Ethernet7"}
	p.SetMembers("Port-Channel1", members...)

	t.Skip("skipping test")

	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestPortChannelSetLacpMode_UnitTest(t *testing.T) {
	p := PortChannel(dummyNode)
	cmds := []string{
		"interface Ethernet5",
		"no channel-group 1",
		"interface Ethernet6",
		"no channel-group 1",
		"interface Ethernet5",
		"channel-group 1 mode active",
		"interface Ethernet6",
		"channel-group 1 mode active",
	}
	p.SetLacpMode("Port-Channel1", "active")

	t.Skip("skipping test")

	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestPortChannelSetLacpModeInvalid_UnitTest(t *testing.T) {
	p := PortChannel(dummyNode)
	if ok := p.SetLacpMode("Port-Channel1", "InvalidParam"); ok {
		t.Fatalf("Passed/Accepted Invalid parameter")
	}
}

func TestPortChannelSetMinimumLinksDefault_UnitTest(t *testing.T) {
	p := PortChannel(dummyNode)
	cmds := []string{
		"interface Port-Channel1",
		"default port-channel min-links",
	}
	p.SetMinimumLinksDefault("Port-Channel1")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
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
		{"", ""},
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
   %s
   %s
   vxlan flood vtep 1.1.1.1 2.2.2.2
   no vxlan vlan flood vtep
   no vxlan learn-restrict vtep
   no vxlan vlan learn-restrict vtep
!
`
	tests := [...]struct {
		in1 string
		in2 string
	}{
		{"", "no vxlan vlan flood vtep"},
		{"vxlan vlan 10 vni 10", ""},
		{"", "vxlan vlan 10 flood vtep 3.3.3.3 4.4.4.4"},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in1, tt.in2)
		conf := VxlanParseVlans(&v, testConfig)
		if conf == nil {
			t.Fatalf("parseVlans()")
		}
	}
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

func TestVxlanInterfaceGet_UnitTest(t *testing.T) {
	//initFixture()
	v := Vxlan(dummyNode)

	keys := []string{
		"name",
		"type",
		"shutdown",
		"description",
		"source_interface",
		"multicast_group",
		"udp_port",
		"flood_list",
	}

	config := v.Get("Vxlan1")

	for _, key := range keys {
		if _, found := config[key]; !found {
			t.Fatalf("Get(Vxlan1) key mismatch expect: %q got %#v", keys, config)
		}
		fmt.Printf("key: %s config:%s\n", key, config[key])
	}
}

func TestVxlanInterfaceSetSourceInterface_UnitTest(t *testing.T) {
	v := Vxlan(dummyNode)
	cmds := []string{
		"interface Vxlan1",
		"default vxlan source-interface",
	}
	tests := [...]struct {
		val  string
		want string
	}{
		{"", "no vxlan source-interface"},
		{"Loopback0", "vxlan source-interface Loopback0"},
	}

	for _, tt := range tests {
		v.SetSourceInterface("Vxlan1", tt.val)
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

func TestVxlanInterfaceSetSourceInterfaceDefault_UnitTest(t *testing.T) {
	v := Vxlan(dummyNode)
	cmds := []string{
		"interface Vxlan1",
		"default vxlan source-interface",
	}
	v.SetSourceInterfaceDefault("Vxlan1")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVxlanInterfaceSetMulticastGroup_UnitTest(t *testing.T) {
	v := Vxlan(dummyNode)
	cmds := []string{
		"interface Vxlan1",
		"default vxlan multicast-group",
	}
	tests := [...]struct {
		val  string
		want string
	}{
		{"", "no vxlan multicast-group"},
		{"239.10.10.10", "vxlan multicast-group 239.10.10.10"},
	}

	for _, tt := range tests {
		v.SetMulticastGroup("Vxlan1", tt.val)
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

func TestVxlanInterfaceSetMulticastGroupDefault_UnitTest(t *testing.T) {
	v := Vxlan(dummyNode)
	cmds := []string{
		"interface Vxlan1",
		"default vxlan multicast-group",
	}
	v.SetMulticastGroupDefault("Vxlan1")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVxlanInterfaceSetUDPPort_UnitTest(t *testing.T) {
	v := Vxlan(dummyNode)
	cmds := []string{
		"interface Vxlan1",
		"default vxlan udp-port",
	}
	tests := [...]struct {
		val  int
		want string
		rc   bool
	}{
		{1, "", false},
		{1023, "", false},
		{1024, "vxlan udp-port 1024", true},
		{8000, "vxlan udp-port 8000", true},
		{65535, "vxlan udp-port 65535", true},
		{65536, "", false},
	}

	for _, tt := range tests {
		v.SetUDPPort("Vxlan1", tt.val)
		cmds[1] = tt.want
		if tt.rc {
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

func TestVxlanInterfaceSetUDPPortDefault_UnitTest(t *testing.T) {
	v := Vxlan(dummyNode)
	cmds := []string{
		"interface Vxlan1",
		"default vxlan udp-port",
	}
	v.SetUDPPortDefault("Vxlan1")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVxlanInterfaceAddVtepGlobalFlood_UnitTest(t *testing.T) {
	v := Vxlan(dummyNode)
	cmds := []string{
		"interface Vxlan1",
		"vxlan flood vtep add 1.1.1.1",
	}
	v.AddVtepGlobalFlood("Vxlan1", "1.1.1.1")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVxlanInterfaceAddVtepLocalFlood_UnitTest(t *testing.T) {
	v := Vxlan(dummyNode)
	cmds := []string{
		"interface Vxlan1",
		"vxlan vlan 10 flood vtep add 1.1.1.1",
	}
	v.AddVtepLocalFlood("Vxlan1", "1.1.1.1", 10)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVxlanInterfaceRemoveVtepGlobalFlood_UnitTest(t *testing.T) {
	v := Vxlan(dummyNode)
	cmds := []string{
		"interface Vxlan1",
		"vxlan flood vtep remove 1.1.1.1",
	}
	v.RemoveVtepGlobalFlood("Vxlan1", "1.1.1.1")
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVxlanInterfaceRemoveVtepLocalFlood_UnitTest(t *testing.T) {
	v := Vxlan(dummyNode)
	cmds := []string{
		"interface Vxlan1",
		"vxlan vlan 10 flood vtep remove 1.1.1.1",
	}
	v.RemoveVtepLocalFlood("Vxlan1", "1.1.1.1", 10)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVxlanInterfaceUpdateVlan_UnitTest(t *testing.T) {
	v := Vxlan(dummyNode)
	cmds := []string{
		"interface Vxlan1",
		"vxlan vlan 10 vni 10",
	}
	v.UpdateVlan("Vxlan1", 10, 10)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVxlanInterfaceRemoveVlan_UnitTest(t *testing.T) {
	v := Vxlan(dummyNode)
	cmds := []string{
		"interface Vxlan1",
		"no vxlan vlan 10 vni",
	}
	v.RemoveVlan("Vxlan1", 10)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
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
