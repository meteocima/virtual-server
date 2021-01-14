package tasks

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/meteocima/virtual-server/config"
	"github.com/meteocima/virtual-server/ctx"
	"github.com/meteocima/virtual-server/testutil"
	"github.com/stretchr/testify/assert"
)

func resultEmit(results chan string, i int) func(vs *ctx.Context) error {
	ID := fmt.Sprintf("TEST%d", i+1)

	return func(vs *ctx.Context) error {
		time.Sleep(time.Duration(i*10) * time.Millisecond)
		results <- ID
		return nil
	}
}

func createsTestTasks(results chan string, count int) []*Task {
	tsks := make([]*Task, count)
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
	Stdout = ioutil.Discard

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
			tsks := createsTestTasks(results, 8)
			parent.AppendChildren(tsks...)
			parent.RunChild(tsks[7])
			parent.RunChild(tsks[2])
			parent.RunChild(tsks[1])
			return nil
		})
		parent.SetMaxParallelism(2)
		parent.Run()

		assert.Equal(t,
			[]string{"TEST3", "TEST2", "TEST8"},
			readResults(t, results, 3),
		)

		parent.Done.AwaitOne()
	})

}
