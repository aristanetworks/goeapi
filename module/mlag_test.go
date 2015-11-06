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
	"testing"
)

/**
 *****************************************************************************
 * Unit Tests
 *****************************************************************************
 **/
func TestMlagParseDomainID_UnitTest(t *testing.T) {
	var m MlagEntity

	shortConfig := `
mac address-table notification host-flap detection window 15
!
mlag configuration
   domain-id %s
   heartbeat-interval 2000
   no local-interface
   no peer-address
   no peer-link
   reload-delay 300
   no reload-delay non-mlag
   no reload-delay mode
   no shutdown
!
no mpls routing`

	tests := [11]struct {
		in   string
		want string
	}{
		{"no domain-id", ""},
	}

	var ts string
	for i := 1; i < len(tests); i++ {
		ts = RandomString(2, 14)
		testConfig := fmt.Sprintf(shortConfig, ts)
		tests[i].in = testConfig
		tests[i].want = ts
	}

	for _, tt := range tests {
		if got := MlagParseDomainID(&m, tt.in); got != tt.want {
			t.Fatalf("parseDomainID() = %q; want %q", got, tt.want)
		}
	}

	if got := MlagParseDomainID(&m, ""); got != "" {
		t.Fatalf("parseDomainID() = %q; want \"\"", got)
	}
}

func TestMlagParseLocalInterface_UnitTest(t *testing.T) {
	var m MlagEntity

	shortConfig := `
mac address-table notification host-flap detection window 15
!
mlag configuration
   no domain-id
   heartbeat-interval 2000
   %s
   no peer-address
   no peer-link
   reload-delay 300
   no reload-delay non-mlag
   no reload-delay mode
   no shutdown
!
no mpls routing`

	tests := []struct {
		in   string
		want string
	}{
		{"no local-interface", ""},
		{"local-interface Vlan2", "Vlan2"},
		{"local-interface Vlan100", "Vlan100"},
	}

	for _, tt := range tests {
		testConf := fmt.Sprintf(shortConfig, tt.in)
		if got := MlagParseLocalInterface(&m, testConf); got != tt.want {
			t.Fatalf("parseLocalInterface() = %q; want %q", got, tt.want)
		}
	}
	if got := MlagParseLocalInterface(&m, ""); got != "" {
		t.Fatalf("parseLocalInterface() = %q; want \"\"", got)
	}
}

func TestMlagParsePeerAddress_UnitTest(t *testing.T) {
	var m MlagEntity

	shortConfig := `
mac address-table notification host-flap detection window 15
!
mlag configuration
   no domain-id
   heartbeat-interval 2000
   no local-interface
   %s
   no peer-link
   reload-delay 300
   no reload-delay non-mlag
   no reload-delay mode
   no shutdown
!
no mpls routing`

	tests := []struct {
		in   string
		want string
	}{
		{"no peer-address", ""},
		{"peer-address 1.1.1.1", "1.1.1.1"},
		{"peer-address 2.2.2.2", "2.2.2.2"},
		{"peer-address 3.3.3.3", "3.3.3.3"},
	}

	for _, tt := range tests {
		testConf := fmt.Sprintf(shortConfig, tt.in)
		if got := MlagParsePeerAddress(&m, testConf); got != tt.want {
			t.Fatalf("parsePeerAddress() = %q; want %q", got, tt.want)
		}
	}
	if got := MlagParsePeerAddress(&m, ""); got != "" {
		t.Fatalf("parsePeerAddress() = %q; want \"\"", got)
	}
}

