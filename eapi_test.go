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
package goeapi

import (
	"regexp"
	"testing"
)

type MyShow struct {
	Cmd              string  `json:"-"`
	ModelName        string  `json:"modelName"`
	InternalVersion  string  `json:"internalVersion"`
	SystemMacAddress string  `json:"systemMacAddress"`
	SerialNumber     string  `json:"serialNumber"`
	MemTotal         int     `json:"memTotal"`
	BootupTimestamp  float64 `json:"bootupTimestamp"`
	MemFree          int     `json:"memFree"`
	Version          string  `json:"version"`
	Architecture     string  `json:"architecture"`
	InternalBuildID  string  `json:"internalBuildId"`
	HardwareRevision string  `json:"hardwareRevision"`
}

func (s *MyShow) SetCmd(cmd string) {
	s.Cmd = cmd
}

func (s *MyShow) GetCmd() string {
	if s.Cmd == "" {
		return "show version"
	}
	return s.Cmd
}

func TestEapiGetHandleNodeInvalid_UnitTest(t *testing.T) {
	var node *Node
	h, err := node.GetHandle("json")
	if err == nil {
		t.Fatal("GetHandle invalid failed")
	}
	n, err := h.getNode()
	if n != nil || err == nil {
		t.Fatal("GetHandle invalid failed")
	}
}

func TestEapiRespHandlerInvalidAddCommandStr_UnitTest(t *testing.T) {
	showdummy := new(MyShow)
	node := &Node{}
	h, _ := node.GetHandle("json")
	h = nil
	err := h.AddCommandStr("", showdummy)
	if err == nil {
		t.Fatal("GetHandle invalid failed")
	}
	err = AddCommand(h, showdummy)
	if err == nil {
		t.Fatal("GetHandle invalid failed")
	}
}

func TestEapiRespHandlerAddCommandStrNull_UnitTest(t *testing.T) {
	showdummy := new(MyShow)
	node := &Node{}
	h, _ := node.GetHandle("json")

	err := h.AddCommandStr("", showdummy)
	if err == nil {
		t.Fatal("GetHandle invalid failed")
	}
}

func TestEapiRespHandlerInvalidAddCommand_UnitTest(t *testing.T) {
	showdummy := new(MyShow)
	node := &Node{}
	h, _ := node.GetHandle("json")
	h = nil
	err := AddCommand(h, showdummy)
	if err == nil {
		t.Fatal("GetHandle invalid failed")
	}
}

func TestEapiRespHandlerGetAllCommands_UnitTest(t *testing.T) {
	showdummy := new(MyShow)
	node := &Node{}
	h, _ := node.GetHandle("json")

	tests := [...]string{
		"show version",
		"show vlan",
		"show arp",
		"show running-config",
		"show interfaces",
	}
	for _, val := range tests {
		h.AddCommandStr(val, showdummy)
	}
	cmds := h.getAllCommands()
	if len(tests) != len(cmds) {
		t.Fatalf("length of tests (%d) doesn't not equal length of cmds (%d)",
			len(tests), len(cmds))
	}
	for idx, val := range cmds {
		if tests[idx] != cmds[idx] {
			t.Fatalf("Got %s expected %s", tests[idx], val)
		}
	}
}

func TestEapiRespHandlerAddCommand_UnitTest(t *testing.T) {
	showdummy := new(MyShow)
	node := &Node{}
	h, _ := node.GetHandle("json")

	tests := [...]struct {
		in   int
		want int
	}{
		{1, 1},
		{4, 5},
		{10, 15},
		{20, 35},
		{28, 63},
		{1, 64},
		{1, 64},
	}
	for _, tt := range tests {
		for i := 0; i < tt.in; i++ {
			h.AddCommand(showdummy)
		}
		if got := h.getCmdLen(); got != tt.want {
			t.Fatalf("Got %d want %d", got, tt.want)
		}
	}
}

func TestEapiRespHandlerClearCommands_UnitTest(t *testing.T) {
	showdummy := new(MyShow)
	node := &Node{}
	h, _ := node.GetHandle("json")

	tests := [...]struct {
		in   int
		want int
	}{
		{1, 1},
		{4, 4},
		{10, 10},
		{35, 35},
		{64, 64},
		{120, 64},
	}
	for _, tt := range tests {
		for i := 0; i < tt.in; i++ {
			h.AddCommand(showdummy)
		}
		if got := h.getCmdLen(); got != tt.want {
			t.Fatalf("Got %d want %d", got, tt.want)
		}
		h.clearCommands()
		if count := h.getCmdLen(); count != 0 {
			t.Fatal("Failed to clear commands from cmd list")
		}
	}

	h.Close()
	if count := h.getCmdLen(); count != 0 {
		t.Fatal("Test4: Failed to clear list after handle Close()")
	}

	h = nil
	if ret := h.getCmdLen(); ret != 0 {
		t.Fatalf("Expected 0 but got %d for inv handle call getCmdLen()", ret)
	}
	h.clearCommands()
}

func TestEapiRespHandlerClodeClearCommands_UnitTest(t *testing.T) {
	showdummy := new(MyShow)
	node := &Node{}

	tests := [...]struct {
		in   int
		want int
	}{
		{1, 1},
		{4, 4},
		{10, 10},
		{35, 35},
		{64, 64},
		{120, 64},
	}
	for _, tt := range tests {
		h, _ := node.GetHandle("json")
		for i := 0; i < tt.in; i++ {
			h.AddCommand(showdummy)
		}
		if got := h.getCmdLen(); got != tt.want {
			t.Fatalf("Got %d want %d", got, tt.want)
		}
		h.Close()
		if count := h.getCmdLen(); count != 0 {
			t.Fatal("Failed to clear commands from cmd list")
		}
	}
}

func TestEapiCall_SystemTest(t *testing.T) {
	showdummy := new(MyShow)
	tests := [...]string{
		"show version",
		"show version",
		"show version",
		"show version",
		"show version",
	}
	for _, dut := range duts {
		h, err := dut.GetHandle("json")
		if err != nil {
			t.Fatalf("GetHandle() failed: Error[%s]", err)
		}
		for _, val := range tests {
			h.AddCommandStr(val, showdummy)
		}
		if err = h.Call(); err != nil {
			t.Fatalf("EapiHandle.Call() failed: Error[%s]", err)
		}
		h.Close()
		h = nil
	}
}

func TestEapiEnable_SystemTest(t *testing.T) {
	showdummy := new(MyShow)
	re := regexp.MustCompile(`^([0-9a-fA-F]{2}[:-]){5}([0-9a-fA-F]{2})$`)
	for _, dut := range duts {
		h, err := dut.GetHandle("json")
		if err != nil {
			t.Fatalf("GetHandle() failed: Error[%s]", err)
		}
		if err = h.Enable(showdummy); err != nil {
			t.Fatalf("EapiHandle.Enable() failed: Error[%s]", err)
		}
		match := re.MatchString(showdummy.SystemMacAddress)

		if !match {
			t.Fatal("failed to find mac address")
		}
		h.Close()
		h = nil
	}
}
