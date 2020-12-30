// Package config allows to load
// the global configuration from a `toml`
// file using the `config.Init` function.
//
// Once Init is called, the configuration
// file path is available as `config.Filename`,
// and the configured hosts as `config.Hosts`
//
// ## Example
//
// __main.go__
//
// ```go
//   import "github.com/meteocima/virtual-server/config"
//
//   func main() {
//     err := config.Init("./config.toml")
//     if err != nil {
//       log.Fatal(err.Error())
//     }
//   }
// ```
//
// __config.toml__
//
// ```
//  [hosts]
//
//  [hosts.localhost]
//  type = 0 #HostTypeOS
//
//
//  [hosts.drihm]
//  type = 1 #HostTypeSSH
//  host = "localhost"
//  port = 2222
//  user = "andrea.parodi"
//  key = "/var/fixtures/private-key"
//
//
//  [hosts.withbackup]
//  type = 1 #HostTypeSSH
//  host = "example.com"
//  backup-hosts = ["local", "drihm"]
//  port = 22
//  user = "andrea.parodi"
//  key = "/var/fixtures/private-key"
//
// ```
//
package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// HostType is an enum that represents
// all possible host type. An host type
// indicates how and where the processes
// is started. An HostType variable can
// have following values:
type HostType int

const (
	// HostTypeOS represents an host that
	// run processes on the local machine
	HostTypeOS HostType = iota

	// HostTypeSSH represents an host that
	// run processes on a remote machine
	// through an SSH connection.
	HostTypeSSH
)

// Host is struct that contains information
// about a host on which to run processes
type Host struct {
	// Contains the type of the host.
	// It can be either `HostTypeOS` or
	// `HostTypeSSH`
	Type HostType
	// Name of the host, written at
	// runtime using the key of the
	// host section in the config file.
	Name string
	// Hostname of the server,
	// used only for SSH type hosts.
	Host string
	// A list of backup hostnames
	// to use in case of failure
	// connecting.
	BackupHosts []string `toml:"backup-hosts"`
	// Tcp port to use
	Port int
	// Username to use to authenticate on
	// the host
	User string
	// Local path of the private SSH
	// key file.
	Key string
}

// Type is a structure which contains the
// configuration data for the running command.
//
// _Used internally and exported only
// to properly unmasharl the `toml` configuration_
type Type struct {
	Hosts map[string]*Host
}

// Hosts contains the configuration public instance
var Hosts map[string]*Host

// Filename contains the absolute path of
// the configuration file used to initialize
// the module
var Filename string

// Init loads the global, public configuration
// from the given file.
func Init(configFile string) error {
	var cfg Type
	_, err := toml.DecodeFile(configFile, &cfg)
	if err != nil {
		return err
	}

	for name, host := range cfg.Hosts {
		host.Name = name
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	Filename, err = filepath.Rel(wd, configFile)
	if err != nil {
		return err
	}

	Hosts = cfg.Hosts
	return nil
}
