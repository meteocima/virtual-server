package connection

import (
	"io"
	"os"

	"github.com/meteocima/virtual-server/vpath"
)

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
}

// FindHost ...
func FindHost(name string) Connection {
	return nil
}
