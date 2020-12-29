# config
--
    import "."


## Usage

```go
var Filename string
```
Filename contains the absolute path of the configuration file used to initialize
the module

```go
var Hosts map[string]*Host
```
Hosts contains the configuration public instance

#### func  Init

```go
func Init(configFile string) error
```
Init loads the global, public configuration from the given file.

#### type Host

```go
type Host struct {
	Type        HostType
	Name        string
	Host        string
	BackupHosts []string `toml:"backup-hosts"`
	Port        int
	User        string
	Key         string
}
```

Host is struct that contains information about a host on which to run processes

#### type HostType

```go
type HostType int
```

HostType is an enum that represents all possible host type. An host type
indicates how and where the processes is started

```go
const (
	// HostTypeOS represents an host that run process on the local machine
	HostTypeOS HostType = iota
	// HostTypeSSH represents an host that run process on a remote machine using SSH
	HostTypeSSH
)
```

#### type Type

```go
type Type struct {
	Hosts map[string]*Host
}
```

Type is a structure which contains the configuration for the running command.
