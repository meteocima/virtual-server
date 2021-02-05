package ctx

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/meteocima/virtual-server/connection"
	"github.com/meteocima/virtual-server/vpath"
)

// Context abstract a set of operations
// on one or multiple FileSystem instances
// that fails or succeed as a whole
type Context struct {
	Err             error
	runningFunction string
	ID              string

	stdout        io.Writer
	stderr        io.Writer
	infoChannel   chan string
	detailChannel chan string
	logCompleted  chan struct{}
	running       bool
	runningLock   *sync.Mutex
	level         LogLevel
}

// New ...
func New(stdout io.Writer, stderr io.Writer) *Context {
	ctx := Context{
		ID:           "ANON",
		stdout:       stdout,
		stderr:       stderr,
		level:        LevelDebug,
		logCompleted: make(chan struct{}),
		runningLock:  &sync.Mutex{},
	}
	ctx.startLogWriter()
	return &ctx
}

// ContextFailed ...
func (ctx *Context) ContextFailed(offendingFunc string, err error) {
	ctx.SetContextFailed("%s: %s: %w", ctx.runningFunction, offendingFunc, err)
}

// SetContextFailed ...
func (ctx *Context) SetContextFailed(format string, args ...interface{}) {
	ctx.Err = fmt.Errorf(format, args...)
}

// setRunningFunction ...
func (ctx *Context) setRunningFunction(msg string, args ...interface{}) func() {
	ctx.runningFunction = fmt.Sprintf(msg, args...)

	ctx.LogDetail("\t⟶\t%s\n", ctx.runningFunction)

	return func() {
		ctx.runningFunction = ""
	}
}

// IsFile ...
func (ctx *Context) IsFile(file vpath.VirtualPath) bool {
	if ctx.Err != nil {
		return false
	}
	defer ctx.setRunningFunction("IsFile `%s`", file.String())()

	conn, err := connection.FindHost(file.Host)
	if err != nil {
		ctx.ContextFailed("connection.FindHost", err)
		return false
	}

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
	defer ctx.setRunningFunction("Exists `%s`", file.String())()

	conn, err := connection.FindHost(file.Host)
	if err != nil {
		ctx.ContextFailed("connection.FindHost", err)
		return false
	}

	_, err = conn.Stat(file)

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
	defer ctx.setRunningFunction("ReadDir `%s`", dir.String())()

	conn, err := connection.FindHost(dir.Host)
	if err != nil {
		ctx.ContextFailed("connection.FindHost", err)
		return nil
	}

	var files vpath.VirtualPathList
	files, err = conn.ReadDir(dir)
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
	defer ctx.setRunningFunction("Copy from `%s` to `%s`", from.String(), to.String())()

	fromConn, err := connection.FindHost(from.Host)
	if err != nil {
		ctx.ContextFailed("connection.FindHost", err)
		return
	}
	toConn, err := connection.FindHost(to.Host)
	if err != nil {
		ctx.ContextFailed("connection.FindHost", err)
		return
	}

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

// OpenWriter ...
func (ctx *Context) OpenWriter(file vpath.VirtualPath, content string) io.WriteCloser {
	if ctx.Err != nil {
		return nil
	}
	defer ctx.setRunningFunction("OpenWriter to `%s`", file.String())()

	toConn, err := connection.FindHost(file.Host)
	if err != nil {
		ctx.ContextFailed("connection.FindHost", err)
		return nil
	}

	writer, err := toConn.OpenWriter(file)
	if err != nil {
		ctx.ContextFailed("toConn.OpenWriter", err)
		return nil
	}

	return writer
}

// WriteString ...
func (ctx *Context) WriteString(file vpath.VirtualPath, content string) {
	if ctx.Err != nil {
		return
	}
	defer ctx.setRunningFunction("WriteString to `%s`", file.String())()

	toConn, err := connection.FindHost(file.Host)
	if err != nil {
		ctx.ContextFailed("connection.FindHost", err)
		return
	}

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
	defer ctx.setRunningFunction("ReadString from `%s`", file.String())()

	conn, err := connection.FindHost(file.Host)
	if err != nil {
		ctx.ContextFailed("connection.FindHost", err)
		return ""
	}

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
	defer ctx.setRunningFunction("Link from %s to %s", from.String(), to.String())()

	conn, err := connection.FindHost(from.Host)
	if err != nil {
		ctx.ContextFailed("connection.FindHost", err)
		return
	}
	err = conn.Link(from, to)
	if err != nil {
		ctx.ContextFailed("conn.Link", err)
	}
}

// MkDir ...
func (ctx *Context) MkDir(dir vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}
	defer ctx.setRunningFunction("MkDir %s", dir.String())()

	conn, err := connection.FindHost(dir.Host)
	if err != nil {
		ctx.ContextFailed("connection.FindHost", err)
		return
	}

	err = conn.MkDir(dir)
	if err != nil {
		ctx.ContextFailed("conn.MkDir", err)
	}
}

