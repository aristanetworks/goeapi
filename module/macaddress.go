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

type ShowMACAddressTable struct {
	MulticastTable struct {
		TableEntries []MACAddressTableEntry `json:"tableEntries"`
	} `json:"multicastTable"`
	UnicastTable struct {
		TableEntries []MACAddressTableEntry `json:"tableEntries"`
	} `json:"unicastTable"`
}

type MACAddressTableEntry struct {
	MACAddress string
	LastMove   float64
	Interface  string
	Moves      int
	EntryType  string
	VlanID     int
}

func (a *ShowMACAddressTable) GetCmd() string {
	return "show mac address-table"
}

func (s *ShowEntity) ShowMACAddressTable() (ShowMACAddressTable, error) {
	var showmacaddresstable ShowMACAddressTable

	handle, err := s.node.GetHandle("json")
	if err != nil {
		return showmacaddresstable, err
	}

	err = handle.AddCommand(&showmacaddresstable)
	if err != nil {
		return showmacaddresstable, err
	}

	if err := handle.Call(); err != nil {
		return showmacaddresstable, err
	}

	handle.Close()
	return showmacaddresstable, nil
}
