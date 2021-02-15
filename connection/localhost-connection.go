package connection

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/meteocima/virtual-server/vpath"
)

// LocalConnection ...
type LocalConnection struct {
	name string
}

// Name ...
func (conn *LocalConnection) Name() string {
	return conn.name
}

// OpenReader ...
func (conn *LocalConnection) OpenReader(file vpath.VirtualPath) (io.ReadCloser, error) {
	freader, err := os.Open(file.Path)
	return freader, err
}

// Glob ...
func (conn *LocalConnection) Glob(pattern vpath.VirtualPath) (vpath.VirtualPathList, error) {

	files, err := filepath.Glob(pattern.Path)
	if err != nil {
		return nil, err
	}
	result := make(vpath.VirtualPathList, len(files))
	for idx, file := range files {
		result[idx] = vpath.New(pattern.Host, file)
	}
	return result, nil
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
	cmd       *exec.Cmd
	completed chan struct{}
	state     int
}

// Kill ...
func (proc *LocalProcess) Kill() error {
	return nil
}

// Wait ...
func (proc *LocalProcess) Wait() (int, error) {
	<-proc.completed
	return proc.state, nil
}

// Run ...
func (conn *LocalConnection) Run(command vpath.VirtualPath, args []string, options RunOptions) (Process, error) {

	cmd := exec.Command(command.Path, args...)

	process := &LocalProcess{
		cmd:       cmd,
		completed: make(chan struct{}),
	}

	if options.Stderr == nil {
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = options.Stderr
	}

	if options.Stdout == nil {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = options.Stdout
	}

	if options.Stdin == nil {
		cmd.Stdin = os.Stdin
	} else {
		cmd.Stdin = options.Stdin
	}

	if options.OutFromLog != nil {
		go copyLines(process, cmd.Stdout, *options.OutFromLog)
	}

	if options.ErrFromLog != nil {
		go copyLines(process, cmd.Stderr, *options.ErrFromLog)
	}

	cmd.Dir = options.Cwd.Path

	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("Run `%s`: Start error: %w", command, err)
	}

	go func() {
		state, err := cmd.Process.Wait()
		if err != nil {
			panic(err)
		}
		process.state = state.ExitCode()
		close(process.completed)
	}()

	return process, nil
}
