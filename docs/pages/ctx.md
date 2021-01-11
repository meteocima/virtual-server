{{ useLayout(".layout.njk") }}
{{ title("CIMA virtual-server") }}
{{ subtitle("ctx package") }}

# [virtual-server](./index) ⟶ {{ meta.subtitle }}





## Usage

#### type Context

```go
type Context struct {
	Err             error
	RunningFunction string
	RunningTask     string
	Log             io.Writer
	DetailLog       io.Writer
}
```

Context abstract a set of operations on one or multiple FileSystem instances
that fails or succeed as a whole

#### func (*Context) ContextFailed

```go
func (ctx *Context) ContextFailed(offendingFunc string, err error)
```
ContextFailed ...

#### func (*Context) Copy

```go
func (ctx *Context) Copy(from, to vpath.VirtualPath)
```
Copy ...

#### func (*Context) Exec

```go
func (ctx *Context) Exec(command vpath.VirtualPath, args []string, options ...connection.RunOptions)
```
Exec ...

#### func (*Context) Exists

```go
func (ctx *Context) Exists(file vpath.VirtualPath) bool
```
Exists ...

#### func (*Context) IsFile

```go
func (ctx *Context) IsFile(file vpath.VirtualPath) bool
```
IsFile ...

#### func (*Context) Link

```go
func (ctx *Context) Link(from, to vpath.VirtualPath)
```
Link ...

#### func (*Context) LogF

```go
func (ctx *Context) LogF(msg string, args ...interface{})
```
LogF ...

#### func (*Context) MkDir

```go
func (ctx *Context) MkDir(dir vpath.VirtualPath)
```
MkDir ...

#### func (*Context) Move

```go
func (ctx *Context) Move(from, to vpath.VirtualPath)
```
Move ...

#### func (*Context) ReadDir

```go
func (ctx *Context) ReadDir(dir vpath.VirtualPath) vpath.VirtualPathList
```
ReadDir ...

#### func (*Context) ReadString

```go
func (ctx *Context) ReadString(file vpath.VirtualPath) string
```
ReadString ...

#### func (*Context) RmDir

```go
func (ctx *Context) RmDir(dir vpath.VirtualPath)
```
RmDir ...

#### func (*Context) RmFile

```go
func (ctx *Context) RmFile(file vpath.VirtualPath)
```
RmFile ...

#### func (*Context) Run

```go
func (ctx *Context) Run(command vpath.VirtualPath, args []string, options ...connection.RunOptions) connection.Process
```
Run ...

#### func (*Context) SetRunning

```go
func (ctx *Context) SetRunning(msg string, args ...interface{}) func()
```
SetRunning ...

#### func (*Context) SetTask

```go
func (ctx *Context) SetTask(msg string, args ...interface{}) func()
```
SetTask ...

#### func (*Context) WriteString

```go
func (ctx *Context) WriteString(file vpath.VirtualPath, content string)
```
WriteString ...
