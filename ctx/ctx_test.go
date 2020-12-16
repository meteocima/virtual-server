package ctx

import (
	"testing"

	"github.com/meteocima/virtual-server/config"
	"github.com/meteocima/virtual-server/testutil"
	"github.com/meteocima/virtual-server/vpath"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := config.Init(testutil.FixtureDir("virt-serv.toml"))
	assert.NoError(t, err)

	ctx := Context{}

	t.Run("Exists", func(t *testing.T) {
		found := ctx.Exists(vpath.New("drihm", "/var/fixtures/drihm"))
		assert.True(t, found)
		assert.NoError(t, ctx.Err)
	})

	assert.PanicsWithValue(
		t,
		"Wrong configuration file \"../fixtures/virt-serv.toml\": unknown host `drum`.",
		func() {
			ctx.Exists(vpath.New("drum", ""))
		},
	)

}
