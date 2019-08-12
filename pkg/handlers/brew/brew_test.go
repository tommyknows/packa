package brew

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommyknows/packa/pkg/cmd"
	"github.com/tommyknows/packa/test/fake"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in  string
		out Formula
		err bool
	}{
		{
			in: "username/repo/vim",
			out: Formula{
				Name: "vim",
				Tap:  "username/repo",
			},
			err: false,
		},
		{
			in: "vim@8.1.0",
			out: Formula{
				Name:    "vim",
				Version: "8.1.0",
			},
			err: false,
		},
		{
			in: "username/repo/vim@8.1.0",
			out: Formula{
				Name:    "vim",
				Tap:     "username/repo",
				Version: "8.1.0",
			},
			err: false,
		},
		{
			in:  "vim@8@8",
			out: Formula{},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			p, err := parse(tt.in)
			assert.Equal(t, tt.out, p)
			assert.Equal(t, tt.err, err != nil)
		})
	}
}

func TestInstall(t *testing.T) {
	b := Handler{
		Formulae: []Formula{
			{
				Name: "somepackage",
			},
		},
	}
	afterInstall := []Formula{
		{
			Name:    "somepackage",
			Version: "newer",
		},
		{Name: "thispackage"},
		{
			Name:    "pkg",
			Version: "version",
		},
		{
			Name: "betterpkg",
			Tap:  "from/tap",
		},
		{
			Name:    "another",
			Tap:     "this/tap",
			Version: "0.0.1",
		},
	}
	afterInstJSON, err := json.Marshal(afterInstall)
	assert.NoError(t, err)

	c := make(chan []string, 20)
	cmd.AddGlobalOptions(fake.NoOp(c, "someoutput"))
	defer cmd.ResetGlobalOptions()
	list, err := b.Install("thispackage", "pkg@version", "from/tap/betterpkg", "this/tap/another@0.0.1", "somepackage@newer")

	close(c)
	var executedCommands [][]string
	for execedCmd := range c {
		executedCommands = append(executedCommands, execedCmd)
	}

	assert.Equal(t, afterInstJSON, []byte(*list))
	assert.NoError(t, err)
	assert.Contains(t, executedCommands, []string{"brew", "install", "thispackage"})
	assert.Contains(t, executedCommands, []string{"brew", "install", "this/tap/another@0.0.1"})
	assert.Contains(t, executedCommands, []string{"brew", "pin", "this/tap/another"})
	assert.Contains(t, executedCommands, []string{"brew", "install", "pkg@version"})
	assert.Contains(t, executedCommands, []string{"brew", "pin", "pkg"})
	assert.Contains(t, executedCommands, []string{"brew", "install", "from/tap/betterpkg"})
	assert.Contains(t, executedCommands, []string{"brew", "install", "somepackage@newer"})
	assert.Contains(t, executedCommands, []string{"brew", "pin", "somepackage"})
	assert.Len(t, executedCommands, 8)
}

func TestRemove(t *testing.T) {
	b := Handler{
		Formulae: []Formula{
			{
				Name:    "somepackage",
				Version: "newer",
			},
			{Name: "thispackage"},
			{
				Name:    "pkg",
				Version: "version",
			},
			{
				Name: "betterpkg",
				Tap:  "from/tap",
			},
		},
	}
	afterRemove := []Formula{
		{Name: "thispackage"},
		{
			Name:    "pkg",
			Version: "version",
		},
	}
	afterRmJSON, err := json.Marshal(afterRemove)
	assert.NoError(t, err)

	c := make(chan []string, 20)
	cmd.AddGlobalOptions(fake.NoOp(c, "someoutput"))
	defer cmd.ResetGlobalOptions()
	list, err := b.Remove("somepackage@newer", "from/tap/betterpkg")

	close(c)
	var executedCommands [][]string
	for execedCmd := range c {
		executedCommands = append(executedCommands, execedCmd)
	}

	assert.Equal(t, afterRmJSON, []byte(*list))
	assert.NoError(t, err)
	assert.Contains(t, executedCommands, []string{"brew", "remove", "somepackage"})
	assert.Contains(t, executedCommands, []string{"brew", "remove", "from/tap/betterpkg"})
	assert.Len(t, executedCommands, 2)
}

func TestUpgrade(t *testing.T) {
	b := Handler{
		Formulae: []Formula{
			{
				Name:    "somepackage",
				Version: "newer",
			},
			{Name: "thispackage"},
			{
				Name:    "pkg",
				Version: "version",
			},
			{
				Name: "betterpkg",
				Tap:  "from/tap",
			},
		},
	}
	afterUpgrade := []Formula{
		{
			Name:    "somepackage",
			Version: "evennewer",
		},
		{Name: "thispackage"},
		{
			Name:    "pkg",
			Version: "version",
		},
		{
			Name: "betterpkg",
			Tap:  "from/tap",
		},
	}
	afterUpJSON, err := json.Marshal(afterUpgrade)
	assert.NoError(t, err)

	c := make(chan []string, 20)
	cmd.AddGlobalOptions(fake.NoOp(c, "someoutput"))
	defer cmd.ResetGlobalOptions()
	list, err := b.Upgrade("somepackage@evennewer", "from/tap/betterpkg")

	close(c)
	var executedCommands [][]string
	for execedCmd := range c {
		executedCommands = append(executedCommands, execedCmd)
	}

	assert.Equal(t, afterUpJSON, []byte(*list))
	assert.NoError(t, err)
	assert.Contains(t, executedCommands, []string{"brew", "unpin", "somepackage"})
	assert.Contains(t, executedCommands, []string{"brew", "upgrade", "somepackage@evennewer"})
	assert.Contains(t, executedCommands, []string{"brew", "pin", "somepackage"})
	assert.Contains(t, executedCommands, []string{"brew", "upgrade", "from/tap/betterpkg"})
	assert.Len(t, executedCommands, 4)

	b = Handler{
		Formulae: []Formula{
			{
				Name:    "somepackage",
				Version: "newer",
			},
			{Name: "thispackage"},
		},
	}

	afterUpgrade = []Formula{
		{
			Name:    "somepackage",
			Version: "newer",
		},
		{Name: "thispackage"},
	}
	afterUpJSON, err = json.Marshal(afterUpgrade)
	assert.NoError(t, err)

	c = make(chan []string, 20)
	cmd.ResetGlobalOptions()
	cmd.AddGlobalOptions(fake.NoOp(c, "someoutput"))
	defer cmd.ResetGlobalOptions()
	list, err = b.Upgrade()

	close(c)
	executedCommands = nil
	for execedCmd := range c {
		executedCommands = append(executedCommands, execedCmd)
	}

	assert.Equal(t, afterUpJSON, []byte(*list))
	assert.NoError(t, err)
	assert.Contains(t, executedCommands, []string{"brew", "upgrade", "thispackage"})
	assert.Len(t, executedCommands, 1)
}
