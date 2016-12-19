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
	"sort"
	"strconv"
	"strings"
	"testing"
)

/**
 *****************************************************************************
 * Unit Tests
 *****************************************************************************
 **/
func compareSlice(slice1 []string, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	sort.Strings(slice1)
	sort.Strings(slice2)

	for idx := range slice1 {
		if slice1[idx] != slice2[idx] {
			return false
		}
	}
	return true
}

func TestVlanFunctions_UnitTest(t *testing.T) {
	isVlanTests := []struct {
		in   string
		want bool
	}{
		{"0", false},
		{"1", true},
		{"64", true},
		{"256", true},
		{"1024", true},
		{"4093", true},
		{"4094", true},
		{"4095", false},
		{"4096", false},
		{"65535", false},
	}

	for _, tt := range isVlanTests {
		if got := isVlan(tt.in); got != tt.want {
			t.Fatalf("isVlan(%q) = %t; want %t", tt.in, got, tt.want)
		}
	}

	findDiffTests := []struct {
		in1  []string
		in2  []string
		want []string
	}{
		{[]string{}, []string{"a", "b", "c"}, []string{}},
		{[]string{"a", "b", "c"}, []string{}, []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, []string{""}, []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, nil, []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}, []string{}},
		{[]string{"a", "b", "c"}, []string{"d", "e", "f"}, []string{"a", "b", "c"}},
		{[]string{"a", "c", "e"}, []string{"c", "d", "f"}, []string{"a", "e"}},
		{[]string{"a", "", "c", "", "e"}, []string{"", "c", "d", "", "f"}, []string{"a", "e"}},
	}

	for _, tt := range findDiffTests {
		got := findDiff(tt.in1, tt.in2)
		if compareSlice(got, tt.want) == false {
			t.Fatalf("findDiff(%q, %q) = %q; want %q", tt.in1, tt.in2, got, tt.want)
		}
	}
}

func TestVlanParseName_UnitTest(t *testing.T) {
	var v VlanEntity
	var shortConf = `
        vlan 10
            %s
            state active
            no private-vlan
            trunk group tg1
            `

	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"name BIGDATA", "BIGDATA"},
		{"name VSAN0200", "VSAN0200"},
		{"name default", "default"},
		{"name test", "test"},
		{"name 42", "42"},
		{"name VLAN0400", "VLAN0400"},
		{"name 12df", "12df"},
		{"name back-end", "back-end"},
	}

	for _, tt := range tests {
		testConf := fmt.Sprintf(shortConf, tt.in)
		if got := VlanParseName(&v, testConf); got != tt.want {
			t.Fatalf("parseName() = %q; want %q", got, tt.want)
		}
	}
	if got := VlanParseName(&v, ""); got != "" {
		t.Fatalf("parseName() = %q; want \"\"", got)
	}
}
func TestVlanParseState_UnitTest(t *testing.T) {
	var v VlanEntity
	var shortConf = `
        vlan 10
            name VSAN0100
            %s
            no private-vlan
            trunk group tg1
            `

	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"state active", "active"},
		{"state suspend", "suspend"},
	}

	for _, tt := range tests {
		testConf := fmt.Sprintf(shortConf, tt.in)
		if got := VlanParseState(&v, testConf); got != tt.want {
			t.Fatalf("parseState() = %q; want %q", got, tt.want)
		}
	}
	if got := VlanParseState(&v, ""); got != "" {
		t.Fatalf("parseState() = %q; want \"\"", got)
	}
}

func TestVlanParseTrunkGroup_UnitTest(t *testing.T) {
	var v VlanEntity
	var shortConf = `
        vlan 10
            name VSAN0100
            state active
            no private-vlan
            trunk group %s
            trunk group %s
            trunk group %s
            trunk group %s
            trunk group %s
            trunk group %s
            trunk group %s
            trunk group %s
            `
	tests := [8]struct {
		in   string
		want string
	}{}

	var tn [len(tests)]string

	// for each test entry
	for idx := range tests {
		// get the random strings
		for i := range tn {
			tn[i] = RandomString(2, 14)
		}
		testConf := fmt.Sprintf(shortConf, tn[0], tn[1], tn[2],
			tn[3], tn[4], tn[5], tn[6], tn[7])
		tests[idx].in = testConf
		tests[idx].want = strings.Join(tn[:], ",")
	}
	for _, tt := range tests {
		got := VlanParseTrunkGroups(&v, tt.in)
		if strings.Compare(got, tt.want) != 0 {
			t.Fatalf("parseTrunkGroups() = %q; want %q", got, tt.want)
		}
	}
}

