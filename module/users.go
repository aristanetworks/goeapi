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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aristanetworks/goeapi"
)

// defaultEncryption
const defaultEncryption = "cleartext"

var encryptionMap = map[string]string{
	"cleartext": "0",
	"md5":       "5",
	"sha512":    "sha512",
}

var usersRegex = regexp.MustCompile(`(?m)username ([^\s]+) privilege (\d+)` +
	`(?m)(?: role ([^\s]+))?` +
	`(?m)(?: (nopassword))?` +
	`(?m)(?: secret (0|5|7|sha512) (.+))?` +
	`(?m).*$\n(?:username ([^\s]+) sshkey (.+)$)?`)

// UserConfig represents a parsed user entry containing:
//
//      "username" : "",
//     "privilege" : "",
//          "role" : "",
//    "nopassword" : "",
//        "secret" : "",
//        "format" : "",
//        "sshkey" : "",
//
//
type UserConfig map[string]string

// UserConfigMap is a mapped entry of UserConfigs containing:
//
//      "username" : {
//            "username" : "",
//           "privilege" : "",
//                "role" : "",
//          "nopassword" : "",
//              "secret" : "",
//              "format" : "",
//              "sshkey" : "",
//      }
//
type UserConfigMap map[string]UserConfig

// UserName returns the username(string) entry for
// this UserConfig
func (u UserConfig) UserName() string {
	return u["username"]
}

// Privilege returns the privilege mode(string) for
// this UserConfig
func (u UserConfig) Privilege() string {
	return u["privilege"]
}

// Role returns the role for this UserConfig
func (u UserConfig) Role() string {
	return u["role"]
}

// Nopassword returns 'true'(string) if nopassowrd is
// configured for this UserConfig
func (u UserConfig) Nopassword() string {
	return u["nopassword"]
}

// Format returns the format(string) for
// this UserConfig
func (u UserConfig) Format() string {
	return u["format"]
}

// Secret returns the login secret passwd type for
// this UserConfig
func (u UserConfig) Secret() string {
	return u["secret"]
}

// SSHKey returns the sshkey(string) for
// this UserConfig
func (u UserConfig) SSHKey() string {
	return u["sshkey"]
}

// UserEntity resource provides configuration for local user resources of
// an EOS node
type UserEntity struct {
	*AbstractBaseEntity
}

// User factory function to initiallize User resource
// given a Node
func User(node *goeapi.Node) *UserEntity {
	return &UserEntity{&AbstractBaseEntity{node}}
}

// isPrivilege Checks value for valid user privilege level.
// True if the value is valid, otherwise False.
func isPrivilege(value int) bool {
	return (value >= 0 && value < 16)
}

// isEqual compaires UserConfigs
func (u UserConfig) isEqual(dest UserConfig) bool {
	if len(u) != len(dest) {
		return false
	}
	for k, v := range u {
		//if dest[k] != v {
		val, found := dest[k]
		if !found || val != v {
			return false
		}
	}
	return true
}

// Get Returns the local user configuration as a UserConfig.
//
// Args:
//  name (string): The user name to return a resource for from the
//                nodes configuration
//
// Returns:
//  UserConfig type
func (u *UserEntity) Get(name string) UserConfig {
	resource, found := u.GetAll()[name]
	if !found {
		return nil
	}
	return resource
}

// GetAll Returns the local user configuration as UserConfig
//
// Returns:
//  UserConfigMap object
func (u *UserEntity) GetAll() UserConfigMap {
	var resources = make(UserConfigMap)

	config := u.Config()
	users := usersRegex.FindAllString(config, -1)

	for _, user := range users {
		result := parseUsername(user)
		if result == nil {
			continue
		}
		resources[result["username"]] = result
	}
	return resources
}

// GetSection Returns the local user configuration as string
//
// Returns:
//  UserConfigMap object
func (u *UserEntity) GetSection() string {

	config := u.Config()
	users := usersRegex.FindAllString(config, -1)
	return strings.Join(users, "")
}

// parseUsername Scans the config block and returns the username
// as a UserConfig object
//
// Args:
//  config (string): The config block to parse
//
// Returns:
//  UserConfig object
func parseUsername(config string) UserConfig {
	var re = regexp.MustCompile(`(?m)username ([^\s]+) privilege (\d+)` +
		`(?: role ([^\s]+))?` +
		`(?: (nopassword))?` +
		`(?: secret (0|5|7|sha512) (.+))?` +
		`.*$\n(?:username ([^\s]+) sshkey (.+)$)?`)

	var resource = make(UserConfig)

	match := re.FindStringSubmatch(config)
	if match == nil {
		return nil
	}
	resource["username"] = match[1]
	resource["privilege"] = match[2]
	resource["role"] = match[3]

	resource["nopassword"] = "false"
	if match[4] == "nopassword" {
		resource["nopassword"] = "true"
	}
	resource["format"] = match[5]
	resource["secret"] = match[6]
	resource["sshkey"] = match[8]

	//for idx, val := range match {
	//    fmt.Printf("val[%d]: %s\n", idx, val)
	//}
	return resource
}

