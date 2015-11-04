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
	"regexp"
	"testing"
)

func TestAclMaskToPrefixlen_UnitTest(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", "32"},
		{"255.255.255.0", "24"},
		{"255.255.0.0", "16"},
		{"255.255.252.0", "22"},
		{"255.0.0.0", "8"},
		{"128.0.0.0", "1"},
		{"255.255.253", "0"},
		{"192.168.0.0", "0"},
	}

	for _, tt := range tests {
		if got := maskToPrefixlen(tt.in); got != tt.want {
			t.Fatalf("maskToPrefixlen(%q) = %q; want %q", tt.in, got, tt.want)
		}
	}

	tests = []struct {
		in, want string
	}{
		{"32", "255.255.255.255"},
		{"24", "255.255.255.0"},
		{"16", "255.255.0.0"},
		{"1", "128.0.0.0"},
		{"8", "255.0.0.0"},
		{"22", "255.255.252.0"},
		{"2", "192.0.0.0"},
		{"0", "0.0.0.0"},
		{"", "255.255.255.255"},
	}

	for _, tt := range tests {
		if got := prefixlenToMask(tt.in); got != tt.want {
			t.Fatalf("prefixlenToMask(%q) = %q; want %q", tt.in, got, tt.want)
		}
	}
}

func TestAclParseEntries_UnitTest(t *testing.T) {
	var a AclEntity

	var config = `
    ip access-list standard test
        no statistics per-entry
        fragment-rules
        10 permit host 1.2.3.4 log
        20 permit 1.2.3.4 255.255.0.0 log
        30 deny any
        40 permit 5.6.7.0/24
        50 permit 16.0.0.0/8
        60 permit any log`

	results := AclEntryMap{
		"10": {
			"action":  "permit",
			"srcaddr": "1.2.3.4",
			"srclen":  "32",
			"log":     "log",
		},
		"20": {
			"action":  "permit",
			"srcaddr": "1.2.3.4",
			"srclen":  "16",
			"log":     "log",
		},
		"30": {
			"action":  "deny",
			"srcaddr": "0.0.0.0",
			"srclen":  "32",
			"log":     "",
		},
		"40": {
			"action":  "permit",
			"srcaddr": "5.6.7.0",
			"srclen":  "24",
			"log":     "",
		},
		"50": {
			"action":  "permit",
			"srcaddr": "16.0.0.0",
			"srclen":  "8",
			"log":     "",
		},
		"60": {
			"action":  "permit",
			"srcaddr": "0.0.0.0",
			"srclen":  "32",
			"log":     "log",
		},
	}

	aclEntries := AclParseEntries(&a, config)

	for seqnum, aclEntry := range results {
		if _, found := aclEntries[seqnum]; !found {
			t.Fatalf("parseEntries(1): entry %s not found in AclEntryMap", seqnum)
		}
		for key, val := range aclEntry {
			if val != aclEntries[seqnum][key] {
				t.Fatalf("parseEntries(2): entry[%s][%s] = %s got %s", seqnum, key,
					val, aclEntries[seqnum][key])
			}
		}
	}
}

func TestAclGet_SystemTest(t *testing.T) {
	cmds := []string{
		"no ip access-list standard test",
		"ip access-list standard test",
	}
	for _, dut := range duts {
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		ret, err := Acl(dut).Get("test")
		if err != nil || ret == nil {
			t.Fatalf("Expecting non-nil value from acl.Get(). Error: %s", err)
		}
	}
}

func TestAclGetters_SystemTest(t *testing.T) {
	cmds := []string{
		"no ip access-list standard test",
		"ip access-list standard test",
		"10 permit host 1.2.3.4 log",
		"20 permit 1.2.3.4 255.255.0.0 log",
		"30 deny any",
		"40 permit 5.6.7.0/24",
		"50 permit 16.0.0.0/8",
		"60 permit any log",
	}
	for _, dut := range duts {
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		ret, err := Acl(dut).Get("test")
		if err != nil {
			t.Fatalf("Expecting non-nil value from acl.Get(). Error: %s", err)
		}
		ret.Name()
		ret.Type()
		entries := ret.Entries()
		for k := range entries {
			entries[k].Action()
			entries[k].SrcAddr()
			entries[k].SrcLen()
			entries[k].Log()
		}
	}
}

func TestAclGetNone_SystemTest(t *testing.T) {
	cmds := []string{
		"no ip access-list standard test",
	}
	for _, dut := range duts {
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		ret, _ := Acl(dut).Get("test")
		if ret != nil {
			t.Fatalf("Expecting nil value from acl.Get().")
		}
	}
}

func TestAclGetAll_SystemTest(t *testing.T) {
	cmds := []string{
		"no ip access-list standard test",
		"ip access-list standard test",
	}
	for _, dut := range duts {
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		ret := Acl(dut).GetAll()
		if _, found := ret["test"]; !found {
			t.Fatalf("Expecting value from acl.GetAll().")
		}
	}
}

func TestAclGetSectionInvalid_SystemTest(t *testing.T) {
	cmds := []string{
		"no ip access-list standard test",
	}

	for _, dut := range duts {
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		acl := Acl(dut)

		if section := acl.GetSection("test"); section != "" {
			t.Fatalf("Invalid acl should return \"\"")
		}
	}
}

