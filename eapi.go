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

// Package goeapi allows for creating connections to EOS eAPI enabled nodes
// using the Connect or ConnectTo function.  Both functions will return an
// instance of a Node object that can be used to send and receive eAPI
// commands.
package goeapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// Request ...
type Request struct {
	Jsonrpc string     `json:"jsonrpc"`
	Method  string     `json:"method"`
	Params  Parameters `json:"params"`
	ID      string     `json:"id"`
}

// Parameters ...
type Parameters struct {
	Version int           `json:"version"`
	Cmds    []interface{} `json:"cmds"`
	Format  string        `json:"format"`
}

// RawJSONRPCResponse ...
type RawJSONRPCResponse struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Result  []json.RawMessage      `json:"result"`
	Error   map[string]interface{} `json:"error"`
	ID      string                 `json:"id"`
}

// JSONRPCResponse ...
type JSONRPCResponse struct {
	Jsonrpc string                   `json:"jsonrpc"`
	Result  []map[string]interface{} `json:"result"`
	ID      string                   `json:"id"`
	Error   *RespError               `json:"error"`
}

// RespError message format breakout
type RespError struct {
	Code    int
	Message string
	Data    interface{}
}

// EapiCommand interface is implemented by any pre-defined response structure
// associated with a command issue toward a node.
type EapiCommand interface {
	GetCmd() string
}

// commandBlock used to map command to an EapiCommand
type commandBlock struct {
	command interface{}
	EapiCommand
}

// Max number of commands to take
const (
	maxCmdBuflen = 64
)

// EapiReqHandle ...
type EapiReqHandle struct {
	node         *Node
	encoding     string
	eapiCommands []commandBlock
	err          error
}

// debugJSON prints out []byte JSON data Indented
func debugJSON(data []byte) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "  ")
	if err == nil {
		fmt.Printf("-------------------- JSON ---------------------------\n")
		fmt.Printf("%s\n", out.Bytes())
		fmt.Printf("-----------------------------------------------------\n")
	}
}

// checkHandle helper function to check the validity of an
// EapiReqHandle
func (handle *EapiReqHandle) checkHandle() error {
	if handle == nil {
		return fmt.Errorf("Invalid EapiReqHandle")
	}
	if handle.node == nil {
		return fmt.Errorf("No connection")
	}
	return nil
}

// getNode returns the Node associated with this EapiReqHandle. Returns error
// if handle is invalid or no Node is associated with this EapiReqHandle.
func (handle *EapiReqHandle) getNode() (*Node, error) {
	if err := handle.checkHandle(); err != nil {
		return nil, err
	}
	return handle.node, nil
}

// AddCommandStr adds a command string with specified EapiCommand type to the
// command block list for this EapiReqHandle.
func AddCommandStr(handle *EapiReqHandle, command string, v EapiCommand) error {
	if err := handle.checkHandle(); err != nil {
		return err
	}
	if command == "" {
		handle.err = fmt.Errorf("Invalid null Command string")
		return handle.err
	}

	if len(handle.eapiCommands) == maxCmdBuflen {
		handle.err = fmt.Errorf("Limit of %d commands reached for AddCommand",
			maxCmdBuflen)
		return handle.err
	}
	cmd := commandBlock{command: command, EapiCommand: v}
	handle.eapiCommands = append(handle.eapiCommands, cmd)
	return nil
}

// AddCommandStr adds a command string with specified EapiCommand type to the
// command block list for this EapiReqHandle.
func (handle *EapiReqHandle) AddCommandStr(command string, v EapiCommand) error {
	return AddCommandStr(handle, command, v)
}

// AddCommand adds a pre-defined EapiCommand type to the command
// block list for this EapiReqHandle.
func AddCommand(handle *EapiReqHandle, v EapiCommand) error {
	command := v.GetCmd()
	return AddCommandStr(handle, command, v)
}

// AddCommand adds a pre-defined EapiCommand type to the command
// block list for this EapiReqHandle.
func (handle *EapiReqHandle) AddCommand(v EapiCommand) error {
	command := v.GetCmd()
	return AddCommandStr(handle, command, v)
}

