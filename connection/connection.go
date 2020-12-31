package connection

import (
	"fmt"
	"io"
	"os"

	"github.com/meteocima/virtual-server/config"
	"github.com/meteocima/virtual-server/vpath"
)

// Process represents a running process
type Process interface {
	Kill() error
	// Stdin, is an io.Writer that will be used
	// to send data to process stdin
	Stdin() io.Writer

	// Stdin, if set, is an io.Reader that will be used
	// to read data from process stdout
	Stdout() io.Reader

	// Stderr, if set, is an io.Reader that will be used
	// to read data from process stderr
	Stderr() io.Reader

	// CombinedOutput returns an io.Reader that reads
	// the combined output and error streams of the process
	CombinedOutput() io.Reader

	// Wait expects the process to terminate
	// and return the exit code.
	Wait() (int, error)
}

// RunOptions ...
type RunOptions struct {
	// OutFromLog if sets, output is read from a file
	// instead of from the process stdout
	OutFromLog vpath.VirtualPath

	// OutFromLog if sets, output is read from a file
	// instead of from the process stderr
	ErrFromLog vpath.VirtualPath

	// Cwd is set the work directory in which the process will be executed.
	Cwd vpath.VirtualPath
	/*
		// Stdin, if set, is an io.Reader that will be used
		// as process Stdin.
		// If nil, a pipe to `Process.Stdin` member is created
		// and used.
		Stdin *io.Reader

		// Stdout, if set, is an io.Writer that will be used
		// as process Stdout
		// If nil, a pipe to `Process.Stdout` member is created
		// and used.
		Stdout *io.Writer

		// Stderr, if set, is an io.Writer that will be used
		// as process Stderr.
		// If nil, a pipe to `Process.Stdout` member is created
		// and used.
		Stderr *io.Writer
	*/
}

/*
// CopyOptions ...
type CopyOptions struct {
}

// MoveOptions ...
type MoveOptions struct {
}
*/

// Connection ...
type Connection interface {
	HostName() string
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
	Run(command vpath.VirtualPath, args []string, options ...RunOptions) (Process, error)
}

var connections = map[string]Connection{}

// NewPath ...
func NewPath(cn Connection, path string, pathArgs ...interface{}) vpath.VirtualPath {
	return vpath.New(cn.HostName(), path, pathArgs...)
}

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
			Name:     host.Name,
			Host:     host.Host,
			Port:     host.Port,
			User:     host.User,
			KeyPath:  host.Key,
			hostName: name,
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
