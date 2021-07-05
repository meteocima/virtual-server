package tasks

import (
	"errors"
	"sync"

	"github.com/meteocima/virtual-server/ctx"
)

// ParentTask is a struct that represents a task that can contains
// children tasks. The struct inherits `Task`, so all its methods
// and fields are avaiable on a ParentTask.
//
// Children task can be appended using the `AppendChildren`
// method of `ParentTask`, during the parent run phase or before.
//
// Parent task completes when all children are completed.
// If any of the children fails, the other ones continue to run,
// but the parent status will be set to failed upon completion.
// This semantic can be changed by calling SetFailFast(true), that
// makes the children `Task` immediately fails when run
// after any children had failed.
//
// It's the parent task responsibility to run children when
// appropriate, using `RunChild` method.
//
// Number of children tasks that can runs concurrently can be
// limited using `SetMaxParallelism` method.
//
// A max parallelism of 0 (the default) means no limit.
// When o is set, children tasks start as soon as you call
// `RunChild` on the parent.
//
// When the max parallelism is greater then 1, children
// tasks that you schedule to run with `RunChild` are delayed
// until enough child tasks will have completed, in order to respect
// the choosen limit.
//
// A max parallelism of 1 (the default) means that
// the children run sequentially, effectively removing any
// parallelism.
type ParentTask struct {
	*Task
	// set of children task. each children could be a Task or another ParentTask
	// so the variable store TaskI interfaces.
	// TODO: I don't remember why it is implemented as a map, maybe to allow for Task removal?
	children map[TaskI]struct{}
	// set of task that are waiting for execution. When a task is tentatively run using
	// RunChild, if it cannot immediately run because it would overflow max parallelism,
	// it's stored here and later retrieved when enough running task completed in order
	// to allow parallelism to be respected
	waitingChildren []TaskI
	// used to implement max parallelism, will be created with a cache.
	runningChild chan struct{}
	failfast     bool
	failed       bool
	// synchronize failed member access
	sem *sync.Mutex
}

// TaskI is an interface implemented by
// ParentTask and Task. It's needed to allow
// ParentTask to contains both ParentTask and
// Task structures.
type TaskI interface {
	Run()
	Status() *TaskStatus
	SetStatus(newStatus *TaskStatus)
	SetCompleted(err error)
	AwaitDone()
}

// AppendChildren add specified children to this
// task. This call is not sinchronized, it's not safe
// to call it concurrently from multiple goroutines.
// Caller should provide sinchronization itself if needed.
func (tsk *ParentTask) AppendChildren(children ...TaskI) {
	for _, child := range children {
		tsk.children[child] = struct{}{}
	}
}

// RunChild schedule specified child task for
// execution. If max parallelism is set to 0, and else if
// running task are below max parallelism allowed, the task
// run immediately. Otherwise, it's put in a queue
// for it to be runner later.
// If the fail fast option is set on the ParentTask,
// and a task already failed, then the child immediately
// fails with an error og "tasks cancelled by parent"
func (tsk *ParentTask) RunChild(child TaskI) {
	if tsk.runningChild == nil {
		child.Run()
		return
	}

	select {
	case tsk.runningChild <- struct{}{}:
		// tsk.runningChild accepted the child in cache,
		// that means max parallelism is respected, and the
		// child is consuming 1 place in runningChild chan cache.
		go func() {
			child.AwaitDone()
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
			child.SetCompleted(errors.New("tasks cancelled by parent"))
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
	tsk.waitingChildren = []TaskI{}
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
		children: map[TaskI]struct{}{},
	}

	tsk.Task = New(ID, wrapRunner(runner, &tsk))

	return &tsk
}

func wrapRunner(runner TaskRunner, tsk *ParentTask) TaskRunner {
	return func(vs *ctx.Context) error {
		err := runner(vs)
		for child := range tsk.children {
			child.AwaitDone()
		}

		return err
	}
}
