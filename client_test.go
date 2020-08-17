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
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestConfigExpandPath_UnitTest(t *testing.T) {
	usr, _ := user.Current()
	homedir := usr.HomeDir
	tests := [...]struct {
		in   string
		want string
	}{
		{"/home/usera/.eapi.conf", "/home/usera/.eapi.conf"},
		{"~/.eapi.conf", homedir + "/.eapi.conf"},
		{"./.eapi.conf", "./.eapi.conf"},
		{"", ""},
		{"a", "a"},
		{"aa", "aa"},
		{"~/", homedir},
		{"~", homedir},
	}
	for _, tt := range tests {
		got, _ := expandPath(tt.in)
		if got != tt.want {
			t.Fatalf("expandPath() got %q;  want %q", got, tt.want)
		}
	}
}

func TestConfigNilEapiConfig_UnitTest(t *testing.T) {
	config := NewEapiConfig()
	config = nil
	if ret := config.Connections(); ret != nil {
		t.Fatalf("got %#v expected nil", config)
	}
}

func TestConfigEapiConfigNilFile_UnitTest(t *testing.T) {
	currEnv := os.Getenv("EAPI_CONF")
	os.Setenv("EAPI_CONF", GetFixture("dontexist.conf"))
	config := NewEapiConfig()
	if len(config.Connections()) != 1 {
		t.Fatalf("got %d expected 1", len(config.Connections()))
	}
	if currEnv != "" {
		os.Setenv("EAPI_CONF", currEnv)
	} else {
		os.Unsetenv("EAPI_CONF")
	}
}

func TestConfigInitInvalid_UnitTest(t *testing.T) {
	config := NewEapiConfigFile(GetFixture("invalid.conf"))
	if len(config.Connections()) != 1 {
		t.Fatalf("got %d expected 1", len(config.Connections()))
	}
}

func TestConfigReadInvalid_UnitTest(t *testing.T) {
	config := NewEapiConfig()
	if err := config.Read(GetFixture("invalid.conf")); err == nil {
		t.Fatalf("expected failure")
	}
}

func TestConfigReadValid_UnitTest(t *testing.T) {
	config := NewEapiConfig()
	if err := config.Read(GetFixture("eapi.conf")); err != nil {
		t.Fatalf("expected failure")
	}
}

func TestConfigLoadConfigFilename_UnitTest(t *testing.T) {
	currEnv := os.Getenv("EAPI_CONF")
	if currEnv != "" {
		fmt.Printf("UNSETTING EAPI_CONF\n")
		os.Unsetenv("EAPI_CONF")
	}
	if os.Getenv("EAPI_CONF") != "" {
		t.Fatalf("Unsetenv failed")
	}
	LoadConfig(GetFixture("eapi.conf"))
	section := ConfigFor("test1")
	if section["host"] != "192.168.1.16" ||
		section["username"] != "eapi" ||
		section["password"] != "password" {
		t.Fatalf("ConfigFor failed: %q", section)
	}
	if currEnv != "" {
		os.Setenv("EAPI_CONF", currEnv)
	} else {
		os.Unsetenv("EAPI_CONF")
	}
}

func TestConfigLoadConfigWithEnv_UnitTest(t *testing.T) {
	currEnv := os.Getenv("EAPI_CONF")

	os.Setenv("EAPI_CONF", GetFixture("eapi.conf"))
	LoadConfig(RandomString(4, 9))
	section := ConfigFor("test1")
	if section["host"] != "192.168.1.16" ||
		section["username"] != "eapi" ||
		section["password"] != "password" {
		t.Fatalf("ConfigFor failed: %q", section)
	}
	if currEnv != "" {
		os.Setenv("EAPI_CONF", currEnv)
	} else {
		os.Unsetenv("EAPI_CONF")
	}
}

func TestConfigLoadConfigCheckSections_UnitTest(t *testing.T) {
	LoadConfig(GetFixture("eapi.conf"))
	if len(configGlobal.File) != 3 {
		t.Fatalf("Incorrect number of sections found: %d", len(configGlobal.File))
	}
}

func TestConfigLoadConfigDefaultConnection_UnitTest(t *testing.T) {
	LoadConfig(GetFixture("invalid.conf"))
	if len(configGlobal.File) != 1 {
		t.Fatalf("Incorrect number of sections found: %d", len(configGlobal.File))
	}
}

