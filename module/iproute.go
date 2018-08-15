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

type ShowIPRoute struct {
	VRFs map[string]Routes `json:"vrfs"`
}

type Routes struct {
	Routes map[string]Route `json:"routes"`
}

type Route struct {
	KernelProgrammed  bool   `json:"kernelProgrammed"`
	DirectlyConnected bool   `json:"directlyConnected"`
	Preference        int    `json:"preference"`
	RouteAction       string `json:"routeAction"`
	Vias              []struct {
		Interface   string `json:"interface"`
		NexthopAddr string `json:"nexthopAddr"`
	} `json:"vias"`
	Metric             int    `json:"metric"`
	HardwareProgrammed bool   `json:"hardwareProgrammed"`
	RouteType          string `json:"routeType"`
}

func (r *ShowIPRoute) GetCmd() string {
	return "show ip route"
}

func (s *ShowEntity) ShowIPRoute() ShowIPRoute {
	handle, _ := s.node.GetHandle("json")
	var showiproute ShowIPRoute
	handle.AddCommand(&showiproute)
	handle.Call()
	handle.Close()
	return showiproute
}
