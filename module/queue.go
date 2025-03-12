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

package module

import "fmt"

// ShowQueueMonitor represents "show queue-monitor length" output
type ShowQueueMonitor struct {
	Cmd                 string
	ReportTime          float64  `json:"reportTime"`
	Warnings            string   `json:"warnings"`
	BytesPerTxmpSegment uint     `json:"bytesPerTxmpSegment"`
	GlobalHitCount      uint     `json:"globalHitCount"`
	LanzEnabled         bool     `json:"lanzEnabled"`
	PlatformName        string   `json:"platformName"`
	EntryList           []*Entry `json:"entryList"`
}

// Entry is queueing event entry from the ShowQueueMonitor output:
// Type       Time                    Intf(TC)           Queue         Duration      Ingress
//
//	Length                      Port-set
//	(bytes)       (usecs)
//
// ---------- ----------------------- --------------- ------------- ---------------- ------------------------------------------------
// P          0:00:03.83243 ago       Et24/1(2)          41904         1000000       Et3/1,4/1,5/1,8/1,9/1,10/1,25/1,26/1,30/1
type Entry struct {
	EntryTimeUsecs              int64    `json:"entryTimeUsecs"`
	GlobalProtectionModeEnabled bool     `json:"globalProtectionModeEnabled"`
	EntryTime                   float64  `json:"entryTime"`
	Interface                   string   `json:"interface"`
	Duration                    uint     `json:"duration"`
	DurationUsecs               uint32   `json:"durationUsecs"`
	EntryType                   string   `json:"entryType"`
	QueueLength                 uint32   `json:"queueLength"`
	TrafficClass                uint     `json:"trafficClass"`
	IngressPortSet              []string `json:"ingressPortSet"`
}

func (l *ShowQueueMonitor) SetCmd(port string, limitBy string, limitValue int) {
	base := "show queue-monitor length"
	if limitBy != "" && limitValue != 0 {
		l.Cmd = fmt.Sprintf("%s %s limit %d %s", base, port, limitValue, limitBy)
	} else {
		l.Cmd = fmt.Sprintf("%s %s", base, port)
	}
}

func (l *ShowQueueMonitor) GetCmd() string {
	return l.Cmd
}

func (s *ShowEntity) ShowQueueMonitorWithLimit(port string, limitBy string, limitValue int) (ShowQueueMonitor, error) {
	showqueuemonitor := ShowQueueMonitor{}
	showqueuemonitor.SetCmd(port, limitBy, limitValue)

	handle, err := s.node.GetHandle("json")
	if err != nil {
		return showqueuemonitor, err
	}

	err = handle.AddCommand(&showqueuemonitor)
	if err != nil {
		return showqueuemonitor, err
	}

	err = handle.Call()
	if err != nil {
		return showqueuemonitor, err
	}

	handle.Close()
	return showqueuemonitor, nil
}

func (s *ShowEntity) ShowQueueMonitor(port string) (ShowQueueMonitor, error) {
	showqueuemonitor := ShowQueueMonitor{}
	showqueuemonitor.SetCmd(port, "", 0)

	handle, err := s.node.GetHandle("json")
	if err != nil {
		return showqueuemonitor, err
	}

	err = handle.AddCommand(&showqueuemonitor)
	if err != nil {
		return showqueuemonitor, err
	}

	err = handle.Call()
	if err != nil {
		return showqueuemonitor, err
	}

	handle.Close()
	return showqueuemonitor, nil
}
