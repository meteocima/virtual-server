package tasks

import (
	"bytes"
	"testing"

	"github.com/meteocima/virtual-server/ctx"
	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	done := make(chan struct{})
	t.Run("Info log and below are logged to stdout", func(t *testing.T) {
		bytesWriter := bytes.Buffer{}
		Stdout = &bytesWriter
		tsk := New("TEST", func(tsk *Task, vs *ctx.Context) error {
			vs.LogInfo("ciao")
			close(done)
			return nil
		})
		tsk.Description = "A task for tests."
		assert.NotNil(t, tsk)

		tsk.Run()

		tsk.Done.AwaitOne()
		assert.Equal(t,
			`INFO: START: TEST: A task for tests.
INFO: ciao
INFO: DONE: TEST
`, bytesWriter.String())

	})
	<-done
}
