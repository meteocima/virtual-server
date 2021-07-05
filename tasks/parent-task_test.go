package tasks

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/meteocima/virtual-server/config"
	"github.com/meteocima/virtual-server/ctx"
	"github.com/meteocima/virtual-server/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resultEmit(results chan string, i int) func(vs *ctx.Context) error {
	ID := fmt.Sprintf("TEST%d", i+1)

	return func(vs *ctx.Context) error {
		fmt.Printf("%s go to sleep for %d ms\n", ID, i*150)
		time.Sleep(time.Duration(i*50) * time.Millisecond)
		fmt.Printf("%s awaked\n", ID)
		results <- ID
		fmt.Printf("%s exiting\n", ID)
		return nil
	}
}

func createsTestTasks(results chan string, count int) []TaskI {
	tsks := make([]TaskI, count)
	for i := 0; i < count; i++ {
		ID := fmt.Sprintf("TEST%d", i+1)
		tsks[i] = New(ID, resultEmit(results, i))
	}
	return tsks
}

func readResults(t *testing.T, results chan string, count int) []string {
	res := make([]string, count)
	for i := 0; i < count; i++ {
		r := <-results
		res[i] = r
	}

	select {
	case <-results:
		t.Fatalf("Expected %d results", count)
	default:
	}

	return res
}

