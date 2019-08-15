package brew

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
	"github.com/tommyknows/packa/pkg/cmd"
	"github.com/tommyknows/packa/test/fake"
)

func TestFilterTaps(t *testing.T) {
	is := is.New(t)

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

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			t.Logf("%v", tt.installedTaps)
			m, s := filterTaps(tt.installedTaps, tt.desiredTaps)
			is.Equal(tt.missingTaps, m)
			is.Equal(tt.spareTaps, s)
		})
	}
}

func TestTapsInstall(t *testing.T) {
	is := is.New(t)

	testTaps := taps{
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
	is.NoErr(err)
	is.Equal([]string{"brew", "tap"}, <-cmds)                    // first command should be brew tap
	is.Equal([]string{"brew", "tap", "my/tap"}, <-cmds)          // second command should add the my/tap tap
	is.Equal([]string{"brew", "untap", "homebrew/cask"}, <-cmds) // third should untap homebrew/cask
}

func TestGetInstalledTaps(t *testing.T) {
	is := is.New(t)

	cmds := make(chan []string, 1)
	installedTaps := `homebrew/cask
this/test
another/here
homebrew/core
`
	cmd.AddGlobalOptions(fake.NoOp(cmds, installedTaps))
	defer cmd.ResetGlobalOptions()
	taps, err := getInstalledTaps()
	is.NoErr(err)
	is.Equal([]string{"homebrew/cask", "this/test", "another/here"}, taps)
	is.Equal([]string{"brew", "tap"}, <-cmds)
}
