package module

import (
	"testing"

	"github.com/aristanetworks/goeapi"
)

func TestShowInterfacesSwitchport_UnitTest(t *testing.T) {
	var dummyNode *goeapi.Node
	var dummyConnection *DummyConnection

	dummyConnection = &DummyConnection{}

	dummyNode = &goeapi.Node{}
	dummyNode.SetConnection(dummyConnection)

	show := Show(dummyNode)
	showInterfacesSwitchport := show.ShowInterfacesSwitchport()

	var scenarios = []struct {
		Interface    string
		AccessVlanID int
		Mode         string
	}{
		{
			Interface:    "Ethernet48",
			AccessVlanID: 1,
			Mode:         "trunk",
		},
	}

	switchports := showInterfacesSwitchport.Switchports

	for _, tt := range scenarios {

		p := switchports[tt.Interface]

		if tt.AccessVlanID != p.SwitchportInfo.AccessVlanID {
			t.Errorf("AccessVlanID does not match: expected %d, got %d", tt.AccessVlanID, p.SwitchportInfo.AccessVlanID)
		}

		if tt.Mode != p.SwitchportInfo.Mode {
			t.Errorf("Mode does not match: expected %s, got %s", tt.Mode, p.SwitchportInfo.Mode)
		}

	}

}
