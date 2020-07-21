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

package goeapi

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/vaughan0/go-ini"
)

//
const (
	RunningConfig = "running-config"
	StartupConfig = "startup-config"
)

// A Node represents a single device for sending and receiving eAPI messages
//
// Node provides an instance for communicating with Arista EOS
// devices.  The Node object provides easy to use methods for sending both
// enable and config commands to the device using a specific transport.  This
// object forms the base for communicating with devices.
type Node struct {
	conn          EapiConnectionEntity
	runningConfig string
	startupConfig string
	autoRefresh   bool
	enablePasswd  string
}

// GetConnection returns the EapiConnectionEntity
// associtated with this Node.
//
// Returns:
//  EapiConnectionEntity
func (n *Node) GetConnection() EapiConnectionEntity {
	if n == nil {
		return nil
	}
	return n.conn

}

// SetConnection sets the EapiConnectionEntity
// associtated with this Node.
func (n *Node) SetConnection(c EapiConnectionEntity) {
	if n == nil {
		return
	}
	n.conn = c

}

// SetAutoRefresh sets the current nodes auto refresh attribute to either
// true or false.
//
// Args:
//  val (bool): If True, the running-config and startup-config are
//              refreshed on config events.  If False, then the config
//              properties must be manually refreshed.
func (n *Node) SetAutoRefresh(val bool) {
	n.autoRefresh = val
}

// EnableAuthentication configures the enable mode authentication
// password present in passwd
//
// Args:
//  passwd (string): The password string in clear text used to
//                   authenticate to exec mode
func (n *Node) EnableAuthentication(passwd string) {
	n.enablePasswd = strings.TrimSpace(passwd)
}

// RunningConfig returns the running configuration for the Arista EOS
// device. A copy is cached locally if one does not already exist.
//
// Returns:
//  String format of the running config
func (n *Node) RunningConfig() string {
	if n.runningConfig != "" {
		return n.runningConfig
	}
	n.runningConfig, _ = n.getConfigText(RunningConfig, "all")
	return n.runningConfig
}

// StartupConfig returns the startup configuration for the Arista EOS
// device. A copy is cached locally if one does not already exist.
//
// Returns:
//  String format of the startup config
func (n *Node) StartupConfig() string {
	if n.startupConfig != "" {
		return n.startupConfig
	}
	n.startupConfig, _ = n.getConfigText(StartupConfig, "")
	return n.startupConfig
}

// Refresh refreshes the config properties.
//
// This method will refresh the runningConfig and startupConfig
// properites.  Since the properties are lazily loaded, this method will
// clear the current internal instance variables.  On the next call the
// instance variables will be repopulated with the current config
func (n *Node) Refresh() {
	n.runningConfig = ""
	n.runningConfig = ""
}

// GetHandle returns the EapiReqHandle for the connection.
//
// Args:
//  n (*Node): Node for which we are aquiring an EapiReqHandle
//  encoding (string): Encoding to be used
//
// Returns:
//  Pointer to an EapiReqHandle or error on failure
func GetHandle(n *Node, encoding string) (*EapiReqHandle, error) {
	if strings.ToLower(encoding) != "json" &&
		strings.ToLower(encoding) != "text" {
		return nil, fmt.Errorf("Invalid encoding specified: %s", encoding)
	}
	if n == nil {
		return nil, fmt.Errorf("Invalid node")
	}
	return &EapiReqHandle{node: n, encoding: encoding}, nil
}

// GetHandle returns the EapiReqHandle for the connection.
//
// Args:
//  encoding (string): Encoding to be used
//
// Returns:
//  Pointer to an EapiReqHandle or error on failure
func (n *Node) GetHandle(encoding string) (*EapiReqHandle, error) {
	return GetHandle(n, encoding)
}

// GetConfig retrieves the config from the node.
//
// The config to retrieve can be specified as either
// the startup-config or the running-config. An error is returned on
// invalid parameter or if the underlying transmit failed.
//
// Args:
//  config (string): Specifies to return either the nodes startup-config
//                or running-config.  The default value is the running-config
//  params (string): A string of keywords to append to the command for
//                retrieving the config.
//  encoding (string): Encoding to be used
// Returns:
//  Will return the config requested or error if failure
func (n *Node) GetConfig(config, params, encoding string) (map[string]interface{}, error) {
	if config != RunningConfig && config != StartupConfig {
		return nil, fmt.Errorf("Invalid config type: %s", config)
	}
	commands := []string{strings.TrimSpace("show " + config + " " + params)}

	result, err := n.RunCommands(commands, encoding)
	if err != nil {
		return nil, err
	}
	return result.Result[0], nil
}