func TestParentTask(t *testing.T) {
	err := config.Init(testutil.FixtureDir("virt-serv.toml"))
	assert.NoError(t, err)
	Stdout = os.Stdout
	Stderr = os.Stdout
	t.Run("when a task runner completed without children it fails", func(t *testing.T) {

		var parent TaskI = NewParent("PARENT1", func(vs *ctx.Context) error {
			return nil
		})

		parent.Run()
		parent.AwaitDone()
		require.NotNil(t, parent.Status())
		assert.True(t, parent.Status().IsFailure())
		require.NotNil(t, parent.Status().Err)
		assert.Equal(t, "Err: task completed without children", parent.Status().String())
		assert.Equal(t, "task completed without children", parent.Status().Err.Error())

	})

	t.Run("change max parallelism when already started panics", func(t *testing.T) {
		results := make(chan string)
		children := createsTestTasks(results, 1)

		var parent *ParentTask

		parent = NewParent("PARENT1", func(vs *ctx.Context) error {
			return nil
		})

		assert.Panics(t, func() {
			parent.AppendChildren(children[0])
			parent.RunChild(children[0])
			parent.SetMaxParallelism(1)
		})

	})

	t.Run("adding same child multiple times panics", func(t *testing.T) {
		results := make(chan string)
		children := createsTestTasks(results, 1)

		var parent *ParentTask

		parent = NewParent("PARENT1", func(vs *ctx.Context) error {
			return nil
		})

		assert.Panics(t, func() {
			parent.AppendChildren(children[0])
			parent.AppendChildren(children[0])
		})
	})

	mkSerialParent := func(results chan string) *ParentTask {
		var parent *ParentTask
		parent = NewParent("PARENT", func(vs *ctx.Context) error {
			var tsks []TaskI
			tsks = append(tsks, New("TEST1", func(vs *ctx.Context) error {
				results <- "TEST1"
				return nil
			}))
			tsks = append(tsks, New("TEST2", func(vs *ctx.Context) error {
				results <- "TEST2"
				return errors.New("ciccio")
			}))
			tsks = append(tsks, New("TEST3", func(vs *ctx.Context) error {
				results <- "TEST3"
				return nil
			}))

			parent.AppendChildren(tsks...)
			parent.RunChild(tsks[0])
			parent.RunChild(tsks[1])
			parent.RunChild(tsks[2])
			return nil
		})

		parent.SetMaxParallelism(1)
		return parent
	}

	t.Run("nested children", func(t *testing.T) {
		results := make(chan string)
		children := createsTestTasks(results, 4)

		var parent1 *ParentTask
		var parent2 *ParentTask
		var grandParent *ParentTask

		parent1 = NewParent("PARENT1", func(vs *ctx.Context) error {
			parent1.AppendChildren(children[1], children[0])
			parent1.RunChild(children[1])
			parent1.RunChild(children[0])
			return nil
		})

		parent2 = NewParent("PARENT2", func(vs *ctx.Context) error {
			parent2.AppendChildren(children[2], children[3])
			parent2.RunChild(children[2])
			parent2.RunChild(children[3])
			return nil
		})

		grandParent = NewParent("GRANDPA", func(vs *ctx.Context) error {
			grandParent.AppendChildren(parent1, parent2)
			grandParent.RunChild(parent2)
			grandParent.RunChild(parent1)
			return nil
		})

		grandParent.SetMaxParallelism(1)
		grandParent.Run()

		assert.Equal(t,
			[]string{"TEST3", "TEST4", "TEST1", "TEST2"},
			readResults(t, results, 4),
		)

		grandParent.Done.AwaitOne()
		fmt.Print("DONE")
	})

	t.Run("stop on first failed task when no max parallelism is set", func(t *testing.T) {
		results := make(chan string)

		var parent *ParentTask
		parent = NewParent("PARENT", func(vs *ctx.Context) error {
			var tsks []TaskI
			tsks = append(tsks, New("TEST1", func(vs *ctx.Context) error {
				results <- "TEST1"
				return errors.New("TEST ERROR")
			}))
			tsks = append(tsks, New("TEST2", func(vs *ctx.Context) error {
				results <- "TEST2"
				return errors.New("ciccio")
			}))
			tsks = append(tsks, New("TEST3", func(vs *ctx.Context) error {
				results <- "TEST3"
				return nil
			}))

			parent.AppendChildren(tsks...)
			parent.RunChild(tsks[0])
			go func() {
				time.Sleep(50 * time.Millisecond)
				parent.RunChild(tsks[1])
				parent.RunChild(tsks[2])
			}()
			return nil
		})
		parent.SetMaxParallelism(0)
		parent.SetFailFast()
		parent.Run()
		assert.Equal(t,
			[]string{"TEST1"},
			readResults(t, results, 1),
		)

		parent.Done.AwaitOne()
		fmt.Print("DONE")
	})
	t.Run("stop on first failed task", func(t *testing.T) {
		results := make(chan string)

		parent := mkSerialParent(results)
		parent.SetFailFast()
		parent.Run()
		assert.Equal(t,
			[]string{"TEST1", "TEST2"},
			readResults(t, results, 2),
		)
		parent.Done.AwaitOne()
	})

	t.Run("continue with other tasks on error", func(t *testing.T) {
		results := make(chan string)

		parent := mkSerialParent(results)

		parent.Run()
		assert.Equal(t,
			[]string{"TEST1", "TEST2", "TEST3"},
			readResults(t, results, 3),
		)
		parent.Done.AwaitOne()
	})

	t.Run("no max parallelism", func(t *testing.T) {
		var parent *ParentTask
		results := make(chan string)
		parent = NewParent("PARENT", func(vs *ctx.Context) error {
			tsks := createsTestTasks(results, 3)
			parent.AppendChildren(tsks...)
			parent.RunChild(tsks[0])
			parent.RunChild(tsks[1])
			parent.RunChild(tsks[2])
			return nil
		})
		parent.Run()

		assert.Equal(t,
			[]string{"TEST1", "TEST2", "TEST3"},
			readResults(t, results, 3),
		)

		parent.Done.AwaitOne()
	})

	t.Run("sequential parallelism", func(t *testing.T) {
		var parent *ParentTask
		results := make(chan string)
		parent = NewParent("PARENT", func(vs *ctx.Context) error {
			tsks := createsTestTasks(results, 3)
			parent.AppendChildren(tsks...)
			parent.RunChild(tsks[2])
			parent.RunChild(tsks[1])
			parent.RunChild(tsks[0])
			return nil
		})
		parent.SetMaxParallelism(1)
		parent.Run()

		assert.Equal(t,
			[]string{"TEST3", "TEST2", "TEST1"},
			readResults(t, results, 3),
		)

		parent.Done.AwaitOne()
	})

	t.Run("limited parallelism", func(t *testing.T) {
		var parent *ParentTask
		results := make(chan string)
		parent = NewParent("PARENT", func(vs *ctx.Context) error {
			tsks := createsTestTasks(results, 4)
			parent.AppendChildren(tsks...)
			parent.RunChild(tsks[1])
			parent.RunChild(tsks[0])
			parent.RunChild(tsks[2])
			parent.RunChild(tsks[3])

			return nil
		})
		parent.SetMaxParallelism(3)
		parent.Run()

		assert.Equal(t,
			[]string{"TEST1", "TEST2", "TEST3", "TEST4"},
			readResults(t, results, 4),
		)

		parent.Done.AwaitOne()
	})

}
