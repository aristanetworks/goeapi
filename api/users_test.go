//
// Copyright (c) 2015, Arista Networks, Inc.
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
package api

import (
	"regexp"
	"testing"
)

func TestUsersFuntions_UnitTest(t *testing.T) {
	isPrivTests := []struct {
		in   int
		want bool
	}{
		{-1, false},
		{0, true},
		{1, true},
		{4, true},
		{12, true},
		{15, true},
		{16, false},
		{4096, false},
		{65535, false},
	}

	for _, tt := range isPrivTests {
		if got := isPrivilege(tt.in); got != tt.want {
			t.Fatalf("isPrivilege(%q) = %t; want %t", tt.in, got, tt.want)
		}
	}
}

func TestUsersIsEqual_UnitTest(t *testing.T) {
	tests := []struct {
		in   UserConfig
		want bool
	}{
		{
			UserConfig{
				"username":   "test",
				"privilege":  "1",
				"role":       "network-admin",
				"nopassword": "true",
				"format":     "5",
				"secret":     "$1$eCurfHLe$JCbuUNM7Xwy6i6/zknYha.",
				"sshkey":     "",
			}, true,
		},
		{
			UserConfig{
				"username":   "test",
				"privilege":  "1",
				"role":       "network-admin",
				"nopassword": "true",
				"secret":     "$1$eCurfHLe$JCbuUNM7Xwy6i6/zknYha.",
			}, false,
		},
		{
			UserConfig{
				"username":   "test",
				"role":       "network-admin",
				"nopassword": "true",
				"format":     "5",
				"secret":     "$1$eCurfHLe$JCbuUNM7Xwy6i6/zknYha.",
				"sshkey":     "",
			}, false,
		},
		{
			UserConfig{
				"username":   "test",
				"privilege":  "1",
				"nopassword": "true",
				"format":     "5",
				"secret":     "$1$eCurfHLe$JCbuUNM7Xwy6i6/zknYha.",
				"sshkey":     "",
			}, false,
		},
		{
			UserConfig{
				"username":  "test",
				"privilege": "1",
				"role":      "network-admin",
				"sshkey":    "",
			}, false,
		},
	}

	conf := UserConfig{
		"username":   "test",
		"privilege":  "1",
		"role":       "network-admin",
		"nopassword": "true",
		"format":     "5",
		"secret":     "$1$eCurfHLe$JCbuUNM7Xwy6i6/zknYha.",
		"sshkey":     "",
	}

	for _, tt := range tests {
		if got := conf.isEqual(tt.in); got != tt.want {
			t.Fatalf("isEqual(%q) = %t; want %t", tt.in, got, tt.want)
		}
	}
}

func TestUsersParseUsername_UnitTest(t *testing.T) {
	tests := []struct {
		in   string
		want UserConfig
	}{
		{
			"username admin privilege 1 role network-admin nopassword\n",
			UserConfig{
				"username":   "admin",
				"privilege":  "1",
				"role":       "network-admin",
				"nopassword": "true",
				"format":     "",
				"secret":     "",
				"sshkey":     "",
			},
		},
		{
			"username eapi privilege 1 secret 5 $1$eCurfHLe$JCbuUNM7Xwy6i6/zknYha.\n",
			UserConfig{
				"username":   "eapi",
				"privilege":  "1",
				"role":       "",
				"nopassword": "false",
				"format":     "5",
				"secret":     "$1$eCurfHLe$JCbuUNM7Xwy6i6/zknYha.",
				"sshkey":     "",
			},
		},
		{
			"username test privilege 10 nopassword\n",
			UserConfig{
				"username":   "test",
				"privilege":  "10",
				"role":       "",
				"nopassword": "true",
				"format":     "",
				"secret":     "",
				"sshkey":     "",
			},
		},
		{
			"username test1 privilege 1 secret 5 $1$o/po05ru$92uegC/GGu3i4MS7MH9AE0\n",
			UserConfig{
				"username":   "test1",
				"privilege":  "1",
				"role":       "",
				"nopassword": "false",
				"format":     "5",
				"secret":     "$1$o/po05ru$92uegC/GGu3i4MS7MH9AE0",
				"sshkey":     "",
			},
		},
		{
			"username test2 privilege 9 role ops secret 5 $1$Kraw0Knu$dfIURYuRCxzDDcyyKnAD1/\n",
			UserConfig{
				"username":   "test2",
				"privilege":  "9",
				"role":       "ops",
				"nopassword": "false",
				"format":     "5",
				"secret":     "$1$Kraw0Knu$dfIURYuRCxzDDcyyKnAD1/",
				"sshkey":     "",
			},
		},
		{
			"username test3 privilege 9 role ops nopassword\n",
			UserConfig{
				"username":   "test3",
				"privilege":  "9",
				"role":       "ops",
				"nopassword": "true",
				"format":     "",
				"secret":     "",
				"sshkey":     "",
			},
		},
	}

	for _, tt := range tests {
		got := parseUsername(tt.in)
		if got == nil {
			t.Fatalf("parseUsername(%q) == nil", tt.in)
		}
		if !got.isEqual(tt.want) {
			t.Fatalf("parseUsername(%q) = %q; want %q", tt.in, got, tt.want)
		}
	}
}

var testSSHKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDKL1UtBALa4CvFUsHUipNym" +
	"A04qCXuAtTwNcMj84bTUzUI+q7mdzRCTLkllXeVxKuBnaTm2PW7W67K5CVpl0" +
	"EVCm6IY7FS7kc4nlnD/tFvTvShy/fzYQRAdM7ZfVtegW8sMSFJzBR/T/Y/sxI" +
	"16Y/dQb8fC3la9T25XOrzsFrQiKRZmJGwg8d+0RLxpfMg0s/9ATwQKp6tPoLE" +
	"4f3dKlAgSk5eENyVLA3RsypWADHpenHPcB7sa8D38e1TS+n+EUyAdb3Yov+5E" +
	"SAbgLIJLd52Xv+FyYi0c2L49ByBjcRrupp4zfXn4DNRnEG4K6GcmswHuMEGZv" +
	"5vjJ9OYaaaaaaa"

func TestUserGet_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
			"username test nopassword",
			"username test sshkey " + testSSHKey,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)
		userConfig := user.Get("test")

		tmpConfig := UserConfig{
			"username":   "test",
			"privilege":  "1",
			"nopassword": "true",
			"role":       "",
			"format":     "",
			"secret":     "",
			"sshkey":     testSSHKey,
		}

		if tmpConfig.isEqual(userConfig) != true {
			t.Fatalf("Unequal configs.")
		}
	}
}

func TestUserGetGetters_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
			"username test nopassword",
			"username test sshkey " + testSSHKey,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)
		userConfig := user.Get("test")

		tmpConfig := UserConfig{
			"username":   "test",
			"privilege":  "1",
			"nopassword": "true",
			"role":       "",
			"format":     "",
			"secret":     "",
			"sshkey":     testSSHKey,
		}

		if tmpConfig.UserName() != userConfig.UserName() ||
			tmpConfig.Privilege() != userConfig.Privilege() ||
			tmpConfig.Role() != userConfig.Role() ||
			tmpConfig.Nopassword() != userConfig.Nopassword() ||
			tmpConfig.Format() != userConfig.Format() ||
			tmpConfig.Secret() != userConfig.Secret() ||
			tmpConfig.SSHKey() != userConfig.SSHKey() {
			t.Fatalf("Unequal configs.")
		}
	}
}

func TestUserGetInvalid_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)
		userConfig := user.Get("test")

		if userConfig != nil {
			t.Fatalf("Invalid Get(name) returns non-nil value")
		}
	}
}

func TestUserGetAll_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
			"username test nopassword",
			"username test sshkey " + testSSHKey,
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)
		userConfigs := user.GetAll()

		tmpConfig := UserConfig{
			"username":   "test",
			"privilege":  "1",
			"nopassword": "true",
			"role":       "",
			"format":     "",
			"secret":     "",
			"sshkey":     testSSHKey,
		}

		if tmpConfig.isEqual(userConfigs["test"]) != true {
			t.Fatalf("Unequal configs.")
		}
	}

}

func TestUserCreate_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)

		config := user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1`, config); found {
			t.Fatalf("\"username test privilege 1\" NOT expected but not seen under "+
				"user section.\n[%s]", config)
		}

		if ok, err := user.Create("test", true, "", ""); !ok || err != nil {
			t.Fatalf("Create of user failed")
		}

		config = user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1`, config); !found {
			t.Fatalf("\"username test privilege 1\" expected but not seen under "+
				"user section.\n[%s]", config)
		}
	}
}

func TestUserCreateInvalid_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)

		config := user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1`, config); found {
			t.Fatalf("\"username test privilege 1\" NOT expected but not seen under "+
				"user section.\n[%s]", config)
		}

		if ok, err := user.Create("test", false, "", ""); ok || err == nil {
			t.Fatalf("Create with nopasswd or secret not specified should fail")
		}
	}
}

func TestUserCreateWithSecret_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)

		config := user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1`, config); found {
			t.Fatalf("\"username test privilege 1\" NOT expected but not seen under "+
				"user section.\n[%s]", config)
		}

		if ok, err := user.CreateWithSecret("test", "password", "cleartext"); !ok || err != nil {
			t.Fatalf("CreateWithSecret of user failed")
		}

		config = user.GetSection()
		if found, _ := regexp.MatchString(`username test`, config); !found {
			t.Fatalf("\"username test privilege 1\" expected but not seen under "+
				"user section.\n[%s]", config)
		}
	}
}

