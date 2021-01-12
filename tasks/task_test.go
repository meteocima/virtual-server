package tasks

import (
	"bytes"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/meteocima/virtual-server/config"
	"github.com/meteocima/virtual-server/ctx"
	"github.com/meteocima/virtual-server/testutil"
	"github.com/meteocima/virtual-server/vpath"
	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	err := config.Init(testutil.FixtureDir("virt-serv.toml"))
	assert.NoError(t, err)

	tests := sync.WaitGroup{}
	tests.Add(3)
	t.Run("Debug and details levels are not logged to stdout", func(t *testing.T) {
		bytesWriter := bytes.Buffer{}
		Stdout = &bytesWriter
		tsk := New("TEST", func(tsk *Task, vs *ctx.Context) error {
			vs.LogInfo("ciao")
			vs.LogDetail("salve")
			vs.LogDebug("urrà")
			assert.Same(t, tsk.Status, Running)
			tests.Done()

			return nil
		})
		tsk.Description = "A task for tests."
		assert.NotNil(t, tsk)

		tsk.Run()

		tsk.Done.AwaitOne()
		assert.Same(t, tsk.Status, DoneOk)
		assert.Equal(t,
			`INFO: START: TEST: A task for tests.
INFO: ciao
INFO: DONE: TEST
`, bytesWriter.String())

	})
	t.Run("All log levels are logged to file", func(t *testing.T) {
		Stdout = ioutil.Discard
		tsk := New("TEST", func(tsk *Task, vs *ctx.Context) error {
			vs.LogInfo("ciao")
			vs.LogDetail("salve")
			vs.LogDebug("urrà")
			tests.Done()
			return nil
		})
		tsk.Description = "A task for tests."
		assert.NotNil(t, tsk)

		tsk.Run()

		tsk.Done.AwaitOne()
		assert.Same(t, tsk.Status, DoneOk)

		contentBuff, err := ioutil.ReadFile("TEST.detailed.log")
		assert.NoError(t, err)
		assert.Equal(t,
			`INFO: START: TEST: A task for tests.
INFO: ciao
DETAIL: salve
DEBUG: urrà
INFO: DONE: TEST
`, string(contentBuff))

	})

	t.Run("Non existent server", func(t *testing.T) {
		Stdout = ioutil.Discard
		tsk := New("TEST", func(tsk *Task, vs *ctx.Context) error {
			vs.Link(vpath.New("peppa", "./bad"), vpath.New("peppa", "./bad"))
			tests.Done()
			return nil
		})
		tsk.Description = "A task for tests."
		assert.NotNil(t, tsk)

		tsk.Run()

		tsk.Done.AwaitOne()

		assert.True(t, tsk.Status.IsFailure())
		assert.Equal(t,
			"Link from peppa:./bad to peppa:./bad: connection.FindHost: wrong configuration file \"../fixtures/virt-serv.toml\": unknown host `peppa`",
			tsk.Status.Err.Error(),
		)
	})

	tests.Wait()
}
