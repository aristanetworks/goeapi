package module

import (
	"testing"

	"github.com/aristanetworks/goeapi"
)

func TestShowARP_UnitTest(t *testing.T) {
	var dummyNode *goeapi.Node
	var dummyConnection *DummyConnection

	dummyConnection = &DummyConnection{}

	dummyNode = &goeapi.Node{}
	dummyNode.SetConnection(dummyConnection)

	show := Show(dummyNode)
	showArp, err := show.ShowARP()
	if err != nil {
		t.Errorf("Error during Show ARP, %s", err)
	}

	var scenarios = []struct {
		HWAddress string
		Address   string
		Interface string
	}{
		{
			HWAddress: "444c.a8aa.bbcc",
			Address:   "10.10.10.10",
			Interface: "Ethernet1/1",
		},
		{
			HWAddress: "444c.a8bb.ccdd",
			Address:   "10.10.10.20",
			Interface: "Ethernet2/1",
		},
		{
			HWAddress: "444c.a8cc.ddee",
			Address:   "10.10.10.30",
			Interface: "Ethernet3/1",
		},
		{
			HWAddress: "444c.a8dd.eeff",
			Address:   "10.10.10.40",
			Interface: "Ethernet4/1",
		},
	}

	neighbors := showArp.IPv4Neighbors

	for i, tt := range scenarios {
		if tt.HWAddress != neighbors[i].HWAddress {
			t.Errorf("HWAddress does not match expected %s, got %s", tt.HWAddress, neighbors[i].HWAddress)
		}

		if tt.Address != neighbors[i].Address {
			t.Errorf("Address does not match expected %s, got %s", tt.Address, neighbors[i].Address)
		}

		if tt.Interface != neighbors[i].Interface {
			t.Errorf("Interface does not match expected %s, got %s", tt.Interface, neighbors[i].Interface)
		}
	}
}
