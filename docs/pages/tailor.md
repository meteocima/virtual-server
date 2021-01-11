{{ useLayout(".layout.njk") }}
{{ title("CIMA virtual-server") }}
{{ subtitle("tailor package") }}

# [virtual-server](./index) ⟶ {{ meta.subtitle }}





## Usage

```go
const (
	ErrNoSource      = "cannot start scroller: no source"
	ErrNoTarget      = "cannot start scroller: no target"
	ErrNegativeLines = "negative number of lines not allowed: %d"
)
```
Error messages ...

#### type Tailor

```go
type Tailor struct {
}
```

Tailor scrolls and filters a ReadSeeker line by line and writes the data into a
Writer.

#### func  New

```go
func New(source io.ReadSeeker, target io.Writer, buffsz int) *Tailor
```
New creates a Tailor for the given source and target.

#### func (*Tailor) Start

```go
func (s *Tailor) Start() chan error
```
Start is the goroutine for reading, filtering and writing. The returned chan
emit an error in case of failure, otherwise it's closed upon tail completion
without emitting nothing.

#### func (*Tailor) Stop

```go
func (s *Tailor) Stop()
```
Stop will causes the tail operation to ends. The log file is readed until EOF
before stopping.
