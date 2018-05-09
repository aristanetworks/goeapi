package module

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aristanetworks/goeapi"
)

type DummyConnection struct {
	goeapi.EapiConnection
	err error
}

func (conn *DummyConnection) Execute(commands []interface{},
	encoding string) (*goeapi.JSONRPCResponse, error) {

	if encoding != "json" {
		return nil, fmt.Errorf("%s encoding not implemented", encoding)
	}

	if conn.err != nil {
		conn.SetError(conn.err)
		return &goeapi.JSONRPCResponse{}, conn.err
	}
	// command 0: enable
	// command 1: show lldp neighbors
	cmd := commands[1].(string)
	fixtureName := strings.Replace(cmd, " ", "_", -1) + ".json"
	// cmd: 'show lldp neighbors' will cause us to look for
	// fixture 'show_lldp_neighbors.json'
	r, err := os.Open(GetFixture(fixtureName))
	if err != nil {
		return nil, fmt.Errorf("Error opening fixture: %s", err)
	}
	defer r.Close()
	return conn.decodeJSONFile(r), nil
}

func (conn *DummyConnection) SetTimeout(uint32) {
}

func (conn *DummyConnection) Error() error {
	return conn.err
}

func (conn *DummyConnection) decodeJSONFile(r io.Reader) *goeapi.JSONRPCResponse {
	dec := json.NewDecoder(r)
	var v goeapi.JSONRPCResponse
	if err := dec.Decode(&v); err != nil {
		panic(err)
	}
	return &v
}

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
