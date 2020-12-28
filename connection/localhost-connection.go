package connection

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"

	"github.com/meteocima/virtual-server/vpath"
)

// LocalConnection ...
type LocalConnection struct{}

// HostName ...
func (conn *LocalConnection) HostName() string {
	return "localhost"
}

// OpenReader ...
func (conn *LocalConnection) OpenReader(file vpath.VirtualPath) (io.ReadCloser, error) {
	freader, err := os.Open(file.Path)
	return freader, err
}

// OpenWriter ...
func (conn *LocalConnection) OpenWriter(file vpath.VirtualPath) (io.WriteCloser, error) {
	fwriter, err := os.OpenFile(file.Path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0664))
	return fwriter, err
}

// ReadDir ...
func (conn *LocalConnection) ReadDir(dir vpath.VirtualPath) (vpath.VirtualPathList, error) {
	files, err := ioutil.ReadDir(dir.Path)
	if err != nil {
		return nil, fmt.Errorf("ReadDir `%s`: ioutil.ReadDir: %w", dir.String(), err)
	}
	filenames := make(vpath.VirtualPathList, len(files))
	for i, f := range files {
		filenames[i] = dir.Join(f.Name())
	}
	sort.Sort(filenames)
	return filenames, nil
}

// Open ...
func (conn *LocalConnection) Open() error { return nil }

// Close ...
func (conn *LocalConnection) Close() error { return nil }

// Stat ...
func (conn *LocalConnection) Stat(path vpath.VirtualPath) (os.FileInfo, error) {
	return os.Stat(path.Path)
}

// Link ...
func (conn *LocalConnection) Link(source, target vpath.VirtualPath) error {
	return os.Symlink(source.Path, target.Path)
}

// MkDir ...
func (conn *LocalConnection) MkDir(dir vpath.VirtualPath) error {
	err := os.MkdirAll(dir.Path, os.FileMode(0775))
	if err != nil {
		return fmt.Errorf("Error: MkDir `%s`: os.MkdirAll: %w", dir.String(), err)
	}
	return nil
}

// RmDir ...
func (conn *LocalConnection) RmDir(dir vpath.VirtualPath) error {
	err := os.RemoveAll(dir.Path)
	if err != nil {
		return fmt.Errorf("RmDir `%s`: os.RemoveAll: %w", dir.String(), err)
	}
	return nil
}

// RmFile ...
func (conn *LocalConnection) RmFile(file vpath.VirtualPath) error {
	err := os.Remove(file.Path)
	if err != nil {
		return fmt.Errorf("RmFile `%s`: os.Remove: %w", file.String(), err)
	}
	return nil
}

// LocalProcess ...
type LocalProcess struct {
	cmd              *exec.Cmd
	stdout           io.Reader
	stderr           io.Reader
	combinedOutput   io.Reader
	streamsCompleted *sync.WaitGroup
}

// CombinedOutput ...
func (proc *LocalProcess) CombinedOutput() io.Reader {
	proc.streamsCompleted.Add(1)
	combined, combinedWriter := io.Pipe()

	done := sync.WaitGroup{}
	done.Add(2)

	go func() {
		io.Copy(combinedWriter, proc.stdout)
		done.Done()
	}()

	go func() {
		io.Copy(combinedWriter, proc.stderr)
		done.Done()
	}()

	go func() {
		done.Wait()
		combinedWriter.Close()
		proc.streamsCompleted.Done()
	}()

	return combined
}

// Kill ...
func (proc *LocalProcess) Kill() error {
	return nil
}

// Stdin ...
func (proc *LocalProcess) Stdin() io.Writer {
	return nil
}

// Stdout ...
func (proc *LocalProcess) Stdout() io.Reader {
	proc.streamsCompleted.Add(1)
	processStdout, processStdoutWriter := io.Pipe()

	go func() {
		io.Copy(processStdoutWriter, proc.stdout)
		processStdoutWriter.Close()
		proc.streamsCompleted.Done()
	}()

	return processStdout
}

// Stderr ...
func (proc *LocalProcess) Stderr() io.Reader {
	proc.streamsCompleted.Add(1)
	processStderr, processStderrWriter := io.Pipe()

	go func() {
		io.Copy(processStderrWriter, proc.stderr)
		processStderrWriter.Close()
		proc.streamsCompleted.Done()
	}()

	return processStderr
}

// Wait ...
func (proc *LocalProcess) Wait() (int, error) {
	proc.streamsCompleted.Wait()
	err := proc.cmd.Wait()
	return proc.cmd.ProcessState.ExitCode(), err
}

// Run ...
func (conn *LocalConnection) Run(command vpath.VirtualPath, args []string, options ...RunOptions) (Process, error) {
	fmt.Println(strings.Repeat("*", 20))
	fmt.Println("EXECUTING", command.Path)
	fmt.Println(strings.Repeat("*", 20))

	cmd := exec.Command(command.Path, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Run `%s`: StdoutPipe error: %w", command, err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("Run `%s`: StderrPipe error: %w", command, err)
	}

	process := &LocalProcess{
		stdout:           stdout,
		stderr:           stderr,
		cmd:              cmd,
		streamsCompleted: &sync.WaitGroup{},
	}

	if len(options) > 0 {
		cmd.Dir = options[0].Cwd.Path
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("Run `%s`: Start error: %w", command, err)
	}

	return process, nil

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
