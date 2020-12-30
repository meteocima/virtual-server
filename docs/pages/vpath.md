{{ useLayout(".layout.njk") }}
{{ title("CIMA virtual-server") }}
{{ subtitle("vpath package") }}

# [virtual-server](./index) ⟶ {{ meta.subtitle }}



Package vpath represents an instance of a virtual path. It is formed by the name
of a Host defined in configuration and by the absolute path within that file
system. The two parts are separated by a colon when representing the
vpath.VirtualPath as a string: host:path an empty string in Host field represent
the localhost Host. an empty string in Path field represent the current
directory (.).

## Usage

```go
var Stderr = &VirtualPath{
	Host: "any",
	Path: "stderr",
}
```
Stderr is a placeholder VirtualPath which represents the `stderr` stream of a
process.

```go
var Stdin = &VirtualPath{
	Host: "any",
	Path: "stdin",
}
```
Stdin is a placeholder VirtualPath which represents the `stdin` stream of a
process.

```go
var Stdout = &VirtualPath{
	Host: "any",
	Path: "stdout",
}
```
Stdout is a placeholder VirtualPath which represents the `stdout` stream of a
process.

#### type VirtualPath

```go
type VirtualPath struct {
	Host string
	Path string
}
```

VirtualPath represents an instance of a virtual path. It is formed by the name
of a Host defined in configuration and by the absolute path within that file
system. The two parts are separated by a colon when representing the
vpath.VirtualPath as a string: host:path an empty string in Host field represent
the localhost Host. an empty string in Path field represent the current
directory (.).

#### func  FromS

```go
func FromS(pathRepr string) VirtualPath
```
FromS returns a new VirtualPath with host and path parsed from pathRepr string
argument.

#### func  Local

```go
func Local(pathFormat string, pathArgs ...interface{}) VirtualPath
```
Local returns a new VirtualPath on localhost with the given path

#### func  New

```go
func New(host string, pathFormat string, pathArgs ...interface{}) VirtualPath
```
New returns an VirtualPath given its host and path parts. The path is built
using pathFormat argument as fmt.Sprintf format string and any pathArgs as
fmt.Sprintf arguments.

#### func (VirtualPath) AddExt

```go
func (vPath VirtualPath) AddExt(newExt string) VirtualPath
```
AddExt returns a new virtual path where the specified extension is appended to
the current filename.

#### func (VirtualPath) Dir

```go
func (vPath VirtualPath) Dir() VirtualPath
```
Dir returns a new VirtualPath formed by the same host and a path that is the
directory path of the original instance.

#### func (VirtualPath) Filename

```go
func (vPath VirtualPath) Filename() string
```
Filename returns the filename (with extension, but without directory path) of
the virtual path

#### func (VirtualPath) Join

```go
func (vPath VirtualPath) Join(pathFormat string, pathArgs ...interface{}) VirtualPath
```
Join returns a new virtual path formed by the same host and a path that is the
joining of the original path and and an additional one. The additional path is
built using pathFormat argument as fmt.Sprintf format string and any pathArgs as
fmt.Sprintf arguments.

#### func (VirtualPath) JoinP

```go
func (vPath VirtualPath) JoinP(other VirtualPath) VirtualPath
```
JoinP returns a new virtual path formed by the same host of the instance and a
path that is the joining of the original path and the path of an additional
VirtualPath

#### func (VirtualPath) ReplaceExt

```go
func (vPath VirtualPath) ReplaceExt(newExt string) VirtualPath
```
ReplaceExt returns a new virtual path where the extension of the file is
replaced with the given one.

#### func (VirtualPath) String

```go
func (vPath VirtualPath) String() string
```
String returns a string representing the virtual path Host and path parts are
separated by a colon: host:path

#### func (VirtualPath) StringRel

```go
func (vPath VirtualPath) StringRel() string
```
StringRel returns a string representing the virtual path Host and path parts are
separated by a colon: host:path

#### func (*VirtualPath) UnmarshalText

```go
func (vPath *VirtualPath) UnmarshalText(data []byte) error
```
UnmarshalText ...

#### type VirtualPathList

```go
type VirtualPathList []VirtualPath
```

VirtualPathList ...

#### func (VirtualPathList) Len

```go
func (list VirtualPathList) Len() int
```
Len is the number of elements in the collection.

#### func (VirtualPathList) Less

```go
func (list VirtualPathList) Less(i, j int) bool
```
Less reports whether the element with index i should sort before the element
with index j.

#### func (VirtualPathList) Swap

```go
func (list VirtualPathList) Swap(i, j int)
```
Swap swaps the elements with indexes i and j.
