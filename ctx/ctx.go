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
	runningFunction string
	infoLog         io.Writer
	detailLog       io.Writer
	level           LogLevel
}

// New ...
func New(infoLog io.Writer, detailLog io.Writer) Context {
	return Context{
		infoLog:   infoLog,
		detailLog: detailLog,
		level:     LevelDebug,
	}
}

// ContextFailed ...
func (ctx *Context) ContextFailed(offendingFunc string, err error) {
	ctx.Err = fmt.Errorf("%s: %s: %w", ctx.runningFunction, offendingFunc, err)
}

// setRunningFunction ...
func (ctx *Context) setRunningFunction(msg string, args ...interface{}) func() {
	ctx.runningFunction = fmt.Sprintf(msg, args...)

	if ctx.infoLog != nil {
		fmt.Fprintf(ctx.infoLog, "\t⟶\t%s\n", ctx.runningFunction)
	}
	return func() {
		ctx.runningFunction = ""
	}
}

/*

// SetTask ...
func (ctx *Context) SetTask(msg string, args ...interface{}) func() {
	ctx.RunningTask = fmt.Sprintf(msg, args...)
	if ctx.infoLog != nil {
		fmt.Fprintf(ctx.infoLog, "\n\n# START: %s\n", ctx.RunningTask)
	}
	return func() {
		if ctx.Err == nil && ctx.infoLog != nil {
			fmt.Fprintf(ctx.infoLog, "# COMPLETED SUCCESSUFULLY: %s\n", ctx.RunningTask)
		}
		ctx.RunningTask = ""
	}
}
*/

// IsFile ...
func (ctx *Context) IsFile(file vpath.VirtualPath) bool {
	if ctx.Err != nil {
		return false
	}
	defer ctx.setRunningFunction("IsFile `%s`", file.StringRel())()

	conn := connection.FindHost(file.Host)

	info, err := conn.Stat(file)

	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		ctx.ContextFailed("connection.Stat", err)
		return false
	}

	return !info.IsDir()
}

// Exists ...
func (ctx *Context) Exists(file vpath.VirtualPath) bool {
	if ctx.Err != nil {
		return false
	}
	defer ctx.setRunningFunction("Exists `%s`", file.StringRel())()

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
	defer ctx.setRunningFunction("ReadDir `%s`", dir.StringRel())()

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
	defer ctx.setRunningFunction("Copy from `%s` to `%s`", from.StringRel(), to.StringRel())()

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

	bufIn := bufio.NewReaderSize(reader, 1024*1024)
	bufOut := bufio.NewWriterSize(writer, 1024*1024)

	_, err = io.Copy(bufOut, bufIn)
	if err != nil {
		ctx.ContextFailed("io.Copy", err)
		return
	}

	err = bufOut.Flush()
	if err != nil {
		ctx.ContextFailed("bufOut.Flush", err)
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
	defer ctx.setRunningFunction("WriteString to `%s`", file.StringRel())()

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
	defer ctx.setRunningFunction("ReadString from `%s`", file.StringRel())()

	conn := connection.FindHost(file.Host)

	reader, err := conn.OpenReader(file)
	if err != nil {
		ctx.ContextFailed("conn.OpenReader", err)
		return ""
	}

	defer reader.Close()

	//bufReader := bufio.NewReader(reader)

	buf, err := ioutil.ReadAll(reader)
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
	defer ctx.setRunningFunction("Link from %s to %s", from.StringRel(), to.StringRel())()

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
	defer ctx.setRunningFunction("MkDir %s", dir.StringRel())()

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
	defer ctx.setRunningFunction("RmDir %s", dir.StringRel())()

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
	defer ctx.setRunningFunction("RmFile %s", file.StringRel())()

	conn := connection.FindHost(file.Host)
	err := conn.RmFile(file)
	if err != nil {
		ctx.ContextFailed("conn.RmFile", err)
	}
}

// Exec ...
func (ctx *Context) Exec(command vpath.VirtualPath, args []string, options ...connection.RunOptions) {
	p := ctx.Run(command, args, options...)
	if p != nil {
		if ctx.detailLog != nil {
			io.Copy(ctx.detailLog, p.CombinedOutput())
		}
		p.Wait()
	}
}

// Run ...
func (ctx *Context) Run(command vpath.VirtualPath, args []string, options ...connection.RunOptions) connection.Process {
	if ctx.Err != nil {
		return nil
	}
	defer ctx.setRunningFunction("Run %s %s", command.StringRel(), strings.Join(args, " "))()

	conn := connection.FindHost(command.Host)
	proc, err := conn.Run(command, args, options...)
	if err != nil {
		ctx.ContextFailed("conn.Run", err)
		return nil
	}

	return proc
}

// LogLevel is a type that represents
// the importance level of a log message
type LogLevel int

const (
	// LevelError identify error messages
	LevelError LogLevel = iota
	// LevelWarning identify Warning messages
	LevelWarning
	// LevelInfo identify Info messages
	LevelInfo
	// LevelDetail identify Detail messages
	LevelDetail
	// LevelDebug identify Debug messages
	LevelDebug
)

func (ll LogLevel) String() string {
	switch ll {
	case LevelError:
		return "ERROR"
	case LevelWarning:
		return "WARNING"
	case LevelInfo:
		return "INFO"
	case LevelDetail:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	default:
		return "WRONGLEVEL"
	}
}

func (ctx *Context) logWrite(msgLevel LogLevel, msgText string, args []interface{}) {
	if msgLevel > ctx.level {
		return
	}
	ErrStream := ctx.infoLog
	if msgLevel >= LevelDetail {
		ErrStream = ctx.detailLog
	}

	fmt.Fprintf(ErrStream, msgLevel.String()+": "+msgText+"\n", args...)
}

// SetLevel set the maximum
// level a message must have to be
// logged.
func (ctx *Context) SetLevel(value LogLevel) {
	ctx.level = value
}

// LogDebug prints a log string if
// the configured log level is
// equal or great than levelDebug
func (ctx *Context) LogDebug(msg string, args ...interface{}) {
	ctx.logWrite(LevelDebug, msg, args)
}

// LogInfo prints a log string if
// the configured log level is
// equal or great than levelInfo
func (ctx *Context) LogInfo(msg string, args ...interface{}) {
	ctx.logWrite(LevelInfo, msg, args)
}

// LogDetail prints a log string if
// the configured log level is
// equal or great than levelDetail
func (ctx *Context) LogDetail(msg string, args ...interface{}) {
	ctx.logWrite(LevelInfo, msg, args)
}

// LogWarning prints a log string if
// the configured log level is
// equal or great than levelWarning
func (ctx *Context) LogWarning(msg string, args ...interface{}) {
	ctx.logWrite(LevelWarning, msg, args)
}

// LogError prints a log string if
// the configured log level is
// equal or great than levelError
func (ctx *Context) LogError(msg string, args ...interface{}) {
	ctx.logWrite(LevelError, msg, args)
}