// getAllCommands iterates through the list of command blocks
// and returns the commands as an array of interfaces.
func (handle *EapiReqHandle) getAllCommands() []interface{} {
	if err := handle.checkHandle(); err != nil {
		return nil
	}
	var tmp []interface{}

	if handle.eapiCommands == nil {
		return nil
	}

	for _, cmd := range handle.eapiCommands {
		tmp = append(tmp, cmd.command)
	}
	return tmp
}

// getCmdLen returns the length of the EapiReqHandle command block
// list.
func (handle *EapiReqHandle) getCmdLen() int {
	if handle == nil {
		return 0
	}
	return len(handle.eapiCommands)
}

// clearCommands frees all entries from the EapiReqHandle command
// block list
func (handle *EapiReqHandle) clearCommands() {
	if handle == nil {
		return
	}
	handle.eapiCommands = nil
}

// Call executes the commands previously added to the command block
// using AddCommand().
//
// Responses from issued commands are stored in the EapiCommand associated
// with that commands response.
//
// Returns:
//  error if handle is invalid, or problem encountered during sending or
//  receiveing.
func (handle *EapiReqHandle) Call() error {
	if err := handle.checkHandle(); err != nil {
		return err
	}
	if handle.err != nil {
		return handle.err
	}

	var cmd interface{}
	if handle.node.enablePasswd != "" {
		cmd = map[string]string{
			"cmd":   "enable",
			"input": handle.node.enablePasswd,
		}
	} else {
		cmd = "enable"
	}
	tmpSlice := []commandBlock{{command: cmd,
		EapiCommand: nil}}
	handle.eapiCommands = append(tmpSlice, handle.eapiCommands...)

	commands := handle.getAllCommands()

	jsonrsp, err := handle.node.conn.Execute(commands, handle.encoding)
        
	if err != nil {
		return err
	}

	err = handle.parseResponse(jsonrsp)
	handle.clearCommands()
	return err
}

// Enable takes an EapiCommand type to issue toward the Node.
// Decoded results are stored in the EapiCommand.
// Returns:
//  error on failure
func (handle *EapiReqHandle) Enable(v EapiCommand) error {
	if err := handle.AddCommand(v); err != nil {
		return err
	}
	return handle.Call()
}

// Close closes the relationship between this EapiReqHandle and
// a Node. Any queued commands are deleted.  User is responsible for
// proper handling of the handle after the close.
func (handle *EapiReqHandle) Close() error {
	return Close(handle)
}

// Close closes the relationship between the EapiReqHandle supplied
// and a Node. Any queued commands are deleted. User is responsible for
// proper handling of the handle after the close.
func Close(handle *EapiReqHandle) error {
	if handle == nil {
		return fmt.Errorf("Invalid EapiReqHandle")
	}
	handle.clearCommands()
	handle.node = nil
	return nil
}

// parseResponse is a speciallized function to parse a JSON response for
// one (or many) command(s) and store the result in the command block
// associated with matching command request.
func (handle *EapiReqHandle) parseResponse(resp *JSONRPCResponse) error {
	var err error

	// check for errors in the JSON response
	if resp.Error != nil {
		err := fmt.Errorf("JSON Error(%d): %s", resp.Error.Code,
			resp.Error.Message)
		return err
	}

	if len(resp.Result) != len(handle.eapiCommands) {
		err := fmt.Errorf("Number of Result entries(%d) does not match"+
			"commands sent(%d)",
			len(resp.Result), len(handle.eapiCommands))
		return err
	}

	for index, result := range resp.Result {
		cmd := handle.eapiCommands[index]
		if cmd.EapiCommand == nil {
			continue
		}

                d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ TagName: "json", Result: cmd.EapiCommand })
                if err != nil {
			return err
                } 

                err = d.Decode(result)
		if err != nil {
			return err
		}
	}
	return err
}

// decodeEapiResponse [private] Used to decode JSON Response into
// structure format defined by type JSONRPCResponse
func decodeEapiResponse(resp *http.Response) (*JSONRPCResponse, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Http error: %s", resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	var v JSONRPCResponse
	if err := dec.Decode(&v); err != nil {
		log.Println(err)
		return nil, err
	}

	if v.Error != nil {
		err := fmt.Errorf("JSON Error(%d): %s", v.Error.Code, v.Error.Message)
		return &v, err
	}
	return &v, nil
}
