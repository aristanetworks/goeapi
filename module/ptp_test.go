package module

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/aristanetworks/goeapi"
)

/**
 *****************************************************************************
 * Unit Tests
 *****************************************************************************
 **/

func TestPTPGet_UnitTest(t *testing.T) {
	ptp := Ptp(dummyNode)
	config := ptp.Get()
	if config.SourceIP() != "1.1.1.1" || config.TTL() != "30" ||
		config.Mode() != "boundary" {
		t.Fatalf("Invalid result from Get(): %#v", config)
	}
}

func TestPTPInterfaces_UnitTest(t *testing.T) {
	ptp := Ptp(&goeapi.Node{})
	i := ptp.Interfaces()
	if i == nil {
		t.Fatalf("No PTPInterfaces")
	}
}

func TestPTPSetMode_UnitTest(t *testing.T) {
	ptp := Ptp(dummyNode)

	cmds := []string{
		"default ptp mode",
	}

	tests := [...]struct {
		mode string
		want string
		rc   bool
	}{
		{"boundary", "ptp mode boundary", true},
		{"e2etransparent", "ptp mode e2etransparent", true},
		{"p2ptransparent", "ptp mode p2ptransparent", true},
		{"gptp", "ptp mode gptp", true},
		{"", "no ptp mode", true},
		{"disabled", "ptp mode disabled", true},
	}

	for _, tt := range tests {
		if got := ptp.SetMode(tt.mode); got != tt.rc {
			t.Fatalf("SetMode(%s) = %t; want %t", tt.mode, got, tt.rc)
		}
		if tt.rc {
			cmds[0] = tt.want
			// first two commands are 'enable', 'configure terminal'
			commands := dummyConnection.GetCommands()[2:]
			for idx, val := range commands {
				if cmds[idx] != val {
					t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
				}
			}
		}
	}
}

func TestPTPSetSourceIp_UnitTest(t *testing.T) {
	ptp := Ptp(dummyNode)

	cmds := []string{
		"default ptp source ip",
	}

	tests := [...]struct {
		mode string
		want string
		rc   bool
	}{
		{"1.1.1.1", "ptp source ip 1.1.1.1", true},
		{"", "no ptp source ip", true},
	}

	for _, tt := range tests {
		if got := ptp.SetSourceIP(tt.mode); got != tt.rc {
			t.Fatalf("SetSourceIP(%s) = %t; want %t", tt.mode, got, tt.rc)
		}
		if tt.rc {
			cmds[0] = tt.want
			// first two commands are 'enable', 'configure terminal'
			commands := dummyConnection.GetCommands()[2:]
			for idx, val := range commands {
				if cmds[idx] != val {
					t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
				}
			}
		}
	}
}

func TestPTPSetTTL_UnitTest(t *testing.T) {
	ptp := Ptp(dummyNode)

	cmds := []string{
		"default ptp ttl",
	}

	tests := [...]struct {
		mode string
		want string
		rc   bool
	}{
		{"30", "ptp ttl 30", true},
		{"300", "ptp ttl 300", true}, // must be false, invalid command for EOS, range is 1 to 2546
		{"", "no ptp ttl", true},
	}

	for _, tt := range tests {
		if got := ptp.SetTTL(tt.mode); got != tt.rc {
			t.Fatalf("SetTTL(%s) = %t; want %t", tt.mode, got, tt.rc)
		}
		if tt.rc {
			cmds[0] = tt.want
			// first two commands are 'enable', 'configure terminal'
			commands := dummyConnection.GetCommands()[2:]
			for idx, val := range commands {
				if cmds[idx] != val {
					t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
				}
			}
		}
	}
}

