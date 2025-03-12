package module

import (
	"reflect"
	"testing"

	"github.com/aristanetworks/goeapi"
)

func TestShowQueueMonitor_UnitTest(t *testing.T) {
	var dummyNode *goeapi.Node
	var dummyConnection *DummyConnection

	dummyConnection = &DummyConnection{}

	dummyNode = &goeapi.Node{}
	dummyNode.SetConnection(dummyConnection)

	show := Show(dummyNode)
	showqueue, _ := show.ShowQueueMonitor("Et24")

	type ShowQueueMonitor struct {
		Cmd                 string
		ReportTime          float64  `json:"report_time"`
		Warnings            string   `json:"warnings"`
		BytesPerTxmpSegment uint     `json:"bytes_per_txmp_segment"`
		GlobalHitCount      uint     `json:"global_hit_count"`
		LanzEnabled         bool     `json:"lanz_enabled"`
		PlatformName        string   `json:"platform_name"`
		EntryList           []*Entry `json:"entry_list"`
	}

	type Entry struct {
		EntryTimeUsecs              int64    `json:"entry_time_usecs"`
		GlobalProtectionModeEnabled bool     `json:"global_protection_mode_enabled"`
		EntryTime                   float64  `json:"entry_time"`
		Interface                   string   `json:"interface"`
		Duration                    uint     `json:"duration"`
		DurationUsecs               uint32   `json:"duration_usecs"`
		EntryType                   string   `json:"entry_type"`
		QueueLength                 uint32   `json:"queue_length"`
		TrafficClass                uint     `json:"traffic_class"`
		IngressPortSet              []string `json:"ingress_port_set"`
	}

	var scenarios = []struct {
		GlobalHitCount uint
		LanzEnabled    bool
		PlatformName   string
		EntryList      []*Entry
	}{
		{
			GlobalHitCount: 0,
			LanzEnabled:    true,
			PlatformName:   "Sand",
			EntryList: []*Entry{
				{
					EntryTimeUsecs:              1568892560445182,
					GlobalProtectionModeEnabled: true,
					EntryTime:                   1568892560.445182,
					Interface:                   "Ethernet24",
					EntryType:                   "P",
					QueueLength:                 73520,
					IngressPortSet: []string{
						"Ethernet4/1",
						"Ethernet8/1",
						"Ethernet26/1",
						"Ethernet30/1",
						"Ethernet10/1",
						"Ethernet25/1",
						"Ethernet5/1",
						"Ethernet9/1",
						"Ethernet3/1",
					},
				},
				{
					EntryTimeUsecs:              1568892557441976,
					GlobalProtectionModeEnabled: true,
					EntryTime:                   1568892557.441976,
					Interface:                   "Ethernet24/1",
					EntryType:                   "P",
					QueueLength:                 489744,
					IngressPortSet: []string{
						"Ethernet24/3",
						"Ethernet13/1",
						"Ethernet35/1",
						"Ethernet23/1",
						"Ethernet33/1",
						"Ethernet19/1",
						"Ethernet14/1",
						"Ethernet18/1",
						"Ethernet24/1",
						"Ethernet34/1",
						"Ethernet24/4",
						"Ethernet24/2",
					},
				},
				{
					EntryTimeUsecs:              1568892549432686,
					GlobalProtectionModeEnabled: true,
					EntryTime:                   1568892549.432686,
					Interface:                   "Ethernet24/1",
					EntryType:                   "P",
					QueueLength:                 46384,
					IngressPortSet: []string{
						"Ethernet4/1",
						"Ethernet8/1",
						"Ethernet26/1",
						"Ethernet30/1",
						"Ethernet10/1",
						"Ethernet25/1",
						"Ethernet5/1",
						"Ethernet9/1",
						"Ethernet3/1",
					},
				},
				{
					EntryTimeUsecs:              1568892541423988,
					GlobalProtectionModeEnabled: true,
					EntryTime:                   1568892541.423988,
					Interface:                   "Ethernet24/1",
					EntryType:                   "P",
					QueueLength:                 247568,
					IngressPortSet: []string{
						"Ethernet24/3",
						"Ethernet13/1",
						"Ethernet35/1",
						"Ethernet23/1",
						"Ethernet33/1",
						"Ethernet19/1",
						"Ethernet14/1",
						"Ethernet18/1",
						"Ethernet24/1",
						"Ethernet34/1",
						"Ethernet24/4",
						"Ethernet24/2",
					},
				},
			},
		},
	}

	entries := showqueue.EntryList

	for _, tt := range scenarios {
		if tt.LanzEnabled != showqueue.LanzEnabled {
			t.Errorf("Lanz is disabled, which does not match expected %t, got %t", tt.LanzEnabled,
				showqueue.LanzEnabled)
		}
		if tt.PlatformName != showqueue.PlatformName {
			t.Errorf("Platform name doesn't match expected %s, got %s", tt.PlatformName, showqueue.PlatformName)
		}
		for idx, entry := range tt.EntryList {
			if entry.EntryTime != entries[idx].EntryTime {
				t.Errorf("Entry time does not match expected %f, got %f",
					entry.EntryTime, entries[idx].EntryTime)
			}
			if entry.EntryType != entries[idx].EntryType {
				t.Errorf("Entry type does not match expected %s, got %s",
					entry.EntryType, entries[idx].EntryType)
			}
			if entry.QueueLength != entries[idx].QueueLength {
				t.Errorf("Queue Length does not match expected %d, got %d",
					entry.QueueLength, entries[idx].QueueLength)
			}
			if !reflect.DeepEqual(entry.IngressPortSet, entries[idx].IngressPortSet) {
				t.Errorf("Ingres Port set does not match expected %v, got %v",
					entry.IngressPortSet, entries[idx].IngressPortSet)
			}
		}
	}

	showqueue, _ = show.ShowQueueMonitorWithLimit("Et24", "samples", 2)
	entries = showqueue.EntryList

	for _, tt := range scenarios {
		if tt.LanzEnabled != showqueue.LanzEnabled {
			t.Errorf("Lanz is disabled, which does not match expected %t, got %t", tt.LanzEnabled,
				showqueue.LanzEnabled)
		}
		if tt.PlatformName != showqueue.PlatformName {
			t.Errorf("Platform name doesn't match expected %s, got %s", tt.PlatformName, showqueue.PlatformName)
		}
		for idx, entry := range tt.EntryList {
			if entry.EntryTime != entries[idx].EntryTime {
				t.Errorf("Entry time does not match expected %f, got %f",
					entry.EntryTime, entries[idx].EntryTime)
			}
			if entry.EntryType != entries[idx].EntryType {
				t.Errorf("Entry type does not match expected %s, got %s",
					entry.EntryType, entries[idx].EntryType)
			}
			if entry.QueueLength != entries[idx].QueueLength {
				t.Errorf("Queue Length does not match expected %d, got %d",
					entry.QueueLength, entries[idx].QueueLength)
			}
			if !reflect.DeepEqual(entry.IngressPortSet, entries[idx].IngressPortSet) {
				t.Errorf("Ingres Port set does not match expected %v, got %v",
					entry.IngressPortSet, entries[idx].IngressPortSet)
			}
		}
	}
}
