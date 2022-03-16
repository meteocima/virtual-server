package tasks

import (
	"errors"
	"fmt"
	"sync"

	"github.com/meteocima/virtual-server/ctx"
	"github.com/tevino/abool"
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
	children map[TaskI]struct{}
	// set of task that are waiting for execution. When a task is tentatively run using
	// RunChild, if it cannot immediately run because it would overflow max parallelism,
	// it's stored here and later retrieved when enough running task completed in order
	// to allow parallelism to be respected
	waitingChildren []TaskI
	// used to implement max parallelism, will be created with a cache.
	// it's nil when max parallelism is set to 0 (no max parallelism)
	runningChild        chan struct{}
	failfast            bool
	failed              *abool.AtomicBool
	someChildrenStarted *abool.AtomicBool

	// synchronizes `waitingChildren` members access
	sem *sync.Mutex
	// synchronizes `children` members access
	lckChildren *sync.Mutex
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
	TaskID() string
}

func (tsk *ParentTask) setFailed(value bool) {
	tsk.failed.SetTo(value)
}

func (tsk *ParentTask) getFailed() bool {
	return tsk.failed.IsSet()
}

func (tsk *ParentTask) popWaitingChild() TaskI {
	tsk.sem.Lock()
	defer tsk.sem.Unlock()
	next := tsk.waitingChildren[0]
	tsk.waitingChildren = tsk.waitingChildren[1:]
	return next
}

func (tsk *ParentTask) pushWaitingChild(child TaskI) {
	tsk.sem.Lock()
	defer tsk.sem.Unlock()
	tsk.waitingChildren = append(tsk.waitingChildren, child)
}

func (tsk *ParentTask) setWaitingChild(value []TaskI) {
	tsk.sem.Lock()
	defer tsk.sem.Unlock()
	tsk.waitingChildren = value
}

func (tsk *ParentTask) hasWaitingChildren() bool {
	tsk.sem.Lock()
	defer tsk.sem.Unlock()
	return len(tsk.waitingChildren) > 0
}

// AppendChildren add specified children to this
// task. This call is not sinchronized, it's not safe
// to call it concurrently from multiple goroutines.
// Caller should provide sinchronization itself if needed.
func (tsk *ParentTask) AppendChildren(children ...TaskI) {
	tsk.lckChildren.Lock()
	defer tsk.lckChildren.Unlock()
	for _, child := range children {

		if _, exists := tsk.children[child]; exists {
			panic(fmt.Sprintf("Task %s already appended", child.TaskID()))
		}
		tsk.children[child] = struct{}{}
	}
}

// TaskID returns ID of the task
func (tsk *ParentTask) TaskID() string {
	return tsk.ID
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
	tsk.someChildrenStarted.Set()

	if tsk.runningChild == nil {
		// no max parallelism, so just run the task.

		if tsk.getFailed() {
			// another child already failed, and failfast option is set.
			child.SetCompleted(errors.New("tasks cancelled by parent"))
			return
		}

		child.Run()
		if tsk.failfast {
			go func() {
				// await for the children to finish,
				// and eventually set the parent task
				// failed
				child.AwaitDone()

				// when failfast option is set, set the whole task has failed.
				if child.Status().IsFailure() {
					tsk.setFailed(true)
				}
			}()
		}

		return
	}

	select {
	case tsk.runningChild <- struct{}{}:
		// tsk.runningChild accepted the child in cache,
		// that means max parallelism is respected, and the
		// child is consuming 1 slot in runningChild chan cache.
		go func() {
			child.AwaitDone()
			// task has done, free 1 slot in runningChild chan cache.
			<-tsk.runningChild

			// when failfast option is set, set the whole task has failed.
			if tsk.failfast && child.Status().IsFailure() {
				tsk.setFailed(true)
			}

			if tsk.hasWaitingChildren() {
				// there is at least one waiting child task,
				// pick and run it.
				tsk.RunChild(tsk.popWaitingChild())
			}
		}()

		if tsk.getFailed() {
			// another child already failed, and failfast option is set.
			child.SetCompleted(errors.New("tasks cancelled by parent"))
			return
		}

		child.Run()

	default:
		// max parallelism reached. Store
		// child task in waiting store, in
		// order to be picked later for execution.
		tsk.pushWaitingChild(child)
	}

}

// SetMaxParallelism sets the maximum allowed number of
// children tasks that can run concurrently.
func (tsk *ParentTask) SetMaxParallelism(count uint) {

	if tsk.someChildrenStarted.IsSet() {
		panic("cannot change parallelism: a children task has already been started")
	}
	tsk.setWaitingChild([]TaskI{})
	if count == 0 {
		tsk.runningChild = nil
		return
	}
	tsk.runningChild = make(chan struct{}, count)
}

// SetFailFast makes the parent fails
// on first child failure.
func (tsk *ParentTask) SetFailFast() {
	tsk.failfast = true
}

// NewParent creates a ParentTask instance that
// can contains children tasks amd returns its
// addrress.
func NewParent(ID string, runner TaskRunner) *ParentTask {
	tsk := ParentTask{
		sem:                 &sync.Mutex{},
		lckChildren:         &sync.Mutex{},
		children:            map[TaskI]struct{}{},
		failed:              abool.New(),
		someChildrenStarted: abool.New(),
	}

	tsk.Task = New(ID, wrapRunner(runner, &tsk))

	return &tsk
}

// wraps a TaskRunner function in order to
// await for completion of all children tasks
// before termination.
func wrapRunner(runner TaskRunner, tsk *ParentTask) TaskRunner {
	return func(vs *ctx.Context) error {
		err := runner(vs)
		tsk.lckChildren.Lock()
		children := map[TaskI]struct{}{}
		for k, v := range tsk.children {
			children[k] = v
		}
		tsk.lckChildren.Unlock()

		if len(children) == 0 {
			return errors.New("task completed without children")
		}
		for child := range children {
			child.AwaitDone()
		}

		return err
	}
}
