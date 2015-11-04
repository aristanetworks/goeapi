# Arista eAPI Golang Library

* JSON-RPC over HTML - not supported in the standard library JSON-RPC package
* Specific structs for switch information


## Requirements
* Arista EOS v4.12 or later
* Arista eAPI enabled for either http or https

## Getting Started
The following steps need to be followed to assure successful configuration of goeapi:

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
* **port** - Configures the port to use for the eAPI connection. (Currently Not Implemented)

_Note:_ See the EOS User Manual found at arista.com for more details on configuring eAPI values.

# Using Arista eAPI with the Go programming language

## Web
web directory is an attempt at creating REST API

## Main
main directory is an example application

