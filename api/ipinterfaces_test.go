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
	"testing"
)

func TestIpInterfaceFunctions_UnitTest(t *testing.T) {
	tests := []struct {
		in   int
		want bool
	}{
		{0, false},
		{67, false},
		{68, true},
		{256, true},
		{1024, true},
		{4096, true},
		{65535, true},
		{65536, false},
	}

	for _, tt := range tests {
		if got := isValidMtu(tt.in); got != tt.want {
			t.Fatalf("isValidMtu(%q) = %t; want %t", tt.in, got, tt.want)
		}
	}
}
func TestIpInterfaceParseAddress_UnitTest(t *testing.T) {
	var i IPInterfaceEntity

	var shortConf = `
     interface Management1
        no description
        no shutdown
        mtu 1500
        no ip local-proxy-arp
        %s
        ip local-proxy-arp
        no ip verify unicast`

	tests := []struct {
		in   string
		want string
	}{
		{"ip address 1.1.1.1/24", "1.1.1.1/24"},
		{"ip address 1.1.1.1/32", "1.1.1.1/32"},
		{"ip address 10.10.10.10/24", "10.10.10.10/24"},
		{"ip address 100.100.100.100/24", "100.100.100.100/24"},
		{"ip address 192.168.1.16/24", "192.168.1.16/24"},
		{"", ""},
	}

	for _, tt := range tests {
		testConf := fmt.Sprintf(shortConf, tt.in)
		if got := IPIntfParseAddress(&i, testConf); got != tt.want {
			t.Fatalf("parseAddress(%q) = %q; want %q", tt.in, got, tt.want)
		}
	}
}

func TestIpInterfaceParseMtu_UnitTest(t *testing.T) {
	var i IPInterfaceEntity

	var shortConf = `
     interface Management1
        no description
        no shutdown
        %s
        no ip local-proxy-arp
        ip address 1.1.1.1/24`

	tests := []struct {
		in   string
		want string
	}{
		{"mtu 1500", "1500"},
		{"mtu 2096", "2096"},
		{"mtu 4096", "4096"},
		{"", ""},
	}

	for _, tt := range tests {
		testConf := fmt.Sprintf(shortConf, tt.in)
		if got := IPIntfParseMtu(&i, testConf); got != tt.want {
			t.Fatalf("parseMtu(%q) = %q; want %q", tt.in, got, tt.want)
		}
	}
}

func (src IPInterfaceConfig) isEqual(dst IPInterfaceConfig) bool {
	for k, v := range src {
		if dst[k] != v {
			return false
		}
	}
	return true
}

func TestIpInterfaceGet_SystemTest(t *testing.T) {
	for _, dut := range duts {
		ipIntf := IPInterface(dut)

		intf := RandomChoice(ipIntf.GetEthInterfaces())
		if intf == "" {
			t.Fatalf("Null Interface")
			continue
		}
		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
			"ip address 99.98.99.99/24",
			"mtu 1800",
		}

		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure on intf %s", intf)
		}

		ret, err := ipIntf.Get(intf)
		if err != nil || ret == nil {
			t.Fatalf("Expecting non-nil value from IPInterface.Get(). Error: %s Intf: %s", err, intf)
		}

		cmds = []string{"default interface " + intf}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure on intf %s", intf)
		}
	}
}

func TestIpInterfaceGetWOIpAddress_SystemTest(t *testing.T) {
	for _, dut := range duts {
		ipIntf := IPInterface(dut)

		intf := RandomChoice(ipIntf.GetEthInterfaces())
		if intf == "" {
			t.Fatalf("Null Interface")
			continue
		}
		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
		}

		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure on intf %s", intf)
		}

		ret, _ := ipIntf.Get(intf)
		if ret["address"] != "" {
			t.Fatalf("Expecting null string value from IPInterface.Get(). interface %s", intf)
		}

		cmds = []string{"default interface " + intf}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure on intf %s", intf)
		}
	}
}

