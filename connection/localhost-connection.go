package connection

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"sync"
	"time"

	"github.com/meteocima/virtual-server/tailor"
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

/*
var tailCfg = tailor.Config{
	Poll:      true,
	Follow:    true,
	MustExist: false,
	ReOpen:    true,
	//Logger:    tailor.DiscardingLogger,
}
*/
func copyLines(proc *LocalProcess, w io.WriteCloser, outLogFile vpath.VirtualPath) {
	var logFile *os.File
	var err error = errors.New("empty")
	for err != nil {
		logFile, err = os.Open(outLogFile.Path)
		if os.IsNotExist(err) {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: copyLines error: (os.Open `%s`\n): %s", outLogFile.Path, err.Error())
			return
		}
	}

	go func() {
		tailProc := tailor.New(logFile, w, 1024)
		errs := tailProc.Start()
		proc.cmd.Wait()
		tailProc.Stop()
		w.Close()
		err := <-errs
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: copyLines error (reading lines from `%s`): %s\n", outLogFile.Path, err.Error())
		}
	}()

}

// Run ...
func (conn *LocalConnection) Run(command vpath.VirtualPath, args []string, options ...RunOptions) (Process, error) {
	//fmt.Println(strings.Repeat("*", 20))
	//fmt.Println("EXECUTING", command.Path)
	//fmt.Println(strings.Repeat("*", 20))

	cmd := exec.Command(command.Path, args...)
	process := &LocalProcess{
		cmd:              cmd,
		streamsCompleted: &sync.WaitGroup{},
	}

	errLogFile := vpath.Stderr
	outLogFile := vpath.Stdout
	if len(options) > 0 {
		if options[0].OutFromLog.Host != "" {
			outLogFile = &options[0].OutFromLog
		}
		if options[0].ErrFromLog.Host != "" {
			errLogFile = &options[0].ErrFromLog
		}
	}

	if outLogFile != vpath.Stdout {

		output, pwrite := io.Pipe()
		process.stdout = output

		go copyLines(process, pwrite, *outLogFile)
	} else {

		output, err := cmd.StdoutPipe()
		if err != nil {
			return nil, fmt.Errorf("Run `%s`: StdoutPipe error: %w", command, err)
		}
		process.stdout = output

	}

	if errLogFile != vpath.Stderr {

		output, pwrite := io.Pipe()
		process.stderr = output

		go copyLines(process, pwrite, *errLogFile)
	} else {
		output, err := cmd.StderrPipe()
		if err != nil {
			return nil, fmt.Errorf("Run `%s`: StderrPipe error: %w", command, err)
		}
		process.stderr = output
	}

	if len(options) > 0 {
		cmd.Dir = options[0].Cwd.Path
	}

	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("Run `%s`: Start error: %w", command, err)
	}

	return process, nil

}