func TestParse_UnitTest(t *testing.T) {
	var p PTPEntity
	shortConfig := `
  no ip policy unresolved-nexthop action drop
  no ip policy mac address aging
  no ip policy match protocol bgp
  !
  power PowerSupply1 input voltage warning low 0.00 readings 10
  power PowerSupply2 input voltage warning low 0.00 readings 10
  !
  power poll-interval 5
  !
  ptp priority1 128
  ptp clock-identity 00:00:00:00:00:00:00:00
  ptp priority2 128
  ptp domain 0
  ptp source ip 1.1.1.1
  ptp mode boundary
  ptp ttl 30
  ptp message-type general dscp 0 default
  ptp message-type event dscp 21 default
  ptp hold-ptp-time 28800
  no ptp forward-v1
  no ptp forward-unicast
  ptp monitor
  no ptp monitor threshold offset-from-master
  no ptp monitor threshold mean-path-delay
  no ptp monitor threshold skew
  !
  ptp hardware-sync interval 1000
  !
  no radius-server key
  radius-server timeout 5
  radius-server retransmit 3
  no radius-server deadtime
  no radius-server attribute 32 include-in-access-req format
  radius-server qos dscp 0
`
	tests := [...]struct {
		in   string
		want string
	}{
		{"mode", "boundary"},
		{"source ip", "1.1.1.1"},
		{"ttl", "30"},
	}

	for _, tt := range tests {
		testConfig := fmt.Sprintf(shortConfig)
		if got := p.parse(testConfig, tt.in); got != tt.want {
			t.Fatalf("parse(config, %s) = %v; want %v", tt.in, got, tt.want)
		}
	}
}

func TestPTPIntfGet_UnitTest(t *testing.T) {
	ptp := PTPInterfaces(dummyNode)
	config := ptp.Get("Ethernet1")
	for _, val := range []string{"enabled", "role", "transport"} {
		if _, found := config[val]; !found {
			t.Fatalf("Get() missing key %s", val)
		}
	}
}

func TestPTPIntfGetAll_UnitTest(t *testing.T) {
	ptp := PTPInterfaces(dummyNode)
	interfaces := ptp.GetAll()
	for _, val := range []string{"enabled", "role", "transport"} {
		if _, found := interfaces["Ethernet1"][val]; !found {
			t.Fatalf("Get() missing key %s", val)
		}
	}
}

func TestPTPIntfGetPTPAdminStatus_UnitTest(t *testing.T) {
	ptp := PTPInterfaces(dummyNode)
	tests := [...]struct {
		in   string
		want bool
	}{
		{"Ethernet1", true},
		{"Ethernet2", false},
		{"Ethernet3", false},
	}
	for _, tt := range tests {
		intf := ptp.Get(tt.in)
		if got, _ := strconv.ParseBool(intf.GetPTPAdminStatus()); got != tt.want {
			t.Fatalf("%s.GetPTPAdminStatus() = %v; want %v", tt.in, got, tt.want)
		}
	}
}

func TestPTPIntfSetEnable_UnitTest(t *testing.T) {
	ptp := PTPInterfaces(dummyNode)
	tests := []struct {
		in    string
		value bool
		want  string
		rc    bool
	}{
		{"Ethernet2", true, "ptp enable", true},
		{"Ethernet3", false, "no ptp enable", true},
	}

	for _, tt := range tests {
		cmds := []string{
			"interface " + tt.in,
			"default ptp enable",
		}

		if ok := ptp.SetEnable(tt.in, tt.value); ok != tt.rc {
			t.Fatalf("Expected status \"%t\" got \"%t\"", tt.rc, ok)
		}
		cmds[1] = tt.want
		// first two commands are 'enable', 'configure terminal'
		commands := dummyConnection.GetCommands()[2:]
		for idx, val := range commands {
			if cmds[idx] != val {
				t.Fatalf("Expected \"%q\" got \"%q\"", cmds, commands)
			}
		}
	}
}

