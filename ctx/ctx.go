package ctx

import (
	"bufio"
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

// Run ...
func (ctx *Context) Run(cwd vpath.VirtualPath, logFile vpath.VirtualPath, command string, args ...string) {
	if ctx.Err != nil {
		return
	}
	/*Logf("\tRun %s %s\n", command, args)
	cmd := exec.Command(command, args...)
	cmd.Dir = ctx.Root.JoinP(cwd).String()

	if logFile != "" {
		err := os.Remove(ctx.Root.JoinP(logFile).String())
		if err != nil && !os.IsNotExist(err) {
			ctx.Err = fmt.Errorf("Run `%s`: Remove error: %w", command, err)
			return
		}
	}

	output, pwrite := io.Pipe()

	var tailProc *tail.Tail
	if logFile == "" {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			ctx.Err = fmt.Errorf("Run `%s`: StdoutPipe error: %w", command, err)
			return
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			ctx.Err = fmt.Errorf("Run `%s`: StderrPipe error: %w", command, err)
			return
		}

		go func() {
			done := sync.WaitGroup{}
			done.Add(2)
			go func() {
				io.Copy(pwrite, stdout)
				done.Done()
			}()

			go func() {
				io.Copy(pwrite, stderr)
				done.Done()
			}()

			done.Wait()

			pwrite.Close()
		}()

	} else {

		tail, err := tail.TailFile(ctx.Root.JoinP(logFile).String(), tail.Config{
			Follow:    true,
			MustExist: false,
			ReOpen:    true,
		})

		if err != nil {
			ctx.Err = fmt.Errorf("Run `%s`: TailFile error: %w", command, err)
			return
		}
		tailProc = tail

		go func() {
			for l := range tail.Lines {
				pwrite.Write([]byte(l.Text + "\n"))
				if l.Err != nil {
					ctx.Err = fmt.Errorf("Run `%s`: TailFile error (lines): %w", command, err)
					break
				}
			}
			pwrite.Close()
		}()

	}

	err := cmd.Start()
	if ctx.Err != nil {
		ctx.Err = fmt.Errorf("Run `%s`: Start error: %w", command, err)
		return
	}

	go func() {
		stdoutBuff := bufio.NewReader(output)
		line, _, err := stdoutBuff.ReadLine()
		for line != nil {
			line, _, err = stdoutBuff.ReadLine()
			if err != nil && err != io.EOF {
				ctx.Err = fmt.Errorf("Run `%s`: ReadLine error: %w", command, err)
			}
			fmt.Println(string(line))
		}
	}()
	err = cmd.Wait()

	if err != nil {
		ctx.Err = fmt.Errorf("Run `%s`: Wait error: %w", command, err)
	}

	if tailProc != nil {
		tailProc.Stop()
	}
	*/
}
