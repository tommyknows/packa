package brew

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommyknows/packa/pkg/cmd"
	"github.com/tommyknows/packa/test/fake"
)

func TestFilterTaps(t *testing.T) {
	tests := []struct {
		installedTaps, desiredTaps []string
		missingTaps, spareTaps     []string
	}{
		{
			[]string{"test/this", "test/yuhere", "test/that"},
			[]string{"test/that", "test/newtap", "test/this"},
			[]string{"test/newtap"},
			[]string{"test/yuhere"},
		},
		{
			[]string{"test/one", "test/two", "test/three"},
			[]string{"test/one", "test/two", "test/three"},
			nil,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			t.Logf("%v", tt.installedTaps)
			m, s := filterTaps(tt.installedTaps, tt.desiredTaps)
			assert.Equal(t, tt.missingTaps, m)
			assert.Equal(t, tt.spareTaps, s)
		})
	}
}

func TestTapsInstall(t *testing.T) {
	testTaps := Taps{
		{
			Name: "my/tap",
			Full: false,
		},
	}

	// allow executing max. 10 commands
	cmds := make(chan []string, 10)
	cmd.AddGlobalOptions(fake.NoOp(cmds, "homebrew/cask\nhomebrew/core"))
	defer cmd.ResetGlobalOptions()
	err := testTaps.sync()

	assert.NoError(t, err)
	assert.Equal(t, []string{"brew", "tap"}, <-cmds)
	assert.Equal(t, []string{"brew", "tap", "my/tap"}, <-cmds)
	assert.Equal(t, []string{"brew", "untap", "homebrew/cask"}, <-cmds)
}

func TestGetInstalledTaps(t *testing.T) {
	cmds := make(chan []string, 1)
	installedTaps := `homebrew/cask
this/test
another/here
homebrew/core`
	cmd.AddGlobalOptions(fake.NoOp(cmds, installedTaps))
	defer cmd.ResetGlobalOptions()
	taps, err := getInstalledTaps()
	assert.NoError(t, err)
	assert.Equal(t, []string{"homebrew/cask", "this/test", "another/here"}, taps)
	assert.Equal(t, []string{"brew", "tap"}, <-cmds)
}
