package tailor

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/meteocima/virtual-server/testutil"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	fakeLog := testutil.FixtureDir("virt-serv.toml")
	//fmt.Println(fakeLog)
	fakeFile, err := os.Open(fakeLog)
	assert.NoError(t, err)
	outbuff := bytes.Buffer{}
	tail := New(fakeFile, &outbuff, 1024)
	assert.NotNil(t, tail)
	errs := tail.Start()
	assert.NotNil(t, errs)
	time.Sleep(100 * time.Millisecond)
	tail.Stop()

	err = <-errs
	assert.NoError(t, err)
	//result := string(outbuff.Bytes())
	expected, err := ioutil.ReadFile(fakeLog)
	assert.NoError(t, err)

	assert.Equal(t, expected, outbuff.Bytes())
}
