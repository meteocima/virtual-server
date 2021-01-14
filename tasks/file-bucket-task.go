package tasks

import (
	"fmt"
	"sync"
	"time"

	"github.com/meteocima/virtual-server/ctx"
	"github.com/meteocima/virtual-server/event"
)

// NewFilesBucketTask creates a Task that listen to one
// or more tasks for files produced.
//
// For each produced file, it creates and run a new task
// using the provided `factory` function.
//
// The task terminates when all given tasks are done
// and also all new created child tasks are done.
func NewFilesBucketTask(factory func(TaskFile) *Task, tasks ...*Task) *Task {
	var tsk *Task

	tsk = New("tskID", func(vs *ctx.Context) error {
		tasksToWait := sync.WaitGroup{}
		waitTask := func(tsk *Task) {
			fmt.Println(tsk.ID, "WAIT")
			tasksToWait.Add(1)
			c := make(chan struct{})
			go func() {
				if tsk.Status() != Scheduled {
					panic(tsk.Status().String() + "You cannot create a FilesBucketTask from already started tasks. " +
						"FilesBucketTask takes care of running the tasks itself.")
				}

				fmt.Println(tsk.ID, "AWAIT COMPLETION")
				close(c)
				tsk.Done.AwaitOne()
				time.Sleep(5 * time.Millisecond)
				fmt.Println(tsk.ID, "COMPLETED")

				tasksToWait.Done()
			}()
			<-c
		}

		for _, tsk := range tasks {
			waitTask(tsk)
			tsk.FileProduced.Listen(func(e *event.Event) {
				newTask := factory(e.Payload.(TaskFile))
				waitTask(newTask)
				newTask.Run()
			})
			tsk.Run()
		}

		fmt.Println(tsk.ID, "WAITING ALL TASK...")

		tasksToWait.Wait()
		fmt.Println(tsk.ID, "ALL TASK DONE...")

		return nil
	})
	return tsk
}
