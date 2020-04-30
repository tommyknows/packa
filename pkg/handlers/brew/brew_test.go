package brew

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"github.com/tommyknows/packa/pkg/cmd"
	"github.com/tommyknows/packa/pkg/output"
	"github.com/tommyknows/packa/test/fake"
)

func TestParse(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name string
		cask bool
		out  formula
		err  bool
	}{
		{
			name: "username/repo/vim",
			out: formula{
				Name: "vim",
				Tap:  "username/repo",
			},
			err: false,
		},
		{
			name: "vim@8.1.0",
			out: formula{
				Name:    "vim",
				Version: "8.1.0",
			},
			err: false,
		},
		{
			name: "username/repo/vim@8.1.0",
			cask: true,
			out: formula{
				Name:    "vim",
				Tap:     "username/repo",
				Version: "8.1.0",
				Cask:    true,
			},
			err: false,
		},
		{
			name: "vim@8@8",
			out:  formula{},
			err:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := parse(tt.name, tt.cask)
			is.Equal(tt.out, p)
			is.Equal(tt.err, err != nil)
		})
	}
}

var (
	isCask    = true
	isNotCask = false
)

func TestInstall(t *testing.T) {
	is := is.New(t)

	// redirect the output logs
	var buf bytes.Buffer
	output.Set(&buf, &buf)
	defer buf.Reset()
	b := Handler{
		Formulae: []formula{
			{
				Name: "somepackage",
			},
		},
		cask: &isNotCask,
	}
	afterInstall := []formula{
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
	is.NoErr(err)

	c := make(chan []string, 20)
	cmd.AddGlobalOptions(fake.NoOp(c, "someoutput"))
	defer cmd.ResetGlobalOptions()
	list, err := b.Install("thispackage", "pkg@version", "from/tap/betterpkg", "this/tap/another@0.0.1", "somepackage@newer")
	is.NoErr(err)
	is.Equal(afterInstJSON, []byte(*list))

	close(c)
	var executedCommands [][]string
	for execedCmd := range c {
		executedCommands = append(executedCommands, execedCmd)
	}

	is.Equal(executedCommands, [][]string{
		{"brew", "install", "thispackage"},
		{"brew", "install", "pkg@version"},
		{"brew", "pin", "pkg"},
		{"brew", "install", "from/tap/betterpkg"},
		{"brew", "install", "this/tap/another@0.0.1"},
		{"brew", "pin", "this/tap/another"},
		{"brew", "install", "somepackage@newer"},
		{"brew", "pin", "somepackage"},
	})
}

func TestUninstall(t *testing.T) {
	is := is.New(t)

	// redirect the output logs
	var buf bytes.Buffer
	output.Set(&buf, &buf)
	defer buf.Reset()

	b := Handler{
		Formulae: []formula{
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
		cask: &isNotCask,
	}
	afterUninstall := []formula{
		{Name: "thispackage"},
		{
			Name:    "pkg",
			Version: "version",
		},
	}
	afterRmJSON, err := json.Marshal(afterUninstall)
	is.NoErr(err)

	c := make(chan []string, 20)
	cmd.AddGlobalOptions(fake.NoOp(c, "someoutput"))
	defer cmd.ResetGlobalOptions()
	list, err := b.Remove("somepackage@newer", "from/tap/betterpkg")
	is.NoErr(err)
	is.Equal(afterRmJSON, []byte(*list))

	close(c)
	var executedCommands [][]string
	for execedCmd := range c {
		executedCommands = append(executedCommands, execedCmd)
	}

	is.Equal(executedCommands, [][]string{
		{"brew", "uninstall", "somepackage"},
		{"brew", "uninstall", "from/tap/betterpkg"},
	})
}

func TestUpgrade(t *testing.T) {
	is := is.New(t)

	// redirect the output logs
	var buf bytes.Buffer
	output.Set(&buf, &buf)
	defer buf.Reset()

	b := Handler{
		Formulae: []formula{
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
		cask: &isNotCask,
	}
	afterUpgrade := []formula{
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
	is.NoErr(err)

	c := make(chan []string, 20)
	cmd.AddGlobalOptions(fake.NoOp(c, "someoutput"))
	defer cmd.ResetGlobalOptions()
	list, err := b.Upgrade("somepackage@evennewer", "from/tap/betterpkg")
	is.NoErr(err)
	is.Equal(afterUpJSON, []byte(*list))

	close(c)
	var executedCommands [][]string
	for execedCmd := range c {
		executedCommands = append(executedCommands, execedCmd)
	}

	is.Equal(executedCommands, [][]string{
		{"brew", "unpin", "somepackage"},
		{"brew", "upgrade", "somepackage@evennewer"},
		{"brew", "pin", "somepackage"},
		{"brew", "upgrade", "from/tap/betterpkg"},
	})

	b = Handler{
		Formulae: []formula{
			{
				Name:    "somepackage",
				Version: "newer",
			},
			{Name: "thispackage"},
		},
	}

	afterUpgrade = []formula{
		{
			Name:    "somepackage",
			Version: "newer",
		},
		{Name: "thispackage"},
	}
	afterUpJSON, err = json.Marshal(afterUpgrade)
	is.NoErr(err)

	c = make(chan []string, 20)
	cmd.ResetGlobalOptions()
	cmd.AddGlobalOptions(fake.NoOp(c, "someoutput"))
	defer cmd.ResetGlobalOptions()
	list, err = b.Upgrade()
	is.NoErr(err)
	is.Equal(afterUpJSON, []byte(*list))

	close(c)
	executedCommands = nil
	for execedCmd := range c {
		executedCommands = append(executedCommands, execedCmd)
	}

	is.Equal(executedCommands, [][]string{{"brew", "upgrade", "thispackage"}})
}

func TestCaskInstall(t *testing.T) {
	is := is.New(t)

	// redirect the output logs
	var buf bytes.Buffer
	output.Set(&buf, &buf)
	defer buf.Reset()
	b := Handler{
		Formulae: []formula{
			{
				Name: "somepackage",
			},
		},
		cask: &isCask,
	}
	afterInstall := []formula{
		{
			Name:    "somepackage",
			Version: "newer",
			Cask:    true,
		},
		{
			Name: "thispackage",
			Cask: true,
		},
		{
			Name:    "pkg",
			Version: "version",
			Cask:    true,
		},
		{
			Name: "betterpkg",
			Tap:  "from/tap",
			Cask: true,
		},
		{
			Name:    "another",
			Tap:     "this/tap",
			Version: "0.0.1",
			Cask:    true,
		},
	}
	afterInstJSON, err := json.Marshal(afterInstall)
	is.NoErr(err)

	c := make(chan []string, 20)
	cmd.AddGlobalOptions(fake.NoOp(c, "someoutput"))
	defer cmd.ResetGlobalOptions()
	list, err := b.Install("thispackage", "pkg@version", "from/tap/betterpkg", "this/tap/another@0.0.1", "somepackage@newer")
	is.NoErr(err)
	is.Equal(afterInstJSON, []byte(*list))

	close(c)
	var executedCommands [][]string
	for execedCmd := range c {
		executedCommands = append(executedCommands, execedCmd)
	}

	is.Equal(executedCommands, [][]string{
		{"brew", "cask", "install", "thispackage"},
		{"brew", "cask", "install", "pkg@version"},
		{"brew", "pin", "pkg"},
		{"brew", "cask", "install", "from/tap/betterpkg"},
		{"brew", "cask", "install", "this/tap/another@0.0.1"},
		{"brew", "pin", "this/tap/another"},
		{"brew", "cask", "install", "somepackage@newer"},
		{"brew", "pin", "somepackage"},
	})
}