func (n *Node) getConfigText(config, params string) (string, error) {
	result, err := n.GetConfig(config, params, "text")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result["output"].(string)), nil
}

// GetSection retrieves the config section from the Node
//
// Args:
//  regex (string):
//  config (string):
//
// Returns:
//  String value of the config section requested.
//  Error returned on failure.
func (n *Node) GetSection(regex string, config string) (string, error) {
	var params string
	if config == "" || config == RunningConfig {
		config = RunningConfig
		params = "all"
	}
	if config != RunningConfig && config != StartupConfig {
		return "", fmt.Errorf("Invalid config type: %s", config)
	}
	sectionRegex, err := regexp.Compile(regex)
	if err != nil {
		return "", fmt.Errorf("Invalid regexp")
	}
	config, err = n.getConfigText(config, params)
	if err != nil || config == "" {
		return "", err
	}

	match := sectionRegex.FindStringIndex(config)
	if match == nil {
		return "", fmt.Errorf("Config section not found %d", match)
	}

	blockStart := match[0]
	lineEnd := match[1]

	blockRegex := regexp.MustCompile(`(?m)^[^\s]`)
	match = blockRegex.FindStringIndex(config[lineEnd:])
	if match == nil {
		return "", fmt.Errorf("Block section/end not found")
	}

	blockEnd := match[0]
	blockEnd = lineEnd + blockEnd

	return config[blockStart:blockEnd], nil
}

// ConfigWithErr the node with the specified commands
//
// This method is used to send configuration commands to the node.
// It will takes a list of strings and prepend the necessary commands
// to put the session into config mode. Returns error if issues arise.
func (n *Node) ConfigWithErr(commands ...string) error {
	commands = append([]string{"configure terminal"}, commands...)
	_, err := n.RunCommands(commands, "json")
	if n.autoRefresh {
		n.Refresh()
	}
	return err
}

// Config the node with the specified commands
//
// This method is used to send configuration commands to the node.
// It will takes a list of strings and prepend the necessary commands
// to put the session into config mode.
func (n *Node) Config(commands ...string) bool {
	err := n.ConfigWithErr(commands...)
	return (err == nil)
}

// Enable issues an array of commands to the node in enable mode
//
// This method will send the commands to the node and evaluate
// the results.  (TODO) If a command fails due to an encoding error,
// then the command set will be re-issued individual with text
// encoding.
//
// Args:
//  commands (string array): The list of commands to send to the node
// Returns:
//  An array of map'd interfaces that includes the response for each
//  command along with the encoding. Error is returned on failure.
func (n *Node) Enable(commands []string) ([]map[string]string, error) {
	for _, cmd := range commands {
		found, _ := regexp.MatchString(`^\s*configure(\s+terminal)?\s*$`, cmd)
		if found {
			return nil, fmt.Errorf("Config mode commands not supported")
		}
	}

	results := make([]map[string]string, len(commands))
	rsp, err := n.RunCommands(commands, "text")
	if err != nil {
		return results, err
	}
	for idx, resp := range rsp.Result {
		results[idx] = make(map[string]string)
		results[idx]["command"] = commands[idx]
		results[idx]["result"] = strings.TrimSpace(resp["output"].(string))
	}
	return results, nil
}

