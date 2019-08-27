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

type ShowInterfacesSwitchport struct {
	Switchports map[string]Switchport `json:"switchports"`
}

type Switchport struct {
	Enabled        bool           `json:"enabled"`
	SwitchportInfo SwitchportInfo `json:"switchportInfo"`
}

type SwitchportInfo struct {
	AccessVlanID           int           `json:"accessVlanId"`
	AccessVlanName         string        `json:"accessVlanName"`
	DynamicAllowedVlans    struct{}      `json:"dynamicAllowedVlans"`
	DynamicTrunkGroups     []interface{} `json:"dynamicTrunkGroups"`
	MacLearning            bool          `json:"macLearning"`
	Mode                   string        `json:"mode"`
	StaticTrunkGroups      []interface{} `json:"staticTrunkGroups"`
	Tpid                   string        `json:"tpid"`
	TpidStatus             bool          `json:"tpidStatus"`
	TrunkAllowedVlans      string        `json:"trunkAllowedVlans"`
	TrunkingNativeVlanID   int           `json:"trunkingNativeVlanId"`
	TrunkingNativeVlanName string        `json:"trunkingNativeVlanName"`
}

func (l *ShowInterfacesSwitchport) GetCmd() string {
	return "show interfaces switchport"
}

func (s *ShowEntity) ShowInterfacesSwitchport() ShowInterfacesSwitchport {
	handle, _ := s.node.GetHandle("json")
	var showInterfacesSwitchport ShowInterfacesSwitchport
	handle.AddCommand(&showInterfacesSwitchport)
	handle.Call()
	handle.Close()
	return showInterfacesSwitchport
}