func TestConfigLoadConfigProperty_UnitTest(t *testing.T) {
	LoadConfig(GetFixture("eapi.conf"))

	validConnections := []string{"test1", "test2", "localhost"}
	sort.Strings(validConnections)

	connections := Connections()
	sort.Strings(connections)

	if len(connections) == len(validConnections) {
		for idx, val := range connections {
			if validConnections[idx] != val {
				t.Fatalf("Got %s expected %s", val, validConnections[idx])
			}
		}
		return
	}
	t.Fatalf("Incorrect length: got %d expected: %d", len(connections), len(validConnections))
}

func TestConfigLoadConfigReplaceHostWithName_UnitTest(t *testing.T) {
	LoadConfig(GetFixture("nohost.conf"))
	section := ConfigFor("test")
	if section["host"] != "test" {
		t.Fatalf("Got %s expected: test", section["host"])
	}
}

// HACK. Reload dut.conf after the above test.
// Don't account for localhost in config since we don't support yet.
func TestConfigLoadRest_UnitTest(t *testing.T) {
	duts = nil
	LoadConfig(GetFixture("dut.conf"))
	conns := Connections()
	for _, name := range conns {
		if name != "localhost" {
			node, _ := ConnectTo(name)
			duts = append(duts, node)
		}
	}
}

func TestClientRunningConfig_UnitTest(t *testing.T) {
	dummyNode.Refresh()
	if config := dummyNode.RunningConfig(); config == "" {
		t.Fatalf("No config returned")
	}
	if config := dummyNode.RunningConfig(); config == "" {
		t.Fatalf("No config returned")
	}
}
func TestClientStartupConfig_UnitTest(t *testing.T) {
	dummyNode.Refresh()
	if config := dummyNode.StartupConfig(); config == "" {
		t.Fatalf("No config returned")
	}
	if config := dummyNode.StartupConfig(); config == "" {
		t.Fatalf("No config returned")
	}
}

func TestClientGetConfig_UnitTest(t *testing.T) {
	tests := [...]struct {
		config string
		params string
		rc     bool
	}{
		{"running-config", "all", true},
		{"startup-config", "", true},
		{"bogus-config", "", false},
	}
	for idx, tt := range tests {
		_, err := dummyNode.getConfigText(tt.config, tt.params)
		if (err == nil) != tt.rc {
			t.Fatalf("Test[%d] Expected %t in eval of (err == nil): err:%#v", idx, tt.rc, err)
		}
	}
}

func TestClientGetConfigConnectionError_UnitTest(t *testing.T) {
	conn := dummyNode.GetConnection().(*DummyEapiConnection)
	conn.setReturnError(true)
	ret, err := dummyNode.getConfigText("running-config", "")
	if ret != "" && err == nil {
		t.Fatalf("Connection error didn't raise issue")
	}
}

func TestClientGetSection_UnitTest(t *testing.T) {
	tests := [...]struct {
		reg    string
		config string
		rc     bool
	}{
		{`(?m)^interface Ethernet1$`, "running-config", true},
		{`(?m)^interface Ethernet1$`, "", true},
		{`(?m)^interface Ethernet1$`, "startup-config", true},
		{`(?=>)^interface Ethernet1$`, "running-config", false},
		{`(?m)^interface Ethernet1$`, "bogus-config", false},
		{`(?m)^interface TokenRing1$`, "running-config", false},
	}
	for idx, tt := range tests {
		_, err := dummyNode.GetSection(tt.reg, tt.config)
		if (err == nil) != tt.rc {
			t.Fatalf("Test[%d] Expected %t in eval of (err == nil): err:%#v", idx, tt.rc, err)
		}
	}
}

func TestClientGetSectionConnectionError_UnitTest(t *testing.T) {
	conn := dummyNode.GetConnection().(*DummyEapiConnection)
	conn.setReturnError(true)
	ret, err := dummyNode.GetSection(`(?m)^interface Ethernet1$`, "startup-config")
	if ret != "" && err == nil {
		t.Fatalf("Connection error didn't raise issue")
	}
}