func TestUserCreateWithSecretInvalidType_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)

		config := user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1`, config); found {
			t.Fatalf("\"username test privilege 1\" NOT expected but not seen under "+
				"user section.\n[%s]", config)
		}

		if ok, err := user.CreateWithSecret("test", "password", "invalidType"); ok || err == nil {
			t.Fatalf("CreateWithSecret using invalid type should fail")
		}
	}
}

func TestUserDelete_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
			"username test nopassword",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)

		config := user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1 nopassword`, config); !found {
			t.Fatalf("\"username test privilege 1\" expected but not seen under "+
				"user section.\n[%s]", config)
		}

		if ok := user.Delete("test"); !ok {
			t.Fatalf("Delete of user failed")
		}

		config = user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1 nopassword`, config); found {
			t.Fatalf("\"username test privilege 1\" expected but not seen under "+
				"user section.\n[%s]", config)
		}
	}
}

func TestUserDefault_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
			"username test nopassword",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)

		config := user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1 nopassword`, config); !found {
			t.Fatalf("\"username test privilege 1\" expected but not seen under "+
				"user section.\n[%s]", config)
		}

		if ok := user.Default("test"); !ok {
			t.Fatalf("Default config for user failed")
		}

		config = user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1 nopassword`, config); found {
			t.Fatalf("\"username test privilege 1\" expected but not seen under "+
				"user section.\n[%s]", config)
		}
	}
}

func TestUserSetPrivWithVal_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
			"username test nopassword",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)

		config := user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1 nopassword`, config); !found {
			t.Fatalf("\"username test privilege 1\" expected but not seen under "+
				"user section.\n[%s]", config)
		}

		if ok, err := user.SetPrivilege("test", 8); !ok || err != nil {
			t.Fatalf("SetPrivilege config for user failed")
		}

		config = user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 8 nopassword`, config); !found {
			t.Fatalf("\"username test privilege 8 nopasswd\" expected but not seen under "+
				"user section.\n[%s]", config)
		}
	}
}

func TestUserSetPrivWithInvalidVal_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
			"username test privilege 8 nopassword",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)

		config := user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 8 nopassword`, config); !found {
			t.Fatalf("\"username test privilege 8\" expected but not seen under "+
				"user section.\n[%s]", config)
		}

		if ok, err := user.SetPrivilege("test", 65535); ok || err == nil {
			t.Fatalf("SetPrivilege config for user failed")
		}

		config = user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 8 nopassword`, config); !found {
			t.Fatalf("\"username test privilege 8 nopasswd\" expected but not seen under "+
				"user section.\n[%s]", config)
		}
	}
}

func TestUserSetRoleWithVal_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
			"username test nopassword",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)

		config := user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1 nopassword`, config); !found {
			t.Fatalf("\"username test privilege 1\" expected but not seen under "+
				"user section.\n[%s]", config)
		}

		if ok := user.SetRole("test", "network-admin"); !ok {
			t.Fatalf("SetRole config for user failed")
		}

		config = user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1 role network-admin nopassword`, config); !found {
			t.Fatalf("\"username test privilege 1 role network-admin nopasswd\" expected but not seen under "+
				"user section.\n[%s]", config)
		}
	}
}

func TestUserSetRoleWithNoVal_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
			"username test role network-admin nopassword",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)

		config := user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1 role network-admin nopassword`, config); !found {
			t.Fatalf("\"username test privilege 1\" expected but not seen under "+
				"user section.\n[%s]", config)
		}

		if ok := user.SetRole("test", ""); !ok {
			t.Fatalf("SetRole config for user failed")
		}

		config = user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1 role network-admin nopassword`, config); found {
			t.Fatalf("\"username test privilege 1 role network-admin nopasswd\" expected but not seen under "+
				"user section.\n[%s]", config)
		}
	}
}

func TestUserSetSSHKeyWithVal_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
			"username test nopassword",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)

		config := user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1 nopassword`, config); !found {
			t.Fatalf("\"username test privilege 1\" expected but not seen under "+
				"user section.\n[%s]", config)
		}

		if ok := user.SetSshkey("test", testSSHKey); !ok {
			t.Fatalf("SetSshkey config for user failed")
		}

		config = user.GetSection()
		configStr := "username test sshkey " + regexp.QuoteMeta(testSSHKey)
		if found, _ := regexp.MatchString(configStr, config); !found {
			t.Fatalf("\"%s\" expected but not seen under "+
				"user section.\n[%s]", configStr, config)
		}
	}
}

func TestUserSetSSHKeyWithNoVal_SystemTest(t *testing.T) {
	for _, dut := range duts {
		cmds := []string{
			"no username test",
			"username test nopassword",
		}
		if ok := dut.Config(cmds...); !ok {
			t.Fatalf("dut.Config() failure")
		}

		user := User(dut)

		config := user.GetSection()
		if found, _ := regexp.MatchString(`username test privilege 1 nopassword`, config); !found {
			t.Fatalf("\"username test privilege 1\" expected but not seen under "+
				"user section.\n[%s]", config)
		}

		if ok := user.SetSshkey("test", ""); !ok {
			t.Fatalf("SetSshkey config for user failed")
		}

		config = user.GetSection()
		configStr := "username test sshkey " + regexp.QuoteMeta(testSSHKey)
		if found, _ := regexp.MatchString(configStr, config); found {
			t.Fatalf("\"%s\" NOT expected but not seen under "+
				"user section.\n[%s]", configStr, config)
		}
	}
}