func TestMlagParsePeerLink_UnitTest(t *testing.T) {
	var m MlagEntity

	shortConfig := `
mac address-table notification host-flap detection window 15
!
mlag configuration
   no domain-id
   heartbeat-interval 2000
   no local-interface
   no peer-address
   %s
   reload-delay 300
   no reload-delay non-mlag
   no reload-delay mode
   no shutdown
!
no mpls routing`

	tests := []struct {
		in   string
		want string
	}{
		{"no peer-link", ""},
		{"peer-link Ethernet1", "Ethernet1"},
		{"peer-link Ethernet10", "Ethernet10"},
		{"peer-link Port-Channel2", "Port-Channel2"},
		{"peer-link Port-Channel20", "Port-Channel20"},
	}

	for _, tt := range tests {
		testConf := fmt.Sprintf(shortConfig, tt.in)
		if got := MlagParsePeerLink(&m, testConf); got != tt.want {
			t.Fatalf("parsePeerLink() = %q; want %q", got, tt.want)
		}
	}
	if got := MlagParsePeerLink(&m, ""); got != "" {
		t.Fatalf("parsePeerLink() = %q; want \"\"", got)
	}
}

func TestMlagParseShutdown_UnitTest(t *testing.T) {
	var m MlagEntity

	shortConfig := `
mac address-table notification host-flap detection window 15
!
mlag configuration
   no domain-id
   heartbeat-interval 2000
   no local-interface
   no peer-address
   no peer-link
   reload-delay 300
   no reload-delay non-mlag
   no reload-delay mode
   %s
!
no mpls routing`

	tests := []struct {
		in   string
		want bool
	}{
		{"no shutdown", false},
		{"shutdown", true},
	}

	for _, tt := range tests {
		testConf := fmt.Sprintf(shortConfig, tt.in)
		if got := MlagParseShutdown(&m, testConf); got != tt.want {
			t.Fatalf("parseShutdown() = %v; want %v", got, tt.want)
		}
	}
	if got := MlagParseShutdown(&m, ""); got != false {
		t.Fatalf("parseShutdown() = %t; want false", got)
	}
}

