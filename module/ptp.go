package module

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/aristanetworks/goeapi"
)

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

// PtpConfig represents the parsed ptp config
type PtpConfig struct {
	sourceIP string
	mode     string
	ttl      string
}

// SourceIP returns source ip address of ptp packets sent by this node
func (p *PtpConfig) SourceIP() string {
	return p.sourceIP
}

// Mode returns ptp mode, either dynamic or master, which this node operates
func (p *PtpConfig) Mode() string {
	return p.mode
}

// TTL returns time to leave number sent in ipv4 header
func (p *PtpConfig) TTL() string {
	return p.ttl
}

// PTPEntity provides a configuration resource for ptp
type PTPEntity struct {
	interfaces *PTPInterfaceEntity
	*AbstractBaseEntity
}

// Interfaces returns the STPInterfaces instance
func (p *PTPEntity) Interfaces() *PTPInterfaceEntity {
	if p.interfaces == nil {
		p.interfaces = PTPInterfaces(p.AbstractBaseEntity.node)
	}
	return p.interfaces
}

// Ptp factory function to initiallize PTPEntity resource
// given a Node
func Ptp(node *goeapi.Node) *PTPEntity {
	return &PTPEntity{AbstractBaseEntity: &AbstractBaseEntity{node}}
}

// Get the PTP Config for the current entity.
// Returns a Ptponfig object
func (p *PTPEntity) Get() *PtpConfig {
	config := p.Config()

	var ptp = new(PtpConfig)
	ptp.sourceIP = p.parse(config, "source ip")
	ptp.mode = p.parse(config, "mode")
	ptp.ttl = p.parse(config, "ttl")

	return ptp
}

// ConfigurePtp configures the PTP Entity with the given
// command. Returns true (bool) if the commands complete
// successfully
func (p *PTPEntity) ConfigurePtp(cmd string) bool {
	config := p.Get()
	if config == nil {
		return false
	}
	commands := []string{
		cmd,
	}
	return p.Configure(commands...)
}

// SetSourceIP configures the source ip using the provided value.
// Returns true(bool) if the commands complete successfully
func (p *PTPEntity) SetSourceIP(value string) bool {
	if value == "" {
		return p.ConfigurePtp("no ptp source ip")
	}
	return p.ConfigurePtp("ptp source ip " + value)
}

// SetMode configures the ptp mode using the provided value.
// Returns true(bool) if the commands complete successfully
func (p *PTPEntity) SetMode(value string) bool {
	if value == "" {
		return p.ConfigurePtp("no ptp mode")
	}
	return p.ConfigurePtp("ptp mode " + value)
}

// SetTTL configures the ptp ttl using the provided value.
// Returns true(bool) if the commands complete successfully
func (p *PTPEntity) SetTTL(value string) bool {
	if value == "" {
		return p.ConfigurePtp("no ptp ttl")
	}
	return p.ConfigurePtp("ptp ttl " + value)
}

// parse parses the given PTP config for the give pattern value
func (p *PTPEntity) parse(config string, pattern string) string {
	str := fmt.Sprintf(`(?m)ptp %s ([^\s]+)`, pattern)
	re := regexp.MustCompile(str)
	match := re.FindStringSubmatch(config)
	if match == nil {
		return ""
	}
	return match[1]
}

// PTPInterfaceConfig represents the parsed PTP interface config
// {
//		"enabled"     : "true",
//		"role"        : "dynamic",
//		"transport"   : "layer2",
// }
type PTPInterfaceConfig map[string]string

// PTPInterfaceCollection is a collection of PTPInterfaceConfigs
// mapped by interface name:
// {
//		"Ethernet49/1" : PTPInterfaceConfig {
//							"enabled"   : "false",
//							"role"      : "dynamic",
//							"transport" : "layer2",
//						  },
//		   "Ethernet1" : ...
// }
type PTPInterfaceCollection map[string]PTPInterfaceConfig

// PTPInterfaceEntity provides a configuration resource for PTP
type PTPInterfaceEntity struct {
	*AbstractBaseEntity
}

// PTPInterfaces factory function to initiallize PTPInterfaceEntity resource
// given a Node
func PTPInterfaces(node *goeapi.Node) *PTPInterfaceEntity {
	return &PTPInterfaceEntity{&AbstractBaseEntity{node}}
}

