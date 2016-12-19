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

// Package api
// This module provides a set of classes that are used to build API modules
// that work with Node objects.
// All API modules will ultimately derive from AbstractBaseEntity which provides
// some common functions to make building API modules easier.

package module

import "github.com/aristanetworks/goeapi"

// AbstractBaseEntity object for all resources to derive from
//
// This AbstractBaseEntity object should not be directly instatiated.  It is
// designed to be implemented by all resource classes to provide common
// methods.
// Attributes:
//  node (Node): The node instance this resource will perform operations
//               against for configuration
//  config (Config): Returns an instance of Config with the nodes
//                   current running configuration
//  error (CommandError): Holds the latest Error exception
//                    instance if raised
type AbstractBaseEntity struct {
	node *goeapi.Node
}

// Config returns the current running configuration
// Returns:
//      String: running config
func (b *AbstractBaseEntity) Config() string {
	return b.node.RunningConfig()
}

// Error returns the current error exception
// Returns:
//      Error: current error
func (b *AbstractBaseEntity) Error() error {
	return b.node.GetConnection().Error()
}

// GetBlock scans the config and returns a block of code
//
// Args:
//  parent (str): The parent string to search the config for and
//                return the block
//
// Returns:
//  A string that represents the block from the config.  If
//  the parent string is not found, then this method will
//  return None.
func (b *AbstractBaseEntity) GetBlock(parent string) (string, error) {
	parent = `(?m)^` + parent + `$`
	return b.node.GetSection(parent, "running-config")
}

// Configure sends the commands list to the node in config mode
//
// This method performs configuration the node using the array of
// commands specified.
//
// Args:
//  commands (list): A list of commands to be sent to the node in
//                   config mode
//
// Returns:
//  True if the commands are executed without exception otherwise
//  False is returned
func (b *AbstractBaseEntity) Configure(commands ...string) bool {
	return b.node.Config(commands...)
}

// CommandBuilder builds a command with keywords
//
// Args:
//  cmd (string): The Command string
//  value (string): The configuration setting to substitute into the command
//                  string.
//  def (bool):    Specifies if command should use default keyword argument
//  enable (bool): Specifies if command is enabled or disabled
// Returns:
//  A command string that can be used to configure the node
func (b *AbstractBaseEntity) CommandBuilder(cmd string, value string, def bool,
	enable bool) string {
	if value != "" {
		cmd = cmd + " " + value
	}
	if def {
		return "default " + cmd
	}
	if !enable {
		return "no " + cmd
	}
	return cmd
}

// ConfigureInterface Configures the specified interface with the commands
//
// Args:
//  name (str): The interface name to configure
//  commands: The commands to configure in the interface
//
// Returns:
//  True if the commands completed successfully
func (b *AbstractBaseEntity) ConfigureInterface(name string, commands ...string) bool {
	var cmd = []string{"interface " + name}
	commands = append(cmd, commands...)
	return b.Configure(commands...)
}
