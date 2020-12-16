package vpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// using simple path
	p := New("timoteo", "/tmp")
	assert.NotNil(t, p)
	assert.Equal(t, "timoteo", p.Host)
	assert.Equal(t, "/tmp", p.Path)

	// using fmt.Sprintf arguments
	p = New("timoteo", "/tmp/%s/", "check")
	assert.NotNil(t, p)
	assert.Equal(t, "timoteo", p.Host)
	assert.Equal(t, "/tmp/check/", p.Path)

	// nil are fields are resolved by new
	p = New("timoteo", "/tmp/%s/", "check")
	assert.NotNil(t, p)
	assert.Equal(t, "timoteo", p.Host)
	assert.Equal(t, "/tmp/check/", p.Path)
}

func TestResolve(t *testing.T) {
	// create instance from nil
	var pnil *VirtualPath
	assert.Nil(t, pnil)
	assert.Panics(t, func() {
		pnil.resolve()
	})

	p := VirtualPath{"", ""}
	assert.Equal(t, "", p.Host)
	assert.Equal(t, "", p.Path)

	p.resolve()
	assert.Equal(t, "localhost", p.Host)
	assert.Equal(t, ".", p.Path)

	// New call resolve
	p2 := New("", "")
	assert.Equal(t, "localhost", p2.Host)
	assert.Equal(t, ".", p2.Path)

}

func TestString(t *testing.T) {
	p2 := New("timoteo", "/tmp")
	assert.Equal(t, "timoteo:/tmp", p2.String())
	p2 = VirtualPath{}
	assert.Equal(t, "localhost:.", p2.String())
	assert.Equal(t, "", p2.Host)
	assert.Equal(t, "", p2.Path)
}

func TestJoin(t *testing.T) {
	p := New("timoteo", "/tmp")
	p2 := p.Join("%sPath/subdir", "other")
	assert.Equal(t, "timoteo:/tmp/otherPath/subdir", p2.String())

	assert.Equal(t, "localhost:adir", VirtualPath{}.Join("adir").String())
}

func TestDir(t *testing.T) {
	p := New("timoteo", "/tmp/afile.txt")
	assert.Equal(t, "timoteo:/tmp", p.Dir().String())
	assert.Equal(t, "localhost:.", VirtualPath{}.Dir().String())
}

func TestFileName(t *testing.T) {
	p := New("timoteo", "/tmp/afile.txt")
	assert.Equal(t, "afile.txt", p.Filename())
	assert.Equal(t, ".", VirtualPath{}.Filename())
}

func TestExt(t *testing.T) {
	p := New("timoteo", "/tmp/afile.txt")
	assert.Equal(t, "afile.txt.ciao", p.AddExt("ciao").Filename())
	assert.Equal(t, "afile.ciao", p.ReplaceExt("ciao").Filename())

	p2 := New("timoteo", "/tmp/afile")
	assert.Equal(t, "afile.ciao", p2.ReplaceExt("ciao").Filename())
	assert.Equal(t, "afile.ciao", p2.AddExt("ciao").Filename())

	p3 := New("timoteo", ".")
	assert.Equal(t, ".ciao", p3.ReplaceExt("ciao").Filename())
	assert.Equal(t, ".ciao", p3.AddExt("ciao").Filename())

	p4 := VirtualPath{}
	assert.Equal(t, ".ciao", p4.ReplaceExt("ciao").Filename())
	assert.Equal(t, ".ciao", p4.AddExt("ciao").Filename())
}

func TestJoinP(t *testing.T) {
	p1 := New("timoteo", "/tmp/caio")
	p2 := New("localhost", "../other")
	assert.Equal(t, "timoteo:/tmp/other", p1.JoinP(p2).String())

	assert.Equal(t, p1.String(), p1.JoinP(VirtualPath{}).String())
	assert.Equal(t, "localhost:.", VirtualPath{}.JoinP(VirtualPath{}).String())
}

func TestFromS(t *testing.T) {
	p1 := FromS("timoteo:/tmp/caio")

	assert.Equal(t, "timoteo", p1.Host)
	assert.Equal(t, "/tmp/caio", p1.Path)

	p2 := FromS(":/tmp/caio")

	assert.Equal(t, "localhost", p2.Host)
	assert.Equal(t, "/tmp/caio", p2.Path)

	p3 := FromS("/tmp/caio")

	assert.Equal(t, "localhost", p3.Host)
	assert.Equal(t, "/tmp/caio", p3.Path)

	p4 := FromS("timoteo:")

	assert.Equal(t, "timoteo", p4.Host)
	assert.Equal(t, ".", p4.Path)

	p5 := FromS(":")

	assert.Equal(t, "localhost", p5.Host)
	assert.Equal(t, ".", p5.Path)
}
