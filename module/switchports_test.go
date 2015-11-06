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
	"strings"
	"testing"
)

var interfaceList = []string{
	"Ethernet1",
	"Ethernet1/1",
	"Port-Channel1",
}

/**
 *****************************************************************************
 * Unit Tests
 *****************************************************************************
 **/
func TestSwitchPortParseMode_UnitTest(t *testing.T) {
	var s SwitchPortEntity

	shortConfig := `
interface Port-Channel10
    no description
    no shutdown
    default load-interval
    logging event link-status use-global
    switchport access vlan 1
    switchport trunk native vlan 1
    switchport trunk allowed vlan 1-4094
    %s
    switchport mac address learning
    no switchport private-vlan mapping
    switchport
    default encapsulation dot1q vlan`

	tests := []struct {
		in   string
		want string
	}{
		{"switchport mode access", "access"},
		{"switchport mode dot1q-tunnel", "dot1q"},
		{"switchport mode trunk", "trunk"},
	}

	for _, tt := range tests {
		testConf := fmt.Sprintf(shortConfig, tt.in)
		if got := SwPortParseMode(&s, testConf); got != tt.want {
			t.Fatalf("parseMode() = %q; want %q", got, tt.want)
		}
	}
	if got := SwPortParseMode(&s, ""); got != "" {
		t.Fatalf("parseMode() = %q; want \"\"", got)
	}
}

func TestSwitchPortParseTrunkGroups_UnitTest(t *testing.T) {
	var s SwitchPortEntity

	shortConfig := `
    interface Ethernet1
        switchport access vlan 1
        switchport trunk native vlan 1
        switchport trunk allowed vlan 1-4094
        switchport mode access
        switchport mac address learning
        switchport trunk group %s
        switchport trunk group %s
        switchport trunk group %s
        switchport trunk group %s
        switchport trunk group %s
        switchport trunk group %s
        switchport trunk group %s
        switchport trunk group %s
        no switchport private-vlan mapping
        switchport
    `

	tests := [8]struct {
		in   string
		want string
	}{}

	var tg [len(tests)]string

	// for each test entry
	for idx := range tests {
		// get the random strings
		for i := range tg {
			tg[i] = RandomString(2, 14)
		}
		testConf := fmt.Sprintf(shortConfig, tg[0], tg[1], tg[2],
			tg[3], tg[4], tg[5], tg[6], tg[7])
		tests[idx].in = testConf
		tests[idx].want = strings.Join(tg[:], ",")
	}
	for _, tt := range tests {
		got := SwPortParseTrunkGroups(&s, tt.in)
		if strings.Compare(got, tt.want) != 0 {
			t.Fatalf("parseTrunkGroups() = %q; want %q", got, tt.want)
		}
	}
	if got := SwPortParseTrunkGroups(&s, ""); got != "" {
		t.Fatalf("parseTrunkGroups() = %q; want \"\"", got)
	}
}