// RunCommands sends the commands over the transport to the device
//
// This method sends the commands to the device using the nodes
// transport.  This is a lower layer function that shouldn't normally
// need to be used, prefering instead to use config() or enable().
//
// Args:
//  commands (array): The ordered list of commands to send to the
//                   device using the transport
//  encoding (string): The encoding method to use for the request and
//                  excpected response. ('json' or 'text')
//
// Returns:
//  This method will return the raw response from the connection
//  which is a JSONRPCResponse object or error on failure.
func (n *Node) RunCommands(commands []string,
	encoding string) (*JSONRPCResponse, error) {
	var cmds []interface{}

	// Check to see if enablePasswd has been set. In the case where
	// enablePassword is provided, the following cmds value format would let
	// you enter exec mode and clear interface counters
	//
	// [ { "cmd": "enable", "input": <enablePasswd> },  "clear counters" ]
	//
	// In these cases we prepend this sequence to the commands.
	if n.enablePasswd != "" {
		cmds = n.prependEnableSequence(commands)
	} else {
		commands = append([]string{"enable"}, commands...)
		cmds = cmdsToInterface(commands)
	}

	result, err := n.conn.Execute(cmds, encoding)
	if err != nil {
		return nil, err
	}
	// pop the result for enable off the result list
	result.Result = append(result.Result[:0], result.Result[1:]...)
	return result, err
}

// prependEnableSequence helper fuction to convert the provided array of
// strings (commands) to type []interface{} and prepends with the entry for
// map[string]interface {"cmd":"enable","info":enablePasswd}
//
// Args:
//  commands (string array): list of commands to convert
//
// Returns:
//  An array of []interface{} if successful.
func (n *Node) prependEnableSequence(commands []string) []interface{} {
	length := len(commands) + 1

	var interfaceSlice []interface{}
	interfaceSlice = make([]interface{}, length)
	interfaceSlice[0] = map[string]interface{}{
		"cmd":   "enable",
		"input": n.enablePasswd,
	}

	for i := 1; i < length; i++ {
		interfaceSlice[i] = commands[i-1]
	}
	return interfaceSlice
}

// cmdsToInterface is a helper fuction that converts a given array
// of strings (commands) to an array of interfaces.
//
// Args:
//  commands (string array): list of commands
//
// Returns:
//  Interface array of converted commands
func cmdsToInterface(commands []string) []interface{} {
	if commands == nil || len(commands) == 0 {
		return nil
	}
	var interfaceSlice []interface{}
	length := len(commands)

	interfaceSlice = make([]interface{}, length)

	for i := 0; i < length; i++ {
		interfaceSlice[i] = commands[i]
	}
	return interfaceSlice

}

var configGlobal = NewEapiConfig()

var configSearchPath = []string{
	"~/.eapi.conf",
	"/mnt/flash/eapi.conf",
}

type fn func(transport string, host string, username string,
	password string, port int) EapiConnectionEntity

// transports provides the method
var transports = map[string]fn{
	"socket":     NewSocketEapiConnection,
	"http_local": NewHTTPLocalEapiConnection,
	"http":       NewHTTPEapiConnection,
	"https":      NewHTTPSEapiConnection,
}

// EapiConfig provides the instance for managing of eapi.conf file.
// We embed ini.File here to use properties of the ini.File type.
type EapiConfig struct {
	// full path to the loaded filename
	filename string
	ini.File
}

// NewEapiConfig creates a new EapiConfig instance and initiates
// the autoload.
func NewEapiConfig() *EapiConfig {
	config := &EapiConfig{}
	config.AutoLoad()
	return config
}

// NewEapiConfigFile creates a new EapiConfig instance with
// the provided file name. After setting the filename, the method
// initiates the autoload for the config file.
//
// Args:
//  filename (string): filename/path of the eapi.conf file.
func NewEapiConfigFile(filename string) *EapiConfig {
	config := &EapiConfig{filename: filename}
	config.AutoLoad()
	return config
}

// AutoLoad loads the eapi.conf file
//
// This method will use the module variable CONFIG_SEARCH_PATH to
// attempt to locate a valid eapi.conf file if a filename is not already
// configured.   This method will load the first eapi.conf file it
// finds and then return.
//
// The CONFIG_SEARCH_PATH can be overridden using an environment variable
// by setting EAPI_CONF.
func (e *EapiConfig) AutoLoad() {
	var searchPath []string
	path := os.Getenv("EAPI_CONF")
	if path == "" {
		if e.filename != "" {
			path = e.filename
		}
	}

	if path != "" {
		searchPath = append(searchPath, path)
	} else {
		searchPath = append(searchPath, configSearchPath...)
	}

	for _, file := range searchPath {
		file, _ = expandPath(file)

		if _, err := os.Stat(file); err == nil {
			e.filename = file
			e.Read(file)
			return
		}
	}
	e.File = make(ini.File)
	e.addDefaultConnection()
	return
}