func TestVlanGetConnectionError_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	conn := dummyNode.GetConnection().(*DummyEapiConnection)
	conn.setReturnError(true)

	config := vlan.Get("10")
	if config != nil {
		t.Fatalf("Get() should return nil on underlying error: config: %#v", config)
	}
}

func TestVlanGet_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)

	config := vlan.Get("10")
	if config == nil || config.Name() != "VLAN0010" || config.State() != "active" ||
		config.TrunkGroups() != "tg1" {
		t.Fatalf("Get() returned invalid data. %#v", config)
	}
}

func TestVlanGetAll_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)

	configs := vlan.GetAll()
	if configs == nil || configs["10"] == nil {
		t.Fatalf("GetAll() returned invalid data. %#v", configs)
	}
}

func TestVlanGetSection_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)

	section := vlan.GetSection("10")
	if section == "" {
		t.Fatalf("GetSection() returned nil")
	}
}

func TestVlanGetSectionConnectionError_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	conn := dummyNode.GetConnection().(*DummyEapiConnection)
	conn.setReturnError(true)

	section := vlan.GetSection("10")
	if section != "" {
		t.Fatalf("GetSection() during connection error should return nil")
	}
}

func TestVlanCreate_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(2, 4094))

	cmds := []string{
		"vlan " + vid,
	}
	vlan.Create(vid)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVlanCreateInvalid_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(4095, 65535))
	if ok := vlan.Create(vid); ok {
		t.Fatalf("Invalid vid(%s) used. Should fail", vid)
	}
}

func TestVlanDelete_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(2, 4094))

	cmds := []string{
		"no vlan " + vid,
	}
	vlan.Delete(vid)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVlanDeleteInvalid_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(4095, 65535))
	if ok := vlan.Delete(vid); ok {
		t.Fatalf("Invalid vid(%s) used. Should fail", vid)
	}
}

func TestVlanDefault_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(2, 4094))

	cmds := []string{
		"default vlan " + vid,
	}
	vlan.Default(vid)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVlanDefaultInvalid_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(4095, 65535))
	if ok := vlan.Default(vid); ok {
		t.Fatalf("Invalid vid(%s) used. Should fail", vid)
	}
}

func TestVlanConfigureVlan_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(2, 4094))
	tg := RandomString(1, 32)

	cmds := []string{
		"name Test",
		"state active",
		"trunk group " + tg,
		"no shutdown",
	}
	vlan.ConfigureVlan(vid, cmds...)
	cmds = append([]string{"vlan " + vid}, cmds...)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVlanSetName_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(2, 4094))
	name := RandomString(1, 32)

	cmds := []string{
		"vlan " + vid,
		"name " + name,
	}
	vlan.SetName(vid, name)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVlanSetNameDefault_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(2, 4094))

	cmds := []string{
		"vlan " + vid,
		"default name",
	}
	vlan.SetNameDefault(vid)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVlanSetState_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(2, 4094))

	cmds := []string{
		"vlan " + vid,
		"default state ",
	}
	tests := []struct {
		state string
		want  string
	}{
		{"", "no state"},
		{"active", "state active"},
		{"suspend", "state suspend"},
	}

	for _, tt := range tests {
		vlan.SetState(vid, tt.state)
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

func TestVlanSetStateDefault_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(2, 4094))

	cmds := []string{
		"vlan " + vid,
		"default state",
	}
	vlan.SetStateDefault(vid)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVlanSetTrunkGroupDefault_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(2, 4094))

	cmds := []string{
		"vlan " + vid,
		"default trunk group",
	}
	vlan.SetTrunkGroupDefault(vid)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVlanSetTrunkGroupAdd_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)

	cmds := []string{
		"vlan 10",
		"trunk group tg2",
	}
	vlan.SetTrunkGroup("10", []string{"tg1", "tg2"})
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVlanSetTrunkGroupAddError_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	conn := dummyNode.GetConnection().(*DummyEapiConnection)
	conn.setReturnError(true)

	if ok := vlan.SetTrunkGroup("10", []string{"tg1", "tg2"}); !ok {
		t.Fatalf("SetTrunkGroup expected to fail during connection error test")
	}
}