func TestSwitchPortParseAccessVlan_UnitTest(t *testing.T) {
	var s SwitchPortEntity

	shortConfig := `
interface Port-Channel10
    no description
    no shutdown
    default load-interval
    logging event link-status use-global
    %s
    switchport trunk native vlan 1
    switchport trunk allowed vlan 1-4094
    switchport mode access
    switchport mac address learning
    no switchport private-vlan mapping
    switchport
    default encapsulation dot1q vlan`

	tests := []struct {
		in   string
		want string
	}{
		{"switchport access vlan 1", "1"},
		{"switchport access vlan 3", "3"},
		{"switchport access vlan 10", "10"},
		{"switchport access vlan 100", "100"},
	}

	for _, tt := range tests {
		testConf := fmt.Sprintf(shortConfig, tt.in)
		if got := SwPortParseAccessVlan(&s, testConf); got != tt.want {
			t.Fatalf("parseAccessVlan() = %q; want %q", got, tt.want)
		}
	}
	if got := SwPortParseAccessVlan(&s, ""); got != "" {
		t.Fatalf("parseAccessVlan() = %q; want \"\"", got)
	}
}
func TestSwitchPortParseTrunkNativeVlan_UnitTest(t *testing.T) {
	var s SwitchPortEntity

	shortConfig := `
interface Port-Channel10
    no description
    no shutdown
    default load-interval
    logging event link-status use-global
    switchport access vlan 1
    %s
    switchport trunk allowed vlan 1-4094
    switchport mode access
    switchport mac address learning
    no switchport private-vlan mapping
    switchport
    default encapsulation dot1q vlan`

	tests := []struct {
		in   string
		want string
	}{
		{"switchport trunk native vlan 1", "1"},
		{"switchport trunk native vlan 3", "3"},
		{"switchport trunk native vlan 10", "10"},
		{"switchport trunk native vlan 100", "100"},
	}

	for _, tt := range tests {
		testConf := fmt.Sprintf(shortConfig, tt.in)
		if got := SwPortParseTrunkNativeVlan(&s, testConf); got != tt.want {
			t.Fatalf("parseTrunkNativeVlan() = %q; want %q", got, tt.want)
		}
	}
	if got := SwPortParseTrunkNativeVlan(&s, ""); got != "" {
		t.Fatalf("parseTrunkNativeVlan() = %q; want \"\"", got)
	}
}
func TestSwitchPortParseTrunkAllowedVlans_UnitTest(t *testing.T) {
	var s SwitchPortEntity

	shortConfig := `
interface Port-Channel10
    no description
    no shutdown
    default load-interval
    logging event link-status use-global
    switchport access vlan 1
    switchport trunk native vlan 1
    %s
    switchport mode access
    switchport mac address learning
    no switchport private-vlan mapping
    switchport
    default encapsulation dot1q vlan`

	tests := []struct {
		in   string
		want string
	}{
		{"switchport trunk allowed vlan none", "none"},
		{"switchport trunk allowed vlan 1-4094", "1-4094"},
		{"switchport trunk allowed vlan 1-2", "1-2"},
		{"switchport trunk allowed vlan 42", "42"},
		{"switchport trunk allowed vlan 1,100", "1,100"},
		{"switchport trunk allowed vlan 1,100-105", "1,100-105"},
	}

	for _, tt := range tests {
		testConf := fmt.Sprintf(shortConfig, tt.in)
		if got := SwPortParseTrunkAllowedVlans(&s, testConf); got != tt.want {
			t.Fatalf("parseTrunkAllowedVlans() = %q; want %q", got, tt.want)
		}
	}
	if got := SwPortParseTrunkAllowedVlans(&s, ""); got != "" {
		t.Fatalf("parseTrunkAllowedVlans() = %q; want \"\"", got)
	}
}

func TestSwitchPortCreate_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"no ip address",
			"switchport",
		}
		sp.Create(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestSwitchPortDelete_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"no switchport",
		}
		sp.Delete(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestSwitchPortDefault_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"no ip address",
			"default switchport",
		}
		sp.Default(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestSwitchPortSetMode_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)
	tests := []struct {
		mode string
		want string
	}{
		{"access", "switchport mode access"},
		{"trunk", "switchport mode trunk"},
	}

	for _, intf := range interfaceList {

		for _, tt := range tests {
			cmds := []string{
				"interface " + intf,
				"switchport mode " + tt.mode,
			}
			sp.SetMode(intf, tt.mode)
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

func TestSwitchPortSetModeDefault_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"default switchport mode",
		}
		sp.SetModeDefault(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestSwitchPortSetAccessVlan_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)
	for _, intf := range interfaceList {
		vid := strconv.Itoa(RandomInt(2, 4094))

		cmds := []string{
			"interface " + intf,
			"switchport access vlan " + vid,
		}
		sp.SetAccessVlan(intf, vid)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestSwitchPortSetAccessVlanDefault_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"default switchport access vlan",
		}
		sp.SetAccessVlanDefault(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestSwitchPortSetTrunkNativeVlan_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)
	for _, intf := range interfaceList {
		vid := strconv.Itoa(RandomInt(2, 4094))

		cmds := []string{
			"interface " + intf,
			"switchport trunk native vlan " + vid,
		}
		sp.SetTrunkNativeVlan(intf, vid)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestSwitchPortSetTrunkNativeVlanDefault_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"default switchport trunk native vlan",
		}
		sp.SetTrunkNativeVlanDefault(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestSwitchPortSetTrunkAllowedVlans_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)
	for _, intf := range interfaceList {
		vid := "1,2,3-5,6,7"

		cmds := []string{
			"interface " + intf,
			"switchport trunk allowed vlan " + vid,
		}
		sp.SetTrunkAllowedVlans(intf, vid)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestSwitchPortSetTrunkAllowedVlansDefault_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"default switchport trunk allowed vlan",
		}
		sp.SetTrunkAllowedVlansDefault(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestSwitchPortSetTrunkGroupDefault_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"default switchport trunk group",
		}
		sp.SetTrunkGroupsDefault(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestSwitchPortAddTrunkGroup_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)
	tg := RandomString(1, 32)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"switchport trunk group " + tg,
		}
		sp.AddTrunkGroup(intf, tg)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestSwitchPortRemoveTrunkGroup_UnitTest(t *testing.T) {
	sp := SwitchPort(dummyNode)
	tg := RandomString(1, 32)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"no switchport trunk group " + tg,
		}
		sp.RemoveTrunkGroup(intf, tg)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

