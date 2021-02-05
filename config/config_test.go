package config

import (
	"testing"

	"github.com/meteocima/virtual-server/testutil"
	"github.com/stretchr/testify/assert"
)

func TestConfigSSHInit(t *testing.T) {
	err := Init(testutil.FixtureDir("sshconfig-virt-serv.toml"))
	assert.NoError(t, err)

	assert.Equal(t, 13, len(Hosts))
	local, timoteo := Hosts["localhost"], Hosts["timoteo"]
	assert.Equal(t, "localhost", local.Name)
	assert.Equal(t, "timoteo", timoteo.Name)

}

func TestInit(t *testing.T) {
	err := Init(testutil.FixtureDir("virt-serv.toml"))
	assert.NoError(t, err)

	assert.Equal(t, 3, len(Hosts))
	local, drihm, withBck := Hosts["localhost"], Hosts["drihm"], Hosts["withbackup"]
	assert.Equal(t, "localhost", local.Name)

	assert.Equal(t, "drihm", drihm.Name)
	assert.Equal(t, "andrea.parodi", drihm.User)
	assert.Equal(t, "withbackup", withBck.Name)
}
