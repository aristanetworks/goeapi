package module

import (
	"errors"
	"testing"

	"github.com/aristanetworks/goeapi"
)

func TestShowMACAddressTable_UnitTest(t *testing.T) {
	var dummyNode *goeapi.Node
	var dummyConnection *DummyConnection

	dummyConnection = &DummyConnection{}

	dummyNode = &goeapi.Node{}
	dummyNode.SetConnection(dummyConnection)

	show := Show(dummyNode)
	showMACAddressTable, err := show.ShowMACAddressTable()
	if err != nil {
		t.Errorf("Error during Show MAC Address Table, %s", err)
	}

	var scenarios = []struct {
		MACAddress string
		Interface  string
		VlanID     int
	}{
		{
			MACAddress: "ab:cd:ef:12:34:56",
			Interface:  "Ethernet1",
			VlanID:     123,
		},
		{
			MACAddress: "ab:cd:ef:78:90:12",
			Interface:  "Ethernet2",
			VlanID:     456,
		},
	}

	unicastTableEntries := showMACAddressTable.UnicastTable.TableEntries

	for i, tt := range scenarios {
		if tt.MACAddress != unicastTableEntries[i].MACAddress {
			t.Errorf("MACAddress does not match expected %s, got %s", tt.MACAddress, unicastTableEntries[i].MACAddress)
		}

		if tt.Interface != unicastTableEntries[i].Interface {
			t.Errorf("Interface does not match expected %s, got %s", tt.Interface, unicastTableEntries[i].Interface)
		}

		if tt.VlanID != unicastTableEntries[i].VlanID {
			t.Errorf("VlanID does not match expected %d, got %d", tt.VlanID, unicastTableEntries[i].VlanID)
		}
	}
}

func TestShowMACAddressTableErrorDuringCall_UnitTest(t *testing.T) {
	dummyConnection := &DummyConnection{err: errors.New("error during connection")}
	dummyNode := &goeapi.Node{}
	dummyNode.SetConnection(dummyConnection)

	show := Show(dummyNode)
	_, err := show.ShowMACAddressTable()
	if err == nil {
		t.Errorf("Error expected during show mac address-table")
	}
}
