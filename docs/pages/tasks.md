{{ useLayout(".layout.njk") }}
{{ title("CIMA virtual-server") }}
{{ subtitle("tasks package") }}

# [virtual-server](./index) ⟶ {{ meta.subtitle }}





## Usage

```go
var Cancelled = &TaskStatus{}
```
Cancelled is the status of a task that won't run, because one of it's
prerequisites has failed

```go
var DoneOk = &TaskStatus{}
```
DoneOk is the status of a task successfully executed

```go
var Running = &TaskStatus{}
```
Running is the status of a running task

```go
var Scheduled = &TaskStatus{}
```
Scheduled is the status of a task scheduled but not yet executed

```go
var Stdout io.Writer
```
Stdout ...

#### func  List

```go
func List(w io.Writer)
```
List ...

#### type MultiWriteCloser

```go
type MultiWriteCloser struct {
	io.Writer
}
```

MultiWriteCloser is a wrapper io.WriteCloser that writes everything it's written
both to the source writer and to stdout

#### func  NewMultiWriteCloser

```go
func NewMultiWriteCloser(closer io.WriteCloser, writers ...io.Writer) *MultiWriteCloser
```
NewMultiWriteCloser creates a new MultiWriteCloser that wraps source

#### func (*MultiWriteCloser) Close

```go
func (mwc *MultiWriteCloser) Close() error
```
Close the source io.WriteCloser

#### type Task

```go
type Task struct {
	StatusChanged *event.Emitter
	Failed        *event.Emitter
	Succeeded     *event.Emitter
	Done          *event.Emitter

	Progress     *event.Emitter
	FileProduced *event.Emitter

	StartedAt   time.Time
	CompletedAt time.Time

	ID string

	Description string
}
```

Task ...

#### func  New

```go
func New(ID string, runner TaskRunner) *Task
```
New ...

#### func (*Task) Run

```go
func (tsk *Task) Run()
```
Run ...

#### func (*Task) SetStatus

```go
func (tsk *Task) SetStatus(newStatus *TaskStatus)
```
SetStatus ...

#### func (*Task) Status

```go
func (tsk *Task) Status() *TaskStatus
```
Status ...

#### type TaskRunner

```go
type TaskRunner func(ctx *ctx.Context) error
```

TaskRunner ...

#### type TaskStatus

```go
type TaskStatus struct {
	Err error
}
```

TaskStatus represents the status of a single task

#### func  Failed

```go
func Failed(err error) *TaskStatus
```
Failed returns the status of a task that failed with an error

#### func (*TaskStatus) IsFailure

```go
func (status *TaskStatus) IsFailure() bool
```
IsFailure returns whether the task status represents a failure

#### func (*TaskStatus) String

```go
func (st *TaskStatus) String() string
```