func TestClientConfig_UnitTest(t *testing.T) {
	cmds := []string{
		"interface Ethernet1",
		"no flowcontrol send",
		"ip address 1.1.1.1/24",
		"no shut",
	}
	dummyNode.SetAutoRefresh(true)
	dummyNode.Config(cmds...)
	commands := dummyConnection.GetCommands()[2:]
	for idx, val := range commands {
		if cmds[idx] != val {
			t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
		}
	}
}

func TestClientEnable_UnitTest(t *testing.T) {
	tests := [...]struct {
		in []string
		rc bool
	}{
		{[]string{"configure terminal"}, false},
		{[]string{"configure    terminal"}, false},
		{[]string{"  configure"}, false},
		{[]string{"configure"}, false},
		{[]string{"show running-config"}, true},
	}
	dummyNode.EnableAuthentication("root")
	for idx, tt := range tests {
		if _, got := dummyNode.Enable(tt.in); (got == nil) != tt.rc {
			dummyNode.EnableAuthentication("")
			t.Fatalf("Test[%d] Expected %t in eval of (err == nil)", idx, tt.rc)
		}
	}
	dummyNode.EnableAuthentication("")
}

func TestClientEnableConnectionError_UnitTest(t *testing.T) {
	conn := dummyNode.GetConnection().(*DummyEapiConnection)
	conn.setReturnError(true)
	cmds := []string{"show running-config"}
	ret, err := dummyNode.Enable(cmds)
	if ret != nil && err == nil {
		t.Fatalf("Connection error didn't raise issue")
	}
}

func TestClientPrependEnableSequence_UnitTest(t *testing.T) {
	cmds := []string{
		"show version",
		"show arp",
		"show interfaces",
	}
	dummyNode.EnableAuthentication("root")
	got := dummyNode.prependEnableSequence(cmds)
	if len(got) != len(cmds)+1 {
		dummyNode.EnableAuthentication("")
		t.Fatalf("Missing or extra entry in prepend: %#v", got)
	}
	cmd := got[0].(map[string]interface{})["cmd"]
	passwd := got[0].(map[string]interface{})["input"]
	if cmd != "enable" || passwd != "root" {
		dummyNode.EnableAuthentication("")
		t.Fatalf("Invalid Entry. Cmd:%s Passwd:%s Got:%#v", cmd, passwd, got[0])
	}
	dummyNode.EnableAuthentication("")
}

func TestClientPrependEnableSequenceInvalidParams_UnitTest(t *testing.T) {
	cmds := []string{
		"show version",
		"show arp",
		"show interfaces",
	}
	got := dummyNode.prependEnableSequence(cmds)
	if len(got) != len(cmds)+1 {
		dummyNode.EnableAuthentication("")
		t.Fatalf("Missing or extra entry in prepend: %#v", got)
	}
	dummyNode.EnableAuthentication("")
}

func TestClientCmdsToInterface_UnitTest(t *testing.T) {
	cmds := []string{
		"show version",
		"show arp",
		"show interfaces",
	}
	got := cmdsToInterface(cmds)
	if len(got) != len(cmds) {
		t.Fatalf("Missing or extra entry in Conversion: %#v", got)
	}
	if got := cmdsToInterface(nil); got != nil {
		t.Fatalf("Should return nil on nil being passed")
	}
	if got := cmdsToInterface([]string{}); got != nil {
		t.Fatalf("Should return nil on empty list being passed")
	}
}

func TestClientEnableSingleResult_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"show version",
		}
		ret, _ := dut.runCommands(cmds, "json", false)
		if len(ret.Result) != 1 {
			t.Fatalf("sizes do not match Result[%d] != 1\n", len(ret.Result))
		}
	}
}
func TestClientEnableMultipleResult_SystemTest(t *testing.T) {
	for _, dut := range duts {
		var cmds []string
		for i := 0; i < RandomInt(10, 200); i++ {
			cmds = append(cmds, "show version")
		}
		ret, _ := dut.runCommands(cmds, "json", false)
		if len(ret.Result) != len(cmds) {
			t.Fatalf("sizes do not match Result[%d] != cmds[%d]\n", len(ret.Result), len(cmds))
		}
	}
}
func TestClientEnableMultiRequests_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"show version",
		}
		for i := 0; i < RandomInt(10, 200); i++ {
			ret, _ := dut.runCommands(cmds, "json", false)
			if len(ret.Result) != 1 {
				t.Fatalf("sizes do not match Result[%d] != 1\n", len(ret.Result))
			}
		}
	}
}
func TestClientConfigSingle_SystemTest(t *testing.T) {
	for _, dut := range duts {

		cmds := []string{
			"hostname " + RandomString(5, 50),
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("Config failure\n")
		}
		name := strings.Split(cmds[len(cmds)-1], " ")[1]
		ret, _ := dut.runCommands([]string{"show hostname"}, "json", false)
		if ret.Result[0]["hostname"] != name {
			t.Fatalf("Expecting %s got %s\n", name, ret.Result[0]["hostname"])
		}
	}
}

