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
	"errors"
	"testing"

	"github.com/aristanetworks/goeapi"
)

func TestShowEnvironmentPower_UnitTest(t *testing.T) {
	var dummyNode *goeapi.Node
	var dummyConnection *DummyConnection

	dummyConnection = &DummyConnection{}

	dummyNode = &goeapi.Node{}
	dummyNode.SetConnection(dummyConnection)

	show := Show(dummyNode)
	showEnvironmentPower, err := show.ShowEnvironmentPower()
	if err != nil {
		t.Errorf("Error during show environment power, %s", err)
	}

	// Test PowerSupplies inside the Environment Power output
	if len(showEnvironmentPower.PowerSupplies) != 2 {
		t.Errorf("2 PowerSupplies expected, %d found", len(showEnvironmentPower.PowerSupplies))
	}

	var powerSupplyScenarios = []struct {
		Number    string
		State     string
		ModelName string
		Capacity  int
	}{
		{
			Number:    "1",
			State:     "ok",
			ModelName: "PWR-500AC-R",
			Capacity:  500,
		},
		{
			Number:    "2",
			State:     "powerLoss",
			ModelName: "PWR-500AC-R",
			Capacity:  500,
		},
	}

	for _, powerSupply := range powerSupplyScenarios {
		if _, ok := showEnvironmentPower.PowerSupplies[powerSupply.Number]; !ok {
			t.Errorf("PowerSupply %s does not exist", powerSupply.Number)
		} else {
			ps := showEnvironmentPower.PowerSupplies[powerSupply.Number]

			if ps.State != powerSupply.State {
				t.Errorf("State does not match expected %s, got %s", ps.State, powerSupply.State)
			}

			if ps.ModelName != powerSupply.ModelName {
				t.Errorf("ModelName does not match expected %s, got %s", ps.ModelName, powerSupply.ModelName)
			}

			if ps.Capacity != powerSupply.Capacity {
				t.Errorf("Capacity does not match expected %d, got %d", ps.Capacity, powerSupply.Capacity)
			}
		}
	}
}

func TestShowEnvironmentPowerErrorDuringCall_UnitTest(t *testing.T) {
	dummyConnection := &DummyConnection{err: errors.New("error during connection")}
	dummyNode := &goeapi.Node{}
	dummyNode.SetConnection(dummyConnection)

	show := Show(dummyNode)
	_, err := show.ShowEnvironmentPower()
	if err == nil {
		t.Errorf("Error expected during show environment power")
	}
}