// Get returns an PTPInterfaceConfig type for a given interface name(string).
func (p *PTPInterfaceEntity) Get(name string) PTPInterfaceConfig {
	parent := `interface\s+` + name
	config, _ := p.GetBlock(parent)

	return PTPInterfaceConfig{
		"enabled":   strconv.FormatBool(p.parsePTPEnabled(config)),
		"role":      p.parsePTPRole(config),
		"transport": p.parsePTPTransport(config),
	}
}

// GetPTPAdminStatus returns "true" if ptp enabled on this interfaces
func (p PTPInterfaceConfig) GetPTPAdminStatus() string {
	return p["enabled"]
}

// GetAll returns a collection of PTPInterfaceConfigs key'd by
// interface name.
func (p *PTPInterfaceEntity) GetAll() PTPInterfaceCollection {
	config := p.Config()

	re := regexp.MustCompile(`(?m)^interface\s(Eth.+|Po.+)$`)
	interfaces := re.FindAllStringSubmatch(config, -1)

	collection := make(PTPInterfaceCollection)

	for _, matchedLine := range interfaces {
		intf := matchedLine[1]
		if tmp := p.Get(intf); tmp != nil {
			collection[intf] = tmp
		}
	}
	return collection
}

// parsePTPEnabled parses the provided interface config returning
// true(bool) if ptp is enabled
func (p *PTPInterfaceEntity) parsePTPEnabled(config string) bool {
	matched, _ := regexp.MatchString(`(?m)^\s+ptp enable`, config)
	return matched
}

// parsePTPRole parses the provided interface config returning the
// ptp role (dynamic, master)
func (p *PTPInterfaceEntity) parsePTPRole(config string) string {
	found, _ := regexp.MatchString("ptp mode master", config)
	if found {
		return "master"
	}
	return "dynamic"
}

// parsePTPTransport parses the provided interface config returning the
// ptp transport (ipv4, layer2)
func (p *PTPInterfaceEntity) parsePTPTransport(config string) string {
	found, _ := regexp.MatchString("ptp transport layer2", config)
	if found {
		return "layer2"
	}
	return "ipv4"
}

// SetEnable enables(true) or disables(false) ptp for the interface
// name(string). Returns true(bool) if configuration successful
func (p *PTPInterfaceEntity) SetEnable(name string, enable bool) bool {
	str := "no ptp enable"
	if enable {
		str = "ptp enable"
	}
	cmd := p.CommandBuilder(str, "", false, true)
	return p.ConfigureInterface(name, cmd)
}

// ShowPTP represents "show ptp" output
type ShowPTP struct {
	PtpMode          string             `json:"ptpMode"`
	PtpClockSummary  PtpClockSummary    `json:"ptpClockSummary"`
	PtpIntfSummaries map[string]PtpIntf `json:"ptpIntfSummaries"`
}

// PtpIntf represents inidividual interface in "show ptp" output
type PtpIntf struct {
	PortState      string `json:"portState"`
	DelayMechanism string `json:"delayMechanism"`
	TransportMode  string `json:"transportMode"`
}

// PtpClockSummary represents common data in "show ptp" output
type PtpClockSummary struct {
	ClockIdentity        string  `json:"clockIdentity"`
	MeanPathDelay        int     `json:"meanPathDelay"`
	StepsRemoved         int     `json:"stepsRemoved"`
	Skew                 float64 `json:"skew"`
	GmClockIdentity      string  `json:"gmClockIdentity"`
	SlavePort            string  `json:"slavePort"`
	NumberOfMasterPorts  int     `json:"numberOfMasterPorts"`
	NumberOfSlavePorts   int     `json:"numberOfSlavePorts"`
	CurrentPtpSystemTime int     `json:"currentPtpSystemTime"`
	OffsetFromMaster     int     `json:"offsetFromMaster"`
	LastSyncTime         int     `json:"lastSyncTime"`
}

func (b *ShowPTP) GetCmd() string {
	return "show ptp"
}

func (s *ShowEntity) ShowPTP() (ShowPTP, error) {
	var showptp ShowPTP
	handle, err := s.node.GetHandle("json")
	if err != nil {
		return showptp, err
	}
	err = handle.AddCommand(&showptp)
	if err != nil {
		return showptp, err
	}
	err = handle.Call()
	if err != nil {
		return showptp, err
	}
	handle.Close()
	return showptp, nil
}