// Connections returns all of the loaded connections names as a list
func (e *EapiConfig) Connections() []string {
	if e == nil {
		return nil
	}
	var connections []string
	for name := range e.File {
		str := strings.Replace(name, "connection:", "", 1)
		connections = append(connections, str)
	}
	return connections
}

// Connections returns all of the loaded connections names as a list
func Connections() []string {
	return configGlobal.Connections()
}

// Read reads the file specified by filename
//
// This method will load the eapi.conf file specified by filename into
// the instance object.  It will also add the default connection localhost
// if it was not defined in the eapi.conf file
// Args:
//  filename (string): The full path to the file to load
func (e *EapiConfig) Read(filename string) error {
	file, err := ini.LoadFile(filename)
	if err != nil {
		return fmt.Errorf("Cant read filename: %s, %#v", filename, err)
	}
	e.File = file

	// for each section
	for name := range e.File {
		if _, found := e.Get(name, "host"); !found {
			e.Section(name)["host"] = strings.Split(name, ":")[1]
		}
	}
	e.addDefaultConnection()
	return nil
}

// Load loads the file specified by filename
//
// This method works in conjunction with the autoload method to load the
// file specified by filename.
//
// Args:
//  filename (string): The full path to the file to be loaded
// Returns:
//  bool: True if successful
func (e *EapiConfig) Load(filename string) bool {
	e.filename = filename
	e.Reload()
	return true
}

// Reload reloades the configuration
//
// This method will reload the configuration instance using the last
// known filename.  Note this method will initially clear the
// configuration and reload all entries.
//
// Returns:
//  bool: True if successful
func (e *EapiConfig) Reload() bool {
	for name := range e.File {
		delete(e.File, name)
	}
	e.AutoLoad()
	return true
}

// GetConnection returns the properties for a connection name
//
// This method will return the settings for the configuration specified
// by name.  Note that the name argument should only be the name.
//
// For instance, give the following eapi.conf file
//
// .. code-block:: ini
//
//  [connection:veos01]
//  transport: http
//
// Args:
//  name (string): The name of the connection to return
//
// Returns:
//  ini.Section object of key/value pairs that represent
//  the node configuration.  If the name provided in the argument
//  is not found, then nil is returned.
func (e *EapiConfig) GetConnection(name string) ini.Section {
	name = "connection:" + name
	section, found := e.File[name]
	if !found {
		return nil
	}
	return section
}

// AddConnection adds a connection to the configuration
//
// This method will add a connection to the configuration.  The connection
// added is only available for the lifetime of the object and is not
// persisted.
//
// Note:
//  If a call is made to load() or reload(), any connections added
//  with this method must be re-added to the config instance
//
// Args:
//  name (string): The name of the connection to add to the config.  The
//              name provided will automatically be prepended with the string
//              connection:
// Returns:
//  bool: True if successful
func (e *EapiConfig) AddConnection(name string) ini.Section {
	return e.Section("connection:" + name)
}

// addDefaultConnection checks the loaded config and adds the
// localhost profile if needed
//
// This method wil load the connection:localhost profile into the client
// configuration if it is not already present.
func (e *EapiConfig) addDefaultConnection() {
	conn := e.GetConnection("localhost")
	if conn == nil {
		e.AddConnection("localhost")["transport"] = "socket"
	}
}

// LoadConfig function method that loads a conf file
//
// This function will load the file specified by filename into the config
// instance.   Its a convenience function that calls load on the config
// instance
//
//Args:
//  filename (string): The full path to the filename to load
func LoadConfig(filename string) {
	configGlobal.Load(filename)
}

// ConfigFor function to get settings for named config
//
// This function will return the settings for a specific connection as
// specified by name.  Its a convenience function that calls get_connection
// on the global config instance
//
// Args:
//  name (string): The name of the connection to return.  The connection
//              name is specified as the string right of the : in the INI file
// Returns:
//  An ini.Section object of key/value pairs that represent the
//  nodes configuration settings from the config instance
func ConfigFor(name string) ini.Section {
	return configGlobal.GetConnection(name)
}

