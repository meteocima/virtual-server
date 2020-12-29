package ctx

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	connection "github.com/meteocima/virtual-server/connection"
	"github.com/meteocima/virtual-server/vpath"
)

// Context abstract a set of operations
// on one or multiple FileSystem instances
// that fails or succeed as a whole
type Context struct {
	Err             error
	RunningFunction string
	RunningTask     string
}

// ContextFailed ...
func (ctx *Context) ContextFailed(offendingFunc string, err error) {
	ctx.Err = fmt.Errorf("Error: %s: %s: %s error: %w", ctx.RunningTask, ctx.RunningFunction, offendingFunc, err)
}

// SetRunning ...
func (ctx *Context) SetRunning(msg string, args ...interface{}) func() {
	ctx.RunningFunction = fmt.Sprintf(msg, args...)
	fmt.Printf("\t⟶\t%s\n", ctx.RunningFunction)
	return func() {
		ctx.RunningFunction = ""
	}
}

// Exists ...
func (ctx *Context) Exists(file vpath.VirtualPath) bool {
	if ctx.Err != nil {
		return false
	}
	defer ctx.SetRunning("Exists `%s`", file.StringRel())()

	conn := connection.FindHost(file.Host)

	_, err := conn.Stat(file)

	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		ctx.ContextFailed("connection.Stat", err)
		return false
	}

	return true
}

// ReadDir ...
func (ctx *Context) ReadDir(dir vpath.VirtualPath) vpath.VirtualPathList {
	if ctx.Err != nil {
		return vpath.VirtualPathList{}
	}
	defer ctx.SetRunning("ReadDir `%s`", dir.StringRel())()

	conn := connection.FindHost(dir.Host)
	var files vpath.VirtualPathList
	files, err := conn.ReadDir(dir)
	if err != nil {
		ctx.ContextFailed("connection.ReadDir", err)
		return nil
	}
	return files
}

// Copy ...
func (ctx *Context) Copy(from, to vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}
	defer ctx.SetRunning("Copy from `%s` to `%s`", from.StringRel(), to.StringRel())()

	fromConn := connection.FindHost(from.Host)
	toConn := connection.FindHost(to.Host)

	reader, err := fromConn.OpenReader(from)
	if err != nil {
		ctx.ContextFailed("fromConn.OpenReader", err)
		return
	}
	defer reader.Close()

	writer, err := toConn.OpenWriter(to)
	if err != nil {
		ctx.ContextFailed("toConn.OpenWriter", err)
		return
	}
	defer writer.Close()

	_, err = io.Copy(writer, reader)
	if err != nil {
		ctx.ContextFailed("io.Copy", err)
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
	defer ctx.SetRunning("WriteString to `%s`", file.StringRel())()

	toConn := connection.FindHost(file.Host)

	writer, err := toConn.OpenWriter(file)
	if err != nil {
		ctx.ContextFailed("toConn.OpenWriter", err)
		return
	}

	defer writer.Close()

	_, err = writer.Write([]byte(content))
	if err != nil {
		ctx.ContextFailed("writer.Write", err)
		return
	}
}

// ReadString ...
func (ctx *Context) ReadString(file vpath.VirtualPath) string {
	if ctx.Err != nil {
		return ""
	}
	defer ctx.SetRunning("ReadString from `%s`", file.StringRel())()

	conn := connection.FindHost(file.Host)

	reader, err := conn.OpenReader(file)
	if err != nil {
		ctx.ContextFailed("conn.OpenReader", err)
		return ""
	}

	defer reader.Close()

	bufReader := bufio.NewReader(reader)

	buf, err := ioutil.ReadAll(bufReader)
	if err != nil {
		ctx.ContextFailed("ioutil.ReadAll", err)
		return ""
	}
	return string(buf)
}

// Link ...
func (ctx *Context) Link(from, to vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}
	defer ctx.SetRunning("Link from %s to %s", from.StringRel(), to.StringRel())()

	conn := connection.FindHost(from.Host)
	err := conn.Link(from, to)
	if err != nil {
		ctx.ContextFailed("conn.Link", err)
	}
}

// MkDir ...
func (ctx *Context) MkDir(dir vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}
	defer ctx.SetRunning("MkDir %s", dir.StringRel())()

	conn := connection.FindHost(dir.Host)
	err := conn.MkDir(dir)
	if err != nil {
		ctx.ContextFailed("conn.MkDir", err)
	}
}

// RmDir ...
func (ctx *Context) RmDir(dir vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}
	defer ctx.SetRunning("RmDir %s", dir.StringRel())()

	conn := connection.FindHost(dir.Host)
	err := conn.RmDir(dir)
	if err != nil {
		ctx.ContextFailed("conn.RmDir", err)
	}
}

// RmFile ...
func (ctx *Context) RmFile(file vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}
	defer ctx.SetRunning("RmFile %s", file.StringRel())()

	conn := connection.FindHost(file.Host)
	err := conn.RmFile(file)
	if err != nil {
		ctx.ContextFailed("conn.RmFile", err)
	}
}

// LogF ...
func (ctx *Context) LogF(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}

// Exec ...
func (ctx *Context) Exec(command vpath.VirtualPath, args []string, options ...connection.RunOptions) {
	p := ctx.Run(command, args, options...)
	if p != nil {
		io.Copy(os.Stderr, p.CombinedOutput())
		p.Wait()
	}
}

// Run ...
func (ctx *Context) Run(command vpath.VirtualPath, args []string, options ...connection.RunOptions) connection.Process {
	if ctx.Err != nil {
		return nil
	}
	defer ctx.SetRunning("Run %s %s", command.StringRel(), strings.Join(args, " "))()

	conn := connection.FindHost(command.Host)
	proc, err := conn.Run(command, args, options...)
	if err != nil {
		ctx.ContextFailed("conn.Run", err)
		return nil
	}

	return proc
}
