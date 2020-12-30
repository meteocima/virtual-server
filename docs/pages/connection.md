{{ useLayout(".layout.njk") }}
{{ title("CIMA virtual-server") }}
{{ subtitle("connection package") }}

# [virtual-server](./index) ⟶ {{ meta.subtitle }}





## Usage

#### func  NewPath

```go
func NewPath(cn Connection, path string, pathArgs ...interface{}) vpath.VirtualPath
```

#### type Connection

```go
type Connection interface {
	HostName() string
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
	Run(command vpath.VirtualPath, args []string, options ...RunOptions) (Process, error)
}
```

Connection ...

#### func  FindHost

```go
func FindHost(name string) Connection
```
FindHost ...

#### type CopyOptions

```go
type CopyOptions struct {
}
```

CopyOptions ...

#### type LocalConnection

```go
type LocalConnection struct{}
```

LocalConnection ...

#### func (*LocalConnection) Close

```go
func (conn *LocalConnection) Close() error
```
Close ...

#### func (*LocalConnection) HostName

```go
func (conn *LocalConnection) HostName() string
```
HostName ...

#### func (*LocalConnection) Link

```go
func (conn *LocalConnection) Link(source, target vpath.VirtualPath) error
```
Link ...

#### func (*LocalConnection) MkDir

```go
func (conn *LocalConnection) MkDir(dir vpath.VirtualPath) error
```
MkDir ...

#### func (*LocalConnection) Open

```go
func (conn *LocalConnection) Open() error
```
Open ...

#### func (*LocalConnection) OpenReader

```go
func (conn *LocalConnection) OpenReader(file vpath.VirtualPath) (io.ReadCloser, error)
```
OpenReader ...

#### func (*LocalConnection) OpenWriter

```go
func (conn *LocalConnection) OpenWriter(file vpath.VirtualPath) (io.WriteCloser, error)
```
OpenWriter ...

#### func (*LocalConnection) ReadDir

```go
func (conn *LocalConnection) ReadDir(dir vpath.VirtualPath) (vpath.VirtualPathList, error)
```
ReadDir ...

#### func (*LocalConnection) RmDir

```go
func (conn *LocalConnection) RmDir(dir vpath.VirtualPath) error
```
RmDir ...

#### func (*LocalConnection) RmFile

```go
func (conn *LocalConnection) RmFile(file vpath.VirtualPath) error
```
RmFile ...

#### func (*LocalConnection) Run

```go
func (conn *LocalConnection) Run(command vpath.VirtualPath, args []string, options ...RunOptions) (Process, error)
```
Run ...

#### func (*LocalConnection) Stat

```go
func (conn *LocalConnection) Stat(path vpath.VirtualPath) (os.FileInfo, error)
```
Stat ...

#### type LocalProcess

```go
type LocalProcess struct {
}
```

LocalProcess ...

#### func (*LocalProcess) CombinedOutput

```go
func (proc *LocalProcess) CombinedOutput() io.Reader
```
CombinedOutput ...

#### func (*LocalProcess) Kill

```go
func (proc *LocalProcess) Kill() error
```
Kill ...

#### func (*LocalProcess) Stderr

```go
func (proc *LocalProcess) Stderr() io.Reader
```
Stderr ...

#### func (*LocalProcess) Stdin

```go
func (proc *LocalProcess) Stdin() io.Writer
```
Stdin ...

#### func (*LocalProcess) Stdout

```go
func (proc *LocalProcess) Stdout() io.Reader
```
Stdout ...

#### func (*LocalProcess) Wait

```go
func (proc *LocalProcess) Wait() (int, error)
```
Wait ...

#### type MoveOptions

```go
type MoveOptions struct {
}
```

MoveOptions ...

#### type Process

```go
type Process interface {
	Kill() error
	// Stdin, is an io.Writer that will be used
	// to send data to process stdin
	Stdin() io.Writer

	// Stdin, if set, is an io.Reader that will be used
	// to read data from process stdout
	Stdout() io.Reader

	// Stderr, if set, is an io.Reader that will be used
	// to read data from process stdout
	Stderr() io.Reader

	// CombinedOutput returns an io.Reader that reads
	// the combined output and error streams of the process
	CombinedOutput() io.Reader

	// Wait expects the process to terminate
	// and return the exit code.
	Wait() (int, error)
}
```