/**
 *****************************************************************************
 * System Tests
 *****************************************************************************
 **/
func getEthInterfaces(s *SwitchPortEntity) []string {
	var re = regexp.MustCompile(`(?m)^interface\s(Eth.+)`)
	config := s.Config()

	interfaces := re.FindAllStringSubmatch(config, -1)

	response := make([]string, len(interfaces))

	for idx, name := range interfaces {
		response[idx] = name[1]
	}
	return response
}

func TestSwitchPortGet_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		config := sp.Get(intf)
		if config == nil {
			t.Fatalf("SwitchPort Get for %s failed.", intf)
		}
		if config.Mode() != "access" ||
			config.AccessVlan() != "1" ||
			config.TrunkNativeVlan() != "1" ||
			config.TrunkAllowedVlans() != "1-4094" ||
			config.Name() != intf ||
			config.TrunkGroups() != "" {
			t.Fatalf("SwitchPort Get for %s failed.", intf)
		}
	}
}

func TestSwitchPortGetReturnsNil_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		config := sp.Get(intf)
		if config != nil {
			t.Fatalf("SwitchPort Get for %s failed.", intf)
		}
	}
}

func TestSwitchPortGetAll_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		configs := sp.GetAll()
		if configs == nil {
			t.Fatalf("SwitchPort Get for %s failed.", intf)
		}
		if configs[intf]["name"] != intf {
			t.Fatalf("SwitchPort Get for %s failed. Got %s", intf, configs[intf]["name"])
		}
	}
}

func TestSwitchPortCreateTrue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.Create(intf); !ok {
			t.Fatalf("SwitchPort Create for %s failed.", intf)
		}

		section := sp.GetSection(intf)
		if found, _ := regexp.MatchString(`no switchport\s*\n`, section); found {
			t.Fatalf("\"no switchport\" NOT expected but not seen under "+
				"Interface %s", intf)
		}

		if ok := dut.Config("default interface " + intf); !ok {
			t.Fatalf("dut.Config() failure")
		}

	}
}

func TestSwitchPortDeleteTrue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.Delete(intf); !ok {
			t.Fatalf("SwitchPort Delete for %s failed.", intf)
		}

		section := sp.GetSection(intf)
		if found, _ := regexp.MatchString(`no switchport\s*\n`, section); !found {
			t.Fatalf("\"no switchport\" expected but not seen under "+
				"Interface %s", intf)
		}

		if ok := dut.Config("default interface " + intf); !ok {
			t.Fatalf("dut.Config() failure")
		}

	}
}

func TestSwitchPortDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.Default(intf); !ok {
			t.Fatalf("SwitchPort Default for %s failed.", intf)
		}

		section := sp.GetSection(intf)
		if found, _ := regexp.MatchString(`no switchport\s*\n`, section); found {
			t.Fatalf("\"no switchport\" NOT expected but not seen under "+
				"Interface %s", intf)
		}

		if ok := dut.Config("default interface " + intf); !ok {
			t.Fatalf("dut.Config() failure")
		}

	}
}

func TestSwitchPortSetMode_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.SetMode(intf, "dot1q-tunnel"); !ok {
			t.Fatalf("SwitchPort SetMode for %s failed.", intf)
		}

		section := sp.GetSection(intf)
		if found, _ := regexp.MatchString(`switchport mode dot1q-tunnel\s*\n`, section); !found {
			t.Fatalf("\"no switchport\" expected but not seen under "+
				"Interface %s", intf)
		}

		if ok := dut.Config("default interface " + intf); !ok {
			t.Fatalf("dut.Config() failure")
		}

	}
}

func TestSwitchPortSetModeDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"switchport mode trunk",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.SetModeDefault(intf); !ok {
			t.Fatalf("SwitchPort Default for %s failed.", intf)
		}

		section := sp.GetSection(intf)
		if found, _ := regexp.MatchString(`switchport mode trunk\s*\n`, section); found {
			t.Fatalf("\"switchport mode trunk\" NOT expected but not seen under "+
				"Interface %s", intf)
		}

		if ok := dut.Config("default interface " + intf); !ok {
			t.Fatalf("dut.Config() failure")
		}

	}
}

func TestSwitchPortSetAccessVlan_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
			"vlan 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.SetAccessVlan(intf, "100"); !ok {
			t.Fatalf("SwitchPort SetAccessVlan for %s failed.", intf)
		}

		section := sp.GetSection(intf)
		if found, _ := regexp.MatchString(`switchport access vlan 100\s*\n`,
			section); !found {
			t.Fatalf("\"switchport access vlan 100\" expected but not seen under "+
				"Interface %s", intf)
		}

		if ok := dut.Config("default interface " + intf); !ok {
			t.Fatalf("dut.Config() failure")
		}
	}
}

func TestSwitchPortSetAccessVlanDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
			"vlan 100",
			"interface " + intf,
			"switchport access vlan 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.SetAccessVlanDefault(intf); !ok {
			t.Fatalf("SwitchPort SetAccessVlanDefault for %s failed.", intf)
		}

		section := sp.GetSection(intf)
		if found, _ := regexp.MatchString(`switchport access vlan 100\s*\n`,
			section); found {
			t.Fatalf("\"switchport access vlan 100\" NOT expected but not seen under "+
				"Interface %s", intf)
		}

		if ok := dut.Config("default interface " + intf); !ok {
			t.Fatalf("dut.Config() failure")
		}
	}
}

func TestSwitchPortSetTrunkNativeVlan_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
			"vlan 100",
			"interface " + intf,
			"switchport mode trunk",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.SetTrunkNativeVlan(intf, "100"); !ok {
			t.Fatalf("SwitchPort SetTrunkNativeVlan for %s failed.", intf)
		}

		section := sp.GetSection(intf)
		if found, _ := regexp.MatchString(`switchport trunk native vlan 100\s*\n`,
			section); !found {
			t.Fatalf("\"switchport trunk native vlan 100\" expected but not seen under "+
				"Interface %s", intf)
		}

		if ok := dut.Config("default interface " + intf); !ok {
			t.Fatalf("dut.Config() failure")
		}
	}
}

func TestSwitchPortSetTrunkNativeVlanDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
			"vlan 100",
			"interface " + intf,
			"switchport mode trunk",
			"switchport trunk native vlan 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.SetTrunkNativeVlanDefault(intf); !ok {
			t.Fatalf("SwitchPort SetTrunkNativeVlanDefault for %s failed.", intf)
		}

		section := sp.GetSection(intf)
		if found, _ := regexp.MatchString(`switchport trunk native vlan 100\s*\n`,
			section); found {
			t.Fatalf("\"switchport trunk native vlan 100\" NOT expected but not seen under "+
				"Interface %s", intf)
		}

		if ok := dut.Config("default interface " + intf); !ok {
			t.Fatalf("dut.Config() failure")
		}
	}
}

func TestSwitchPortSetTrunkAllowedVlans_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
			"vlan 100",
			"interface " + intf,
			"switchport mode trunk",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.SetTrunkAllowedVlans(intf, "1,10,100"); !ok {
			t.Fatalf("SwitchPort SetTrunkAllowedVlans for %s failed.", intf)
		}

		section := sp.GetSection(intf)
		if found, _ := regexp.MatchString(`switchport trunk allowed vlan 1,10,100\s*\n`,
			section); !found {
			t.Fatalf("\"switchport trunk allowed vlan 1,10,100\" expected but not seen under "+
				"Interface %s", intf)
		}

		if ok := dut.Config("default interface " + intf); !ok {
			t.Fatalf("dut.Config() failure")
		}
	}
}

