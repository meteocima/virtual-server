package ctx

import (
	"io/ioutil"
	"testing"

	"github.com/meteocima/virtual-server/config"
	"github.com/meteocima/virtual-server/testutil"
	"github.com/meteocima/virtual-server/vpath"
	"github.com/stretchr/testify/assert"
)

const sOut = "THIS IS A TEST COMMAND\n"

func TestNew(t *testing.T) {
	err := config.Init(testutil.FixtureDir("virt-serv.toml"))
	assert.NoError(t, err)
	drihmFixt := vpath.New("drihm", "/var/fixtures/")
	ctx := Context{}

	t.Run("Exists", func(t *testing.T) {
		found := ctx.Exists(drihmFixt.Join("drihm"))
		assert.True(t, found)
		assert.NoError(t, ctx.Err)
	})

	t.Run("Run", func(t *testing.T) {
		testcmd := drihmFixt.Join("testcmd")
		process := ctx.Run(testcmd, []string{"/var/fixtures/"})

		out, err := ioutil.ReadAll(process.Stdout())
		assert.NoError(t, err)
		assert.Equal(t, sOut, string(out))

		assert.NoError(t, ctx.Err)

	})

	t.Run("ReadDir", func(t *testing.T) {
		files := ctx.ReadDir(drihmFixt.Join("new-dir"))
		assert.Equal(t, 4, len(files))
		assert.Equal(t, "drihm:/var/fixtures/new-dir/file1.txt", files[0].String())
		assert.NoError(t, ctx.Err)
	})

	t.Run("ReadString", func(t *testing.T) {
		content := ctx.ReadString(drihmFixt.Join("ciao.txt"))
		assert.Equal(t, "ciao\n", content)
		assert.NoError(t, ctx.Err)
	})

	t.Run("Copy", func(t *testing.T) {
		ctx.Copy(
			drihmFixt.Join("ciao.txt"),
			vpath.Local("/tmp/hi"),
		)
		assert.NoError(t, ctx.Err)
		actual := ctx.ReadString(vpath.Local("/tmp/hi"))
		assert.Equal(t, "ciao\n", actual)

		ctx.RmFile(drihmFixt.Join("added"))
		ctx.Err = nil

		ctx.Copy(
			vpath.Local("/tmp/hi"),
			drihmFixt.Join("added"),
		)
		assert.NoError(t, ctx.Err)
		actual2 := ctx.ReadString(drihmFixt.Join("added"))
		assert.Equal(t, "ciao\n", actual2)
	})

	t.Run("MkDir", func(t *testing.T) {
		dir := drihmFixt.Join("created-by-tests")
		ctx.RmDir(dir)
		ctx.Err = nil

		assert.False(t, ctx.Exists(dir))
		ctx.MkDir(dir)

		assert.True(t, ctx.Exists(dir))
		assert.NoError(t, ctx.Err)
	})

	t.Run("RmDir", func(t *testing.T) {
		dir := drihmFixt.Join("removed-by-tests")
		ctx.MkDir(dir)
		ctx.Err = nil

		assert.True(t, ctx.Exists(dir))
		ctx.RmDir(dir)
		assert.False(t, ctx.Exists(dir))

		assert.NoError(t, ctx.Err)
	})

	t.Run("RmFile", func(t *testing.T) {
		file := drihmFixt.Join("file-created-by-tests")
		ctx.WriteString(file, "something")
		ctx.Err = nil

		assert.True(t, ctx.Exists(file))

		ctx.RmFile(file)
		assert.False(t, ctx.Exists(file))

		assert.NoError(t, ctx.Err)
	})

	t.Run("WriteString", func(t *testing.T) {
		file := drihmFixt.Join("file-created-by-tests")
		ctx.RmFile(file)
		ctx.Err = nil
		assert.False(t, ctx.Exists(file))

		ctx.WriteString(file, "something")

		assert.True(t, ctx.Exists(file))
		content := ctx.ReadString(file)
		assert.Equal(t, "something", content)

		ctx.RmFile(file)

		assert.NoError(t, ctx.Err)
	})

	t.Run("Move", func(t *testing.T) {
		file := drihmFixt.Join("tbmoved")
		ctx.WriteString(file, "something")

		ctx.Move(
			file,
			vpath.Local("/tmp/tbmoved"),
		)
		assert.NoError(t, ctx.Err)

		actual := ctx.ReadString(vpath.Local("/tmp/tbmoved"))
		assert.Equal(t, "something", actual)

		ctx.RmFile(vpath.Local("/tmp/tbmoved"))
		assert.False(t, ctx.Exists(file))
		assert.NoError(t, ctx.Err)
	})

	t.Run("Bad hostname or config", func(t *testing.T) {
		ctx.Exists(vpath.New("drum", ""))
		assert.Equal(t, ctx.Err.Error(), "Exists `drum:.`: connection.FindHost: wrong configuration file \"../fixtures/virt-serv.toml\": unknown host `drum`")
	})

}