Process represents a running process

#### type RunOptions

```go
type RunOptions struct {
	// OutFromLog if sets, output is read from a file
	// instead of from the process stdout
	OutFromLog vpath.VirtualPath

	// OutFromLog if sets, output is read from a file
	// instead of from the process stderr
	ErrFromLog vpath.VirtualPath

	// Cwd is set the work directory in which the process will be executed.
	Cwd vpath.VirtualPath

	// Stdin, if set, is an io.Reader that will be used
	// as process Stdin.
	// If nil, a pipe to `Process.Stdin` member is created
	// and used.
	Stdin *io.Reader

	// Stdout, if set, is an io.Writer that will be used
	// as process Stdout
	// If nil, a pipe to `Process.Stdout` member is created
	// and used.
	Stdout *io.Writer

	// Stderr, if set, is an io.Writer that will be used
	// as process Stderr.
	// If nil, a pipe to `Process.Stdout` member is created
	// and used.
	Stderr *io.Writer
}
```

RunOptions ...

#### type SSHConnection

```go
type SSHConnection struct {
	BackupHosts []string
	Name        string

	Host    string
	Port    int
	User    string
	KeyPath string
}
```

SSHConnection ...

#### func (*SSHConnection) Close

```go
func (conn *SSHConnection) Close() error
```
Close ...

#### func (*SSHConnection) HostName

```go
func (conn *SSHConnection) HostName() string
```
HostName ...

#### func (*SSHConnection) Link

```go
func (conn *SSHConnection) Link(source, target vpath.VirtualPath) error
```
Link ...

#### func (*SSHConnection) MkDir

```go
func (conn *SSHConnection) MkDir(dir vpath.VirtualPath) error
```
MkDir ...

#### func (*SSHConnection) Open

```go
func (conn *SSHConnection) Open() error
```
Open ...

#### func (*SSHConnection) OpenReader

```go
func (conn *SSHConnection) OpenReader(file vpath.VirtualPath) (io.ReadCloser, error)
```
OpenReader ...

#### func (*SSHConnection) OpenWriter

```go
func (conn *SSHConnection) OpenWriter(file vpath.VirtualPath) (io.WriteCloser, error)
```
OpenWriter ...

#### func (*SSHConnection) ReadDir

```go
func (conn *SSHConnection) ReadDir(dir vpath.VirtualPath) (vpath.VirtualPathList, error)
```
ReadDir ...

#### func (*SSHConnection) RmDir

```go
func (conn *SSHConnection) RmDir(dir vpath.VirtualPath) error
```
RmDir ...

#### func (*SSHConnection) RmFile

```go
func (conn *SSHConnection) RmFile(file vpath.VirtualPath) error
```
RmFile ...

#### func (*SSHConnection) Run

```go
func (conn *SSHConnection) Run(command vpath.VirtualPath, args []string, options ...RunOptions) (Process, error)
```
Run ...

#### func (*SSHConnection) Stat

```go
func (conn *SSHConnection) Stat(path vpath.VirtualPath) (os.FileInfo, error)
```
Stat ...

#### type SSHProcess

```go
type SSHProcess struct {
}
```

SSHProcess ...

#### func (*SSHProcess) CombinedOutput

```go
func (proc *SSHProcess) CombinedOutput() io.Reader
```
CombinedOutput ...

#### func (*SSHProcess) Kill

```go
func (proc *SSHProcess) Kill() error
```
Kill ...

#### func (*SSHProcess) Stderr

```go
func (proc *SSHProcess) Stderr() io.Reader
```
Stderr ...

#### func (*SSHProcess) Stdin

```go
func (proc *SSHProcess) Stdin() io.Writer
```
Stdin ...

#### func (*SSHProcess) Stdout

```go
func (proc *SSHProcess) Stdout() io.Reader
```
Stdout ...

#### func (*SSHProcess) Wait

```go
func (proc *SSHProcess) Wait() (int, error)
```
Wait ...
