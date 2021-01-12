{{ useLayout(".layout.njk") }}
{{ title("CIMA virtual-server") }}
{{ subtitle("ctx package") }}

# [virtual-server](./index) ⟶ {{ meta.subtitle }}





## Usage

#### type Context

```go
type Context struct {
	Err error
}
```

Context abstract a set of operations on one or multiple FileSystem instances
that fails or succeed as a whole

#### func  New

```go
func New(infoLog io.Writer, detailLog io.Writer) Context
```
New ...

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

#### func (*Context) LogDebug

```go
func (ctx *Context) LogDebug(msg string, args ...interface{})
```
LogDebug prints a log string if the configured log level is equal or great than
levelDebug

#### func (*Context) LogDetail

```go
func (ctx *Context) LogDetail(msg string, args ...interface{})
```
LogDetail prints a log string if the configured log level is equal or great than
levelDetail

#### func (*Context) LogError

```go
func (ctx *Context) LogError(msg string, args ...interface{})
```
LogError prints a log string if the configured log level is equal or great than
levelError

#### func (*Context) LogInfo

```go
func (ctx *Context) LogInfo(msg string, args ...interface{})
```
LogInfo prints a log string if the configured log level is equal or great than
levelInfo

#### func (*Context) LogWarning

```go
func (ctx *Context) LogWarning(msg string, args ...interface{})
```
LogWarning prints a log string if the configured log level is equal or great
than levelWarning

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

#### func (*Context) SetLevel

```go
func (ctx *Context) SetLevel(value LogLevel)
```
SetLevel set the maximum level a message must have to be logged.

#### func (*Context) WriteString

```go
func (ctx *Context) WriteString(file vpath.VirtualPath, content string)
```
WriteString ...

#### type LogLevel

```go
type LogLevel int
```

LogLevel is a type that represents the importance level of a log message

```go
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
```

#### func (LogLevel) String

```go
func (ll LogLevel) String() string
```