//Create Creates a new user on the local system.
//
// Args:
//  name (string):     The name of the user to craete
//  nopassword (bool): Configures the user to be able to authenticate
//                     without a password challenage
//  secret (string):   The secret (password) to assign to this user
//  encryption (string): Specifies how the secret is encoded.  Valid
//                           values are "cleartext", "md5", "sha512".
//                           The default is "cleartext"
// Returns:
//  True if the operation was successful otherwise False
func (u *UserEntity) Create(name string, nopassword bool, secret string,
	encryption string) (bool, error) {
	if secret != "" {
		return u.CreateWithSecret(name, secret, encryption)
	} else if nopassword {
		return u.CreateWithNoPassword(name), nil
	} else {
		return false, fmt.Errorf("either \"nopassword\" or \"secret\" must be" +
			" specified to create a user")
	}
}

// CreateWithSecret Creates a new user on the local node
//
// Args:
//  name (string):       The name of the user to craete
//  secret (string):     The secret (password) to assign to this user
//  encryption (string): Specifies how the secret is encoded.  Valid
//                      values are "cleartext", "md5", "sha512".  The
//                      default is "cleartext"
// Returns:
//  True if the operation was successful otherwise False
func (u *UserEntity) CreateWithSecret(name string, secret string,
	encryption string) (bool, error) {
	var enc string

	enc, found := encryptionMap[encryption]
	if !found {
		return false, fmt.Errorf("encryption must be one of \"cleartext\", " +
			"\"md5\" or \"sha512\"")
	}
	cmd := "username " + name + " secret " + enc + " " + secret
	return u.Configure(cmd), nil
}

// CreateWithNoPassword Creates a new user on the local node
//
// Args:
//  name (string): The name of the user to create
// Returns:
//  True if the operation was successful otherwise False
func (u *UserEntity) CreateWithNoPassword(name string) bool {
	var cmd = "username " + name + " nopassword"
	return u.Configure(cmd)
}

// Delete Deletes the local username from the config
//
// Args:
//  name (string): The name of the user to delete
// Returns:
//  True if the operation was successful otherwise False
func (u *UserEntity) Delete(name string) bool {
	var cmd = "no username " + name
	return u.Configure(cmd)
}

// Default Configures the local username using the default keyword
//
// Args:
//  name (string): The name of the user to configure
// Returns:
//  True if the operation was successful otherwise False
func (u *UserEntity) Default(name string) bool {
	var cmd = "default username " + name
	return u.Configure(cmd)
}

// SetPrivilege Configures the user privilege value in EOS
//
// Args:
//  name (string): The name of the user to craete
//  value (int):   The privilege value to assign to the user.  Valid
//                 values are in the range of 0 to 15
// Returns:
//  True if the operation was successful otherwise False
func (u *UserEntity) SetPrivilege(name string, value int) (bool, error) {
	if !isPrivilege(value) {
		return false, fmt.Errorf("priviledge value must be between 0 and 15")
	}
	var cmd = "username " + name + " privilege " + strconv.Itoa(value)
	return u.Configure(cmd), nil
}

// SetRole Configures the user role vale in EOS
//
// Args:
//  name (string):  The name of the user to craete
//  value (string): The value to configure for the user role
// Returns:
//  True if the operation was successful otherwise False
func (u *UserEntity) SetRole(name string, value string) bool {
	var cmd = "username " + name
	if value != "" {
		cmd = cmd + " role " + value
	} else {
		cmd = "default " + cmd + " role"
	}
	return u.Configure(cmd)
}

// SetSshkey Configures the user sshkey
//
// Args:
//  name (string):  The name of the user to add the sshkey to
//  value (string): The value to configure for the sshkey.
// Returns:
//  True if the operation was successful otherwise False
func (u *UserEntity) SetSshkey(name string, value string) bool {
	var cmd = "username " + name
	if value != "" {
		cmd = cmd + " sshkey " + value
	} else {
		cmd = "no " + cmd + " sshkey"
	}
	return u.Configure(cmd)
}
