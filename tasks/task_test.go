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
		var tsk *Task
		tsk = New("TEST", func(vs *ctx.Context) error {
			vs.LogInfo("ciao")
			vs.LogDetail("salve")
			vs.LogDebug("urrà")
			assert.NoError(t, MustBeEqual(tsk.Status(), Running))
			tests.Done()

			return nil
		})
		tsk.Description = "A task for tests."
		assert.NotNil(t, tsk)
		assert.Equal(t, "TEST", tsk.ID)

		tsk.Run()

		tsk.Done.AwaitOne()
		assert.NoError(t, MustBeEqual(tsk.Status(), DoneOk))
		assert.Equal(t,
			`INFO: TEST: START: A task for tests.
INFO: TEST: ciao
INFO: TEST: DONE
`, bytesWriter.String())

	})
	t.Run("All log levels are logged to file", func(t *testing.T) {
		Stdout = ioutil.Discard
		var tsk *Task
		tsk = New("TEST", func(vs *ctx.Context) error {
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
		assert.NoError(t, MustBeEqual(tsk.Status(), DoneOk))

		contentBuff, err := ioutil.ReadFile("TEST.detailed.log")
		assert.NoError(t, err)
		assert.Equal(t,
			`INFO: TEST: START: A task for tests.
INFO: TEST: ciao
DETAIL: TEST: salve
DEBUG: TEST: urrà
INFO: TEST: DONE
`, string(contentBuff))

	})

	t.Run("Non existent server", func(t *testing.T) {
		Stdout = ioutil.Discard
		var tsk *Task
		tsk = New("TEST", func(vs *ctx.Context) error {
			vs.Link(vpath.New("peppa", "./bad"), vpath.New("peppa", "./bad"))
			tests.Done()
			return nil
		})
		tsk.Description = "A task for tests."
		assert.NotNil(t, tsk)

		tsk.Run()

		tsk.Done.AwaitOne()

		assert.True(t, tsk.Status().IsFailure())
		assert.Equal(t,
			"Link from peppa:./bad to peppa:./bad: connection.FindHost: wrong configuration file \"../fixtures/virt-serv.toml\": unknown host `peppa`",
			tsk.Status().Err.Error(),
		)
	})

	tests.Wait()
}
