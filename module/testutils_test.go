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

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

var fixturesPath = ""

func init() {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fixturesPath = path.Join(workingDir, "../testdata/fixtures")

	rand.Seed(time.Now().UTC().UnixNano())
}

const charBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString randomly creates a string within the
// range of minchar to maxchar
func RandomString(minchar int, maxchar int) string {
	size := RandomInt(minchar, maxchar)
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[i] = charBytes[rand.Intn(len(charBytes)-1)]
	}
	return string(bytes)
}

// RandomIPAddress randomly creates a ip address
func RandomIPAddress() string {
	return strconv.Itoa(RandomInt(1, 100)) +
		"." + strconv.Itoa(RandomInt(1, 100)) +
		"." + strconv.Itoa(RandomInt(1, 100)) +
		"." + strconv.Itoa(RandomInt(1, 100))
}

// RandomInt randomly creates a int between
// min and max
func RandomInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// RandomVlan gets a random vlan number
func RandomVlan() int {
	return RandomInt(1, 4094)
}

// RandomChoice chooses a random element in the list
func RandomChoice(list []string) string {
	if len(list) == 0 {
		return ""
	}
	return list[RandomInt(0, len(list))]
}

// GetFixturesPath aquires the global fixtures path
func GetFixturesPath() string {
	return fixturesPath
}

// GetFixture returns the full path to filenmae within
// fixtures
func GetFixture(filename string) string {
	return path.Join(fixturesPath, filename)
}

// LoadFixtureFile reads the fixtures file into a string.
func LoadFixtureFile(file string) string {
	data, err := ioutil.ReadFile(GetFixture(file))
	if err != nil {
		panic(err)
	}
	return string(data)
}