func TestIpInterfaceGetInterfaceGetAll_SystemTest(t *testing.T) {
	for _, dut := range duts {
		ret := IPInterface(dut).GetAll()
		if ret == nil {
			t.Fatalf("Expecting non-nil value from IPInterface.GetAll().")
		}
		if _, found := ret["Management1"]; !found {
			t.Fatalf("Management interface not found")
		}
	}
}

func TestIpInterfaceCreate_SystemTest(t *testing.T) {
	for _, dut := range duts {
		ipIntf := IPInterface(dut)

		intf := RandomChoice(ipIntf.GetEthInterfaces())
		if intf == "" {
			t.Fatalf("Null Interface")
			continue
		}
		cmds := []string{
			"default interface " + intf,
		}

		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ret, _ := ipIntf.Get(intf); ret != nil {
			t.Fatalf("Expecting nil from IPInterface.Get().")
		}

		if ok := ipIntf.Create(intf); !ok {
			t.Fatalf("Failure. IPInterface.Create().")
		}

		if ret, _ := ipIntf.Get(intf); ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}
	}
}

func TestIpInterfaceDelete_SystemTest(t *testing.T) {
	for _, dut := range duts {
		ipIntf := IPInterface(dut)

		intf := RandomChoice(ipIntf.GetEthInterfaces())
		if intf == "" {
			t.Fatalf("Null Interface")
			continue
		}
		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
			"ip address 99.98.99.99/24",
			"mtu 1800",
		}

		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ret, _ := ipIntf.Get(intf); ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}

		if ok := ipIntf.Delete(intf); !ok {
			t.Fatalf("Failure. IPInterface.Create().")
		}

		if ret, _ := ipIntf.Get(intf); ret != nil {
			t.Fatalf("Expecting nil from IPInterface.Get().")
		}
	}
}

func TestIpInterfaceSetAddress_SystemTest(t *testing.T) {
	for _, dut := range duts {
		ipIntf := IPInterface(dut)

		intf := RandomChoice(ipIntf.GetEthInterfaces())
		if intf == "" {
			t.Fatalf("Null Interface")
			continue
		}
		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
			"ip address 99.98.99.99/24",
			"mtu 1800",
		}

		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ret, _ := ipIntf.Get(intf); ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}

		if ok := ipIntf.SetAddress(intf, "1.1.1.1/24"); !ok {
			t.Fatalf("Failure. IPInterface.SetAddress(%s,1.1.1.1/24)", intf)
		}

		ret, _ := ipIntf.Get(intf)
		if ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}

		if ret["address"] != "1.1.1.1/24" {
			t.Fatalf("Got %s, Expecting 1.1.1.1/24", ret["address"])
		}

		cmds = []string{"default interface " + intf}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure on intf %s", intf)
		}
	}
}

func TestIpInterfaceSetAddressNo_SystemTest(t *testing.T) {
	for _, dut := range duts {
		ipIntf := IPInterface(dut)

		intf := RandomChoice(ipIntf.GetEthInterfaces())
		if intf == "" {
			t.Fatalf("Null Interface")
			continue
		}
		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
			"ip address 99.98.99.99/24",
			"mtu 1800",
		}

		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ret, _ := ipIntf.Get(intf); ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}

		if ok := ipIntf.SetAddress(intf, ""); !ok {
			t.Fatalf("Failure. IPInterface.SetAddress(%s,\"\") should fail", intf)
		}

		ret, _ := ipIntf.Get(intf)
		if ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}

		if ret["address"] != "" {
			t.Fatalf("Got %s, Expecting \"\" intf %s", ret["address"], intf)
		}

		cmds = []string{"default interface " + intf}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure on intf %s", intf)
		}
	}
}

func TestIpInterfaceSetAddressDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		ipIntf := IPInterface(dut)

		intf := RandomChoice(ipIntf.GetEthInterfaces())
		if intf == "" {
			t.Fatalf("Null Interface")
			continue
		}
		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
			"ip address 99.98.99.99/24",
			"mtu 1800",
		}

		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ret, _ := ipIntf.Get(intf); ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}

		if ok := ipIntf.SetAddressDefault(intf); !ok {
			t.Fatalf("Failure. IPInterface.SetAddressDefault(%s)", intf)
		}

		ret, _ := ipIntf.Get(intf)
		if ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}

		if ret["address"] != "" {
			t.Fatalf("Got %s, Expecting \"\"", ret["address"])
		}

		cmds = []string{"default interface " + intf}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure on intf %s", intf)
		}
	}
}

func TestIpInterfaceSetMtu_SystemTest(t *testing.T) {
	for _, dut := range duts {
		ipIntf := IPInterface(dut)

		intf := RandomChoice(ipIntf.GetEthInterfaces())
		if intf == "" {
			t.Fatalf("Null Interface")
			continue
		}
		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
			"mtu 1800",
			"ip address 99.98.99.99/24",
		}

		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ret, _ := ipIntf.Get(intf); ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}

		if ok := ipIntf.SetMtu(intf, 2400); !ok {
			t.Fatalf("Failure. IPInterface.SetMtu(%s,2400).", intf)
		}

		ret, _ := ipIntf.Get(intf)
		if ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}

		if ret["mtu"] != "2400" {
			t.Fatalf("Got %s, Expecting 2400", ret["mtu"])
		}

		cmds = []string{"default interface " + intf}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure on intf %s", intf)
		}
	}
}

func TestIpInterfaceSetMtuInvalid_SystemTest(t *testing.T) {
	for _, dut := range duts {
		ipIntf := IPInterface(dut)

		intf := RandomChoice(ipIntf.GetEthInterfaces())
		if intf == "" {
			t.Fatalf("Null Interface")
			continue
		}
		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
			"mtu 1800",
			"ip address 99.98.99.99/24",
		}

		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ret, _ := ipIntf.Get(intf); ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}

		if ok := ipIntf.SetMtu(intf, 1); ok {
			t.Fatalf("Failure. Should disallow invalid MTU size. %s", intf)
		}

		cmds = []string{"default interface " + intf}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure on intf %s", intf)
		}
	}
}

func TestIpInterfaceSetMtuDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		ipIntf := IPInterface(dut)

		intf := RandomChoice(ipIntf.GetEthInterfaces())
		if intf == "" {
			t.Fatalf("Null Interface")
			continue
		}
		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
			"mtu 1800",
			"ip address 99.98.99.99/24",
		}

		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		if ret, _ := ipIntf.Get(intf); ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}

		if ok := ipIntf.SetMtuDefault(intf); !ok {
			t.Fatalf("Failure. IPInterface.SetMtuDefault(%s).", intf)
		}

		ret, _ := ipIntf.Get(intf)
		if ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}

		if ret["mtu"] == "1800" {
			t.Fatalf("Got %s, Expecting default value", ret["mtu"])
		}

		cmds = []string{"default interface " + intf}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure on intf %s", intf)
		}
	}
}

func TestIpInterfaceGetters_SystemTest(t *testing.T) {
	for _, dut := range duts {
		ipIntf := IPInterface(dut)

		intf := RandomChoice(ipIntf.GetEthInterfaces())
		if intf == "" {
			t.Fatalf("Null Interface")
			continue
		}
		cmds := []string{
			"default interface " + intf,
			"interface " + intf,
			"no switchport",
			"ip address 99.98.99.99/24",
			"mtu 1800",
		}

		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		ret, _ := ipIntf.Get(intf)
		if ret == nil {
			t.Fatalf("Expecting non-nil from IPInterface.Get().")
		}

		if ret.Name() != intf ||
			ret.Address() != "99.98.99.99/24" ||
			ret.Mtu() != "1800" {
			t.Fatalf("Values for %s incorrect. [%#v]", intf, ret)
		}

		cmds = []string{"default interface " + intf}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure on intf %s", intf)
		}
	}
}