func TestMlagSetDomainID_UnitTest(t *testing.T) {
	mlag := Mlag(dummyNode)
	cmds := []string{
		"mlag configuration",
		"default domain-id",
	}
	tests := []struct {
		in   string
		want string
	}{
		{"", "no domain-id"},
		{"test", "domain-id test"},
		{"test4", "domain-id test4"},
	}

	for _, tt := range tests {
		mlag.SetDomainID(tt.in)
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

func TestMlagSetDomainIDDefault_UnitTest(t *testing.T) {
	mlag := Mlag(dummyNode)

	cmds := []string{
		"mlag configuration",
		"default domain-id",
	}
	mlag.SetDomainIDDefault()
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestMlagSetLocalInterface_UnitTest(t *testing.T) {
	mlag := Mlag(dummyNode)
	cmds := []string{
		"mlag configuration",
		"default local-interface",
	}
	tests := []struct {
		in   string
		want string
	}{
		{"", "no local-interface"},
		{"Ethernet1", "local-interface Ethernet1"},
		{"Port-Channel1", "local-interface Port-Channel1"},
	}

	for _, tt := range tests {
		mlag.SetLocalInterface(tt.in)
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

func TestMlagSetLocalInterfaceDefault_UnitTest(t *testing.T) {
	mlag := Mlag(dummyNode)

	cmds := []string{
		"mlag configuration",
		"default local-interface",
	}
	mlag.SetLocalInterfaceDefault()
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestMlagSetPeerAddress_UnitTest(t *testing.T) {
	mlag := Mlag(dummyNode)
	cmds := []string{
		"mlag configuration",
		"default peer-address",
	}
	tests := []struct {
		in   string
		want string
	}{
		{"", "no peer-address"},
		{"1.1.1.1", "peer-address 1.1.1.1"},
		{"2.2.2.2", "peer-address 2.2.2.2"},
	}

	for _, tt := range tests {
		mlag.SetPeerAddress(tt.in)
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

func TestMlagSetPeerAddressDefault_UnitTest(t *testing.T) {
	mlag := Mlag(dummyNode)

	cmds := []string{
		"mlag configuration",
		"default peer-address",
	}
	mlag.SetPeerAddressDefault()
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestMlagSetPeerLink_UnitTest(t *testing.T) {
	mlag := Mlag(dummyNode)
	cmds := []string{
		"mlag configuration",
		"default peer-link",
	}
	tests := []struct {
		in   string
		want string
	}{
		{"", "no peer-link"},
		{"Ethernet1", "peer-link Ethernet1"},
		{"Port-Channel1", "peer-link Port-Channel1"},
	}

	for _, tt := range tests {
		mlag.SetPeerLink(tt.in)
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

func TestMlagSetPeerLinkDefault_UnitTest(t *testing.T) {
	mlag := Mlag(dummyNode)

	cmds := []string{
		"mlag configuration",
		"default peer-link",
	}
	mlag.SetPeerLinkDefault()
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestMlagSetShutdown_UnitTest(t *testing.T) {
	mlag := Mlag(dummyNode)
	cmds := []string{
		"mlag configuration",
		"default shutdown",
	}
	tests := []struct {
		in   bool
		want string
	}{
		{false, "no shutdown"},
		{true, "shutdown"},
	}

	for _, tt := range tests {
		mlag.SetShutdown(tt.in)
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

func TestMlagSetShutdownDefault_UnitTest(t *testing.T) {
	mlag := Mlag(dummyNode)

	cmds := []string{
		"mlag configuration",
		"default shutdown",
	}
	mlag.SetShutdownDefault()
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestMlagID_UnitTest(t *testing.T) {
	mlag := Mlag(dummyNode)
	cmds := []string{
		"mlag configuration",
		"default domain-id",
	}
	tests := []struct {
		in   string
		want string
	}{
		{"", "no domain-id"},
		{"test", "domain-id test"},
		{"test4", "domain-id test4"},
	}

	for _, tt := range tests {
		mlag.SetDomainID(tt.in)
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

func TestMlagSetMlagIDDefault_UnitTest(t *testing.T) {
	mlag := Mlag(dummyNode)

	cmds := []string{
		"interface Port-Channel1",
		"default mlag",
	}
	mlag.SetMlagIDDefault("Port-Channel1")
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
func TestMlagGet_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no interface Port-Channel1-2000",
			"default mlag configuration",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		gconfig := GlobalMlagConfig{
			"domain_id":       "",
			"local_interface": "",
			"peer_address":    "",
			"peer_link":       "",
			"shutdown":        "false",
		}
		iconfig := InterfaceMlagConfig{}
		configTemp := MlagConfig{config: gconfig, interfaces: iconfig}

		mlag := Mlag(dut)
		config := mlag.Get()
		if config.isEqual(configTemp) == false {
			t.Fatalf("Unequal configs seen")
		}

	}
}

func TestMlagGetter_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no interface Port-Channel1-2000",
			"default mlag configuration",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		gconfig := GlobalMlagConfig{
			"domain_id":       "",
			"local_interface": "",
			"peer_address":    "",
			"peer_link":       "",
			"shutdown":        "false",
		}
		iconfig := InterfaceMlagConfig{}
		configTemp := MlagConfig{config: gconfig, interfaces: iconfig}

		mlag := Mlag(dut)
		config := mlag.Get()
		if configTemp.DomainID() != config.DomainID() ||
			configTemp.LocalInterface() != config.LocalInterface() ||
			configTemp.PeerAddress() != config.PeerAddress() ||
			configTemp.PeerLink() != config.PeerLink() ||
			configTemp.Shutdown() != config.Shutdown() ||
			configTemp.InterfaceConfig("Port-Channel5") != config.InterfaceConfig("Port-Channel5") {
			t.Fatalf("Unequal configs. Got [%#v] Want [%#v]", config, configTemp)
		}

	}
}

func TestMlagSetDomainId_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"default mlag configuration",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("no domain-id", section); !found {
			t.Fatalf("\"no domain-id\" expected but not seen in "+
				"mlag configuration.\n[%s]", section)
		}

		if ok := mlag.SetDomainID("test"); !ok {
			t.Fatalf("Failure setting domain ID")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("domain-id test", section); !found {
			t.Fatalf("\"no domain-id\" expected but not seen under "+
				"mlag section.\n[%s]", section)
		}
	}
}
func TestMlagSetDomainIdNoValue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"mlag configuration",
			"domain-id test",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("domain-id test", section); !found {
			t.Fatalf("\"no domain-id\" expected but not seen in "+
				"mlag configuration.\n[%s]", section)
		}

		if ok := mlag.SetDomainID(""); !ok {
			t.Fatalf("Failure setting domain ID")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("no domain-id", section); !found {
			t.Fatalf("\"no domain-id\" expected but not seen under "+
				"mlag section.\n[%s]", section)
		}
	}
}
func TestMlagSetDomainIdDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"mlag configuration",
			"domain-id test",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("domain-id test", section); !found {
			t.Fatalf("\"domain-id test\" expected but not seen under mlag "+
				"configuration.\n[%s]", section)
		}

		if ok := mlag.SetDomainIDDefault(); !ok {
			t.Fatalf("Failure setting domain ID to default config")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("no domain-id", section); !found {
			t.Fatalf("\"no domain-id\" expected but not seen under "+
				"mlag section.\n[%s]", section)
		}
	}
}
func TestMlagSetLocalInt_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"interface Vlan1234",
			"default mlag configuration",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("no local-interface", section); !found {
			t.Fatalf("\"no local-interface\" exepected but not seen in Mlag "+
				"configuration.\n[%s]", section)
		}

		if ok := mlag.SetLocalInterface("Vlan1234"); !ok {
			t.Fatalf("Failure setting local-interface")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("local-interface Vlan1234", section); !found {
			t.Fatalf("\"no local-interface\" expected but not seen under "+
				"mlag section.\n[%s]", section)
		}
	}
}
func TestMlagSetLocalIntNoValue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"interface Vlan1234",
			"default mlag configuration",
			"mlag configuration",
			"local-interface Vlan1234",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("local-interface Vlan1234", section); !found {
			t.Fatalf("\"local-interface Vlan1234\" expected but not seen in Mlag "+
				"configuration.\n[%s]", section)
		}

		if ok := mlag.SetLocalInterface(""); !ok {
			t.Fatalf("Failure setting local-interface")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("no local-interface", section); !found {
			t.Fatalf("\"no local-interface\" expected but not seen under "+
				"mlag section.\n[%s]", section)
		}
	}
}
func TestMlagSetLocalIntDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"interface Vlan1234",
			"mlag configuration",
			"local-interface Vlan1234",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("local-interface Vlan1234", section); !found {
			t.Fatalf("\"local-interface Vlan1234\" expected but not seen in Mlag "+
				"configuration.\n[%s]", section)
		}

		if ok := mlag.SetLocalInterfaceDefault(); !ok {
			t.Fatalf("Failure setting default config mode for local-interface")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("no local-interface", section); !found {
			t.Fatalf("\"no local-interface\" expected but not seen under "+
				"mlag section.\n[%s]", section)
		}
	}
}
func TestMlagSetPeerAddress_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"default mlag configuration",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("no peer-address", section); !found {
			t.Fatalf("\"no peer-address\" expected but not seen under Mlag "+
				"configuration.\n[%s]", section)
		}

		if ok := mlag.SetPeerAddress("1.2.3.4"); !ok {
			t.Fatalf("Failure setting address for peer-address")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("peer-address 1.2.3.4", section); !found {
			t.Fatalf("\"peer-address 1.2.3.4\" expected but not seen under "+
				"mlag section.\n[%s]", section)
		}
	}
}
func TestMlagSetPeerAddressNoValue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"interface Vlan1234",
			"ip address 1.2.3.1/24",
			"mlag configuration",
			"peer-address 1.2.3.4",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("peer-address 1.2.3.4", section); !found {
			t.Fatalf("\"peer-address 1.2.3.4\" expected but not seen under "+
				"Mlag configuration.\n[%s]", section)
		}

		if ok := mlag.SetPeerAddress(""); !ok {
			t.Fatalf("Failure setting no config mode for peer-address")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("no peer-address", section); !found {
			t.Fatalf("\"no peer-address\" expected but not seen under "+
				"mlag section.\n[%s]", section)
		}
	}
}
func TestMlagSetPeerAddressDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"interface Vlan1234",
			"ip address 1.2.3.1/24",
			"mlag configuration",
			"peer-address 1.2.3.4",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("peer-address 1.2.3.4", section); !found {
			t.Fatalf("\"peer-address 1.2.3.4\" expected but not seen under "+
				"Mlag configuration.\n[%s]", section)
		}

		if ok := mlag.SetPeerAddressDefault(); !ok {
			t.Fatalf("Failure setting default mode for peer-address")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("no peer-address", section); !found {
			t.Fatalf("\"no peer-address\" expected but not seen under "+
				"mlag section.\n[%s]", section)
		}
	}
}
func TestMlagPeerLink_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"default mlag configuration",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("no peer-link", section); !found {
			t.Fatalf("\"no peer-link\" expected but not seen in Mlag "+
				"configuration.\n[%s]", section)
		}

		if ok := mlag.SetPeerLink("Ethernet1"); !ok {
			t.Fatalf("Failure setting Ethernet1 for peer-link")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("peer-link Ethernet1", section); !found {
			t.Fatalf("\"peer-link Ethernet1\" expected but not seen under "+
				"mlag section.\n[%s]", section)
		}
	}
}
func TestMlagPeerLinkPortChannel_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"default mlag configuration",
			"interface Port-Channel5",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("no peer-link", section); !found {
			t.Fatalf("\"no peer-link\" expected but not seen in Mlag "+
				"configuration.\n[%s]", section)
		}

		if ok := mlag.SetPeerLink("Port-Channel5"); !ok {
			t.Fatalf("Failure setting Port-Channel5 for peer-link")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("peer-link Port-Channel5", section); !found {
			t.Fatalf("\"peer-link Port-Channel5\" expected but not seen "+
				"under mlag section.\n[%s]", section)
		}
	}
}
func TestMlagPeerLinkNoValue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"mlag configuration",
			"peer-link Ethernet1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("peer-link Ethernet1", section); !found {
			t.Fatalf("\"peer-link Ethernet1\" expected but not seen in Mlag "+
				"configuration.\n[%s]", section)
		}

		if ok := mlag.SetPeerLink(""); !ok {
			t.Fatalf("Failure setting no value for peer-link")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("no peer-link", section); !found {
			t.Fatalf("\"no peer-link\" expected but not seen under "+
				"mlag section.\n[%s]", section)
		}
	}
}
func TestMlagPeerLinkDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"mlag configuration",
			"peer-link Ethernet1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("peer-link Ethernet1", section); !found {
			t.Fatalf("\"peer-link Ethernet1\" expected but not seen in Mlag "+
				"configuration.\n[%s]", section)
		}

		if ok := mlag.SetPeerLinkDefault(); !ok {
			t.Fatalf("Failure SetPeerLinkDefault for peer-link")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("no peer-link", section); !found {
			t.Fatalf("\"no peer-link\" expected but not seen under "+
				"mlag section.\n[%s]", section)
		}
	}
}
func TestMlagShutdownTrue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"default mlag configuration",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("no shutdown", section); !found {
			t.Fatalf("\"no shutdown\" expected but not seen under Mlag "+
				"configuration.\n[%s]", section)
		}

		if ok := mlag.SetShutdown(true); !ok {
			t.Fatalf("Failure setting default for mlag")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("shutdown", section); !found {
			t.Fatalf("\"shutdown\" expected but not seen under Mlag "+
				"configuration.\n[%s]", section)
		}
	}
}
func TestMlagShutdownFalse_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"mlag configuration",
			"shutdown",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("shutdown", section); !found {
			t.Fatalf("\"shutdown\" expected but not seen under Mlag "+
				"configuration.\n[%s]", section)
		}

		if ok := mlag.SetShutdown(false); !ok {
			t.Fatalf("Failure setting default for mlag")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("no shutdown", section); !found {
			t.Fatalf("\"no shutdown\" expected but not seen under Mlag "+
				"configuration.\n[%s]", section)
		}
	}
}

func TestMlagShutdownDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"mlag configuration",
			"shutdown",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section := mlag.GetSection()
		if found, _ := regexp.MatchString("shutdown", section); !found {
			t.Fatalf("\"shutdown\" expected but not seen under Mlag "+
				"configuration.\n[%s]", section)
		}

		if ok := mlag.SetShutdownDefault(); !ok {
			t.Fatalf("Failure SetShutdownDefault for mlag")
		}

		section = mlag.GetSection()
		if found, _ := regexp.MatchString("no shutdown", section); !found {
			t.Fatalf("\"no shutdown\" expected but not seen under Mlag "+
				"configuration.\n[%s]", section)
		}
	}
}
func TestMlagSetMlagId_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no interface Port-Channel10",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section, _ := mlag.GetBlock("interface Port-Channel10")
		if found, _ := regexp.MatchString("no mlag", section); found {
			t.Fatalf("\"no mlag\" expected but not seen under "+
				"Port-Channel10 section.\n[%s]", section)
		}

		if ok := mlag.SetMlagID("Port-Channel10", "100"); !ok {
			t.Fatalf("Failure setting default Mlag ID for Port-Channel10")
		}

		section, _ = mlag.GetBlock("interface Port-Channel10")
		if found, _ := regexp.MatchString("mlag 100", section); !found {
			t.Fatalf("\"mlag 100\" expected but not seen under "+
				"Port-Channel10 section.\n[%s]", section)
		}
	}
}
func TestMlagSetMlagIdNoValue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no interface Port-Channel10",
			"interface Port-Channel10",
			"mlag 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section, _ := mlag.GetBlock("interface Port-Channel10")
		if found, _ := regexp.MatchString("mlag 100", section); !found {
			t.Fatalf("Port-Channel10 does not have [mlag 100] config.\n[%s]", section)
		}

		if ok := mlag.SetMlagID("Port-Channel10", ""); !ok {
			t.Fatalf("Invalid setting Mlag ID for Port-Channel10")
		}

		section, _ = mlag.GetBlock("interface Port-Channel10")
		if found, _ := regexp.MatchString("no mlag", section); !found {
			t.Fatalf("\"no mlag\" expected but not seen under "+
				"Port-Channel10 section.\n[%s]", section)
		}
	}
}
func TestMlagSetMlagIdDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no interface Port-Channel10",
			"interface Port-Channel10",
			"mlag 100",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		mlag := Mlag(dut)

		section, _ := mlag.GetBlock("interface Port-Channel10")
		if found, _ := regexp.MatchString("mlag 100", section); !found {
			t.Fatalf("Mlag section does not have [mlag 100] config.\n[%s]", section)
		}

		if ok := mlag.SetMlagIDDefault("Port-Channel10"); !ok {
			t.Fatalf("Failure setting default Mlag ID for Port-Channel10")
		}

		section, _ = mlag.GetBlock("interface Port-Channel10")
		if found, _ := regexp.MatchString("no mlag", section); !found {
			t.Fatalf("\"no mlag\" expected but not seen under "+
				"Port-Channel10 section.\n[%s]", section)
		}
	}
}
