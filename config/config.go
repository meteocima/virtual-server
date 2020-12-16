package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// HostType is an enum that represents
// all possible host type. An host type
// indicates how and where the processes is started
type HostType int

const (
	// HostTypeOS represents an host that run process on the local machine
	HostTypeOS HostType = iota
	// HostTypeSSH represents an host that run process on a remote machine using SSH
	HostTypeSSH
)

// Host is struct that contains information
// about a host on which to run processes
type Host struct {
	Type        HostType
	Name        string
	Host        string
	BackupHosts []string `toml:"backup-hosts"`
	Port        int
	User        string
	Key         string
}

// Type is a structure which contains the
// configuration for the running command.
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