func TestClientConfigMultiple_SystemTest(t *testing.T) {
	for _, dut := range duts {
		var cmds []string
		for i := 0; i < RandomInt(10, 200); i++ {
			cmds = append(cmds, "hostname "+RandomString(5, 50))
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatal("Config failure\n")
		}
		// just check last
		name := strings.Split(cmds[len(cmds)-1], " ")[1]
		ret, _ := dut.runCommands([]string{"show hostname"}, "json", false)
		if ret.Result[0]["hostname"] != name {
			t.Fatalf("Expecting %s got %s\n", name, ret.Result[0]["hostname"])
		}
	}
}

func TestClientConnectInvalidTransport_UnitTest(t *testing.T) {
	if _, err := Connect("invalid", "hostname", "username", "passwd", 10); err == nil {
		t.Fatalf("No error seen for invalid transport type")
	}
}

func TestClientConnectTimeout_UnitTest(t *testing.T) {
	// 1.1.1.2 is assumed to be an unreachable bogus address
	node, err := Connect("http", "1.1.1.2", "admin", "admin", 80)
	if err != nil {
		t.Fatal("Error in connect.")
	}
	node.GetConnection().SetTimeout(5)

	if _, err = node.getConfigText("running-config", "all"); err == nil {
		t.Fatal("Should timeout and return error")
	}
}

func TestClientConnectDefaultPort_UnitTest(t *testing.T) {
	tests := [...]struct {
		proto string
		port  int
		descr string
	}{
		{"http", UseDefaultPortNum, "Http default port"},
		{"http", 80, "Http port 80"},
		{"https", UseDefaultPortNum, "Https default port"},
		{"https", 443, "Https port 443"},
		{"http_local", UseDefaultPortNum, "Http_local default port"},
		{"http_local", 8080, "Http_local port 8080"},
	}

	for idx, tt := range tests {
		// 1.1.1.2 is assumed to be an unreachable bogus address
		_, err := Connect(tt.proto, "1.1.1.2", "admin", "admin", tt.port)
		if err != nil {
			t.Fatalf("Test[%d] %s. Error in connect: err:%#v", idx, tt.descr, err)
		}
	}
}

func TestClientNodeEnablePasswd_UnitTest(t *testing.T) {
	node := &Node{}
	node.EnableAuthentication("root")
	if node.enablePasswd != "root" {
		t.Fatal("EnableAuthentication failed to set")
	}
	node.EnableAuthentication("")
	if node.enablePasswd != "" {
		t.Fatal("EnableAuthentication failed to set")
	}
}

func TestClientNodeAutoRefresh_UnitTest(t *testing.T) {
	for _, dut := range duts {
		dut.SetAutoRefresh(true)
		if dut.autoRefresh == false {
			t.Fatal("SetAutoRefresh(true) failed to set")
		}
		dut.SetAutoRefresh(false)
		if dut.autoRefresh == true {
			t.Fatal("SetAutoRefresh(false) failed to set")
		}
	}
}

func TestClientNodeGetRunningConfig_SystemTest(t *testing.T) {
	re := regexp.MustCompile(`^!\s+Command: show running-config`)
	for _, dut := range duts {
		dut.Refresh()
		config := dut.RunningConfig()
		if found := re.MatchString(config); !found {
			t.Fatal("Failed to obtain running-config")
		}

		config = dut.RunningConfig()
		if found := re.MatchString(config); !found {
			t.Fatal("Failed to obtain non-cached running-config")
		}
	}
}

