package tasks

import (
	"fmt"

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
}

// AppendChildren ...
func (tsk *ParentTask) AppendChildren(children ...*Task) {
	for _, child := range children {
		tsk.children[child] = struct{}{}
	}
}

// RunChild ...
func (tsk *ParentTask) RunChild(child *Task) {
	if tsk.runningChild != nil {
		select {
		case tsk.runningChild <- struct{}{}:
			fmt.Println("TASK ENTER", child.ID)
			child.Run()
			go func() {
				child.Done.AwaitOne()
				<-tsk.runningChild
				if len(tsk.waitingChildren) > 0 {
					fmt.Println("TASK RECOVERED", child.ID)
					next := tsk.waitingChildren[0]
					tsk.waitingChildren = tsk.waitingChildren[1:]
					tsk.RunChild(next)
				}
			}()
		default:
			fmt.Println("TASK WAIT", child.ID)
			tsk.waitingChildren = append(tsk.waitingChildren, child)
		}

		return
	}
	child.Run()

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

// NewParent creates a ParentTask instance that
// can contains children tasks amd returns its
// addrress
func NewParent(ID string, runner TaskRunner) *ParentTask {
	tsk := ParentTask{

		children: map[*Task]struct{}{},
	}

	tsk.Task = New(ID, wrapRunner(runner, &tsk))

	return &tsk
}

func wrapRunner(runner TaskRunner, tsk *ParentTask) TaskRunner {
	return func(vs *ctx.Context) error {
		err := runner(vs)
		/*for child := range tsk.children {
			child.Done.AwaitOne()
		}*/
		return err
	}
}
