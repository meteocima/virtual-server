package connection

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/meteocima/virtual-server/vpath"
	"github.com/stretchr/testify/assert"
)

func exists(t *testing.T, conn Connection, file vpath.VirtualPath) bool {
	_, err := conn.Stat(file)

	if os.IsNotExist(err) {
		return false
	}
	assert.NoError(t, err)
	return true
}

func CheckMkDir(conn Connection) func(t *testing.T) {
	return func(t *testing.T) {
		dirPath := vpath.VirtualPath{Path: "/tmp/check-dir"}
		conn.RmDir(dirPath)
		assert.False(t, exists(t, conn, dirPath))

		err := conn.MkDir(dirPath)
		assert.NoError(t, err)

		defer conn.RmDir(dirPath)

		assert.True(t, exists(t, conn, dirPath))

	}
}

func CheckRmDir(conn Connection) func(t *testing.T) {
	return func(t *testing.T) {
		dirPath := vpath.VirtualPath{Path: "/tmp/check-dir"}
		conn.MkDir(dirPath)
		assert.True(t, exists(t, conn, dirPath))

		err := conn.RmDir(dirPath)
		assert.NoError(t, err)
		assert.False(t, exists(t, conn, dirPath))
	}
}

func CheckStat(conn Connection) func(t *testing.T) {
	return func(t *testing.T) {
		info, err := conn.Stat(vpath.VirtualPath{Path: "/tmp"})
		assert.NoError(t, err)
		assert.Equal(t, "tmp", info.Name())

		info, err = conn.Stat(vpath.VirtualPath{Path: "/timpa/tompa"})
		assert.Error(t, err)
		assert.Equal(t, nil, info)
	}
}

func CheckOpenReader(conn Connection) func(t *testing.T) {
	return func(t *testing.T) {
		reader, err := conn.OpenReader(vpath.VirtualPath{Path: "/var/fixtures/ciao.txt"})
		assert.NoError(t, err)
		assert.NotNil(t, reader)

		buf, err := ioutil.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, "ciao\n", string(buf))
		err = reader.Close()
		assert.NoError(t, err)
	}
}

func writeFile(t *testing.T, conn Connection, file vpath.VirtualPath, value string) {
	writer, err := conn.OpenWriter(file)
	assert.NoError(t, err)
	_, err = writer.Write([]byte(value))
	assert.NoError(t, err)
	writer.Close()
}

func CheckReadDir(conn Connection) func(t *testing.T) {
	return func(t *testing.T) {
		dir := vpath.VirtualPath{Path: "/var/fixtures/new-dir"}
		conn.RmDir(dir)
		defer conn.RmDir(dir)
		err := conn.MkDir(dir)
		assert.NoError(t, err)

		writeFile(t, conn, dir.Join("file1.txt"), "1")
		writeFile(t, conn, dir.Join("file2.txt"), "2")
		writeFile(t, conn, dir.Join("file3.txt"), "3")
		writeFile(t, conn, dir.Join("file4.txt"), "4")

		files, err := conn.ReadDir(dir)
		assert.NoError(t, err)

		assert.Equal(t, vpath.VirtualPathList{
			dir.Join("file1.txt"),
			dir.Join("file2.txt"),
			dir.Join("file3.txt"),
			dir.Join("file4.txt"),
		}, files)
	}
}

func CheckRmFile(conn Connection) func(t *testing.T) {
	return func(t *testing.T) {
		file := vpath.VirtualPath{Path: "/var/fixtures/somefile"}
		writeFile(t, conn, file, "a test line")

		assert.True(t, exists(t, conn, file))

		err := conn.RmFile(file)
		assert.NoError(t, err)

		assert.False(t, exists(t, conn, file))
	}
}

func CheckOpenWriter(conn Connection) func(t *testing.T) {
	return func(t *testing.T) {
		dir := vpath.VirtualPath{Path: "/var/fixtures/tmp"}
		conn.RmDir(dir)
		defer conn.RmDir(dir)
		err := conn.MkDir(dir)
		assert.NoError(t, err)

		writer, err := conn.OpenWriter(dir.Join("afile"))
		assert.NoError(t, err)
		_, err = writer.Write([]byte("a test line"))
		assert.NoError(t, err)
		writer.Close()

		reader, err := conn.OpenReader(dir.Join("afile"))
		assert.NoError(t, err)
		assert.NotNil(t, reader)

		buf, err := ioutil.ReadAll(reader)
		assert.NoError(t, err)
		assert.Equal(t, "a test line", string(buf))
		err = reader.Close()
		assert.NoError(t, err)
	}
}

func CheckRun(conn Connection) func(t *testing.T) {
	return func(t *testing.T) {
		fixtures := vpath.VirtualPath{Path: "/var/fixtures/", Host: conn.HostName()}
		process, err := conn.Run(fixtures.Join("testcmd"), []string{"/var/fixtures/"})

		assert.NotNil(t, process)
		assert.NoError(t, err)

		outReader := process.Stdout()
		assert.NotNil(t, outReader)
		out, err := ioutil.ReadAll(outReader)
		assert.NoError(t, err)
		s := string(out)
		fmt.Println(s)
		assert.Equal(t, "THIS IS A TEST COMMAND\nTHIS IS AN ERROR COMMAND\n", s)
	}
}

func DoAllChecks(t *testing.T, conn Connection) {
	//t.Run("CheckStat", CheckStat(conn))
	//t.Run("CheckMkDir", CheckMkDir(conn))
	//t.Run("CheckRmDir", CheckRmDir(conn))
	//t.Run("CheckOpenReader", CheckOpenReader(conn))
	//t.Run("CheckOpenWriter", CheckOpenWriter(conn))
	//t.Run("CheckRmFile", CheckRmFile(conn))
	//t.Run("CheckReadDir", CheckReadDir(conn))
	t.Run("CheckRun", CheckRun(conn))
}

func TestLocalHost(t *testing.T) {
	osConn := LocalConnection{}
	err := osConn.Open()
	assert.NoError(t, err)
	DoAllChecks(t, &osConn)
	t.Run("CheckStat", CheckStat(&osConn))
	assert.NoError(t, osConn.Close())

}

/*
func TestSSH(t *testing.T) {
	conn := SSHConnection{
		Host:    "drihm",
		Port:    22,
		User:    "andrea.parodi",
		KeyPath: "/var/fixtures/private-key",
	}

	err := conn.Open()
	assert.NoError(t, err)
	DoAllChecks(t, &conn)
	assert.NoError(t, conn.Close())

}
*/
