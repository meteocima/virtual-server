package tasks

import (
	"io"
	"os"
	"time"

	"github.com/meteocima/virtual-server/ctx"
	"github.com/meteocima/virtual-server/event"
)

// Task ...
type Task struct {
	Status        *TaskStatus
	StatusChanged *event.Emitter
	Failed        *event.Emitter
	Succeeded     *event.Emitter
	Done          *event.Emitter

	Progress     *event.Emitter
	FileProduced *event.Emitter

	StartedAt   time.Time
	CompletedAt time.Time

	ID          string
	detailedLog io.WriteCloser
	infoLog     io.WriteCloser
	Description string
	runner      TaskRunner
}

// TaskRunner ...
type TaskRunner func(tsk *Task, ctx *ctx.Context) error

var tasks = map[string]Task{}

// Stdout ...
var Stdout io.Writer

// Run ...
func (tsk *Task) Run() {
	go func() {

		infoLog := openTaskLog(tsk.ID + ".info.log")
		detailedLog := openTaskLog(tsk.ID + ".detailed.log")

		tsk.infoLog = NewLoggedWriteCloser(infoLog, detailedLog, Stdout)
		tsk.detailedLog = detailedLog
		vs := ctx.New(tsk.infoLog, tsk.detailedLog)

		vs.LogInfo("START: %s: %s", tsk.ID, tsk.Description)

		tsk.Status = Running
		err := tsk.runner(tsk, &vs)
		if err == nil && vs.Err != nil {
			err = vs.Err
		}

		if err != nil {
			vs.LogError("%s: %s", tsk.ID, err.Error())
		} else {
			vs.LogInfo("DONE: %s", tsk.ID)
		}

		infoLog.Close()
		detailedLog.Close()
		delete(tasks, tsk.ID)

		tsk.Done.Invoke(err)
		if err != nil {
			tsk.Failed.Invoke(err)
			tsk.Status = Failed(err)
		} else {
			tsk.Succeeded.Invoke(nil)
			tsk.Status = DoneOk
		}
	}()
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
		Status: Scheduled,
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

	return &t
}
