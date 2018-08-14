package module

import (
	"testing"

	"github.com/aristanetworks/goeapi"
)

func TestShowLLDPNeighbors_UnitTest(t *testing.T) {
	var dummyNode *goeapi.Node
	var dummyConnection *DummyConnection

	dummyConnection = &DummyConnection{}

	dummyNode = &goeapi.Node{}
	dummyNode.SetConnection(dummyConnection)

	show := Show(dummyNode)
	showLldpNeighbors := show.ShowLLDPNeighbors()

	var scenarios = []struct {
		NeighborDevice string
		NeighborPort   string
		Port           string
	}{
		{
			NeighborDevice: "testsw1.aristanetworks.com",
			NeighborPort:   "Ethernet48",
			Port:           "Ethernet3",
		},
		{
			NeighborDevice: "testsw1.aristanetworks.com",
			NeighborPort:   "Ethernet47",
			Port:           "Ethernet4",
		},
		{
			NeighborDevice: "testsw1.aristanetworks.com",
			NeighborPort:   "Ethernet14",
			Port:           "Ethernet14",
		},
		{
			NeighborDevice: "testsw1.aristanetworks.com",
			NeighborPort:   "Ethernet33",
			Port:           "Management1",
		},
	}

	neighbors := showLldpNeighbors.LLDPNeighbors

	for i, tt := range scenarios {
		if tt.Port != neighbors[i].Port {
			t.Errorf("Port does not match expected %s, got %s", tt.Port, neighbors[i].Port)
		}

		if tt.NeighborPort != neighbors[i].NeighborPort {
			t.Errorf("NeighborPort does not match expected %s, got %s", tt.NeighborPort, neighbors[i].NeighborPort)
		}

		if tt.NeighborDevice != neighbors[i].NeighborDevice {
			t.Errorf("NeighborDevice does not match expected %s, got %s", tt.NeighborDevice, neighbors[i].NeighborDevice)
		}
	}

	//fmt.Printf("%+v\n", showLldpNeighbors)
}
