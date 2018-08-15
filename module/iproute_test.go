package module

import (
	"testing"

	"github.com/aristanetworks/goeapi"
)

func TestShowIPRoute_UnitTest(t *testing.T) {
	var dummyNode *goeapi.Node
	var dummyConnection *DummyConnection

	dummyConnection = &DummyConnection{}

	dummyNode = &goeapi.Node{}
	dummyNode.SetConnection(dummyConnection)

	show := Show(dummyNode)
	showIPRoute := show.ShowIPRoute()

	var scenarios = []struct {
		Route             string
		DirectlyConnected bool
		RouteType         string
	}{
		{
			Route:             "10.1.2.0/20",
			DirectlyConnected: true,
			RouteType:         "connected",
		},
		{
			Route:             "0.0.0.0/0",
			DirectlyConnected: false,
			RouteType:         "static",
		},
	}

	defaultVrfRoutes := showIPRoute.VRFs["default"].Routes

	for _, tt := range scenarios {
		r := defaultVrfRoutes[tt.Route]

		if tt.DirectlyConnected != r.DirectlyConnected {
			t.Errorf("DirectlyConnected does not match expected %t, got %t", tt.DirectlyConnected, r.DirectlyConnected)
		}

		if tt.RouteType != r.RouteType {
			t.Errorf("RouteType does not match expected %s, got %s", tt.RouteType, r.RouteType)
		}
	}
}