// RmDir ...
func (ctx *Context) RmDir(dir vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}
	defer ctx.setRunningFunction("RmDir %s", dir.String())()

	conn, err := connection.FindHost(dir.Host)
	if err != nil {
		ctx.ContextFailed("connection.FindHost", err)
		return
	}

	err = conn.RmDir(dir)
	if err != nil {
		ctx.ContextFailed("conn.RmDir", err)
	}
}

// RmFile ...
func (ctx *Context) RmFile(file vpath.VirtualPath) {
	if ctx.Err != nil {
		return
	}
	defer ctx.setRunningFunction("RmFile %s", file.String())()

	conn, err := connection.FindHost(file.Host)
	if err != nil {
		ctx.ContextFailed("connection.FindHost", err)
		return
	}
	err = conn.RmFile(file)
	if err != nil {
		ctx.ContextFailed("conn.RmFile", err)
	}
}

// Exec ...
func (ctx *Context) Exec(command vpath.VirtualPath, args []string, options *connection.RunOptions) {
	if options == nil {
		options = &connection.RunOptions{}
	}

	/*devNull, err := os.Open(os.DevNull)
	if err != nil {
		ctx.ContextFailed("os.Open(os.DevNull)", err)
	}*/

	options.Stdout = ctx.stdout
	options.Stderr = ctx.stderr

	ctx.LogInfo("START %s %s", command.String(), strings.Join(args, " "))
	p := ctx.Run(command, args, *options)
	p.Wait()
	if ctx.Err == nil {
		ctx.LogInfo("COMPLETED OK %s", command.String())
	}
}

// Run ...
func (ctx *Context) Run(command vpath.VirtualPath, args []string, options connection.RunOptions) connection.Process {
	if ctx.Err != nil {
		return nil
	}
	defer ctx.setRunningFunction("Run %s %s", command.String(), strings.Join(args, " "))()

	fmt.Println("find host ", command.Host)
	conn, err := connection.FindHost(command.Host)
	if err != nil {
		fmt.Println("err host ", err)

		ctx.ContextFailed("connection.FindHost", err)
		return nil
	}

	fmt.Println("using host ", conn)

	proc, err := conn.Run(command, args, options)
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
		return "DETAIL"
	case LevelDebug:
		return "DEBUG"
	default:
		return "WRONGLEVEL"
	}
}

// Close ...
func (ctx *Context) Close() {
	ctx.runningLock.Lock()
	ctx.running = false
	ctx.runningLock.Unlock()
	<-ctx.logCompleted
	close(ctx.infoChannel)
	close(ctx.detailChannel)
}

func (ctx *Context) startLogWriter() {
	ctx.infoChannel = make(chan string, 1024)
	ctx.detailChannel = make(chan string, 1024)
	ctx.runningLock.Lock()
	ctx.running = true
	ctx.runningLock.Unlock()
	go func() {
		defer close(ctx.logCompleted)
		for {
			var chunk string
			select {
			case chunk = <-ctx.infoChannel:
				fmt.Fprintf(ctx.stdout, chunk)
			case chunk = <-ctx.detailChannel:
				fmt.Fprintf(ctx.stderr, chunk)
			case <-time.After(100 * time.Millisecond):
				ctx.runningLock.Lock()
				running := ctx.running
				ctx.runningLock.Unlock()
				if !running {

					return
				}
			}

		}

	}()
}

func (ctx *Context) logWrite(msgLevel LogLevel, msgText string, args []interface{}) {
	if msgLevel > ctx.level {
		return
	}
	channel := ctx.infoChannel
	if msgLevel >= LevelDetail {
		channel = ctx.detailChannel
	}

	channel <- fmt.Sprintf(msgLevel.String()+": "+ctx.ID+": "+msgText+"\n", args...)
}

// OutPrintf ...
func (ctx *Context) OutPrintf(format string, args ...interface{}) {
	channel := ctx.infoChannel
	channel <- fmt.Sprintf(format, args...)
}

// ErrPrintf ...
func (ctx *Context) ErrPrintf(format string, args ...interface{}) {
	channel := ctx.detailChannel
	channel <- fmt.Sprintf(format, args...)
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
	ctx.logWrite(LevelDetail, msg, args)
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
