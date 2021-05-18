package tasks

import (
	"errors"
	"sync"

	"github.com/meteocima/virtual-server/ctx"
	"github.com/meteocima/virtual-server/event"
)

// ParentTask is a struct that represents a task that can contains
// children tasks. The struct inherits `Task`, so all its methods
// and fields are avaiable on a ParentTask.
//
// Children task can be appended using the `AppendChildren`
// method of `ParentTask`, during the parent run phase or before.
//
// Parent task completes when all children are completed.
// If any of the the, fails, the other ones continue to run,
// but the parent status will be set to failed
// upon completion.
//
// It's the parent task responsibility to run children when
// appropriate, using `RunChild` method of `ParentTask`
//
// Number of children tasks that can runs concurrently can be
// limited using `SetMaxParallelism` method.
//
// A max parallelism of 0 (the default) means no limit,
// children tasks start as soon as you call `RunChild`.
//
// When the max parallelism is greater then 1, children
// tasks that you schedule to run with `RunChild` are delayed
// until enough child tasks completes to respect the limit.
//
// A max parallelism of 1 (the default) means that
// the children run sequentially, effectively removing any
// parallelism.
type ParentTask struct {
	*Task
	children        map[*Task]struct{}
	waitingChildren []*Task
	runningChild    chan struct{}
	failfast        bool
	failed          bool
	sem             *sync.Mutex
}

// AppendChildren ...
func (tsk *ParentTask) AppendChildren(children ...*Task) {
	for _, child := range children {
		tsk.children[child] = struct{}{}
	}
}

// RunChild ...
func (tsk *ParentTask) RunChild(child *Task) {
	if tsk.runningChild == nil {
		child.Run()
		return
	}

	select {
	case tsk.runningChild <- struct{}{}:
		//fmt.Println("TASK ENTER", child.ID)

		go func() {
			child.Done.AwaitOne()
			//fmt.Println("TASK COMPLETED", child.ID)
			<-tsk.runningChild

			if tsk.failfast && child.Status().IsFailure() {
				tsk.sem.Lock()
				tsk.failed = true
				tsk.sem.Unlock()
			}

			//fmt.Println("TASK RECOVERING", child.ID)

			tsk.sem.Lock()

			if len(tsk.waitingChildren) == 0 {
				tsk.sem.Unlock()
				return
			}

			next := tsk.waitingChildren[0]
			tsk.waitingChildren = tsk.waitingChildren[1:]
			tsk.sem.Unlock()

			tsk.RunChild(next)
		}()

		tsk.sem.Lock()
		failed := tsk.failed
		tsk.sem.Unlock()

		if failed {
			err := errors.New("tasks cancelled by parent")
			child.Failed.Invoke(err)

			child.SetStatus(Failed(err))
			child.Done.Invoke(err)

			event.CloseEmitters(
				&child.StatusChanged,
				&child.Failed,
				&child.Succeeded,
				&child.Done,
				&child.Progress,
				&child.FileProduced,
			)

			registry.RemoveTask(child.ID)
		} else {
			child.Run()
		}
	default:
		//fmt.Println("TASK WAIT", child.ID)
		tsk.sem.Lock()
		tsk.waitingChildren = append(tsk.waitingChildren, child)
		tsk.sem.Unlock()
	}

}

// SetMaxParallelism ...
func (tsk *ParentTask) SetMaxParallelism(count uint) {
	tsk.waitingChildren = []*Task{}
	if count == 0 {
		tsk.runningChild = nil
		return
	}
	tsk.runningChild = make(chan struct{}, count)
}

// SetFailFast make the parent fails
// on first child failures.
func (tsk *ParentTask) SetFailFast() {
	tsk.failfast = true
}

// NewParent creates a ParentTask instance that
// can contains children tasks amd returns its
// addrress
func NewParent(ID string, runner TaskRunner) *ParentTask {
	tsk := ParentTask{
		sem:      &sync.Mutex{},
		children: map[*Task]struct{}{},
	}

	tsk.Task = New(ID, wrapRunner(runner, &tsk))

	return &tsk
}

func wrapRunner(runner TaskRunner, tsk *ParentTask) TaskRunner {
	return func(vs *ctx.Context) error {
		err := runner(vs)
		for child := range tsk.children {
			child.Done.AwaitOne()
		}

		return err
	}
}