func TestClientNodeGetStartupConfig_SystemTest(t *testing.T) {
	re := regexp.MustCompile(`^!\s+Command: show startup-config`)
	for _, dut := range duts {
		dut.Refresh()
		config := dut.StartupConfig()
		if found := re.MatchString(config); !found {
			t.Fatal("Failed to obtain startup-config")
		}

		config = dut.StartupConfig()
		if found := re.MatchString(config); !found {
			t.Fatal("Failed to obtain non-cached startup-config")
		}
	}
}

func TestClientNodeGetConfigInvalid_SystemTest(t *testing.T) {
	for _, dut := range duts {
		if _, err := dut.getConfigText("bogus-config", ""); err == nil {
			t.Fatal("Failed to return error on bogus-config")
		}
	}
}

func TestClientNodeGetSection_SystemTest(t *testing.T) {
	re := regexp.MustCompile(`^interface\s+Management1`)
	regStr := `(?m)^interface\s+Management1$`
	invalidRegexp := `(?=>)interface\s+Management1$`
	for _, dut := range duts {
		section, err := dut.GetSection(regStr, "bogus-config")
		if err == nil {
			t.Fatalf("GetSection should return error on bogus-config")
		}
		section, err = dut.GetSection(regStr, "")
		if found := re.MatchString(section); !found {
			t.Fatalf("Failed to obtain section from running-config err: %s", err)
		}
		section, _ = dut.GetSection(regStr, "running-config")
		if found := re.MatchString(section); !found {
			t.Fatalf("Failed to obtain section from running-config")
		}
		section, _ = dut.GetSection(regStr, "startup-config")
		if found := re.MatchString(section); !found {
			t.Fatalf("Failed to obtain section from startup-config")
		}
		_, err = dut.GetSection(invalidRegexp, "startup-config")
		if err == nil {
			t.Fatalf("Invalid regexp didn't fail")
		}
	}
}

func TestClientNodeEnableInvalidConfigCommands_SystemTest(t *testing.T) {
	tests := [...]struct {
		in []string
	}{
		{[]string{"configure terminal"}},
		{[]string{"configure    terminal"}},
		{[]string{"  configure"}},
		{[]string{"configure"}},
	}
	for _, dut := range duts {
		for _, tt := range tests {
			if _, got := dut.Enable(tt.in); got == nil {
				t.Fatalf("Should have failed %s", tt.in)
			}
		}
	}
}

func TestClientNodeEnableValid_SystemTest(t *testing.T) {
	re := regexp.MustCompile(`Internal build version`)
	cmds := []string{"show version"}
	for _, dut := range duts {
		crap, _ := dut.Enable(cmds)
		if found := re.MatchString(crap[0]["result"]); !found {
			t.Fatal("Failed to obtain build version")
		}
	}
}

func TestClientHandleEncoding_UnitTest(t *testing.T) {
	node := &Node{}

	if _, err := node.GetHandle("json"); err != nil {
		t.Fatal("GetHandle json")
	}
	if _, err := node.GetHandle("text"); err != nil {
		t.Fatal("GetHandle text")
	}
	if _, err := node.GetHandle("crap"); err == nil {
		t.Fatal("GetHandle crap")
	}
	if _, err := node.GetHandle("JsOn"); err != nil {
		t.Fatal("GetHandle JsOn")
	}
}

func TestClientHandleInvalid_UnitTest(t *testing.T) {
	var node *Node
	if _, err := node.GetHandle("json"); err == nil {
		t.Fatal("GetHandle invalid failed")
	}
}

func TestClientHandleClose_UnitTest(t *testing.T) {
	node := &Node{}
	handle, _ := node.GetHandle("json")
	handle.Close()

	if err := handle.Call(); err == nil {
		t.Fatal("No error for Call() after Close()")
	}
}

func TestClientNodeGetConnectionInvalid_UnitTest(t *testing.T) {
	var node *Node
	conn := node.GetConnection()
	if conn != nil {
		t.Fatal("Should not return valid")
	}
}

func TestClientNodeGetConnection_UnitTest(t *testing.T) {
	for _, dut := range duts {
		conn := dut.GetConnection()
		if conn == nil {
			t.Fatal("Failed to obtain connection")
		}
	}
}
