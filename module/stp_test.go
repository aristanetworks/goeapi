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

	"github.com/aristanetworks/goeapi"
)

/**
 *****************************************************************************
 * Unit Tests
 *****************************************************************************
 **/
func TestSTPInterfaces_UnitTest(t *testing.T) {
	stp := Stp(&goeapi.Node{})
	i := stp.Interfaces()
	if i == nil {
		t.Fatalf("No STPInterface Created")
	}
}

func TestSTPInstances_UnitTest(t *testing.T) {
	stp := Stp(&goeapi.Node{})
	i := stp.Instances()
	if i == nil {
		t.Fatalf("No STPInstances Created")
	}
}

func TestSTPSetMode_UnitTest(t *testing.T) {
	stp := Stp(dummyNode)

	cmds := []string{
		"default spanning-tree mode",
	}

	tests := [...]struct {
		mode string
		want string
		rc   bool
	}{
		{"mstp", "spanning-tree mode mstp", true},
		{"Invalid", "", false},
		{"none", "spanning-tree mode none", true},
		{"", "no spanning-tree mode", true},
	}

	for _, tt := range tests {
		if got := stp.SetMode(tt.mode); got != tt.rc {
			t.Fatalf("SetMode(%s) = %t; want %t", tt.mode, got, tt.rc)
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

func TestSTPIntfIsValidInterface_UnitTest(t *testing.T) {
	tests := [...]struct {
		in   string
		want bool
	}{
		{"Ethernet1", true},
		{"Ethernet2", true},
		{"Management1", false},
		{"Port-Channel10", true},
		{"Token-Ring8", false},
		{"Ethernet10/1", true},
		{"Vlan10", false},
	}

	for _, tt := range tests {
		if got := isValidStpInterface(tt.in); got != tt.want {
			t.Fatalf("isValidStpInterface(%s) = %v; want %v", tt.in, got, tt.want)
		}
	}
}

func TestSTPIntfParseBPDUGuard_UnitTest(t *testing.T) {
	var s STPInterfaceEntity

	shortConfig := `
 interface Ethernet1
    no description
    no shutdown
    default load-interval
    logging event link-status use-global
    no spanning-tree portfast
    spanning-tree portfast auto
    no spanning-tree link-type
    %s
    no spanning-tree bpdufilter
    no spanning-tree cost
    spanning-tree port-priority 128
    no spanning-tree guard
    no spanning-tree bpduguard rate-limit
    logging event spanning-tree use-global
    switchport tap native vlan 1
    no switchport tap identity
    switchport tap allowed vlan 1-4094
`
	tests := [...]struct {
		in   string
		want bool
	}{
		{"no spanning-tree bpduguard", false},
		{"spanning-tree bpduguard enable", true},
		{"spanning-tree bpduguard rate-limit enable", false},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := STPIntfParseBPDUGuard(&s, testConfig); got != tt.want {
			t.Fatalf("parseBPDUGuard(%s) = %v; want %v", tt.in, got, tt.want)
		}
	}
}

func TestSTPIntfParsePortFast_UnitTest(t *testing.T) {
	var s STPInterfaceEntity

	shortConfig := `
 interface Ethernet1
    no description
    no shutdown
    default load-interval
    logging event link-status use-global
    %s
    spanning-tree portfast auto
    no spanning-tree link-type
    no spanning-tree bpduguard
    no spanning-tree bpdufilter
    no spanning-tree cost
    spanning-tree port-priority 128
    no spanning-tree guard
    no spanning-tree bpduguard rate-limit
    logging event spanning-tree use-global
    switchport tap native vlan 1
    no switchport tap identity
    switchport tap allowed vlan 1-4094
`
	tests := [...]struct {
		in   string
		want bool
	}{
		{"no spanning-tree portfast", false},
		{"spanning-tree portfast", true},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		if got := STPIntfParsePortFast(&s, testConfig); got != tt.want {
			t.Fatalf("parseParsePortFast(%s) = %v; want %v", tt.in, got, tt.want)
		}
	}

}

func TestSTPIntfParsePortFastType_UnitTest(t *testing.T) {
	var s STPInterfaceEntity

	shortConfig := `
 interface Ethernet1
    no description
    no shutdown
    default load-interval
    logging event link-status use-global
    %s
    spanning-tree portfast auto
    no spanning-tree link-type
    no spanning-tree bpduguard
    no spanning-tree bpdufilter
    no spanning-tree cost
    spanning-tree port-priority 128
    no spanning-tree guard
    no spanning-tree bpduguard rate-limit
    logging event spanning-tree use-global
    switchport tap native vlan 1
    no switchport tap identity
    switchport tap allowed vlan 1-4094
`
	tests := [...]struct {
		in   string
		want string
	}{
		{"no spanning-tree portfast", "normal"},
		{"spanning-tree portfast", "edge"},
		{"spanning-tree portfast network", "network"},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig, tt.in)
		got := STPIntfParsePortFastType(&s, testConfig)
		if got != tt.want {
			t.Fatalf("parsePortFastType(%s) = %s; want %s", tt.in, got, tt.want)
		}
	}

}

func TestSTPIntfGetKeysReturned_UnitTest(t *testing.T) {
	stp := STPInterfaces(dummyNode)
	config := stp.Get("Ethernet1")
	for _, val := range []string{"bpduguard", "portfast", "portfast_type"} {
		if _, found := config[val]; !found {
			t.Fatalf("Get() missing key %s", val)
		}
	}
}

func TestSTPIntfSetPortfastType_UnitTest(t *testing.T) {
	stp := STPInterfaces(dummyNode)
	tests := []struct {
		value string
		want  string
		rc    bool
	}{
		{"", "", false},
		{"network", "spanning-tree portfast ", true},
		{"edge", "spanning-tree portfast ", true},
		{"normal", "spanning-tree portfast ", true},
		{"InvalidType", "", false},
	}

	for _, intf := range interfaceList {

		for _, tt := range tests {
			cmds := []string{
				"interface " + intf,
				"spanning-tree portfast " + tt.value,
			}
			if tt.value == "edge" {
				cmds = append(cmds, "spanning-tree portfast auto")
			}
			if ok := stp.SetPortfastType(intf, tt.value); ok != tt.rc {
				t.Fatalf("Expected status \"%t\" got \"%t\"", tt.rc, ok)
			}
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
}

func TestSTPIntfSetPortfast_UnitTest(t *testing.T) {
	stp := STPInterfaces(dummyNode)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"default spanning-tree portfast",
		}

		tests := []struct {
			enable bool
			want   string
		}{
			{true, "spanning-tree portfast"},
			{false, "no spanning-tree portfast"},
		}

		for _, tt := range tests {
			stp.SetPortfast(intf, tt.enable)
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

func TestSTPIntfSetPortfastInvalidIntf_UnitTest(t *testing.T) {
	stp := STPInterfaces(dummyNode)
	intf := RandomString(8, 15)

	if ok := stp.SetPortfast(intf, true); ok {
		t.Fatalf("Invalid interface should return false")
	}
}

func TestSTPIntfSetPortfastDefault_UnitTest(t *testing.T) {
	stp := STPInterfaces(dummyNode)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"default spanning-tree portfast",
		}
		stp.SetPortfastDefault(intf)
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestSTPIntfSetBPDUGuard_UnitTest(t *testing.T) {
	stp := STPInterfaces(dummyNode)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"default spanning-tree bpduguard",
		}

		tests := []struct {
			enable bool
			want   string
		}{
			{true, "spanning-tree bpduguard enable"},
			{false, "spanning-tree bpduguard disable"},
		}

		for _, tt := range tests {
			stp.SetBPDUGuard(intf, tt.enable)
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

func TestSTPIntfSetBPDUGuardDefault_UnitTest(t *testing.T) {
	stp := STPInterfaces(dummyNode)

	for _, intf := range interfaceList {

		cmds := []string{
			"interface " + intf,
			"default spanning-tree bpduguard",
		}
		stp.SetBPDUGuardDefault(intf)
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
func getValidInterfaces(s *STPInterfaceEntity) []string {
	var re = regexp.MustCompile(`(?m)^interface\s(Eth.+|Po.+)`)
	config := s.Config()

	interfaces := re.FindAllStringSubmatch(config, -1)

	response := make([]string, len(interfaces))

	for idx, name := range interfaces {
		response[idx] = name[1]
	}
	return response
}

func TestSTPEntitySetMode_SystemTest(t *testing.T) {
	for _, dut := range duts {
		stp := Stp(dut)

		types := []string{"mstp", "none", ""}

		for _, tt := range types {
			if ok := stp.SetMode(tt); !ok {
				t.Fatalf("SetMode(%s) failed", tt)
			}
		}
	}
}

func TestSTPEntitySetModeInvalid_SystemTest(t *testing.T) {
	for _, dut := range duts {
		stp := Stp(dut)

		if ok := stp.SetMode("InvalidMode"); ok {
			t.Fatalf("SetMode(InvalidMode) should fail")
		}
	}
}

func TestSTPIntfGet_SystemTest(t *testing.T) {
	for _, dut := range duts {
		stp := STPInterfaces(dut)

		cmds := []string{
			"default interface Ethernet1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		result := stp.Get("Ethernet1")
		if result == nil {
			t.Fatalf("Get(Ethernet1) failed")
		}
	}
}

func TestSTPIntfGetAll_SystemTest(t *testing.T) {
	for _, dut := range duts {
		stp := STPInterfaces(dut)

		cmds := []string{
			"default interface Et1-2",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		collection := stp.GetAll()
		if collection == nil {
			t.Fatalf("GetAll() failed")
		}
	}
}

func TestSTPIntfSetBPDUGuardTrue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		stp := STPInterfaces(dut)
		intf := RandomChoice(getValidInterfaces(stp))

		cmds := []string{
			"default interface " + intf,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := stp.SetBPDUGuard(intf, true); !ok {
			t.Fatalf("SetBPDUGuard(%s, true) failed", intf)
		}

		parent := `interface\s+` + intf
		config, _ := stp.GetBlock(parent)
		if found := STPIntfParseBPDUGuard(stp, config); !found {
			t.Fatalf("\"spanning-tree bpduguard enable\" expected but not seen under "+
				"%s section.\n[%s]", intf, config)
		}
	}
}

func TestSTPIntfSetBPDUGuardFalse_SystemTest(t *testing.T) {
	for _, dut := range duts {
		stp := STPInterfaces(dut)
		intf := RandomChoice(getValidInterfaces(stp))

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"spanning-tree bpduguard enable",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		parent := `interface\s+` + intf
		config, _ := stp.GetBlock(parent)
		if found := STPIntfParseBPDUGuard(stp, config); !found {
			t.Fatalf("\"spanning-tree bpduguard enable\" expected but not seen under "+
				"%s section.\n[%s]", intf, config)
		}

		if ok := stp.SetBPDUGuard(intf, false); !ok {
			t.Fatalf("SetBPDUGuard(%s, false) failed", intf)
		}

		config, _ = stp.GetBlock(parent)
		if found := STPIntfParseBPDUGuard(stp, config); found {
			t.Fatalf("\"spanning-tree bpduguard enable\" NOT expected but not seen under "+
				"%s section.\n[%s]", intf, config)
		}
	}
}

func TestSTPIntfSetBPDUGuardDefualt_SystemTest(t *testing.T) {
	for _, dut := range duts {
		stp := STPInterfaces(dut)
		intf := RandomChoice(getValidInterfaces(stp))

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"spanning-tree bpduguard enable",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := stp.SetBPDUGuardDefault(intf); !ok {
			t.Fatalf("SetBPDUGuardDefault(%s) failed", intf)
		}

		parent := `interface\s+` + intf
		config, _ := stp.GetBlock(parent)
		if found := STPIntfParseBPDUGuard(stp, config); found {
			t.Fatalf("\"spanning-tree bpduguard enable\" NOT expected but not seen under "+
				"%s section.\n[%s]", intf, config)
		}
	}
}

func TestSTPIntfSetPortFastTrue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		stp := STPInterfaces(dut)
		intf := RandomChoice(getValidInterfaces(stp))

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := stp.SetPortfast(intf, true); !ok {
			t.Fatalf("SetPortfast(%s, true) failed", intf)
		}

		parent := `interface\s+` + intf
		config, _ := stp.GetBlock(parent)
		if found := STPIntfParsePortFast(stp, config); !found {
			t.Fatalf("\"spanning-tree portfast\" expected but not seen under "+
				"%s section.\n[%s]", intf, config)
		}
	}
}

func TestSTPIntfSetPortFastFalse_SystemTest(t *testing.T) {
	for _, dut := range duts {
		stp := STPInterfaces(dut)
		intf := RandomChoice(getValidInterfaces(stp))

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"spanning-tree portfast",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := stp.SetPortfast(intf, false); !ok {
			t.Fatalf("SetPortfast(%s, false) failed", intf)
		}

		parent := `interface\s+` + intf
		config, _ := stp.GetBlock(parent)
		if found := STPIntfParsePortFast(stp, config); found {
			t.Fatalf("\"spanning-tree portfast\" NOT expected but not seen under "+
				"%s section.\n[%s]", intf, config)
		}
	}
}

func TestSTPIntfSetPortFastDefualt_SystemTest(t *testing.T) {
	for _, dut := range duts {
		stp := STPInterfaces(dut)
		intf := RandomChoice(getValidInterfaces(stp))

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := stp.SetPortfastDefault(intf); !ok {
			t.Fatalf("SetPortfastDefault(%s) failed", intf)
		}

		parent := `interface\s+` + intf
		config, _ := stp.GetBlock(parent)
		if found := STPIntfParsePortFast(stp, config); found {
			t.Fatalf("\"spanning-tree portfast\" NOT expected but not seen under "+
				"%s section.\n[%s]", intf, config)
		}
	}
}

func TestSTPIntfInvalidInteface_SystemTest(t *testing.T) {
	for _, dut := range duts {
		stp := STPInterfaces(dut)

		if ok := stp.SetPortfastDefault("Token-Ring5"); ok {
			t.Fatalf("SetPortfastDefault(Token-Ring5) should fail with invalid interface")
		}
	}
}

func TestSTPIntfSetPortFastTypes_SystemTest(t *testing.T) {
	for _, dut := range duts {
		stp := STPInterfaces(dut)
		intf := RandomChoice(getValidInterfaces(stp))

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		types := []string{"edge", "network", "normal"}

		for _, tt := range types {

			if ok := stp.SetPortfastType(intf, tt); !ok {
				t.Fatalf("SetPortfastType(%s, %s) failed", intf, tt)
			}

			parent := `interface\s+` + intf
			config, _ := stp.GetBlock(parent)
			if got := STPIntfParsePortFastType(stp, config); got != tt {
				t.Fatalf("Expected %s but got %s section.\n[%s]", tt, got, config)
			}
		}
	}
}

func TestSTPIntfSetPortFastTypesInvalid_SystemTest(t *testing.T) {
	for _, dut := range duts {
		stp := STPInterfaces(dut)
		intf := RandomChoice(getValidInterfaces(stp))

		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ok := stp.SetPortfastType(intf, "InvalidType"); ok {
			t.Fatalf("SetPortfastType(%s, InvalidType) passed but should fail", intf)
		}
	}
}
