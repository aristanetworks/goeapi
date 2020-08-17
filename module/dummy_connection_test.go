package module

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aristanetworks/goeapi"
)

type DummyConnection struct {
	goeapi.EapiConnection
	err error
}

func (conn *DummyConnection) Execute(commands []interface{},
	encoding string, streaming bool) (*goeapi.JSONRPCResponse, error) {

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
