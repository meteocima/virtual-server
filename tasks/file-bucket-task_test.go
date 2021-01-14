package tasks

import (
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"github.com/meteocima/virtual-server/config"
	"github.com/meteocima/virtual-server/ctx"
	"github.com/meteocima/virtual-server/testutil"
	"github.com/meteocima/virtual-server/vpath"
	"github.com/stretchr/testify/assert"
)

func TestBucketTask(t *testing.T) {
	err := config.Init(testutil.FixtureDir("virt-serv.toml"))
	assert.NoError(t, err)

	t.Run("Works", func(t *testing.T) {
		//bytesWriter := bytes.Buffer{}
		Stdout = ioutil.Discard

		var tsk1 *Task
		var tsk2 *Task

		tsk1 = New("TEST1", func(vs *ctx.Context) error {
			tsk1.FileProduced.Invoke(TaskFile{vpath.Local("/test1"), 42.1})
			return nil
		})

		tsk2 = New("TEST2", func(vs *ctx.Context) error {
			time.Sleep(4 * time.Millisecond)
			tsk2.FileProduced.Invoke(TaskFile{vpath.Local("/test2"), 42.2})
			return nil
		})

		results := []TaskFile{}
		resultsLock := sync.Mutex{}
		mkPostProcTask := func(file TaskFile) *Task {
			return New("POSTPROC-"+file.Path.Filename(), func(vs *ctx.Context) error {
				resultsLock.Lock()
				results = append(results, file)
				resultsLock.Unlock()
				return nil
			})
		}
		bucket := NewFilesBucketTask(mkPostProcTask, tsk1, tsk2)

		assert.NotNil(t, bucket)
		assert.Equal(t, "tskID", bucket.ID)

		bucket.Run()

		bucket.Done.AwaitOne()

		assert.NoError(t, MustBeEqual(bucket.Status(), DoneOk))
		assert.Equal(t, 2, len(results))
		if results[0].Meta == 42.1 {
			assert.Equal(t, results, []TaskFile{
				{vpath.Local("/test1"), 42.1},
				{vpath.Local("/test2"), 42.2},
			})
		} else {
			assert.Equal(t, results, []TaskFile{
				{vpath.Local("/test2"), 42.2},
				{vpath.Local("/test1"), 42.1},
			})
		}
	})

}