func TestSwitchPortSetTrunkAllowedVlansDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))

		cmds := []string{
			"default interface " + intf,
			"vlan 100",
			"interface " + intf,
			"switchport mode trunk",
			"switchport trunk allowed vlan 1,10,100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.SetTrunkAllowedVlansDefault(intf); !ok {
			t.Fatalf("SwitchPort SetTrunkAllowedVlansDefault for %s failed.", intf)
		}

		section := sp.GetSection(intf)
		if found, _ := regexp.MatchString(`switchport trunk allowed vlan 1,10,100\s*\n`,
			section); found {
			t.Fatalf("\"switchport trunk allowed vlan 1,10,100\" NOT expected but not seen under "+
				"Interface %s", intf)
		}

		if ok := dut.Config("default interface " + intf); !ok {
			t.Fatalf("dut.Config() failure")
		}
	}
}

func TestSwitchPortSetTrunkGroups_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))
		tg1 := RandomString(1, 32)
		tg2 := RandomString(1, 10)
		tg3 := RandomString(1, 10)
		cmds := []string{
			"default interface " + intf,
			"vlan 100",
			"interface " + intf,
			"switchport mode trunk",
			"switchport trunk group " + tg1,
			"switchport trunk group " + tg2,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure: Commands: %q", cmds)
		}

		confTgs := []string{tg1, tg3}

		if ok := sp.SetTrunkGroups(intf, confTgs); !ok {
			t.Fatalf("Failure Setting Trunk Group to %q on intf %s", confTgs, intf)
		}

		section := sp.GetSection(intf)
		str := "switchport trunk group " + tg2
		if found, _ := regexp.MatchString(str, section); found {
			t.Fatalf("\"%s\" NOT expected but seen under "+
				"Interface %s", tg2, intf)
		}

		// if ok := dut.Config("default interface "+ intf); !ok {
		//    t.Fatalf("dut.Config() failure")
		// }
	}
}

func TestSwitchPortSetTrunkGroupsDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)

		intf := RandomChoice(getEthInterfaces(sp))
		tg1 := RandomString(1, 32)
		cmds := []string{
			"default interface " + intf,
			"vlan 100",
			"interface " + intf,
			"switchport mode trunk",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.SetTrunkGroups(intf, []string{tg1}); !ok {
			t.Fatalf("Failure Setting Trunk Group to %s on intf %s", tg1, intf)
		}

		section := sp.GetSection(intf)
		str := "switchport trunk group " + tg1
		if found, _ := regexp.MatchString(str, section); !found {
			t.Fatalf("\"%s\" expected but not seen under "+
				"Interface %s", str, intf)
		}

		if ok := sp.SetTrunkGroupsDefault(intf); !ok {
			t.Fatalf("Failure Setting Trunk Group default intf %s", intf)
		}

		section = sp.GetSection(intf)
		str = "switchport trunk group " + tg1
		if found, _ := regexp.MatchString(str, section); found {
			t.Fatalf("\"%s\" NOT expected but not seen under "+
				"Interface %s", str, intf)
		}

		// if ok := dut.Config("default interface "+ intf); !ok {
		//    t.Fatalf("dut.Config() failure")
		// }
	}
}

func TestSwitchPortAddTrunkGroup_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)
		intf := RandomChoice(getEthInterfaces(sp))
		tg := RandomString(1, 32)

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport trunk group",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.AddTrunkGroup(intf, tg); !ok {
			t.Fatalf("Failure Adding Trunk Group to %s on intf %s", tg, intf)
		}

		section := sp.GetSection(intf)
		str := "switchport trunk group " + tg
		if found, _ := regexp.MatchString(str, section); !found {
			t.Fatalf("\"%s\" expected but not seen under "+
				"Interface %s", str, intf)
		}

		if ok := dut.Config("default interface " + intf); !ok {
			t.Fatalf("dut.Config() failure")
		}
	}
}

func TestSwitchPortRemoveTrunkGroup_SystemTest(t *testing.T) {
	for _, dut := range duts {
		sp := SwitchPort(dut)
		intf := RandomChoice(getEthInterfaces(sp))
		tg := RandomString(1, 32)

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"switchport trunk group " + tg,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := sp.RemoveTrunkGroup(intf, tg); !ok {
			t.Fatalf("Failure Removing Trunk Group to %s on intf %s", tg, intf)
		}

		section := sp.GetSection(intf)
		str := "switchport trunk group " + tg
		if found, _ := regexp.MatchString(str, section); found {
			t.Fatalf("\"%s\" NOT expected but not seen under "+
				"Interface %s", str, intf)
		}

		if ok := dut.Config("default interface " + intf); !ok {
			t.Fatalf("dut.Config() failure")
		}
	}
}