func TestShowPTP_UnitTest(t *testing.T) {
	var dummyNode *goeapi.Node
	var dummyConnection *DummyConnection

	dummyConnection = &DummyConnection{}

	dummyNode = &goeapi.Node{}
	dummyNode.SetConnection(dummyConnection)

	show := Show(dummyNode)
	showPTP, _ := show.ShowPTP()

	type ClockSummary struct {
		ClockIdentity        string
		GmClockIdentity      string
		SlavePort            string
		NumberOfMasterPorts  int
		NumberOfSlavePorts   int
		CurrentPtpSystemTime int
		OffsetFromMaster     int
	}

	type IntfSummary struct {
		InterfaceName  string
		PortState      string
		DelayMechanism string
		TransportMode  string
	}

	var scenarios = []struct {
		PtpMode          string
		PtpClockSummary  ClockSummary
		PtpIntfSummaries []IntfSummary
	}{
		{
			PtpMode: "ptpBoundaryClock",
			PtpClockSummary: ClockSummary{
				ClockIdentity:        "0x01:01:01:ff:ff:01:a0:01",
				GmClockIdentity:      "0x01:01:01:ff:fe:01:a0:02",
				SlavePort:            "Ethernet49/1",
				NumberOfMasterPorts:  1,
				NumberOfSlavePorts:   1,
				CurrentPtpSystemTime: 1557315926,
				OffsetFromMaster:     -10,
			},
			PtpIntfSummaries: []IntfSummary{
				{
					InterfaceName:  "Ethernet49/1",
					PortState:      "psSlave",
					DelayMechanism: "e2e",
					TransportMode:  "ipv4",
				},
				{
					InterfaceName:  "Ethernet53/1",
					PortState:      "psMaster",
					DelayMechanism: "e2e",
					TransportMode:  "ipv4",
				},
			},
		},
	}

	interfaces := showPTP.PtpIntfSummaries

	for _, tt := range scenarios {
		if tt.PtpMode != showPTP.PtpMode {
			t.Errorf("Ptp mode does not match expected %s, got %s", tt.PtpMode,
				showPTP.PtpMode)
		}
		if tt.PtpClockSummary.ClockIdentity != showPTP.PtpClockSummary.ClockIdentity {
			t.Errorf("Ptp local clock identity does not match expected %s, got %s",
				tt.PtpClockSummary.ClockIdentity,
				showPTP.PtpClockSummary.ClockIdentity)
		}
		if tt.PtpClockSummary.GmClockIdentity != showPTP.PtpClockSummary.GmClockIdentity {
			t.Errorf("Ptp grand master clock identity does not match expected %s, got %s",
				tt.PtpClockSummary.GmClockIdentity,
				showPTP.PtpClockSummary.GmClockIdentity)
		}
		if tt.PtpClockSummary.SlavePort != showPTP.PtpClockSummary.SlavePort {
			t.Errorf("Ptp slave port not match expected %s, got %s",
				tt.PtpClockSummary.SlavePort, showPTP.PtpClockSummary.SlavePort)
		}
		if tt.PtpClockSummary.OffsetFromMaster != showPTP.PtpClockSummary.OffsetFromMaster {
			t.Errorf("Ptp offset from master not match expected %d, got %d",
				tt.PtpClockSummary.OffsetFromMaster,
				showPTP.PtpClockSummary.OffsetFromMaster)
		}
		if tt.PtpClockSummary.CurrentPtpSystemTime != showPTP.PtpClockSummary.CurrentPtpSystemTime {
			t.Errorf("Ptp time not match expected %d, got %d",
				tt.PtpClockSummary.CurrentPtpSystemTime,
				showPTP.PtpClockSummary.CurrentPtpSystemTime)
		}
		if tt.PtpClockSummary.NumberOfMasterPorts != showPTP.PtpClockSummary.NumberOfMasterPorts {
			t.Errorf("Number of master ports does not match expected %d, got %d",
				tt.PtpClockSummary.NumberOfMasterPorts,
				showPTP.PtpClockSummary.NumberOfMasterPorts)
		}
		if tt.PtpClockSummary.NumberOfSlavePorts != showPTP.PtpClockSummary.NumberOfSlavePorts {
			t.Errorf("Number of slave ports does not match expected %d, got %d",
				tt.PtpClockSummary.NumberOfSlavePorts,
				showPTP.PtpClockSummary.NumberOfSlavePorts)
		}
		for _, intf := range tt.PtpIntfSummaries {
			if intf.PortState != interfaces[intf.InterfaceName].PortState {
				t.Errorf("State of the port %s does not match expected %s, got %s",
					intf.InterfaceName, intf.PortState,
					interfaces[intf.InterfaceName].PortState)
			}
			if intf.DelayMechanism != interfaces[intf.InterfaceName].DelayMechanism {
				t.Errorf("Delay mechanism of the port %s does not match expected %s, got %s",
					intf.InterfaceName, intf.DelayMechanism,
					interfaces[intf.InterfaceName].DelayMechanism)
			}
			if intf.TransportMode != interfaces[intf.InterfaceName].TransportMode {
				t.Errorf("Transport mode of the port %s does not match expected %s, got %s",
					intf.InterfaceName, intf.TransportMode,
					interfaces[intf.InterfaceName].TransportMode)
			}
		}
	}
}
