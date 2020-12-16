package vpath

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

// VirtualPath represents an instance of a virtual path.
// It is formed by the name of a Host defined in configuration
// and by the absolute path within that file system.
// The two parts are separated by a colon when
// representing the vpath.VirtualPath as a string: host:path
// an empty string in Host field represent the localhost Host.
// an empty string in Path field represent the current directory (.).
type VirtualPath struct {
	Host string
	Path string
}

// private method that resolve
// nil fields to their corresponding
// values
func (vPath *VirtualPath) resolve() {
	if vPath == nil {
		panic("vPath is nil")
	}

	if vPath.Host == "" {
		vPath.Host = "localhost"
	}

	if vPath.Path == "" {
		vPath.Path = "."
	}
}

// New returns an VirtualPath given its
// host and path parts.
// The path is built using pathFormat argument as fmt.Sprintf format string
// and any pathArgs as fmt.Sprintf arguments.
func New(host string, pathFormat string, pathArgs ...interface{}) VirtualPath {
	p := VirtualPath{host, fmt.Sprintf(pathFormat, pathArgs...)}
	p.resolve()
	return p
}

// Local returns a new VirtualPath on localhost
// with the given path
func Local(pathFormat string, pathArgs ...interface{}) VirtualPath {
	return New("localhost", pathFormat, pathArgs...)
}

// FromS returns a new VirtualPath
// with host and path parsed from
// pathRepr string argument.
func FromS(pathRepr string) VirtualPath {
	parts := strings.SplitN(pathRepr, ":", 2)
	if len(parts) == 1 {
		return Local(parts[0])
	}

	return New(parts[0], parts[1])
}

// String returns a string representing the virtual path
// Host and path parts are separated by a colon: host:path
func (vPath VirtualPath) String() string {
	vPath.resolve()
	return vPath.Host + ":" + vPath.Path
}

// Join returns a new virtual path formed
// by the same host and a path that is the joining of the original path and
// and an additional one. The additional path is built
// using pathFormat argument as fmt.Sprintf format string
// and any pathArgs as fmt.Sprintf arguments.
func (vPath VirtualPath) Join(pathFormat string, pathArgs ...interface{}) VirtualPath {
	vPath.resolve()
	additionalPath := fmt.Sprintf(pathFormat, pathArgs...)
	newPath := path.Join(vPath.Path, additionalPath)
	return New(vPath.Host, newPath)
}

// JoinP returns a new virtual path formed
// by the same host of the instance and a path that is the joining of
// the original path and the path of an additional VirtualPath
func (vPath VirtualPath) JoinP(other VirtualPath) VirtualPath {
	vPath.resolve()
	newPath := path.Join(vPath.Path, other.Path)
	return New(vPath.Host, newPath)
}

// Dir returns a new VirtualPath formed
// by the same host and a path
// that is the directory path of the original
// instance.
func (vPath VirtualPath) Dir() VirtualPath {
	vPath.resolve()
	return New(vPath.Host, path.Dir(vPath.Path))
}

// Filename returns the filename (with extension,
//	but without directory path) of the virtual path
func (vPath VirtualPath) Filename() string {
	return path.Base(vPath.Path)
}

func rev(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ReplaceExt returns a new virtual path where
// the extension of the file is replaced
// with the given one.
func (vPath VirtualPath) ReplaceExt(newExt string) VirtualPath {
	vPath.resolve()
	if vPath.Path == "." {
		return New(vPath.Host, "."+newExt)
	}
	ext := filepath.Ext(vPath.Path)
	revExt := rev(ext)
	revPt := rev(vPath.Path)
	revNewExt := rev("." + newExt)

	return New(vPath.Host, rev(strings.Replace(revPt, revExt, revNewExt, 1)))
}

// AddExt returns a new virtual path where
// the specified extension is appended
// to the current filename.
func (vPath VirtualPath) AddExt(newExt string) VirtualPath {
	vPath.resolve()
	if vPath.Path == "." {
		return New(vPath.Host, "."+newExt)
	}
	return New(vPath.Host, vPath.Path+"."+newExt)
}