// ConnectTo Creates a Node instance based on an entry from the config
//
// This function will retrieve the settings for the specified connection
// from the config and return a Node instance.  The configuration must
// be loaded prior to calling this function.
//
// Args:
//  name (string): The name of the connection to load from the config.  The
//              name argument should be the connection name (everything
//              right of the colon from the INI file)
// Returns:
//  This function will return an instance of Node with the settings
//  from the config instance.
func ConnectTo(name string) (*Node, error) {
	section := ConfigFor(name)
	if section == nil {
		return nil, fmt.Errorf("Connection profile not found in config")
	}
	host := section["host"]
	username := section["username"]
	passwd := section["password"]
	transport := section["transport"]
	enablepwd := section["enablepwd"]
	var port = UseDefaultPortNum
	_, ok := section["port"]
	if ok {
		port, _ = strconv.Atoi(section["port"])
	}
	node, err := Connect(transport, host, username, passwd, port)
	if err != nil {
		return nil, err
	}
	node.EnableAuthentication(enablepwd)
	return node, nil
}

// Connect establishes a connection (using the supplied settings) and creates a Node instance.
//
// This function will create a connection to an Arista EOS node using
// the arguments.  All arguments are optional with default values.
//
// Args:
//  transport (string): Specifies the type of connection transport to use.
//                   Valid values for the connection are socket, http_local,
//                   http, and https.  https is the default.
//  host (string): The IP addres or DNS host name of the connection device.
//              The default value is 'localhost'
//  username (string): The username to pass to the device to authenticate
//                  the eAPI connection.   The default value is 'admin'
//  password (string): The password to pass to the device to authenticate
//                  the eAPI connection.  The default value is ''
//  port (int): The TCP port of the endpoint for the eAPI connection.  If
//              this keyword is not specified, the default value is
//              automatically determined by the transport type.
//              (http=80, https=443)
// Returns:
//  An instance of Node object for the specified transport.
func Connect(transport string, host string, username string, passwd string,
	port int) (*Node, error) {
	conn, err := Connection(transport, host, username, passwd, port)
	if err != nil {
		return nil, err
	}
	return &Node{conn: conn, autoRefresh: true}, nil
}

// Connection creates a connection using the supplied settings
//
// This function will create a connection to an Arista EOS node using
// the arguments.  All arguments are optional with default values.
//
// Args:
//  transport (string): Specifies the type of connection transport to use.
//                   Valid values for the connection are socket, http_local,
//                   http, and https.  https is the default.
//  host (string): The IP addres or DNS host name of the connection device.
//              The default value is 'localhost'
//  username (string): The username to pass to the device to authenticate
//                  the eAPI connection.   The default value is 'admin'
//  password (string): The password to pass to the device to authenticate
//                  the eAPI connection.  The default value is ''
//  port (int): The TCP port of the endpoint for the eAPI connection.  If
//              this keyword is not specified, the default value is
//              automatically determined by the transport type.
//              (http=80, https=443)
// Returns:
//  An instance of EapiConnectionEntity object for the specified transport.
func Connection(transport string, host string, username string, passwd string,
	port int) (EapiConnectionEntity, error) {
	if transport == "" {
		transport = "https"
	}
	if host == "" {
		host = "localhost"
	}
	if username == "" {
		username = "admin"
	}

	var transFunc fn
	var found bool
	if transFunc, found = transports[transport]; !found {
		return nil, fmt.Errorf("Invalid transport specified: %s", transport)
	}

	obj := transFunc(transport, host, username, passwd, port)
	return obj, nil
}

// expandPath expands out the '~' if specified within the path
//
// Args:
//  path (string): path
//
// Returns:
//  String with newly expanded path or "" with error
func expandPath(path string) (string, error) {
	var homeDir string

	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	// Get current User home dir
	if usr, err := user.Current(); err == nil {
		homeDir = usr.HomeDir
	} else {
		// user.Current() requires cgo and currently this is disabled
		// during cross-compiling. If an error is returned, then attempt
		// to get the home dir from the environment.
		homeDir = os.Getenv("HOME")
	}
	// return nothing if no homeDir found.
	if homeDir == "" {
		return "", fmt.Errorf("No home dir found")
	}
	if len(path) < 2 {
		return homeDir, nil
	}
	return filepath.Join(homeDir, path[1:]), nil
}
