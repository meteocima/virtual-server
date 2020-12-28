package ctx

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	connection "github.com/meteocima/virtual-server/connection"
	"github.com/meteocima/virtual-server/vpath"
)

// Context abstract a set of operations
// on one or multiple FileSystem instances
// that fails or succeed as a whole
type Context struct {
	Err error
}

// Exists ...
func (ctx *Context) Exists(file vpath.VirtualPath) bool {
	if ctx.Err != nil {
		return false
	}
	conn := connection.FindHost(file.Host)

	_, err := conn.Stat(file)

	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		ctx.Err = err
		return false
	}

	return true
}

// ReadDir ...
func (ctx *Context) ReadDir(dir vpath.VirtualPath) vpath.VirtualPathList {
	if ctx.Err != nil {
		return vpath.VirtualPathList{}
	}
	conn := connection.FindHost(dir.Host)
	var files vpath.VirtualPathList
	files, ctx.Err = conn.ReadDir(dir)
	return files
}

// Copy ...
func (ctx *Context) Copy(from, to vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}

	fromConn := connection.FindHost(from.Host)
	toConn := connection.FindHost(to.Host)

	reader, err := fromConn.OpenReader(from)
	if err != nil {
		ctx.Err = err
		return
	}
	defer reader.Close()

	writer, err := toConn.OpenWriter(to)
	if err != nil {
		ctx.Err = err
		return
	}
	defer writer.Close()

	_, err = io.Copy(writer, reader)
	if err != nil {
		ctx.Err = err
		return
	}
}

// Move ...
func (ctx *Context) Move(from, to vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}

	ctx.Copy(from, to)
	ctx.RmFile(from)
}

// WriteString ...
func (ctx *Context) WriteString(file vpath.VirtualPath, content string) {
	if ctx.Err != nil {
		return
	}

	toConn := connection.FindHost(file.Host)

	writer, err := toConn.OpenWriter(file)
	if err != nil {
		ctx.Err = err
		return
	}

	defer writer.Close()

	_, ctx.Err = writer.Write([]byte(content))
}

// ReadString ...
func (ctx *Context) ReadString(file vpath.VirtualPath) string {
	if ctx.Err != nil {
		return ""
	}
	conn := connection.FindHost(file.Host)

	reader, err := conn.OpenReader(file)
	if err != nil {
		ctx.Err = err
		return ""
	}

	defer reader.Close()

	bufReader := bufio.NewReader(reader)

	buf, err := ioutil.ReadAll(bufReader)
	if err != nil {
		ctx.Err = err
		return ""
	}
	return string(buf)
}

// Link ...
func (ctx *Context) Link(from, to vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}
	conn := connection.FindHost(from.Host)
	ctx.Err = conn.Link(from, to)
}

// MkDir ...
func (ctx *Context) MkDir(dir vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}
	conn := connection.FindHost(dir.Host)
	ctx.Err = conn.MkDir(dir)
}

// RmDir ...
func (ctx *Context) RmDir(dir vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}
	conn := connection.FindHost(dir.Host)
	ctx.Err = conn.RmDir(dir)
}

// RmFile ...
func (ctx *Context) RmFile(file vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}
	conn := connection.FindHost(file.Host)
	ctx.Err = conn.RmFile(file)
}

// LogF ...
func (ctx *Context) LogF(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}

// Exec ...
func (ctx *Context) Exec(command vpath.VirtualPath, args []string, options ...connection.RunOptions) {
	p := ctx.Run(command, args, options...)

	io.Copy(os.Stdout, p.Stdout())

	p.Wait()
}

// Run ...
func (ctx *Context) Run(command vpath.VirtualPath, args []string, options ...connection.RunOptions) connection.Process {
	if ctx.Err != nil {
		return nil
	}

	conn := connection.FindHost(command.Host)
	proc, err := conn.Run(command, args, options...)
	if err != nil {
		ctx.Err = err
		return nil
	}

	return proc
}
