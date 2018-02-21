# Arista Go eAPI Library [![Build Status](https://travis-ci.org/aristanetworks/goeapi.svg?branch=master)](https://travis-ci.org/aristanetworks/goeapi) [![codecov.io](http://codecov.io/github/aristanetworks/goeapi/coverage.svg?branch=master)](http://codecov.io/github/aristanetworks/goeapi?branch=master) [![GoDoc](https://godoc.org/github.com/aristanetworks/goeapi?status.png)](https://godoc.org/github.com/aristanetworks/goeapi)


#### Table of Contents

1. [Overview](#overview)
    * [Requirements](#requirements)
2. [Installation](#installation)
3. [Upgrading](#upgrading)
4. [Getting Started](#getting-started)
    * [Example eapi.conf File](#example-eapiconf-file)
    * [Using goeapi](#using-goeapi)
5. [Building Local Documentation](#building-documention)
6. [Testing](#testing)
7. [Contributing](#contributing)
8. [License](#license)


# Overview

The Go Client for eAPI provides a native Go implementation for programming Arista EOS network devices using Golang.  The Go client provides the ability to build native applications in Go that can communicate with EOS remotely over a HTTP/S transport (off-box).  It uses a standard INI-style configuration file to specifiy one or more connection profiles.

The goeapi implemenation also provides an API layer for building native Go objects that allow for configuration and state extraction of EOS nodes.  The API layer provides a consistent implementation for working with EOS configuration resources.  The implementation of the API layer is highly extensible and can be used as a foundation for building custom data models.

The libray is freely provided to the open source community for building robust applications using Arista EOS eAPI.  Support is provided as best effort through Github iusses.

## Requirements
* Arista EOS v4.12 or later
* Arista eAPI enabled for either http or https
* Go 1.5+

## Installation
First, it is assumed you have and are working in a standard [Go](https://www.golang.org) workspace, as described in http://golang.org/doc/code.html, with proper [GOPATH](https://golang.org/doc/code.html#GOPATH) set. Go 1.5+ is what's recommended for using goeapi. To download and install goeapi:

```console
$ go get github.com/aristanetworks/goeapi
```

After setting up Go and installing goeapi, any required build tools can be installed by bootstrapping your environment via:

```console
$ make bootstrap
```

# Upgrading
  ```
  $ go get -u github.com/aristanetworks/goeapi
  ```

# Getting Started
The following steps need to be followed to assure successful configuration of goeapi.

1. EOS Command API must be enabled

    To enable EOS Command API from configuration mode, configure proper protocol under management api, and
    then verify:
    ```
        Switch# configure terminal
        Switch(config)# management api http-commands
        Switch(config-mgmt-api-http-cmds)# protocol ?
          http         Configure HTTP server options
          https        Configure HTTPS server options
          unix-socket  Configure Unix Domain Socket
        Switch(config-mgmt-api-http-cmds)# protocol http
        Switch(config-mgmt-api-http-cmds)# end

        Switch# show management api http-commands
        Enabled:            Yes
        HTTPS server:       running, set to use port 443
        HTTP server:        running, set to use port 80
        Local HTTP server:  shutdown, no authentication, set to use port 8080
        Unix Socket server: shutdown, no authentication
        ...
```

2. Create configuration file with proper node properties. (*See eapi.conf file examples below*)

    **Note:** The default search path for the conf file is ``~/.eapi.conf``
    followed by ``/mnt/flash/eapi.conf``.   This can be overridden by setting
    ``EAPI_CONF=<path file conf file>`` in your environment.

## Example eapi.conf File
Below is an example of an eAPI conf file.  The conf file can contain more than
one node.  Each node section must be prefaced by **connection:\<name\>** where
\<name\> is the name of the connection.

The following configuration options are available for defining node entries:

* **host** - The IP address or FQDN of the remote device.  If the host parameter is omitted then the connection name is used
* **username** - The eAPI username to use for authentication (only required for http or https connections)
* **password** - The eAPI password to use for authentication (only required for http or https connections)
* **enablepwd** - The enable mode password if required by the destination node
* **transport** - Configures the type of transport connection to use.  The default value is _https_.  Valid values are:
  * http
  * https
  * socket
* **port** - Configures the port to use for the eAPI connection. (Currently Not Implemented)

_Note:_ See the EOS User Manual found at arista.com for more details on configuring eAPI values.

# Using Goeapi
Once goeapi has been installed and your .eapi.config file is setup correctly, you are now ready to try it out. Here is a working example of .eapi.config file and go program:

```sh
$ cat ~/.eapi.config
```
```console
[connection:arista1]
host=arista1
username=admin
password=root
enablepwd=passwd
transport=https
```
```sh
$ cat example1.go
```
```go
package main

import (
        "fmt"

        "github.com/aristanetworks/goeapi"
        "github.com/aristanetworks/goeapi/module"
)

func main() {
        // connect to our device
        node, err := goeapi.ConnectTo("Arista1")
        if err != nil {
                panic(err)
        }
        // get the running config and print it
        conf := node.RunningConfig()
        fmt.Printf("Running Config:\n%s\n", conf)

        // get api system module
        sys := module.System(node)
        // change the host name to "Ladie"
        if ok := sys.SetHostname("Ladie"); !ok {
                fmt.Printf("SetHostname Failed\n")
        }
        // get system info
        sysInfo := sys.Get()
        fmt.Printf("\nSysinfo: %#v\n", sysInfo.HostName())
}
```
goeapi provides a way for users to directly couple a command with a predefined response. The underlying api will issue the command and the response stored in the defined type. For example, lets say the configured vlan ports are needed for some form of processing. If we know the JSON response for the command composed like the following:
(from Arista Command API Explorer):
```json
 {
    "jsonrpc": "2.0",
    "result": [
       {
          "sourceDetail": "",
          "vlans": {
             "2": {
                "status": "active",
                "name": "VLAN0002",
                "interfaces": {
                   "Port-Channel10": {
                      "privatePromoted": false
                   },
                   "Ethernet2": {
                      "privatePromoted": false
                   },
                   "Ethernet1": {
                      "privatePromoted": false
                   },
                   "Port-Channel5": {
                      "privatePromoted": false
                   }
                },
                "dynamic": false
             },
          }
       }
    ],
    "id": "CapiExplorer-123"
 }
```
We can then build our Go structures based on the response format and couple our show command with the type:
```go
type MyShowVlan struct {
        SourceDetail string          `json:"sourceDetail"`
        Vlans        map[string]Vlan `json:"vlans"`
}

type Vlan struct {
        Status     string               `json:"status"`
        Name       string               `json:"name"`
        Interfaces map[string]Interface `json:"interfaces"`
        Dynamic    bool                 `json:"dynamic"`
}

type Interface struct {
        Annotation      string `json:"annotation"`
        PrivatePromoted bool   `json:"privatePromoted"`
}

func (s *MyShowVlan) GetCmd() string {
        return "show vlan configured-ports"
}
```
Since the command ``show vlan configured-ports`` is coupled with the response structure, the underlying api knows to issue the command and the response needs to be filled in. The resulting code looks like:
```go
package main

import (
        "fmt"

        "github.com/aristanetworks/goeapi"
)

type MyShowVlan struct {
        SourceDetail string          `json:"sourceDetail"`
        Vlans        map[string]Vlan `json:"vlans"`
}

type Vlan struct {
        Status     string               `json:"status"`
        Name       string               `json:"name"`
        Interfaces map[string]Interface `json:"interfaces"`
        Dynamic    bool                 `json:"dynamic"`
}

type Interface struct {
        Annotation      string `json:"annotation"`
        PrivatePromoted bool   `json:"privatePromoted"`
}

func (s *MyShowVlan) GetCmd() string {
        return "show vlan configured-ports"
}

func main() {
        node, err := goeapi.ConnectTo("dut")
        if err != nil {
                panic(err)
        }

        sv := &MyShowVlan{}

        handle, _ := node.GetHandle("json")
        handle.AddCommand(sv)
        if err := handle.Call(); err != nil {
                panic(err)
        }

        for k, v := range sv.Vlans {
                fmt.Printf("Vlan:%s\n", k)
                fmt.Printf("  Name  : %s\n", v.Name)
                fmt.Printf("  Status: %s\n", v.Status)
        }
}
```
Also, if several commands/responses have been defined, goeapi supports command stacking to batch issue all at once:
```go
    ...
	handle, _ := node.GetHandle("json")
	handle.AddCommand(showVersion)
	handle.AddCommand(showVlan)
	handle.AddCommand(showHostname)
	handle.AddCommand(showIp)
	if err := handle.Call(); err != nil {
		panic(err)
	}
	fmt.Printf("Version           : %s\n", showVersion.Version)
	fnt.Printf("Hostname          : %s\n", showHostname.Hostname)
    ...
```
There are several go example's using goeapi (as well as example .eapi.config file) provided in the examples directory.

# Building Local Documentation

Documentation can be generated locally in plain text via:
```sh
$ godoc github.com/aristanetworks/goeapi
```
Or you can run the local godoc server and view the html version of the documentation by pointing your browser at http://localhost:6060
```sh
$ make doc
    or
$ godoc -http=:6060 -index
```

# Testing

The goeapi library provides various tests. To run System specific tests, you will need to
update the ``dut.conf`` file (found in testutils/fixtures) to include the device level specifics
for your setup. The switch used for testing should have at least interfaces Ethernet1-7.

* For running System tests, issue the following from the root of the goeapi directory:
```sh
$ make systest
    or
$ go test ./... -run SystemTest$
```
* Similarly, Unit tests can be run via:
```sh
$ make unittest
    or
$ go test ./... -run UnitTest$
```
Verbose mode can be specified as a flag to provide additional information:
```sh
$ make GOTEST_FLAGS=-v test
```

Note: Test cases for XXX.go files live in respective XXX_test.go files and have the following function signature:
* Unit Tests: ``TestXXX_UnitTest(t *testing.T){...``
* System Tests: ``TestXXX_SystemTest(t *testing.T){...``

Any tests written must conform to this standard.

# Contributing

Contributing pull requests are gladly welcomed for this repository.  Please
note that all contributions that modify the library behavior require
corresponding test cases otherwise the pull request will be rejected.

# License

Copyright (c) 2015-2016, Arista Networks
Inc. All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

* Neither the name of Arista Networks nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL ARISTA NETWORKS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


