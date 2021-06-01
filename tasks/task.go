package tasks

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/meteocima/virtual-server/ctx"
	"github.com/meteocima/virtual-server/event"
	"github.com/meteocima/virtual-server/vpath"
)

// TaskFile ...
type TaskFile struct {
	Path vpath.VirtualPath
	Meta interface{}
}

// Task ...
type Task struct {
	status        *TaskStatus
	StatusChanged *event.Emitter
	Failed        *event.Emitter
	Succeeded     *event.Emitter
	Done          *event.Emitter

	Progress     *event.Emitter
	FileProduced *event.Emitter

	StartedAt   time.Time
	CompletedAt time.Time

	ID          string
	stdout      io.WriteCloser
	stderr      io.WriteCloser
	Description string
	runner      TaskRunner
}

// TaskRunner ...
type TaskRunner func(ctx *ctx.Context) error

// TaskRegistry ...
type TaskRegistry struct {
	tasks    map[string]*Task
	taskLock sync.Mutex
}

// AllTasks ...
func (reg *TaskRegistry) AllTasks() []*Task {
	reg.taskLock.Lock()
	defer reg.taskLock.Unlock()
	res := make([]*Task, len(reg.tasks))
	idx := 0
	for _, task := range reg.tasks {
		res[idx] = task
		idx++
	}
	return res
}

// RemoveTask ...
func (reg *TaskRegistry) RemoveTask(ID string) {
	reg.taskLock.Lock()
	defer reg.taskLock.Unlock()
	delete(reg.tasks, ID)
}

// AddTask ...
func (reg *TaskRegistry) AddTask(tsk *Task) {
	reg.taskLock.Lock()
	defer reg.taskLock.Unlock()
	reg.tasks[tsk.ID] = tsk
}

var registry = &TaskRegistry{
	tasks:    map[string]*Task{},
	taskLock: sync.Mutex{},
}

// List ...
func List(w io.Writer) {
	for _, task := range registry.AllTasks() {
		fmt.Fprintf(w, "%s: %s [%s]\n", task.ID, task.Description, task.Status().String())
	}
}

// Stdout ...
var Stdout io.Writer

// SetStatus ...
func (tsk *Task) SetStatus(newStatus *TaskStatus) {
	if tsk.status == newStatus {
		return
	}

	tsk.status = newStatus
	tsk.StatusChanged.Invoke(newStatus)
}

// Status ...
func (tsk *Task) Status() *TaskStatus {
	return tsk.status
}

// Run ...
func (tsk *Task) Run() {
	go func() {

		//infoLog := openTaskLog(tsk.ID + ".info.log")
		stderr := openTaskLog(tsk.ID + ".log")

		tsk.stdout = NewMultiWriteCloser(stderr, Stdout)
		tsk.stderr = stderr
		vs := ctx.New(os.Stdin, tsk.stdout, tsk.stderr)
		vs.ID = tsk.ID
		vs.LogInfo("START: %s", tsk.Description)

		tsk.SetStatus(Running)
		err := tsk.runner(vs)
		if err == nil && vs.Err != nil {
			err = vs.Err
		}

		if err != nil {
			vs.LogError(err.Error())
		} else {
			vs.LogInfo("DONE")
		}

		vs.Close()
		stderr.Close()

		tsk.SetCompleted(err)
	}()
}

// AwaitDone ...
func (tsk *Task) AwaitDone() {
	tsk.Done.AwaitOne()
}

// SetCompleted ...
func (tsk *Task) SetCompleted(err error) {
	if err != nil {
		tsk.Failed.Invoke(err)
		tsk.SetStatus(Failed(err))
	} else {
		tsk.Succeeded.Invoke(nil)
		tsk.SetStatus(DoneOk)
	}
	//fmt.Printf("Invoke Done %v\n", tsk.Done)
	tsk.Done.Invoke(err)

	event.CloseEmitters(
		&tsk.StatusChanged,
		&tsk.Failed,
		&tsk.Succeeded,
		&tsk.Done,
		&tsk.Progress,
		&tsk.FileProduced,
	)

	registry.RemoveTask(tsk.ID)
}

func openTaskLog(path string) io.WriteCloser {
	logFile, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		panic(err)
	}
	return logFile
}

// New ...
func New(ID string, runner TaskRunner) *Task {
	t := Task{
		status: Scheduled,
		ID:     ID,
		runner: runner,
	}

	event.InitSource(
		&t,
		&t.StatusChanged,
		&t.Failed,
		&t.Succeeded,
		&t.Done,
		&t.Progress,
		&t.FileProduced,
	)

	registry.AddTask(&t)
	return &t
}
