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

#### func  List

```go
func List(w io.Writer)
```
List ...

#### func  Run

```go
func Run(tsk Task)
```
Run ...

#### type SimulationTaskStatus

```go
type SimulationTaskStatus struct {
	FinalDewetraDelivery        *TaskStatus
	VdADelivery                 []*TaskStatus
	ArpalDelivery               []*TaskStatus
	ArpaPiemonteDelivery        []*TaskStatus
	ContinuumDelivery           []*TaskStatus
	ArpaPiemonteIndexesDelivery []*TaskStatus

	AUXDownloadDomain1       []*TaskStatus
	AUXPostProcessDomain1    []*TaskStatus
	AUXDownloadDomain3       []*TaskStatus
	AUXPostProcessDomain3    []*TaskStatus
	OUTPostProcess           []*TaskStatus
	OUTPostProcessedDownload []*TaskStatus
	FinalZTDScript           *TaskStatus
}
```

SimulationTaskStatus is the status of a single WRF run

#### func  NewSimulationTaskStatus

```go
func NewSimulationTaskStatus(totHours int) SimulationTaskStatus
```
NewSimulationTaskStatus returns a new SimulationTaskStatus instance initialized
for a run of totHours hours

#### type Task

```go
type Task interface {
	fmt.Stringer
	Run() error
	ID() string
	InfoLogFilePath() string
	DetailedLogFilePath() string
	DetailedLog() io.WriteCloser
	InfoLog() io.WriteCloser
}
```

Task ...

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
