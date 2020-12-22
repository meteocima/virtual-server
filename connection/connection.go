package connection

import (
	"fmt"
	"io"
	"os"

	"github.com/meteocima/virtual-server/config"
	"github.com/meteocima/virtual-server/vpath"
)

// RunOptions ...
type RunOptions struct {
	// OutFromLog if sets, output is read from a file
	// instead of from the process stdout
	OutFromLog *vpath.VirtualPath

	// OutFromLog if sets, output is read from a file
	// instead of from the process stderr
	ErrFromLog *vpath.VirtualPath

	// Cwd is set the work directory in which the process will be executed.
	Cwd *vpath.VirtualPath

	// Stdin, if set, is an io.Reader that will be used
	// as process Stdin
	Stdin *io.Reader

	// Stdout, if set, is an io.Writer that will be used
	// as process Stdout
	Stdout *io.Writer

	// Stderr, if set, is an io.Writer that will be used
	// as process Stderr
	Stderr *io.Reader
}

type CopyOptions struct {
}

type MoveOptions struct {
}

// Connection ...
type Connection interface {
	Open() error
	Close() error
	OpenReader(file vpath.VirtualPath) (io.ReadCloser, error)
	OpenWriter(file vpath.VirtualPath) (io.WriteCloser, error)
	ReadDir(dir vpath.VirtualPath) (vpath.VirtualPathList, error)
	Stat(path vpath.VirtualPath) (os.FileInfo, error)
	MkDir(dir vpath.VirtualPath) error
	RmDir(dir vpath.VirtualPath) error
	RmFile(file vpath.VirtualPath) error
	Link(source, target vpath.VirtualPath) error
	Run(command string, args []string, options ...RunOptions) error
}

var connections = map[string]Connection{}

// FindHost ...
func FindHost(name string) Connection {
	fail := func(msg string, args ...interface{}) {
		panic(fmt.Sprintf("Wrong configuration file \"%s\": ", config.Filename) + fmt.Sprintf(msg, args...))
	}

	cn, ok := connections[name]
	if ok {
		return cn
	}

	host, ok := config.Hosts[name]
	if !ok {
		fail("unknown host `%s`.", name)
	}

	if host.Type == config.HostTypeOS {
		cn = &LocalConnection{}
	} else if host.Type == config.HostTypeSSH {
		cn = &SSHConnection{
			Name:    host.Name,
			Host:    host.Host,
			Port:    host.Port,
			User:    host.User,
			KeyPath: host.Key,
		}
	} else {
		fail("unknown connection type %d for host `%s`.", host.Type, name)
	}

	connections[name] = cn

	err := cn.Open()
	if err != nil {
		fail("cannot connect to host `%s`: %w", name, err)
	}

	return cn
}