func TestAclCreate_SystemTest(t *testing.T) {
	cmds := []string{
		"no ip access-list standard test",
	}
	for _, dut := range duts {
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		ret, _ := Acl(dut).Get("test")
		if ret != nil {
			t.Fatalf("Expecting nil value from acl.Get().")
		}

		if ok := Acl(dut).Create("test"); !ok {
			t.Fatalf("Create() failure")
		}

		if ret, err := Acl(dut).Get("test"); ret == nil {
			t.Fatalf("Expecting non-nil value from acl.Get(). Err: %q", err)
		}
	}
}

func TestAclDelete_SystemTest(t *testing.T) {
	cmds := []string{
		"ip access-list standard test",
	}
	for _, dut := range duts {
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		ret, _ := Acl(dut).Get("test")
		if ret == nil {
			t.Fatalf("Expecting non-nil value from acl.Get().")
		}

		if ok := Acl(dut).Delete("test"); !ok {
			t.Fatalf("Delete() failure")
		}

		if ret, _ = Acl(dut).Get("test"); ret != nil {
			t.Fatalf("Expecting nil value from acl.Get().")
		}
	}
}

func TestAclDefault_SystemTest(t *testing.T) {
	cmds := []string{
		"ip access-list standard test",
	}
	for _, dut := range duts {
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		ret, _ := Acl(dut).Get("test")
		if ret == nil {
			t.Fatalf("Expecting non-nil value from acl.Get().")
		}

		if ok := Acl(dut).Default("test"); !ok {
			t.Fatalf("Default() failure")
		}

		if ret, _ = Acl(dut).Get("test"); ret != nil {
			t.Fatalf("Expecting nil value from acl.Get().")
		}
	}
}

func TestAclUpdateEntry_SystemTest(t *testing.T) {
	cmds := []string{
		"no ip access-list standard test",
		"ip access-list standard test",
	}

	for _, dut := range duts {
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		acl := Acl(dut)

		ret, _ := acl.Get("test")
		if ret == nil {
			t.Fatalf("Expecting non-nil value from acl.Get().")
		}

		section := acl.GetSection("test")
		if found, _ := regexp.MatchString("10 permit any log", section); found {
			t.Fatalf("Acl section has existing config [10 permit any log].")
		}
		if ok := acl.UpdateEntry("test", "10", "permit", "0.0.0.0", "0", true); !ok {
			t.Fatalf("acl.UpdateEntry failure.")
		}
		section = acl.GetSection("test")
		if found, _ := regexp.MatchString("10 permit any log", section); !found {
			t.Fatalf("Config not seen in section after UpdateEntry().")
		}
	}
}

func TestAclUpdateEntryExisting_SystemTest(t *testing.T) {
	cmds := []string{
		"no ip access-list standard test",
		"ip access-list standard test",
		"10 permit any log",
	}

	for _, dut := range duts {
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		acl := Acl(dut)

		ret, _ := acl.Get("test")
		if ret == nil {
			t.Fatalf("Expecting non-nil value from acl.Get().")
		}

		section := acl.GetSection("test")
		if found, _ := regexp.MatchString("10 permit any log", section); !found {
			t.Fatalf("Acl section is missing [10 permit any log].")
		}
		if ok := acl.UpdateEntry("test", "10", "deny", "0.0.0.0", "0", true); !ok {
			t.Fatalf("acl.UpdateEntry failure.")
		}
		section = acl.GetSection("test")
		if found, _ := regexp.MatchString("10 deny any log", section); !found {
			t.Fatalf("Config not seen in section after UpdateEntry().")
		}
	}
}

func TestAclAddEntry_SystemTest(t *testing.T) {
	cmds := []string{
		"no ip access-list standard test",
		"ip access-list standard test",
	}

	for _, dut := range duts {
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		acl := Acl(dut)

		ret, _ := acl.Get("test")
		if ret == nil {
			t.Fatalf("Expecting non-nil value from acl.Get().")
		}

		section := acl.GetSection("test")
		if found, _ := regexp.MatchString("10 permit any log", section); found {
			t.Fatalf("Acl section has existing config [10 permit any log].")
		}
		if ok := acl.AddEntry("test", "permit", "0.0.0.0", "0", true); !ok {
			t.Fatalf("acl.AddEntry failure.")
		}
		section = acl.GetSection("test")
		if found, _ := regexp.MatchString("10 permit any log", section); !found {
			t.Fatalf("Config not seen in section after AddEntry().")
		}
	}
}

func TestAclRemoveEntry_SystemTest(t *testing.T) {
	cmds := []string{
		"no ip access-list standard test",
		"ip access-list standard test",
		"10 permit any log",
	}

	for _, dut := range duts {
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}
		acl := Acl(dut)

		ret, _ := acl.Get("test")
		if ret == nil {
			t.Fatalf("Expecting non-nil value from acl.Get().")
		}

		section := acl.GetSection("test")
		if found, _ := regexp.MatchString("10 permit any log", section); !found {
			t.Fatalf("Acl section is missing [10 permit any log].")
		}
		if ok := acl.RemoveEntry("test", "10"); !ok {
			t.Fatalf("acl.RemoveEntry failure.")
		}
		section = acl.GetSection("test")
		if found, _ := regexp.MatchString("10 permit any log", section); found {
			t.Fatalf("Config seen in section after RemoveEntry().")
		}
	}
}