func TestVlanSetTrunkGroupDel_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)

	cmds := []string{
		"vlan 10",
		"no trunk group tg1",
	}
	vlan.SetTrunkGroup("10", []string{"tg2"})
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVlanAddTrunkGroup_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(2, 4094))
	tg := RandomString(1, 32)

	cmds := []string{
		"vlan " + vid,
		"trunk group " + tg,
	}
	vlan.AddTrunkGroup(vid, tg)
	// first two commands are 'enable', 'configure terminal'
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestVlanRemoveTrunkGroup_UnitTest(t *testing.T) {
	vlan := Vlan(dummyNode)
	vid := strconv.Itoa(RandomInt(2, 4094))
	tg := RandomString(1, 32)

	cmds := []string{
		"vlan " + vid,
		"no trunk group " + tg,
	}
	vlan.RemoveTrunkGroup(vid, tg)
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
func TestVlanGet_SystemTest(t *testing.T) {
	vlanTmp := VlanConfig{
		"vlan_id":      "1",
		"name":         "default",
		"state":        "active",
		"trunk_groups": "",
	}
	for _, dut := range duts {
		cmds := []string{
			"no vlan 1",
			"vlan 1",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		vlanConf := vlan.Get("1")
		if vlanConf == nil {
			t.Fatalf("Failure during Get()")
		}

		for k, v := range vlanConf {
			if v != vlanTmp[k] {
				t.Fatalf("Entry %s: Expected %s but got %s", k, vlanTmp[k], v)
			}
		}
	}
}

func TestVlanGetInvalid_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		cmds := []string{
			"no vlan " + strconv.Itoa(vid),
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		vlanConf := vlan.Get(strconv.Itoa(vid))
		if vlanConf != nil {
			t.Fatalf("Get(%d) of invalid vlan should return nil", vid)
		}
	}
}

func TestVlanGetAll_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no vlan 1-4094",
			"vlan 1",
			"vlan 2",
			"vlan 3",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		vlanConfigs := vlan.GetAll()

		for v := 1; v < 4; v++ {
			if _, found := vlanConfigs[strconv.Itoa(v)]; !found {
				t.Fatalf("Expected entry for vlan %d but not found", v)
			}
		}
	}
}

func TestVlanGetSectionInvalid_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		cmds := []string{
			"no vlan " + strconv.Itoa(vid),
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		if ok := vlan.GetSection(strconv.Itoa(vid)); ok != "" {
			t.Fatalf("GetSection() for invalid vlan should return \"\"")
		}
	}
}

func TestVlanCreateRetTrue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		cmds := []string{
			"no vlan " + strconv.Itoa(vid),
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		if ok := vlan.Create(strconv.Itoa(vid)); !ok {
			t.Fatalf("Failure during Create()")

		}
	}
}

func TestVlanCreateRetFalse_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vlan := Vlan(dut)

		if ok := vlan.Create("5000"); ok {
			t.Fatalf("Expected failure during create of invalid Vlan ID")
		}
	}
}

func TestVlanDeleteRetTrue_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		cmds := []string{
			"vlan " + strconv.Itoa(vid),
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		if ok := vlan.Delete(strconv.Itoa(vid)); !ok {
			t.Fatalf("Failure during Create()")
		}
	}
}

func TestVlanDeleteRetFalse_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vlan := Vlan(dut)

		if ok := vlan.Delete("5000"); ok {
			t.Fatalf("Expected failure for Delete of invalid Vlan ID")
		}
	}
}

func TestVlanSetDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		cmds := []string{
			"vlan " + strconv.Itoa(vid),
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		if ok := vlan.Default(strconv.Itoa(vid)); !ok {
			t.Fatalf("Failure during setting Default()")

		}
	}
}

func TestVlanSetDefaultInvalid_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(4095, 65535)
		vlan := Vlan(dut)

		if ok := vlan.Default(strconv.Itoa(vid)); ok {
			t.Fatalf("Should not allow Default() on invalid Vlan %d", vid)
		}
	}
}

func TestVlanSetName_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		name := RandomString(1, 20)
		cmds := []string{
			"no vlan " + strconv.Itoa(vid),
			"vlan " + strconv.Itoa(vid),
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		if ok := vlan.SetName(strconv.Itoa(vid), name); !ok {
			t.Fatalf("Failure during SetName()")

		}
	}
}

func TestVlanSetNameDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		cmds := []string{
			"no vlan " + strconv.Itoa(vid),
			"vlan " + strconv.Itoa(vid),
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		if ok := vlan.SetNameDefault(strconv.Itoa(vid)); !ok {
			t.Fatalf("Failure during SetNameDefault()")

		}
	}
}

func TestVlanSetStateActive_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		cmds := []string{
			"no vlan " + strconv.Itoa(vid),
			"vlan " + strconv.Itoa(vid),
			"state suspend",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		if ok := vlan.SetState(strconv.Itoa(vid), "active"); !ok {
			t.Fatalf("Failure during SetState()")

		}
	}
}

func TestVlanSetStateSuspend_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		cmds := []string{
			"no vlan " + strconv.Itoa(vid),
			"vlan " + strconv.Itoa(vid),
			"state active",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		if ok := vlan.SetName(strconv.Itoa(vid), "suspend"); !ok {
			t.Fatalf("Failure during SetState()")

		}
	}
}

func TestVlanSetStateDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		cmds := []string{
			"no vlan " + strconv.Itoa(vid),
			"vlan " + strconv.Itoa(vid),
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		if ok := vlan.SetStateDefault(strconv.Itoa(vid)); !ok {
			t.Fatalf("Failure during SetStateDefault()")

		}
	}
}

func TestVlanSetTrunkGroups_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		tg1 := RandomString(1, 10)
		tg2 := RandomString(1, 10)
		tg3 := RandomString(1, 10)
		cmds := []string{
			"no vlan " + strconv.Itoa(vid),
			"vlan " + strconv.Itoa(vid),
			"no trunk group",
			"trunk group " + tg1,
			"trunk group " + tg2,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		confTgs := []string{tg1, tg3}

		if ok := vlan.SetTrunkGroup(strconv.Itoa(vid), confTgs); !ok {
			t.Fatalf("Failure Setting Trunk Group to default")

		}

		show := Show(dut)
		tgs := show.ShowTrunkGroups()

		if len(tgs.TrunkGroups[strconv.Itoa(vid)].Names) != len(confTgs) {
			t.Fatalf("Tg lists not equal: [%#v] [%#v]\n",
				tgs.TrunkGroups[strconv.Itoa(vid)].Names,
				confTgs)
		}

		var found bool
		for _, v1 := range confTgs {
			found = false
			for _, v2 := range tgs.TrunkGroups[strconv.Itoa(vid)].Names {
				if v2 == v1 {
					found = true
				}
			}
			if !found {
				t.Fatalf("Could not find element")
			}
		}
	}
}

func TestVlanSetTrunkGroupsDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		tg := RandomString(1, 32)
		cmds := []string{
			"no vlan " + strconv.Itoa(vid),
			"vlan " + strconv.Itoa(vid),
			"no trunk group",
			"trunk group " + tg,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		if ok := vlan.SetTrunkGroupDefault(strconv.Itoa(vid)); !ok {
			t.Fatalf("Failure Setting Trunk Group to default")

		}

		section := vlan.GetSection(strconv.Itoa(vid))
		tgStr := "trunk group " + tg
		if found, _ := regexp.MatchString(tgStr, section); found {
			t.Fatalf("\"%s\" NOT expected but not seen under "+
				"%s section.\n[%s]", tgStr, cmds[1], section)
		}
	}
}
func TestVlanAddTrunkGroup_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		tg := RandomString(1, 32)
		cmds := []string{
			"no vlan " + strconv.Itoa(vid),
			"vlan " + strconv.Itoa(vid),
			"no trunk group",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		if ok := vlan.AddTrunkGroup(strconv.Itoa(vid), tg); !ok {
			t.Fatalf("Failure Adding Trunk Group")

		}

		section := vlan.GetSection(strconv.Itoa(vid))
		tgStr := "trunk group " + tg
		if found, _ := regexp.MatchString(tgStr, section); !found {
			t.Fatalf("\"%s\" expected but not seen under "+
				"%s section.\n[%s]", tgStr, cmds[1], section)
		}
	}
}

func TestVlanRemoveTrunkGroup_SystemTest(t *testing.T) {
	for _, dut := range duts {
		vid := RandomInt(2, 4094)
		tg := RandomString(1, 32)
		cmds := []string{
			"no vlan " + strconv.Itoa(vid),
			"vlan " + strconv.Itoa(vid),
			"no trunk group",
			"trunk group " + tg,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		vlan := Vlan(dut)

		if ok := vlan.RemoveTrunkGroup(strconv.Itoa(vid), tg); !ok {
			t.Fatalf("Failure Removing Trunk Group")

		}

		section := vlan.GetSection(strconv.Itoa(vid))
		if found, _ := regexp.MatchString(cmds[3], section); found {
			t.Fatalf("\"%s\" is NOT expected but not seen under "+
				"%s section.\n[%s]", cmds[3], cmds[1], section)
		}
	}
}
